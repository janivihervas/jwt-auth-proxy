package authproxy

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pkg/errors"

	"github.com/janivihervas/authproxy/internal/mock"

	"github.com/gorilla/sessions"

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

func TestMiddleware_getAccessTokenFromSession_MockStore(t *testing.T) {
	var (
		mockStore = &mock.SessionStore{}
		m         = &Middleware{
			&Config{
				SessionStore: mockStore,
			},
		}
		r           = httptest.NewRequest(http.MethodGet, "/", nil)
		accessToken = "foo"
		ctx         = context.Background()
	)

	t.Run("Get fails", func(t *testing.T) {
		mockStore.ErrGet = errors.New("foo")
		s, err := m.getAccessTokenFromSession(ctx, r)
		assert.Error(t, err)
		assert.Equal(t, "", s)
	})

	t.Run("State not in session", func(t *testing.T) {
		mockStore.ErrGet = nil
		mockStore.Session = &sessions.Session{
			Values: map[interface{}]interface{}{},
		}
		s, err := m.getAccessTokenFromSession(ctx, r)
		assert.Error(t, err)
		assert.Equal(t, "", s)
	})

	t.Run("State in session but wrong type", func(t *testing.T) {
		mockStore.ErrGet = nil
		mockStore.Session = &sessions.Session{
			Values: map[interface{}]interface{}{
				sessionName: 666,
			},
		}
		s, err := m.getAccessTokenFromSession(ctx, r)
		assert.Error(t, err)
		assert.Equal(t, "", s)
	})

	t.Run("State in session but access token is empty", func(t *testing.T) {
		mockStore.ErrGet = nil
		mockStore.Session = &sessions.Session{
			Values: map[interface{}]interface{}{
				sessionName: State{
					AccessToken: "",
				},
			},
		}
		s, err := m.getAccessTokenFromSession(ctx, r)
		assert.Error(t, err)
		assert.Equal(t, "", s)
	})

	t.Run("State in session and access token is not empty", func(t *testing.T) {
		mockStore.ErrGet = nil
		mockStore.Session = &sessions.Session{
			Values: map[interface{}]interface{}{
				sessionName: State{
					AccessToken: accessToken,
				},
			},
		}
		s, err := m.getAccessTokenFromSession(ctx, r)
		assert.NoError(t, err)
		assert.Equal(t, accessToken, s)
	})
}

func TestMiddleware_getAccessTokenFromSession_RealStore(t *testing.T) {
	var (
		m = &Middleware{
			&Config{
				SessionStore: sessions.NewCookieStore(
					[]byte(strings.Repeat("x", 32)),
					[]byte(strings.Repeat("y", 32)),
				),
			},
		}
		r           = httptest.NewRequest(http.MethodGet, "/", nil)
		accessToken = "foo"
		ctx         = context.Background()
	)

	t.Run("No session in request", func(t *testing.T) {
		s, err := m.getAccessTokenFromSession(ctx, r)
		assert.Error(t, err)
		assert.Equal(t, "", s)
	})

	t.Run("Session in request but access token is empty", func(t *testing.T) {
		r = httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		session, err := m.SessionStore.Get(r, sessionName)
		assert.NoError(t, err)

		session.Values[sessionName] = State{
			AccessToken: "",
		}
		err = session.Save(r, w)
		assert.NoError(t, err)

		var cookie *http.Cookie
		for _, c := range w.Result().Cookies() {
			if c.Name == sessionName {
				cookie = c
				break
			}
		}
		assert.NotNil(t, cookie)

		r.AddCookie(cookie)
		s, err := m.getAccessTokenFromSession(ctx, r)
		assert.Error(t, err)
		assert.Equal(t, "", s)
	})

	t.Run("Session in request and access token is not empty", func(t *testing.T) {
		r = httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		session, err := m.SessionStore.Get(r, sessionName)
		assert.NoError(t, err)

		session.Values[sessionName] = State{
			AccessToken: accessToken,
		}
		err = session.Save(r, w)
		assert.NoError(t, err)

		var cookie *http.Cookie
		for _, c := range w.Result().Cookies() {
			if c.Name == sessionName {
				cookie = c
				break
			}
		}
		assert.NotNil(t, cookie)

		r.AddCookie(cookie)
		s, err := m.getAccessTokenFromSession(ctx, r)
		assert.NoError(t, err)
		assert.Equal(t, accessToken, s)
	})

}
