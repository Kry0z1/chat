package test

import (
	"sync"
	"testing"
	"time"

	"github.com/Kry0z1/chat/pkg/chat"
	"github.com/stretchr/testify/require"
)

func TestBasic(t *testing.T) {
	topic := chat.NewTopic("topic", 1)

	sender := topic.RegisterUser("a", 10000)

	numRecievers := 10
	recievers := make([]chat.User, numRecievers)
	pipes := make([]chan chat.Message, numRecievers)

	numMessages := 10
	wg := sync.WaitGroup{}
	wg.Add(numRecievers + 1)

	for i := range numRecievers {
		recievers[i] = topic.RegisterUser(string(byte(i)), 10000)
		pipes[i] = make(chan chat.Message, 10000)
		go func() {
			for range numMessages {
				pipes[i] <- <-recievers[i].Recieve()
			}
			wg.Done()
		}()
	}

	go func() {
		time.Sleep(time.Millisecond * 100)
		for j := range numMessages {
			sender.Publish(j)
		}
		wg.Done()
	}()

	wg.Wait()

	for i := range numRecievers {
		for j := range numMessages {
			require.Equal(t, j, (<-pipes[i]).Content)
		}
	}

	for i := range numRecievers {
		recievers[i].Close()
	}

	sender.Close()
}
