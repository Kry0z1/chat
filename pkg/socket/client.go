package socket

import (
	"net"
)

func Connect(address string, handler func(net.Conn) error) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	return handler(conn)
}
