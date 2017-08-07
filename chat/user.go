package chat

import (
	"golang.org/x/net/websocket"
	"log"
)

type User struct {
	id     uint64
	name   string
	ws     *websocket.Conn
	server *Server
}

func NewUser(id uint64, name string, ws *websocket.Conn, server *Server) *User {

	return &User{
		id:     id,
		name:   name,
		ws:     ws,
		server: server,
	}
}

func (u *User) SetName(newName string) {

	u.name = newName
}

func (u *User) Listen() {

	go u.listenWrite()
	u.listenRead()
}

func (u *User) listenWrite() {

	for {
		select {}
	}
}

func (u *User) listenRead() {

	for {
		select {
		default:
			var pack Pack
			if err := websocket.JSON.Receive(u.ws, &pack); err != nil {
				log.Panic(err) //todo err == io.EOF
			} else {
				_ = pack //todo !!!!!
				log.Printf("user message: %#v", pack)
				//todo u.server. send all users pack
			}
		}
	}
}

type Users struct {
	users map[uint64]*User
	maxId uint64
}

func NewUsers() *Users {

	return &Users{
		users: make(map[uint64]*User),
		maxId: 0,
	}
}

func (uu *Users) Add(name string, ws *websocket.Conn, server *Server) (*User, error) {

	uu.maxId++
	user := NewUser(uu.maxId, name, ws, server)
	//todo повторение имени

	uu.users[uu.maxId] = user

	return user, nil
}
