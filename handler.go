package authproxy

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/janivihervas/authproxy/internal/random"
	"golang.org/x/oauth2"

	"github.com/janivihervas/authproxy/jwt"
)

func (m *Middleware) defaultHandler(w http.ResponseWriter, r *http.Request) {
	for _, regexp := range m.skipAuthenticationRegex {
		if regexp.MatchString(r.URL.Path) {
			m.Logger.Printf("path %s matched regexp %s, skipping authentication", r.URL.Path, regexp.String())
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
		accessTokenStr, err = m.refreshAccessToken(ctx, w)
		if err != nil {
			m.Logger.Printf("couldn't refresh access token: %+v", err)
		}
	}

	var validationErr error
	for i := 0; i < 3; i++ {
		_, validationErr = jwt.ParseAccessToken(ctx, accessTokenStr)
		if validationErr == jwt.ErrTokenExpired {
			accessTokenStr, err = m.refreshAccessToken(ctx, w)
			if err != nil {
				m.Logger.Printf("couldn't refresh access token, previous access token was expired: %+v", err)
			}
			continue
		}
		break
	}

	if validationErr != nil {
		m.Logger.Printf("couldn't validate access token, move to unauthorized flow: %+v", validationErr)
		m.unauthorized(ctx, w, r)
		return
	}

	m.Next.ServeHTTP(w, r)
}

func (m *Middleware) unauthorized(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	err := m.clearSessionAndAccessToken(ctx, w, r)
	if err != nil {
		m.Logger.Printf("couldn't clear session and/or access token: %+v", err)
	}

	state, err := getStateFromContext(ctx)
	if err != nil {
		m.Logger.Printf("couldn't get session from context: %+v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	requestState := random.String(32)
	state.AuthRequestStates[requestState] = authRequestState{
		ExpiresAt:   time.Now().Add(authRequestStateExpiration),
		OriginalURL: r.URL.String(),
	}

	opts := []oauth2.AuthCodeOption{
		oauth2.AccessTypeOffline,
	}
	opts = append(opts, m.AdditionalAuthURLParameters...)
	redirectURL := m.AuthClient.AuthCodeURL(requestState, opts...)
	for _, regexp := range m.skipRedirectToLoginRegex {
		if regexp.MatchString(r.URL.Path) {
			m.Logger.Printf("path %s matched regexp %s, skipping redirection to login page", r.URL.Path, regexp.String())
			m.unauthorizedResponse(w, redirectURL)
			return
		}
	}

	m.Logger.Println("redirecting to login page")
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

type unauthorizedResponse struct {
	StatusCode  int    `json:"statusCode"`
	RedirectURL string `json:"redirectUrl"`
}

func (m *Middleware) unauthorizedResponse(w http.ResponseWriter, redirectURL string) {
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
		m.Logger.Printf("writing json for unauthorized response failed, falling back to plain text: %+v", err)
		w.Header().Set("Content-Type", "text/plain")
		_, err = w.Write([]byte(redirectURL))
		if err != nil {
			m.Logger.Printf("writing plain text for unauthorized response failed: %+v", err)
		}
	}
}

// nolint:funlen
func (m *Middleware) authorizeCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		m.Logger.Printf("authorize callback: received non-GET request: %s", r.Method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var (
		ctx      = r.Context()
		code     = r.URL.Query().Get("code")
		stateNew = r.URL.Query().Get("state")
		now      = time.Now()
	)

	if code == "" {
		m.Logger.Println("authorize callback: no code")
		http.Error(w, "no code in request query", http.StatusBadRequest)
		return
	}

	if stateNew == "" {
		m.Logger.Println("authorize callback: no state")
		http.Error(w, "no state in request query", http.StatusBadRequest)
		return
	}

	state, err := getStateFromContext(ctx)
	if err != nil {
		m.Logger.Printf("authorize callback: %+v", err)
		http.Error(w, "no session in request", http.StatusInternalServerError)
		return
	}

	// nolint:govet
	authRequestState, found := state.AuthRequestStates[stateNew]
	if !found {
		m.Logger.Println("authorize callback: unknown authorization state", stateNew)
		http.Error(w, "unknown authorization state", http.StatusBadRequest)
		return
	}

	// Clean previous state from session
	delete(state.AuthRequestStates, stateNew)
	state.clearExpiredStates()

	if now.After(authRequestState.ExpiresAt) {
		m.Logger.Printf("authorize callback: authorization state was expired at %v", authRequestState.ExpiresAt)
		http.Error(w, "authorization state was expired", http.StatusBadRequest)
		return
	}

	opts := []oauth2.AuthCodeOption{
		oauth2.AccessTypeOffline,
	}
	opts = append(opts, m.AdditionalAuthURLParameters...)
	tokens, err := m.AuthClient.Exchange(ctx, code, opts...)
	if err != nil {
		m.Logger.Printf("couldn't exchange authorization code for tokens: %+v", err)
		http.Error(w, "couldn't exchange authorization code for tokens", http.StatusBadRequest)
		return
	}

	_, err = jwt.ParseAccessToken(ctx, tokens.AccessToken)
	if err != nil {
		m.Logger.Println("access token from exchange was invalid")
		http.Error(w, "access token from exchange was invalid", http.StatusBadRequest)
		return
	}

	state.AccessToken = tokens.AccessToken
	http.SetCookie(w, createAccessTokenCookie(tokens.AccessToken))

	if tokens.RefreshToken != "" {
		state.RefreshToken = tokens.RefreshToken
	}

	url := authRequestState.OriginalURL

	if url == "" {
		url = "/"
	}

	m.Logger.Printf("authorize callback successful, redirecting to %s", url)
	http.Redirect(w, r, url, http.StatusSeeOther)
}
