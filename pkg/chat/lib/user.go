package chat

// User channel is working like queue:
//
// If buffer is full - first message is deleted
type User struct {
	channel chan Message
	name    string

	topic Topic
}

// Post message in name of user
func (u User) Publish(content interface{}) {
	msg := Message{Content: content, Username: u.name}
	u.topic.Publish(msg)
}

// Remove user from recieving
func (u User) Close() {
	u.topic.RemoveUser(u)
	close(u.channel)
}

// Returns channel from which messages are recieved
func (u User) Recieve() chan Message {
	return u.channel
}

// Returns title of current topic
func (u User) Topic() string {
	return u.topic.Title()
}
