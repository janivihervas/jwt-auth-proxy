package middleware

import (
	"net/http"

	"github.com/janivihervas/jwt-auth-proxy/pkg/oidc/provider"
)

const (
	accessTokenName  = "access_token"
	authHeaderName   = "Authorization"
	authHeaderPrefix = "Bearer"
)

// New ...
func New(provider provider.OIDC, next http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/oauth2/callback", func(writer http.ResponseWriter, r *http.Request) {

	})

	mux.HandleFunc("/", func(writer http.ResponseWriter, r *http.Request) {
		accessToken, err := r.Cookie("access_token")
		if err != nil || accessToken.Value == "" {

		}
	})
	return mux
}
