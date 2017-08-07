package chat

import (
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

const websocketPath = "/websocket"

type Server struct {
	users *Users
}

func NewServer() *Server {
	return &Server{
		users: NewUsers(),
	}
}

func (s *Server) Listen() {

	onConnected := func(ws *websocket.Conn) {
		log.Printf("connected")
		// add user to server
		if user, err := s.users.Add("userQwe", ws, s); err != nil {
			log.Panic(err)
		} else {
			user.Listen()
		}
	}

	http.Handle(websocketPath, websocket.Handler(onConnected))

	for {
		select {}
	}
}
