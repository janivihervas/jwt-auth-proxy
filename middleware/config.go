package middleware

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
)

// Valid returns a nil error if the config is valid
func (c Config) Valid() error {
	var errorMsg string

	if c.AuthClient == nil {
		errorMsg = errorMsg + "AuthClient is nil\n"
	}
	if c.AuthenticationCallbackPath == "" || c.AuthenticationCallbackPath[0:] != "/" {
		errorMsg = errorMsg + "AuthenticationCallbackPath is empty of missing \"/\" from the start: " + c.AuthenticationCallbackPath + "\n"
	}
	if c.SessionStorage == nil {
		errorMsg = errorMsg + "SessionStorage is nil\n"
	}
	if c.CookieHashKey == nil || len(c.CookieHashKey) == 0 {
		errorMsg = errorMsg + "CookieHashKey is empty\n"
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

	if errorMsg != "" {
		return errors.New("Config is not valid:\n\n" + errorMsg)
	}

	return nil
}
