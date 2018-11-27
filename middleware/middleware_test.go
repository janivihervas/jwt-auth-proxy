package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janivihervas/jwt-auth-proxy/internal/mock"

	"github.com/janivihervas/jwt-auth-proxy/internal/http/upstream"
)

func TestNewWithMockRedirect(t *testing.T) {
	p := &mock.Provider{}
	h := New(p, upstream.Echo{})
	server := httptest.NewServer(h)
	defer server.Close()

	t.Run("Invalid access token redirects to provider's authorization page", func(t *testing.T) {
		noToken, err := http.NewRequest(http.MethodGet, server.URL, nil)
		assert.NoError(t, err)

		nonEmptyCookie, err := http.NewRequest(http.MethodGet, server.URL, nil)
		assert.NoError(t, err)
		nonEmptyCookie.AddCookie(&http.Cookie{
			Name:  accessTokenName,
			Value: "foo",
		})

		nonEmptyHeader, err := http.NewRequest(http.MethodGet, server.URL, nil)
		assert.NoError(t, err)
		nonEmptyHeader.Header.Add(authHeaderName, authHeaderPrefix+" bar")

		cases := []*http.Request{
			noToken,
			nonEmptyCookie,
			nonEmptyHeader,
		}
		p.ValidationError = errors.New("error")

		for i, req := range cases {
			resp, err := http.DefaultClient.Do(req)
			assert.NoError(t, err, strconv.Itoa(i), "executing request")

			assert.Equal(t, http.StatusSeeOther, resp.StatusCode, strconv.Itoa(i), "status code")
			assert.Equal(t, p.AuthorizationUrl(), resp.Header.Get("Location"), strconv.Itoa(i), "Location header")
			assert.NoError(t, resp.Body.Close(), strconv.Itoa(i), "closing body")
		}
		p.ValidationError = nil
	})

	t.Run("POST with malformed body returns 500", func(t *testing.T) {
		t.Skip("Not implemented")
	})

	t.Run("Expired access token redirects to provider's authorization page", func(t *testing.T) {
		t.Run("No refresh token", func(t *testing.T) {
			t.Skip("Not implemented")
		})

		t.Run("Updating with refresh token fails", func(t *testing.T) {
			t.Skip("Not implemented")
		})
	})

	t.Run("Redirect stores the original request", func(t *testing.T) {
		t.Skip("Not implemented")
	})

}
