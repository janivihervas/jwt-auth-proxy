package oidc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
}

func (c Client) TokenRequest(ctx context.Context, code string) (TokenResponse, error) {
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

func (c Client) RefreshTokens(ctx context.Context, refreshToken string) (TokenResponse, error) {
	form := url.Values{
		"client_id":     []string{c.ClientID},
		"client_secret": []string{c.ClientSecret},
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{refreshToken},
		"scope":         []string{strings.Join(c.Scope, " ")},
	}

	return c.tokenRequest(ctx, form)
}

func (c Client) tokenRequest(ctx context.Context, form url.Values) (TokenResponse, error) {
	var response TokenResponse

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
		return response, errors.Wrapf(err, "oidc: token request failed; not 200 status code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response, errors.Wrap(err, "oidc: couldn't decode token response")
	}

	// TODO: validate tokens

	return response, nil
}
