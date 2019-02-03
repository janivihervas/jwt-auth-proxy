package middleware

import (
	"net/http"

	"github.com/gorilla/securecookie"
)

const (
	accessTokenName   = "access_token"
	authHeaderName    = "Authorization"
	authHeaderPrefix  = "Bearer"
	sessionCookieName = "oidc_session"
)

// Middleware for authentication requests
type Middleware struct {
	Config
}

// ServeHTTP will authenticate the request and forward it to the next http.Handler
func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := m.setupAccessTokenAndSession(r.Context(), w, r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	m.mux.ServeHTTP(w, r)
}

// NewMiddleware creates a new authentication middleware
func NewMiddleware(config Config) (*Middleware, error) {
	if err := config.Valid(); err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	config.cookieStore = securecookie.New(config.CookieHashKey, config.CookieEncryptKey)
	config.mux = mux

	m := &Middleware{
		Config: config,
	}
	mux.HandleFunc("/oauth2/callback", m.authorizeCallback)
	mux.HandleFunc("/", m.defaultHandler)

	return m, nil
}
