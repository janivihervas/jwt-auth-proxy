package middleware

import (
	"net/http"
	"strings"

	"github.com/janivihervas/oidc-go/session"

	"github.com/gorilla/securecookie"

	"github.com/janivihervas/oidc-go"
)

const (
	accessTokenName   = "access_token"
	authHeaderName    = "Authorization"
	authHeaderPrefix  = "Bearer"
	sessionCookieName = "oidc_session"
)

type redirectFunc func(r *http.Request) bool

type middleware struct {
	mux            *http.ServeMux
	client         oidc.Client
	cookieStore    *securecookie.SecureCookie
	sessionStorage session.Storage
	next           http.Handler
	redirect       redirectFunc
}

func (m *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

func New(client oidc.Client, sessionStorage session.Storage, next http.Handler) http.Handler {
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
		redirect: func(r *http.Request) bool {
			// TODO
			return r.Method == http.MethodGet
		},
	}
	mux.HandleFunc("/oauth2/callback", m.authorizeCallback)
	mux.HandleFunc("/", m.defaultHandler)

	return m
}
