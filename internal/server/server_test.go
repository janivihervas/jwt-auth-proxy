package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janivihervas/authproxy/upstream"
)

func TestRunHTTP(t *testing.T) {
	go func() {
		err := RunHTTP("30000", upstream.Echo{}, log.New(ioutil.Discard, "", log.LstdFlags))
		if err != nil {
			panic(err)
		}
	}()

	var err error

	for i := 0; i < 5; i++ {
		var resp *http.Response
		resp, err = http.Get("http://localhost:30000/foo")
		if err != nil {
			continue
		}
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NoError(t, resp.Body.Close())
	}

	if err != nil {
		t.Fatal("Server didn't start", err)
	}
}
