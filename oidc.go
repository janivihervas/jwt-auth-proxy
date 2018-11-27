package oidc

import (
	"errors"
	"net/http"
)

const (
	ScopeOpenID        = "openid"
	ScopeOfflineAccess = "offline_access"

	ResponseModeFormPost = "form_post"
)

var (
	ErrNoCodeInResponse  = errors.New("oidc: no code in authentication response")
	ErrNoStateInResponse = errors.New("oidc: no state in authentication response")
)

type Client struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scope        []string
	Endpoint     Endpoint
	HTTPClient   *http.Client
}

type Endpoint struct {
	AuthURL  string
	TokenURL string
}
