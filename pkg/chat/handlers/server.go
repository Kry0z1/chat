package chat

import (
	"encoding/json"
	"errors"
	"log"
	"net"

	chat "github.com/Kry0z1/chat/pkg/chat/lib"
)

var (
	ErrBadJson        = errors.New("Couldn't parse json")
	ErrUnknownCommand = errors.New("Unknown command")
)

type user struct {
	Username string `json:"username"`
}

type Request struct {
	Command string `json:"command"`
	Content string `json:"content"`

	internalError error
}

type Response struct {
	Topic    string `json:"topic"`
	Username string `json:"username"`
	Content  string `json:"content"`
	Error    error  `json:"error"`
}

func ServerHandler() func(net.Conn) {
	topics := make(map[string]chat.Topic)
	topics["global"] = chat.NewTopic("global", 32)

	return func(c net.Conn) {
		var usr user
		for {
			err := json.NewDecoder(c).Decode(&usr)
			if err != nil {
				json.NewEncoder(c).Encode(Response{
					Error:   ErrBadJson,
					Content: err.Error(),
				})
			} else {
				break
			}
		}

		curUser := topics["global"].RegisterUser(usr.Username, 10_000)

		listenConn := make(chan Request)
		go connListener(listenConn, c)

		for {
			select {
			case msg := <-listenConn:
				if msg.internalError != nil {
					json.NewEncoder(c).Encode(Response{
						Error: msg.internalError,
					})
					break
				}
				switch msg.Command {
				case "topic":
					curUser.Close()

					topic, ok := topics[msg.Content]
					if !ok {
						topics[msg.Content] = chat.NewTopic(msg.Content, 16)
						topic = topics[msg.Content]
					}

					curUser = topic.RegisterUser(usr.Username, 10_000)
				case "publish":
					curUser.Publish(msg.Content)
				case "close":
					curUser.Close()
					close(listenConn)
					c.Close()
					return
				default:
					json.NewEncoder(c).Encode(Response{
						Error: ErrUnknownCommand,
					})
				}
			case msg := <-curUser.Recieve():
				err := json.NewEncoder(c).Encode(Response{
					Topic:    curUser.Topic(),
					Username: msg.Username,
					Content:  msg.Content.(string),
				})
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

// listens on conn c, forwards parsed requests to pipe
func connListener(pipe chan Request, c net.Conn) {
	var req Request
	for {
		err := json.NewDecoder(c).Decode(&req)
		if err != nil {
			pipe <- Request{
				internalError: ErrBadJson,
			}
		}
		pipe <- req
	}
}
