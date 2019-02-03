package authproxy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware_getAccessTokenFromCookie(t *testing.T) {
	var (
		m           = &Middleware{}
		r           = httptest.NewRequest(http.MethodGet, "/", nil)
		accessToken = "foo"
		cookie      = createAccessTokenCookie(accessToken)
		ctx         = context.Background()
	)

	assert.Equal(t, "", m.getAccessTokenFromCookie(ctx, r))

	r.AddCookie(cookie)
	assert.Equal(t, accessToken, m.getAccessTokenFromCookie(ctx, r))
}

func TestMiddleware_getAccessTokenFromHeader(t *testing.T) {
	var (
		m           = &Middleware{}
		r           = httptest.NewRequest(http.MethodGet, "/", nil)
		accessToken = "foo"
		ctx         = context.Background()
	)

	assert.Equal(t, "", m.getAccessTokenFromHeader(ctx, r))

	r.Header.Set(authHeaderName, "")
	assert.Equal(t, "", m.getAccessTokenFromHeader(ctx, r))

	r.Header.Set(authHeaderName, authHeaderPrefix)
	assert.Equal(t, "", m.getAccessTokenFromHeader(ctx, r))

	r.Header.Set(authHeaderName, authHeaderPrefix+" "+accessToken)
	assert.Equal(t, accessToken, m.getAccessTokenFromHeader(ctx, r))
}
