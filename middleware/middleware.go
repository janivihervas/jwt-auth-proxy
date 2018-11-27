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

func (m *middleware) refreshAccessToken(writer http.ResponseWriter, r *http.Request) error {
	return nil
}

// New ...
func New(client oidc.Client, next http.Handler) http.Handler {
	mux := http.NewServeMux()
	m := &middleware{
		mux:    mux,
		client: client,
		// TODO
		cookieStore: securecookie.New(
			[]byte(strings.Repeat("x", 32)),
			[]byte(strings.Repeat("y", 32)),
		),
		next: next,
	}
	mux.HandleFunc("/oauth2/callback", m.authorizeCallback)

	mux.HandleFunc("/", func(writer http.ResponseWriter, r *http.Request) {

	})
	return m
}
