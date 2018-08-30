package main

import (
	"log"
	"net/http"

	"github.com/janivihervas/jwt-auth-proxy/pkg/http/upstream"
)

func main() {
	log.Println("Starting echo server...")
	log.Fatal(http.ListenAndServe("127.0.0.1:3000", upstream.Echo{}))
}
