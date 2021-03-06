package upstream

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/stretchr/testify/assert"
)

func TestEcho_ServeHTTP(t *testing.T) {
	cupaloy := cupaloy.New(
		cupaloy.CreateNewAutomatically(true),
		cupaloy.FailOnUpdate(false),
		cupaloy.ShouldUpdate(func() bool {
			return true
		}),
	)

	test := func(t *testing.T, w *httptest.ResponseRecorder, r *http.Request) {
		t.Helper()

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		b, err := ioutil.ReadAll(w.Result().Body)
		assert.NoError(t, err)

		cupaloy.SnapshotT(t, string(b))
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
