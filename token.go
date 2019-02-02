package oidc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// Token response from requesting tokens
type Token struct {
	// AccessToken is used for authentication
	AccessToken string `json:"access_token"`
	// IDToken is used for identification
	IDToken string `json:"id_token"`
	// RefreshToken is used for refreshing access and ID tokens
	RefreshToken string `json:"refresh_token"`
}

// TokenRequest for requesting new tokens. Code parameter is from the AuthenticationResponse
func (c Client) TokenRequest(ctx context.Context, code string) (Token, error) {
	form := url.Values{
		"client_id":     []string{c.ClientID},
		"client_secret": []string{c.ClientSecret},
		"grant_type":    []string{"authorization_code"},
		"code":          []string{code},
		"redirect_uri":  []string{c.RedirectURL},
		"scope":         []string{strings.Join(c.Scope, " ")},
	}

	return c.tokenRequest(ctx, form)
}

// RefreshTokens will refresh the tokens with provided refresh token
func (c Client) RefreshTokens(ctx context.Context, refreshToken string) (Token, error) {
	form := url.Values{
		"client_id":     []string{c.ClientID},
		"client_secret": []string{c.ClientSecret},
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{refreshToken},
		"scope":         []string{strings.Join(c.Scope, " ")},
	}

	return c.tokenRequest(ctx, form)
}

func (c Client) tokenRequest(ctx context.Context, form url.Values) (Token, error) {
	var response Token

	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	req, err := http.NewRequest(http.MethodPost, c.Endpoint.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return response, errors.Wrap(err, "oidc: couldn't create a token request")
	}
	req = req.WithContext(ctx)

	resp, err := httpClient.Do(req)
	if err != nil {
		return response, errors.Wrap(err, "oidc: token request failed")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return response, errors.Wrapf(err, "oidc: token request failed; non-200 status code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response, errors.Wrap(err, "oidc: couldn't decode token response")
	}

	// TODO: validate tokens

	return response, nil
}
