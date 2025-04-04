package internal

type distributor interface {
	publish(Message)
	removeUser(User)
	registerUser(User)
	close()
}

type myDistributor struct {
	broadcasters []broadcaster
	brQueue      []int

	messages chan Message

	userToBr map[string]int
}

func (d myDistributor) publish(msg Message) {
	for _, b := range d.broadcasters {
		b.pass(msg)
	}
}

func (d myDistributor) removeUser(user User) {
	brIndex := d.userToBr[user.name]
	d.brQueue = append(d.brQueue, brIndex)
	d.broadcasters[brIndex].removeUser(user)
}

func (d myDistributor) registerUser(user User) {
	if len(d.brQueue) == 0 {
		d.brQueue = make([]int, len(d.broadcasters))
		for i := range d.broadcasters {
			d.brQueue[i] = i
		}
	}

	brIndex := d.brQueue[0]

	// This is kinda fine because taking slice is free
	// and sometimes queue will be reallocated so
	// not so much memory will be grabbed
	d.brQueue = d.brQueue[:1]

	d.broadcasters[brIndex].registerUser(user)
	d.userToBr[user.name] = brIndex
}

func (d myDistributor) close() {
	for _, br := range d.broadcasters {
		br.close()
	}
	close(d.messages)
}

func newDistributor(broadcasterCount int) myDistributor {
	result := myDistributor{
		broadcasters: make([]broadcaster, broadcasterCount),
		brQueue:      make([]int, 0),
		messages:     make(chan Message),
		userToBr:     make(map[string]int),
	}

	for i := range result.broadcasters {
		result.broadcasters[i] = newBroadcaster()
		go result.broadcasters[i].start()
	}

	return result
}
