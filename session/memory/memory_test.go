package memory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/janivihervas/oidc-go/session"
)

func TestMemory(t *testing.T) {
	var (
		store = New()
		key   = []byte("key")
		state = session.State{
			ID: key,
		}
		ctx = context.Background()
	)

	_, err := store.Get(ctx, key)
	assert.Error(t, err)
	assert.Equal(t, err, session.ErrNotFound)

	err = store.Save(ctx, key, state)
	assert.NoError(t, err)

	s, err := store.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, key, s.ID)

	err = store.Delete(ctx, key)
	assert.NoError(t, err)

	_, err = store.Get(ctx, key)
	assert.Error(t, err)
	assert.Equal(t, err, session.ErrNotFound)
}
