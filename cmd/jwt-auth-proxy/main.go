package main

import (
	"os"

	"github.com/janivihervas/oidc-go/internal/server"

	"github.com/janivihervas/oidc-go/upstream"
)

func main() {
	err := server.RunHTTP(os.Getenv("PORT"), upstream.Echo{})
	if err != nil {
		panic(err)
	}
}
