// Package session for storing state
package session

import (
	"context"
	"errors"
)

type State struct {
	ID           []byte
	State        string
	RefreshToken string
	OriginalURL  string
}

var (
	ErrNotFound = errors.New("oidc: session not found")
)

type Storage interface {
	Get(ctx context.Context, key []byte) (State, error)
	Save(ctx context.Context, key []byte, session State) error
	Delete(ctx context.Context, key []byte) error
}
