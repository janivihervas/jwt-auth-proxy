package jwt

import "errors"

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

func ParseAccessToken(accessTokenStr string) (AccessToken, error) {
	// TODO
	return AccessToken{}, nil
}
