package oidc

import "context"

type Session struct {
	ID           []byte
	State        string
	RefreshToken string
}

type SessionStorage interface {
	Get(ctx context.Context, key []byte) (Session, error)
	Save(ctx context.Context, key []byte, session Session) error
}
