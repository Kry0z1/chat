package socket

import (
	"fmt"
	"net"
)

func Listen(port int, handler func(c net.Conn), verbose bool) error {
	var log func(...interface{})
	if verbose {
		log = func(i ...interface{}) { fmt.Println(i) }
	} else {
		log = func(i ...interface{}) {}
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	log("Listening on port", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		log("Got new connection")

		go handler(conn)
	}
}
