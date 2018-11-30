package storage

import (
	"context"
	"sync"

	"github.com/janivihervas/oidc-go"
)

func NewMemory() *Memory {
	return &Memory{
		memory: make(map[string]oidc.Session),
	}
}

type Memory struct {
	mutex  sync.RWMutex
	memory map[string]oidc.Session
}

func (mem *Memory) Get(ctx context.Context, key []byte) (oidc.Session, error) {
	mem.mutex.RLock()
	session, ok := mem.memory[string(key)]
	mem.mutex.RUnlock()

	if !ok {
		return oidc.Session{}, oidc.ErrNoSessionFound
	}

	return session, nil
}

func (mem *Memory) Save(ctx context.Context, key []byte, session oidc.Session) error {
	mem.mutex.Lock()
	mem.memory[string(key)] = session
	mem.mutex.Unlock()

	return nil
}
