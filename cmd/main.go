package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/Kry0z1/chat/pkg/socket"
)

var (
	mode    = flag.String("m", "server", "run as server or as client")
	address = flag.String("a", "localhost:8000", "address to connect to in client mode")
	port    = flag.Int("p", 8000, "port to listen on in server mode")
	verbose = flag.Bool("v", false, "run server in verbose mode")
)

func main() {
	flag.Parse()

	if *mode == "server" {
		log.Fatal(socket.Listen(*port, func(c net.Conn) {
			io.Copy(c, c)
			c.Close()
		}, *verbose))
	} else if *mode == "client" {
		log.Fatal(socket.Connect(*address, func(c net.Conn) error {
			reader := bufio.NewReader(os.Stdin)
			for {
				fmt.Print("Text to send: ")
				text, _ := reader.ReadString('\n')
				fmt.Fprintf(c, "%s", text+"\n")
				message, _ := bufio.NewReader(c).ReadString('\n')
				fmt.Print("Message from server: " + message)
			}
		}))
	} else {
		log.Fatalf("Unknown mode: %s", *mode)
	}
}
