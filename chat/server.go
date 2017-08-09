package chat

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

const websocketPath = "/websocket"
const timeFormat = "02.01.2006 15:04:05.000"
const quantityPastMessages = 10

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Server struct {
	users                *Users
	messages             *Messages
	chanSendAllMessage   chan *Message // послать всем сообщение
	chanSendPastMessages chan *User    // послать последние сообщения
	chanSendAllListUsers chan *Users
}

func NewServer() *Server {
	return &Server{
		users:                NewUsers(),
		messages:             NewMessages(),
		chanSendAllMessage:   make(chan *Message),
		chanSendPastMessages: make(chan *User),
		chanSendAllListUsers: make(chan *Users),
	}
}

func (s *Server) Listen() {

	http.HandleFunc(websocketPath, s.ConnectUser)

	for {
		select {
		case message := <-s.chanSendAllMessage:
			s.users.SendAllMessage(message)

		case user := <-s.chanSendPastMessages:
			user.chanSendMessages <- s.messages.GetPast()

		case users := <-s.chanSendAllListUsers:
			s.users.SendAllListUsers(users)
		}
	}
}

func (s *Server) ConnectUser(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Panic(err)
	}
	// add user to server
	if user, err := s.users.Add(ws, s); err != nil {
		log.Panic(err)
	} else {
		log.Printf("подключился пользователь %d", user.id)
		user.Listen()
		s.chanSendPastMessages <- user
	}
}

func (s *Server) DisconnectUser(user *User) {

	log.Printf("пользователь %d отключился", user.id)
	s.users.Delete(user.id)
}
