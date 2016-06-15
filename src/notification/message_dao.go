package notification

// create a message
func NewMessage() *Message {
	message := &Message{}
	return message
}

// load unsent messages from storage
func LoadMessages() []Message {
	// initilize all message marked as persisted
	return nil
}

// persist a message into storage
func (message *Message) Persist() error {
	message.Persisted = true
	return nil
}

// remove a message permanantly from storage
func (message *Message) Remove() error {
	return nil
}
