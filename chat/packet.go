package chat

import "encoding/json"

type jsonPacket struct {
	Command string          `json:"command"`
	Data    json.RawMessage `json:"data"`
}

// from user
type jsonMessage struct {
	Message string `json:"message"`
}

// from user
type jsonSetProfile struct {
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
}

// to user
type jsonNewMessage struct {
	Author          uint64 `json:"author"`
	AuthorName      string `json:"author_name,omitempty"`
	AuthorColorName string `json:"author_color_name,omitempty"`
	Message         string `json:"message"`
	Date            string `json:"date"`
}

// to user
type jsonUser struct {
	Id        uint64 `json:"id"`
	Name      string `json:"name"`
	NameColor string `json:"name_color"`
}

// to user
type jsonError struct {
	Code    uint16 `json:"code"`
	Message string `json:"message"`
}
