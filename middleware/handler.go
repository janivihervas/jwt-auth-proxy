package middleware

import (
	"log"
	"net/http"

	"github.com/janivihervas/oidc-go/jwt"
)

func (m *middleware) defaultHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("default", r.URL.String())

	var (
		accessTokenStr string
		ctx            = r.Context()
	)

	accessTokenStr, err := extractAccessToken(ctx, r)
	if err != nil {
		accessTokenStr = m.refreshAccessToken(ctx, w, r)
	}

	var validationErr error
	for i := 0; i < 3; i++ {
		_, validationErr = jwt.ParseAccessToken(ctx, accessTokenStr)
		if validationErr == jwt.ErrTokenExpired {
			accessTokenStr = m.refreshAccessToken(ctx, w, r)
			continue
		}
		break
	}

	if validationErr != nil {
		m.unauthorized(ctx, w, r)
		return
	}

	m.next.ServeHTTP(w, r)
}
