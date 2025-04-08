package chat

type Message struct {
	Username string
	Content  interface{}
}

type Topic interface {
	Title() string

	// Registers user in this topic
	//
	// All the messages sent by ALL users will
	// be recieved by returned one
	RegisterUser(name string, bufSize int) User

	// Close topic - remove all the users
	Close()

	// Sent message to everyone in topic
	Publish(Message)

	// Remove user from listening to messages
	// and removing the ability to send them
	RemoveUser(user User)
}

type myTopic struct {
	title string

	d Distributor
}

func (t myTopic) Publish(msg Message) {
	t.d.publish(msg)
}

func (t myTopic) RemoveUser(user User) {
	t.d.removeUser(user)
}

func (t myTopic) RegisterUser(name string, bufSize int) User {
	user := User{
		channel: make(chan Message, bufSize),
		name:    name,
		topic:   t,
	}

	t.d.registerUser(user)

	return user
}

func (t myTopic) Title() string {
	return t.title
}

func (t myTopic) Close() {
	t.d.close()
}

func NewTopic(title string, broadcasterCount int) Topic {
	return &myTopic{
		title: title,
		d:     newDistributor(broadcasterCount),
	}
}
