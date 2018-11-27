package oidc

import "net/url"

func (c Client) AuthenticationRequestURL(responseMode string, state string) string {
	return c.Endpoint.AuthURL + "?" + url.Values{
		"client_id":     []string{c.ClientID},
		"response_type": []string{"code"},
		"redirect_uri":  []string{c.RedirectURL},
		"response_mode": []string{responseMode},
		"scope":         c.Scope,
		"state":         []string{state},
	}.Encode()
}

type AuthenticationResponse struct {
	Code  string
	State string
}

func (c Client) AuthenticationResponseForm(form url.Values) (AuthenticationResponse, error) {
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
