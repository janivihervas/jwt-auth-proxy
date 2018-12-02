package jwt

import (
	"context"
	"errors"
)

var (
	ErrTokenExpired = errors.New("oidc: token is expired")
)

type AccessToken struct {
}

func (at AccessToken) Valid() error {
	return nil
}

type Token interface {
	Valid() error
}

func ParseAccessToken(ctx context.Context, accessTokenStr string) (AccessToken, error) {
	// TODO
	if accessTokenStr == "" {
		return AccessToken{}, errors.New("empty")
	}
	return AccessToken{}, nil
}
