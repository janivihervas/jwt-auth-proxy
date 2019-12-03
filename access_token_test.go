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
	)

	s, err := m.getAccessTokenFromCookie(r)
	assert.Error(t, err)
	assert.Equal(t, "", s)

	r.AddCookie(createAccessTokenCookie(""))
	s, err = m.getAccessTokenFromCookie(r)
	assert.Error(t, err)
	assert.Equal(t, "", s)

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(createAccessTokenCookie(accessToken))
	s, err = m.getAccessTokenFromCookie(r)
	assert.NoError(t, err)
	assert.Equal(t, accessToken, s)
}

func TestMiddleware_getAccessTokenFromHeader(t *testing.T) {
	var (
		m           = &Middleware{}
		r           = httptest.NewRequest(http.MethodGet, "/", nil)
		accessToken = "foo"
	)

	s, err := m.getAccessTokenFromHeader(r)
	assert.Error(t, err)
	assert.Equal(t, "", s)

	r.Header.Set(authHeaderName, "")
	s, err = m.getAccessTokenFromHeader(r)
	assert.Error(t, err)
	assert.Equal(t, "", s)

	r.Header.Set(authHeaderName, authHeaderPrefix)
	s, err = m.getAccessTokenFromHeader(r)
	assert.Error(t, err)
	assert.Equal(t, "", s)

	r.Header.Set(authHeaderName, authHeaderPrefix+" "+accessToken)
	s, err = m.getAccessTokenFromHeader(r)
	assert.NoError(t, err)
	assert.Equal(t, accessToken, s)
}

func TestMiddleware_getAccessTokenFromSession(t *testing.T) {
	var (
		m = &Middleware{
			&Config{},
		}
		accessToken = "foo"
	)

	t.Run("State not in context", func(t *testing.T) {
		s, err := m.getAccessTokenFromSession(context.Background())
		assert.Error(t, err)
		assert.Equal(t, "", s)
	})

	t.Run("State in context but nil", func(t *testing.T) {
		var state *sessionState
		ctx := context.WithValue(context.Background(), ctxStateKey, state)
		s, err := m.getAccessTokenFromSession(ctx)
		assert.Error(t, err)
		assert.Equal(t, "", s)
	})

	t.Run("State in context, access token empty", func(t *testing.T) {
		state := &sessionState{}
		ctx := context.WithValue(context.Background(), ctxStateKey, state)
		s, err := m.getAccessTokenFromSession(ctx)
		assert.Error(t, err)
		assert.Equal(t, "", s)
	})

	t.Run("State in session and access token is not empty", func(t *testing.T) {
		state := &sessionState{
			AccessToken: accessToken,
		}
		ctx := context.WithValue(context.Background(), ctxStateKey, state)
		s, err := m.getAccessTokenFromSession(ctx)
		assert.NoError(t, err)
		assert.Equal(t, accessToken, s)
	})
}
