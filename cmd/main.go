package main

import (
	"flag"
	"log"

	"github.com/Kry0z1/chat/cmd/socket"
)

var (
	mode    = flag.String("m", "server", "run as server or as client")
	address = flag.String("a", "localhost:8000", "address to connect to in client mode")
	port    = flag.Int("p", 8000, "port to listen on in server mode")
)

func main() {
	flag.Parse()

	if *mode == "server" {
		socket.Listen(*port)
	} else if *mode == "client" {
		socket.Connect(*address)
	} else {
		log.Fatalf("Unknown mode: %s", *mode)
	}
}
