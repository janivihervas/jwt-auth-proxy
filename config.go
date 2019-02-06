package authproxy

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"

	"github.com/gorilla/sessions"

	"golang.org/x/oauth2"

	"github.com/pkg/errors"
)

// Config for Middleware
type Config struct {
	// Next http.Handler
	Next http.Handler
	// AuthClient for handling authentication flow
	AuthClient *oauth2.Config
	// SessionStore for persisting session state
	SessionStore sessions.Store
	// SkipAuthenticationRegex for skipping authentication on these paths
	SkipAuthenticationRegex []string
	// SkipRedirectToLoginRegex for skipping redirecting user to auth provider's login page.
	// If a path matches one of these, a response with status code 401 or 403 with
	// JSON with redirectUrl field will be returned. Use this to prevent the middleware redirecting
	// API requests to the login page.
	SkipRedirectToLoginRegex []string

	// Logger
	Logger *log.Logger

	mux                      *http.ServeMux
	callbackPath             string
	skipAuthenticationRegex  []*regexp.Regexp
	skipRedirectToLoginRegex []*regexp.Regexp
}

// Valid returns a nil error if the config is valid
func (c *Config) Valid() error {
	if c == nil {
		return errors.New("Config is nil")
	}
	var errorMsg string

	if c.AuthClient == nil {
		errorMsg = errorMsg + "AuthClient is nil\n"
	} else {
		u, err := url.Parse(c.AuthClient.RedirectURL)
		if err != nil {
			errorMsg = errorMsg + "Can't parse AuthClient.RedirectURL:" + err.Error() + "\n"
		} else {
			c.callbackPath = u.Path
		}
	}
	if c.SessionStore == nil {
		errorMsg = errorMsg + "SessionStore is nil\n"
	}
	for i, s := range c.SkipAuthenticationRegex {
		r, err := regexp.Compile(s)
		if err != nil {
			errorMsg = errorMsg + fmt.Sprintf("SkipAuthenticationRegex[%d] (\"%s\") is invalid regex: %+v\n", i, s, err)
		} else {
			c.skipAuthenticationRegex = append(c.skipAuthenticationRegex, r)
		}
	}

	for i, s := range c.SkipRedirectToLoginRegex {
		r, err := regexp.Compile(s)
		if err != nil {
			errorMsg = errorMsg + fmt.Sprintf("SkipRedirectToLoginRegex[%d] (\"%s\") is invalid regex: %+v\n", i, s, err)
		} else {
			c.skipRedirectToLoginRegex = append(c.skipRedirectToLoginRegex, r)
		}
	}

	if c.Logger == nil {
		c.Logger = log.New(ioutil.Discard, "", log.LstdFlags)
	}

	if errorMsg != "" {
		return errors.New("Config is not valid:\n\n" + errorMsg)
	}

	return nil
}
