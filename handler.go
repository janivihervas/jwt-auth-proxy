package authproxy

import (
	"context"
	"encoding/json"
	"fmt"
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
		m.Logger.Printf("couldn't validate access token, unauthorized: %+v", validationErr)
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
		m.Logger.Printf("couldn't type case session or session is empty")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	authRequestState := random.String(32)
	state.AuthRequestState = authRequestState

	session.Values[sessionName] = state
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
	session.Values[sessionName] = state
	err = m.SessionStore.Save(r, w, session)
	if err != nil {
		m.Logger.Printf("couldn't save session: %+v", err)
	}

	m.Logger.Println("redirecting to login page")
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
	if r.Method != http.MethodGet {
		m.Logger.Printf("received non-GET request to authorize callback: %s", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var (
		ctx      = r.Context()
		code     = r.URL.Query().Get("code")
		stateNew = r.URL.Query().Get("state")
	)

	if code == "" {
		m.Logger.Println("no code in authorize callback")
		http.Error(w, "no code in request query", http.StatusBadRequest)
		return
	}

	if stateNew == "" {
		m.Logger.Println("no state in authorize callback")
		http.Error(w, "no state in request query", http.StatusBadRequest)
		return
	}

	session, err := m.SessionStore.Get(r, sessionName)
	if err != nil {
		m.Logger.Printf("couldn't get session: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	state, ok := session.Values[sessionName].(State)
	if !ok {
		m.Logger.Println("couldn't type case session or session is empty")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	stateOld := state.AuthRequestState

	// Clean previous state from session
	state.AuthRequestState = ""

	session.Values[sessionName] = state
	err = m.SessionStore.Save(r, w, session)
	if err != nil {
		m.Logger.Printf("couldn't save session: %+v", err)
	}

	if stateNew != stateOld {
		http.Error(w, fmt.Sprintf("states are not the same: %s != %s", stateNew, stateOld), http.StatusBadRequest)
		return
	}

	tokens, err := m.AuthClient.Exchange(ctx, code, oauth2.AccessTypeOffline)
	if err != nil {
		m.Logger.Printf("couldn't exchange authorization code for tokens: %+v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	_, err = jwt.ParseAccessToken(ctx, tokens.AccessToken)
	if err != nil {
		m.Logger.Println("access token from exchange was invalid")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	state.AccessToken = tokens.AccessToken
	http.SetCookie(w, createAccessTokenCookie(tokens.AccessToken))

	if tokens.RefreshToken != "" {
		state.RefreshToken = tokens.RefreshToken
	}

	session.Values[sessionName] = state
	err = m.SessionStore.Save(r, w, session)
	if err != nil {
		m.Logger.Printf("couldn't save session: %+v", err)
	}

	url := state.OriginalURL

	// Clear previous redirect url from session
	state.OriginalURL = ""

	session.Values[sessionName] = state
	err = m.SessionStore.Save(r, w, session)
	if err != nil {
		m.Logger.Printf("couldn't save session: %+v", err)
	}

	if url == "" {
		url = "/"
	}

	m.Logger.Printf("redirecting to %s", url)
	http.Redirect(w, r, url, http.StatusSeeOther)
}
