package authproxy

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/janivihervas/authproxy/internal/random"
	"golang.org/x/oauth2"

	"github.com/janivihervas/authproxy/jwt"
)

func (m *Middleware) defaultHandler(w http.ResponseWriter, r *http.Request) {
	for _, regexp := range m.skipAuthenticationRegex {
		if regexp.MatchString(r.URL.Path) {
			m.Next.ServeHTTP(w, r)
			return
		}
	}

	var (
		accessTokenStr string
		ctx            = r.Context()
	)

	accessTokenStr, err := m.getAccessToken(ctx, r)
	if err != nil {
		accessTokenStr, err = m.refreshAccessToken(ctx, w, r)
		if err != nil {
			m.Logger.Printf("Couldn't refresh access token: %+v", err)
		}
	}

	var validationErr error
	for i := 0; i < 3; i++ {
		_, validationErr = jwt.ParseAccessToken(ctx, accessTokenStr)
		if validationErr == jwt.ErrTokenExpired {
			accessTokenStr, err = m.refreshAccessToken(ctx, w, r)
			if err != nil {
				m.Logger.Printf("Couldn't refresh access token: %+v", err)
			}
			continue
		}
		break
	}

	if validationErr != nil {
		m.unauthorized(ctx, w, r)
		return
	}

	m.Next.ServeHTTP(w, r)
}

func (m *Middleware) unauthorized(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	err := m.clearSessionAndAccessToken(ctx, w, r)
	if err != nil {
		m.Logger.Printf("%+v", err)
	}

	ctx = r.Context()

	session, err := m.SessionStore.Get(r, sessionName)
	if err != nil {
		m.Logger.Printf("couldn't get session: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	state, ok := session.Values[sessionName].(State)
	if !ok {
		m.Logger.Printf("couldn't type cast session")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	authRequestState := random.String(32)
	state.AuthRequestState = authRequestState

	err = m.SessionStore.Save(r, w, session)
	if err != nil {
		m.Logger.Printf("couldn't save session: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	redirectURL := m.AuthClient.AuthCodeURL(state.AuthRequestState, oauth2.AccessTypeOffline)
	for _, regexp := range m.skipRedirectToLoginRegex {
		if regexp.MatchString(r.URL.Path) {
			m.unauthorizedResponse(ctx, w, redirectURL)
			return
		}
	}

	state.OriginalURL = r.URL.String()
	err = m.SessionStore.Save(r, w, session)
	if err != nil {
		m.Logger.Printf("couldn't save session: %+v", err)
	}

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

type unauthorizedResponse struct {
	StatusCode  int    `json:"statusCode"`
	RedirectURL string `json:"redirectURL"`
}

func (m *Middleware) unauthorizedResponse(ctx context.Context, w http.ResponseWriter, redirectURL string) {
	resp := unauthorizedResponse{
		StatusCode:  http.StatusUnauthorized,
		RedirectURL: redirectURL,
	}
	w.WriteHeader(http.StatusUnauthorized)

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	err := encoder.Encode(resp)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		_, err = w.Write([]byte(redirectURL))
		if err != nil {
			m.Logger.Printf("falling back to plain text response failed: %+v", err)
		}
	}
}

func (m *Middleware) authorizeCallback(w http.ResponseWriter, r *http.Request) {
	//if r.Method != http.MethodGet {
	//	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	//	return
	//}
	//
	//var (
	//	ctx      = r.Context()
	//	code     = r.URL.Query().Get("code")
	//	stateNew = r.URL.Query().Get("state")
	//)
	//
	//if code == "" {
	//	http.Error(w, "no code in response", http.StatusBadRequest)
	//	return
	//}
	//
	//if stateNew == "" {
	//	http.Error(w, "no state in response", http.StatusBadRequest)
	//	return
	//}
	//
	//state, err := m.getStateFromContext(ctx)
	//if err != nil {
	//	m.Logger.Printf("%+v", err)
	//	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	//	return
	//}
	//
	//stateOld := state.AuthRequestState
	//
	//// Clean previous state from session
	//state.AuthRequestState = ""
	//err = m.SessionStore.Save(ctx, state.ID, state)
	//if err != nil {
	//	m.Logger.Printf("%+v", err)
	//	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	//	return
	//}
	//
	//if stateNew != stateOld {
	//	http.Error(w, fmt.Sprintf("states are not the same: %s != %s", stateNew, stateOld), http.StatusBadRequest)
	//	return
	//}
	//
	//tokens, err := m.AuthClient.Exchange(ctx, code, oauth2.AccessTypeOffline)
	//if err != nil {
	//	m.Logger.Printf("%+v", err)
	//	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	//	return
	//}
	//
	//state.RefreshToken = tokens.RefreshToken
	//err = m.SessionStore.Save(ctx, state.ID, state)
	//if err != nil {
	//	m.Logger.Printf("%+v", err)
	//	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	//	return
	//}
	//
	//http.SetCookie(w, createAccessTokenCookie(tokens.AccessToken))
	//
	//url := state.OriginalURL
	//
	//// Clear previous redirect url from session
	//state.OriginalURL = ""
	//err = m.SessionStore.Save(ctx, state.ID, state)
	//if err != nil {
	//	m.Logger.Printf("%+v", err)
	//	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	//}
	//
	//if url == "" {
	//	url = "/"
	//}
	//
	//http.Redirect(w, r, url, http.StatusSeeOther)
}
