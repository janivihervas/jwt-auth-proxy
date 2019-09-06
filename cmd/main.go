package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app     = kingpin.New("authproxy", "Authentication proxy")
	debug   = app.Flag("debug", "Enable debug mode.").Bool()
	timeout = app.Flag("timeout", "Timeout waiting for ping.").Default("5s").Envar("PING_TIMEOUT").Short('t').Duration()
	ip      = app.Arg("ip", "IP address to ping.").Required().IP()
	count   = app.Arg("count", "Number of packets to send").Int()
)

func main() {
	app.Version("0.0.1")
	kingpin.MustParse(app.Parse(os.Args[1:]))

	app.GetFlag("").IsSet
	fmt.Printf("Would ping: %s with timeout %s and count %d\n", *ip, *timeout, count)
}
