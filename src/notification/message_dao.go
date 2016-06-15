package notification

// create a message
func NewMessage() *Message {
	message := &Message{}
	return message
}

// load unsent messages from storage
func LoadMessages() []Message {
	return nil
}

// persist a message into storage
func (message *Message) Persist() error {
	return nil
}

// remove a message permanantly from storage
func (message *Message) Remove() error {
	return nil
}
