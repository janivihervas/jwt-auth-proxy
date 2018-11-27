package main

import (
	"os"

	"github.com/janivihervas/oidc-go/pkg/http/server"

	"github.com/janivihervas/oidc-go/pkg/http/upstream"
)

func main() {
	err := server.RunHTTP(os.Getenv("PORT"), upstream.Echo{})
	if err != nil {
		panic(err)
	}
}
