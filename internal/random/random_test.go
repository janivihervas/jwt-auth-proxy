package random

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBytes(t *testing.T) {
	length := rand.Intn(100)
	b1 := Bytes(length)
	b2 := Bytes(length)

	assert.NotEqual(t, b1, b2)
	assert.Equal(t, length, len(b1))
	assert.Equal(t, length, len(b2))
}

func TestString(t *testing.T) {
	length := rand.Intn(100)
	s1 := String(length)
	s2 := String(length)

	assert.NotEqual(t, s1, s2)
	assert.Equal(t, length, len(s1))
	assert.Equal(t, length, len(s2))
}
