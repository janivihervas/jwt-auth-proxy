package main

import (
	"os"

	"github.com/janivihervas/jwt-auth-proxy/pkg/http/server"

	"github.com/janivihervas/jwt-auth-proxy/pkg/http/upstream"
)

func main() {
	err := server.RunHTTP(os.Getenv("PORT"), upstream.Echo{})
	if err != nil {
		panic(err)
	}
}
