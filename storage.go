package oidc

import (
	"context"
	"errors"
)

type Session struct {
	ID           []byte
	State        string
	RefreshToken string
}

var (
	ErrNoSessionFound = errors.New("oidc: session wasn't found")
)

type SessionStorage interface {
	Get(ctx context.Context, key []byte) (Session, error)
	Save(ctx context.Context, key []byte, session Session) error
}
