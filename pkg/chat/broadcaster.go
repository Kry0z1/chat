package chat

import "sync"

type Broadcaster interface {
	start()
	registerUser(User)
	removeUser(User)
	pass(Message)
	close()
}

type myBroadcaster struct {
	msgChan chan Message

	// Map like this because User doesn't want to be hashable
	users map[string]chan Message

	usersLock *sync.Mutex
}

// Expects to run in goroutine - else blocks
func (b myBroadcaster) start() {
	for msg := range b.msgChan {
		b.usersLock.Lock()
		for user := range b.users {
			ch := b.users[user]
			select {
			case ch <- msg:
			default:
				<-ch
				ch <- msg
			}
		}
		b.usersLock.Unlock()
	}
}

func (b myBroadcaster) registerUser(user User) {
	b.usersLock.Lock()
	b.users[user.name] = user.channel
	b.usersLock.Unlock()
}

func (b myBroadcaster) removeUser(user User) {
	b.usersLock.Lock()
	delete(b.users, user.name)
	b.usersLock.Unlock()
}

func (b myBroadcaster) close() {
	close(b.msgChan)
}

func (b myBroadcaster) pass(msg Message) {
	b.msgChan <- msg
}

func newBroadcaster() Broadcaster {
	return myBroadcaster{
		msgChan:   make(chan Message),
		users:     make(map[string]chan Message),
		usersLock: &sync.Mutex{},
	}
}
