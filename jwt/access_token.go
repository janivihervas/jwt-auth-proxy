package jwt

import (
	"context"
	"errors"
)

var (
	// ErrTokenExpired is returned if the token is expired
	ErrTokenExpired = errors.New("jwt: token is expired")
)

// AccessToken contains the fields from the JWT
type AccessToken struct {
}

// Valid returns non-nil error if any of the fields are invalid or the token is expired
func (token AccessToken) Valid() error {
	return nil
}

// ParseAccessToken and validate it
func ParseAccessToken(ctx context.Context, accessTokenStr string) (AccessToken, error) {
	// TODO
	if accessTokenStr == "" {
		return AccessToken{}, errors.New("empty")
	}
	return AccessToken{}, nil
}
