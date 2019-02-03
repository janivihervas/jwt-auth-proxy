package middleware

import (
	"net/http"
	"regexp"

	"golang.org/x/oauth2"

	"github.com/janivihervas/oidc-go/session"

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

// Config for Middleware
type Config struct {
	// Next http.Handler
	Next http.Handler
	// AuthClient for handling authentication flow
	AuthClient *oauth2.Config
	// AuthenticationCallbackPath handles the authentication callback. E.g. /callback
	AuthenticationCallbackPath string
	// SessionStorage for persisting session state
	SessionStorage session.Storage
	// CookieHashKey for validating session cookie signature. It is recommended to use a key with 32 or 64 bytes.
	CookieHashKey []byte
	// CookieEncryptKey for encrypting session cookie, optional. Valid key lengths are 16, 24, or 32 bytes.
	CookieEncryptKey []byte
	// SkipAuthenticationRegex for skipping authentication on these paths
	SkipAuthenticationRegex []string
	// SkipRedirectToLoginRegex for skipping redirecting user to auth provider's login page.
	// If a path matches one of these, a response with status code 401 or 403 with
	// JSON with redirectUrl field will be returned. Use this to prevent the middleware redirecting
	// API requests to the login page.
	SkipRedirectToLoginRegex []string

	mux                      *http.ServeMux
	cookieStore              *securecookie.SecureCookie
	skipAuthenticationRegex  []*regexp.Regexp
	skipRedirectToLoginRegex []*regexp.Regexp
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
