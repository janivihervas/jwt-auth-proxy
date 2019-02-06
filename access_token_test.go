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

	s, err := m.getAccessTokenFromCookie(ctx, r)
	assert.Error(t, err)
	assert.Equal(t, "", s)

	r.AddCookie(cookie)
	s, err = m.getAccessTokenFromCookie(ctx, r)
	assert.NoError(t, err)
	assert.Equal(t, accessToken, s)
}

func TestMiddleware_getAccessTokenFromHeader(t *testing.T) {
	var (
		m           = &Middleware{}
		r           = httptest.NewRequest(http.MethodGet, "/", nil)
		accessToken = "foo"
		ctx         = context.Background()
	)

	s, err := m.getAccessTokenFromHeader(ctx, r)
	assert.Error(t, err)
	assert.Equal(t, "", s)

	r.Header.Set(authHeaderName, "")
	s, err = m.getAccessTokenFromHeader(ctx, r)
	assert.Error(t, err)
	assert.Equal(t, "", s)

	r.Header.Set(authHeaderName, authHeaderPrefix)
	s, err = m.getAccessTokenFromHeader(ctx, r)
	assert.Error(t, err)
	assert.Equal(t, "", s)

	r.Header.Set(authHeaderName, authHeaderPrefix+" "+accessToken)
	s, err = m.getAccessTokenFromHeader(ctx, r)
	assert.NoError(t, err)
	assert.Equal(t, accessToken, s)
}
