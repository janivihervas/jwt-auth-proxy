package main

import (
	"os"

	"github.com/janivihervas/authproxy/internal/server"

	"github.com/janivihervas/authproxy/upstream"
)

func main() {
	err := server.RunHTTP(os.Getenv("PORT"), upstream.Echo{})
	if err != nil {
		panic(err)
	}
}
