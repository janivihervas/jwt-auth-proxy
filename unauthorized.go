package authproxy

import (
	"context"
	"encoding/json"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/janivihervas/authproxy/internal/random"
)

func (m *Middleware) unauthorized(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	m.clearSessionAndAccessToken(ctx, w, r)

	session, err := m.getSession(ctx, w, r, true)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if session.AuthRequestState == "" {
		state := random.String(32)
		session.AuthRequestState = state

		err = m.sessionStorage.Save(ctx, session.ID, session)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	redirectURL := m.authClient.AuthCodeURL(session.AuthRequestState, oauth2.AccessTypeOffline)
	//redirectURL := m.authClient.AuthenticationRequestURL(session.AuthRequestState, oidc.ResponseModeFormPost)
	if !m.redirect(r) {
		m.unauthorizedResponse(ctx, w, redirectURL)
		return
	}

	session.OriginalURL = r.URL.String()
	err = m.sessionStorage.Save(ctx, session.ID, session)
	if err != nil {
		// log err
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

type unauthorizedResponse struct {
	StatusCode int    `json:"statusCode"`
	RedirectTo string `json:"redirectTo"`
}

func (m *Middleware) unauthorizedResponse(ctx context.Context, w http.ResponseWriter, redirectURL string) {
	resp := unauthorizedResponse{
		StatusCode: http.StatusUnauthorized,
		RedirectTo: redirectURL,
	}
	w.WriteHeader(http.StatusUnauthorized)

	b, err := json.Marshal(resp)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		_, err = w.Write([]byte(redirectURL))
		if err != nil {
			// log err
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		// log err
	}
}
