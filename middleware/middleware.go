package middleware

import (
	"net/http"
	"strings"

	"github.com/gorilla/securecookie"

	"github.com/janivihervas/oidc-go"
)

const (
	accessTokenName   = "access_token"
	authHeaderName    = "Authorization"
	authHeaderPrefix  = "Bearer"
	sessionCookieName = "oidc_session"
)

type middleware struct {
	mux            *http.ServeMux
	client         oidc.Client
	cookieStore    *securecookie.SecureCookie
	sessionStorage oidc.SessionStorage
	next           http.Handler
}

func (m *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

// New ...
func New(client oidc.Client, sessionStorage oidc.SessionStorage, next http.Handler) http.Handler {
	mux := http.NewServeMux()
	m := &middleware{
		mux:    mux,
		client: client,
		// TODO
		cookieStore: securecookie.New(
			[]byte(strings.Repeat("x", 32)),
			[]byte(strings.Repeat("y", 32)),
		),
		next:           next,
		sessionStorage: sessionStorage,
	}
	mux.HandleFunc("/oauth2/callback", m.authorizeCallback)
	mux.HandleFunc("/", m.defaultHandler)

	return m
}
