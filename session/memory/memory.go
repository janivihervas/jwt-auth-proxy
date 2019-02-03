// Package memory provides an in-memory session storage
package memory

import (
	"context"
	"sync"

	"github.com/janivihervas/jwt-auth-proxy/session"
)

// New creates a new in-memory session storage.
// Use only for testing purposes
func New() *Memory {
	return &Memory{
		memory: make(map[string]session.State),
	}
}

// Memory for storing session state in memory.
// Use only for testing purposes
type Memory struct {
	mutex  sync.RWMutex
	memory map[string]session.State
}

// Get a state from session with a key
func (mem *Memory) Get(ctx context.Context, key []byte) (session.State, error) {
	mem.mutex.RLock()
	s, ok := mem.memory[string(key)]
	mem.mutex.RUnlock()

	if !ok {
		return session.State{}, session.ErrNotFound
	}

	return s, nil
}

// Save a state into session
func (mem *Memory) Save(ctx context.Context, key []byte, state session.State) error {
	mem.mutex.Lock()
	mem.memory[string(key)] = state
	mem.mutex.Unlock()

	return nil
}

// Delete a session
func (mem *Memory) Delete(ctx context.Context, key []byte) error {
	mem.mutex.Lock()
	delete(mem.memory, string(key))
	mem.mutex.Unlock()

	return nil
}
