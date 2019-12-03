package authproxy

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/sessions"

	"golang.org/x/oauth2"

	"github.com/pkg/errors"
)

// Config for Middleware
type Config struct {
	// Next http.Handler
	Next http.Handler
	// CallbackPath to handle authorization callback
	CallbackPath string
	// AuthClient for handling authentication flow
	AuthClient *oauth2.Config
	// AdditionalAuthURLParameters for providers who require additional authorization parameters,
	// like Auth0 requires to set and "audience" parameter in order to receive a JWT access token
	AdditionalAuthURLParameters []oauth2.AuthCodeOption
	// SessionStore for persisting session state
	SessionStore sessions.Store
	// SkipAuthenticationRegex for skipping authentication on these paths
	SkipAuthenticationRegex []string
	// SkipRedirectToLoginRegex for skipping redirecting user to auth provider's login page.
	// If a path matches one of these, a response with status code 401 with
	// JSON with redirectUrl field will be returned. Use this to prevent the middleware redirecting
	// API requests to the login page.
	SkipRedirectToLoginRegex []string

	// Logger, optional
	Logger *log.Logger

	mux                      *http.ServeMux
	skipAuthenticationRegex  []*regexp.Regexp
	skipRedirectToLoginRegex []*regexp.Regexp
}

// Valid returns a nil error if the config is valid
func (c *Config) Valid() error {
	if c == nil {
		return errors.New("Config is nil")
	}
	var errorMsg string

	if c.CallbackPath == "" {
		errorMsg += "CallbackPath is empty\n"
	}
	if c.CallbackPath == "/" {
		errorMsg += "CallbackPath is can't be '/'\n"
	}
	if c.AuthClient == nil {
		errorMsg += "AuthClient is nil\n"
	}
	if c.SessionStore == nil {
		errorMsg += "SessionStore is nil\n"
	}
	for i, s := range c.SkipAuthenticationRegex {
		r, err := regexp.Compile(s)
		if err != nil {
			errorMsg += fmt.Sprintf("SkipAuthenticationRegex[%d] (\"%s\") is invalid regex: %+v\n", i, s, err)
		} else {
			c.skipAuthenticationRegex = append(c.skipAuthenticationRegex, r)
		}
	}

	for i, s := range c.SkipRedirectToLoginRegex {
		r, err := regexp.Compile(s)
		if err != nil {
			errorMsg += fmt.Sprintf("SkipRedirectToLoginRegex[%d] (\"%s\") is invalid regex: %+v\n", i, s, err)
		} else {
			c.skipRedirectToLoginRegex = append(c.skipRedirectToLoginRegex, r)
		}
	}

	if c.Logger == nil {
		c.Logger = log.New(ioutil.Discard, "", log.LstdFlags)
	}

	if errorMsg != "" {
		return errors.New("config is not valid:\n\n" + errorMsg)
	}

	return nil
}
