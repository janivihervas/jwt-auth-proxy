package main

import (
	"flag"
	"log"
	"os"

	"github.com/janivihervas/authproxy/internal/server"
	"github.com/janivihervas/authproxy/upstream"
)

func main() {
	var port = "3000"

	if e := os.Getenv("PORT"); e != "" {
		port = e
	}

	var portFlag string
	flag.StringVar(&portFlag, "port", "3000", "port to run the server")
	flag.Parse()

	if portFlag != "" {
		port = portFlag
	}

	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime | log.LUTC)
	err := server.RunHTTP(port, upstream.Echo{}, logger)
	if err != nil {
		panic(err)
	}
}
