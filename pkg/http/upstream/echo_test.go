package upstream

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/janivihervas/cupaloy"
	"github.com/stretchr/testify/assert"
)

func TestEcho_ServeHTTP(t *testing.T) {
	test := func(t *testing.T, w *httptest.ResponseRecorder, r *http.Request) {
		t.Helper()

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		b, err := ioutil.ReadAll(w.Result().Body)
		assert.NoError(t, err)

		err = cupaloy.SnapshotMulti(strings.Replace(t.Name(), "/", "-", -1), string(b))
		assert.NoError(t, err)
	}

	echo := Echo{}

	t.Run("Just request info", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		echo.ServeHTTP(w, r)

		test(t, w, r)
	})

	t.Run("With query", func(t *testing.T) {
		w := httptest.NewRecorder()
		q := url.Values{
			"foo": []string{"bar"},
		}.Encode()
		r := httptest.NewRequest(http.MethodGet, "/queries?"+q, nil)

		echo.ServeHTTP(w, r)

		test(t, w, r)
	})

	t.Run("With headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/headers", nil)
		h := make(http.Header)
		h.Add("foo", "bar")
		r.Header = h

		echo.ServeHTTP(w, r)

		test(t, w, r)
	})

	t.Run("With cookies", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/headers", nil)
		r.AddCookie(&http.Cookie{
			Name:  "foo",
			Value: "bar",
		})

		echo.ServeHTTP(w, r)

		test(t, w, r)
	})

	t.Run("With body", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/body", bytes.NewBufferString("something that is in the body"))

		echo.ServeHTTP(w, r)

		test(t, w, r)
	})
}
