package main

import "testing"
import "github.com/stretchr/testify/assert"

func TestHello(t *testing.T) {
	assert.Equal(t, "Hello world!", Hello())
}
