// Package oidc provides a client to handle OpenID Connect authentication flow
package oidc

import (
	"net/http"

	"golang.org/x/oauth2"
)

const (
	// ScopeOpenID will return basic information
	ScopeOpenID        = "openid"
	ScopeOfflineAccess = "offline_access"
)

// Client for handling OIDC authentication flow
type Client struct {
	// ClientID of the authentication app
	ClientID string
	// ClientSecret of the authentication app
	ClientSecret string
	// RedirectURL
	RedirectURL string
	Scope       []string
	Endpoint    oauth2.Endpoint
	HTTPClient  *http.Client
}
