package authproxy

import (
	"testing"

	"github.com/janivihervas/authproxy/internal/mock"
	"golang.org/x/oauth2"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/stretchr/testify/assert"
)

func TestConfig_Valid(t *testing.T) {
	cupaloy := cupaloy.New(
		cupaloy.CreateNewAutomatically(true),
		cupaloy.FailOnUpdate(false),
		cupaloy.ShouldUpdate(func() bool {
			return true
		}),
	)

	var c *Config

	err := c.Valid()
	assert.Error(t, err)
	err = cupaloy.SnapshotMulti("0", err.Error())
	assert.NoError(t, err)

	c = &Config{
		SkipRedirectToLoginRegex: []string{"["},
		SkipAuthenticationRegex:  []string{"["},
	}
	err = c.Valid()
	assert.Error(t, err)
	err = cupaloy.SnapshotMulti("1", err.Error())
	assert.NoError(t, err)

	c.CallbackPath = "/"
	err = c.Valid()
	assert.Error(t, err)
	err = cupaloy.SnapshotMulti("2", err.Error())
	assert.NoError(t, err)

	c.CallbackPath = "/callback"
	err = c.Valid()
	assert.Error(t, err)
	err = cupaloy.SnapshotMulti("3", err.Error())
	assert.NoError(t, err)

	c.AuthClient = &oauth2.Config{}
	err = c.Valid()
	assert.Error(t, err)
	err = cupaloy.SnapshotMulti("4", err.Error())
	assert.NoError(t, err)

	c.SessionStore = &mock.SessionStore{}
	err = c.Valid()
	assert.Error(t, err)
	err = cupaloy.SnapshotMulti("5", err.Error())
	assert.NoError(t, err)

	c.SkipAuthenticationRegex = []string{}
	err = c.Valid()
	assert.Error(t, err)
	assert.Equal(t, 0, len(c.skipAuthenticationRegex))
	err = cupaloy.SnapshotMulti("6", err.Error())
	assert.NoError(t, err)

	c.SkipAuthenticationRegex = []string{".*"}
	err = c.Valid()
	assert.Error(t, err)
	assert.Equal(t, 1, len(c.skipAuthenticationRegex))
	err = cupaloy.SnapshotMulti("7", err.Error())
	assert.NoError(t, err)

	c.SkipRedirectToLoginRegex = []string{}
	err = c.Valid()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(c.skipRedirectToLoginRegex))

	c.SkipRedirectToLoginRegex = []string{".*"}
	err = c.Valid()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(c.skipRedirectToLoginRegex))
}
