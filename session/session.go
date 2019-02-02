// Package session for storing state
package session

import (
	"context"
	"errors"
)

// State of the current user
type State struct {
	ID           []byte
	State        string
	RefreshToken string
	OriginalURL  string
}

var (
	// ErrNotFound is returned by the storage when the session is not found
	ErrNotFound = errors.New("oidc/session: session not found")
)

// Storage for storing a session state
type Storage interface {
	Get(ctx context.Context, key []byte) (State, error)
	Save(ctx context.Context, key []byte, state State) error
	Delete(ctx context.Context, key []byte) error
}
