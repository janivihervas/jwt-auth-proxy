package oidc

import (
	"context"
	"errors"
	"net/url"
	"strings"
)

const (
	// ResponseModeFormPost is used to define the response mode of the OIDC Provider to the authentication request
	ResponseModeFormPost = "form_post"
)

// AuthenticationRequestURL that the client should be redirected to if no token is present.
// Response mode is the mode the OIDC provider will response to this request.
// Check the documentation of the OIDC provider to see all allowed values. For convenience,
// ResponseModeFormPost is provided for responding with a POST request and a form
func (c Client) AuthenticationRequestURL(state string, responseMode string) string {
	return c.Endpoint.AuthURL + "?" + url.Values{
		"client_id":     []string{c.ClientID},
		"response_type": []string{"code"},
		"redirect_uri":  []string{c.RedirectURL},
		"response_mode": []string{responseMode},
		"scope":         []string{strings.Join(c.Scope, " ")},
		"state":         []string{state},
	}.Encode()
}

// AuthenticationResponse form that is parsed from the callback
type AuthenticationResponse struct {
	Code  string
	State string
}

var (
	// ErrNoCodeInResponse is returned when code field is empty in the authentication response
	ErrNoCodeInResponse = errors.New("oidc: no code in authentication response")
	// ErrNoStateInResponse is returned when state field is empty in the authentication response
	ErrNoStateInResponse = errors.New("oidc: no state in authentication response")
)

// ParseAuthenticationResponseForm from the callback
func (c Client) ParseAuthenticationResponseForm(ctx context.Context, form url.Values) (AuthenticationResponse, error) {
	response := AuthenticationResponse{
		Code:  form.Get("code"),
		State: form.Get("state"),
	}

	if response.Code == "" {
		return response, ErrNoCodeInResponse
	}

	if response.State == "" {
		return response, ErrNoStateInResponse
	}

	return response, nil
}
