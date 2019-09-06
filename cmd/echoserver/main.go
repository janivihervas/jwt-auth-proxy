package main

import (
	"os"

	"github.com/janivihervas/authproxy/internal/server"
	"github.com/janivihervas/authproxy/upstream"
)

func main() {
	var port = "3000"

	if e := os.Getenv("PORT"); e != "" {
		port = e
	}

	err := server.RunHTTP(port, upstream.Echo{})
	if err != nil {
		panic(err)
	}
}
