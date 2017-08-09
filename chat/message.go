package chat

import "time"

type Message struct {
	id      uint64
	author  *User
	message string
	date    time.Time
}

func NewMessage(id uint64, author *User, message string) *Message {
	return &Message{
		id:      id,
		author:  author,
		message: message,
		date:    time.Now(),
	}
}

type Messages struct {
	messages []*Message
	maxId    uint64
}

func NewMessages() *Messages {
	return &Messages{
		messages: make([]*Message, 0),
		maxId:    0,
	}
}

func (mm *Messages) Add(author *User, messageString string) (*Message, error) {

	mm.maxId++
	message := NewMessage(mm.maxId, author, messageString)

	mm.messages = append(mm.messages, message)

	return message, nil
}

func (mm *Messages) GetPast() []*Message {

	if len(mm.messages) == 0 {
		return []*Message{}
	} else if len(mm.messages) > quantityPastMessages {
		return mm.messages[len(mm.messages)-quantityPastMessages:]
	} else {
		return mm.messages
	}
}
