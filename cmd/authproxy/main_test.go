package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/janivihervas/authproxy/internal/server"
	"github.com/janivihervas/authproxy/upstream"

	"github.com/stretchr/testify/assert"
)

//func Test_parseAdditionalParameters(t *testing.T) {
//	assert.Empty(t, parseAdditionalParameters(""))
//	assert.Empty(t, parseAdditionalParameters("efewfwefwe"))
//	assert.Empty(t, parseAdditionalParameters("&"))
//	assert.Empty(t, parseAdditionalParameters("dqowi="))
//	assert.Empty(t, parseAdditionalParameters("=ofiofwen"))
//	assert.Empty(t, parseAdditionalParameters("dqowi:ofiofwen"))
//	assert.Empty(t, parseAdditionalParameters("dqowi=ofiofwen=foo"))
//
//	assert.Equal(
//		t,
//		[]oauth2.AuthCodeOption{oauth2.SetAuthURLParam("foo", "bar")},
//		parseAdditionalParameters("foo=bar"),
//	)
//	assert.Equal(
//		t,
//		[]oauth2.AuthCodeOption{oauth2.SetAuthURLParam("foo", "bar")},
//		parseAdditionalParameters("foo=bar&roeigneroi"),
//	)
//
//	assert.Equal(
//		t,
//		[]oauth2.AuthCodeOption{
//			oauth2.SetAuthURLParam("foo", "bar"),
//			oauth2.SetAuthURLParam("bar", "baz"),
//		},
//		parseAdditionalParameters("foo=bar&bar=baz"),
//	)
//}

func TestAuthProxy(t *testing.T) {
	go func() {
		err := server.RunHTTP("3000", upstream.Echo{}, log.New(ioutil.Discard, "", log.LstdFlags))
		if err != nil {
			panic(err)
		}
	}()

	var startErr error

	for i := 0; i < 5; i++ {
		var resp *http.Response
		resp, startErr = http.Get("http://localhost:3000/foo")
		if startErr == nil {
			assert.NoError(t, resp.Body.Close())
			break
		}
	}

	if startErr != nil {
		t.Fatal("Couldn't start test server", startErr)
	}
}
