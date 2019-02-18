package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func Test_parseAdditionalParameters(t *testing.T) {
	assert.Empty(t, parseAdditionalParameters(""))
	assert.Empty(t, parseAdditionalParameters("efewfwefwe"))
	assert.Empty(t, parseAdditionalParameters("&"))
	assert.Empty(t, parseAdditionalParameters("dqowi="))
	assert.Empty(t, parseAdditionalParameters("=ofiofwen"))
	assert.Empty(t, parseAdditionalParameters("dqowi:ofiofwen"))
	assert.Empty(t, parseAdditionalParameters("dqowi=ofiofwen=foo"))

	assert.Equal(
		t,
		[]oauth2.AuthCodeOption{oauth2.SetAuthURLParam("foo", "bar")},
		parseAdditionalParameters("foo=bar"),
	)
	assert.Equal(
		t,
		[]oauth2.AuthCodeOption{oauth2.SetAuthURLParam("foo", "bar")},
		parseAdditionalParameters("foo=bar&roeigneroi"),
	)

	assert.Equal(
		t,
		[]oauth2.AuthCodeOption{
			oauth2.SetAuthURLParam("foo", "bar"),
			oauth2.SetAuthURLParam("bar", "baz"),
		},
		parseAdditionalParameters("foo=bar&bar=baz"),
	)
}
