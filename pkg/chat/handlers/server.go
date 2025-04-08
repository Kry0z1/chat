package chat

import (
	"encoding/json"
	"errors"
	"net"

	chat "github.com/Kry0z1/chat/pkg/chat/lib"
)

var (
	ErrBadJson        = errors.New("Couldn't parse json")
	ErrUnknownCommand = errors.New("Unknown command")
)

type UserIn struct {
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
	Error    string `json:"error"`
}

func ServerHandler() func(net.Conn) {
	topics := make(map[string]chat.Topic)
	topics["global"] = chat.NewTopic("global", 32)

	return func(c net.Conn) {
		var usr UserIn
		for {
			err := json.NewDecoder(c).Decode(&usr)
			if err != nil {
				json.NewEncoder(c).Encode(Response{
					Error: ErrBadJson.Error(),
				})
			} else {
				break
			}
		}

		json.NewEncoder(c).Encode(Response{
			Username: usr.Username,
		})

		curUser := topics["global"].RegisterUser(usr.Username, 10_000)

		listenConn := make(chan Request)
		go connListener(listenConn, c)

		for {
			select {
			case msg := <-listenConn:
				if msg.internalError != nil {
					json.NewEncoder(c).Encode(Response{
						Error: msg.internalError.Error(),
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
					c.Close()
					return
				default:
					json.NewEncoder(c).Encode(Response{
						Error: ErrUnknownCommand.Error(),
					})
				}
			case msg := <-curUser.Recieve():
				json.NewEncoder(c).Encode(Response{
					Topic:    curUser.Topic(),
					Username: msg.Username,
					Content:  msg.Content.(string),
				})
			}
		}
	}
}

// listens on conn c, forwards parsed requests to pipe
func connListener(pipe chan Request, c net.Conn) {
	var req Request
	dec := json.NewDecoder(c)
	for {
		err := dec.Decode(&req)
		if err != nil {
			pipe <- Request{
				internalError: ErrBadJson,
			}
			continue
		}
		pipe <- req
		if req.Command == "close" {
			return
		}
	}
}
