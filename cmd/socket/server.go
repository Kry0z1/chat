package socket

import (
	"fmt"
	"io"
	"log"
	"net"
)

func Listen(port int) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("listening on port", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(c net.Conn) {
			io.Copy(c, c)
			c.Close()
		}(conn)
	}
}
