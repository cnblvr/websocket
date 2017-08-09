package chat

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type User struct {
	id                uint64
	name              string
	colorName         string
	ws                *websocket.Conn
	server            *Server
	chanSendMessage   chan *Message
	chanSendMessages  chan []*Message
	chanSendListUsers chan *Users
	chanSendError     chan *Error
}

func NewUser(id uint64, ws *websocket.Conn, server *Server) *User {

	return &User{
		id:                id,
		name:              "",
		colorName:         "black",
		ws:                ws,
		server:            server,
		chanSendMessage:   make(chan *Message),
		chanSendMessages:  make(chan []*Message),
		chanSendListUsers: make(chan *Users),
		chanSendError:     make(chan *Error),
	}
}

func (u *User) SetName(newName string) *Error {

	if u.server.users.AvailableName(newName) {
		u.name = newName
	} else {
		return NewError(5)
	}
	return nil
}

func (u *User) SetColorName(newColor string) {

	u.colorName = newColor
}

func (u *User) Listen() {

	go u.listenWrite()
	go u.listenRead()
}

func (u *User) listenWrite() {

	for {
		select {
		case msg := <-u.chanSendMessage:
			jdata := jsonNewMessage{
				Author:  msg.author.id,
				Message: msg.message,
				Date:    msg.date.Format(timeFormat),
			}
			if data, err := json.Marshal(jdata); err != nil {
				log.Panic(err)
			} else {
				u.ws.WriteJSON(jsonPacket{
					Command: "new_message",
					Data:    data,
				})
			}

		case messages := <-u.chanSendMessages:
			if len(messages) == 0 {
				break
			}
			var jdata []jsonNewMessage
			for _, msg := range messages {
				jdata = append(jdata, jsonNewMessage{
					Author:          msg.author.id,
					AuthorName:      msg.author.name,
					AuthorColorName: msg.author.colorName,
					Message:         msg.message,
					Date:            msg.date.Format(timeFormat),
				})
			}
			if data, err := json.Marshal(jdata); err != nil {
				log.Panic(err)
			} else {
				u.ws.WriteJSON(jsonPacket{
					Command: "past_messages",
					Data:    data,
				})
			}

		case users := <-u.chanSendListUsers:
			jdata := []jsonUser{}
			for _, user := range users.GetList() {
				jdata = append(jdata, jsonUser{
					Id:        user.id,
					Name:      user.name,
					NameColor: user.colorName,
				})
			}
			if data, err := json.Marshal(jdata); err != nil {
				log.Panic(err)
			} else {
				u.ws.WriteJSON(jsonPacket{
					Command: "list_users",
					Data:    data,
				})
			}

		case e := <-u.chanSendError:
			jdata := jsonError{
				Code:    e.code,
				Message: e.message,
			}
			if data, err := json.Marshal(jdata); err != nil {
				log.Panic(err)
			} else {
				u.ws.WriteJSON(jsonPacket{
					Command: "error",
					Data:    data,
				})
			}
		}
	}
}

func (u *User) listenRead() {

	for {
		select {
		default:
			var jpack jsonPacket
			err := u.ws.ReadJSON(&jpack)
			if _, ok := err.(*websocket.CloseError); ok {
				u.server.DisconnectUser(u)
				return
			} else {
				switch jpack.Command {
				case "send":
					if u.name == "" {
						u.chanSendError <- NewError(1)
						break
					}
					var jmessage jsonMessage
					if err := json.Unmarshal([]byte(jpack.Data), &jmessage); err != nil {
						log.Panic(err)
					} else {
						log.Printf("%d: %s", u.id, jmessage.Message)
						// добавление сообщения в список сервера
						if message, err := u.server.messages.Add(u, jmessage.Message); err != nil {
							log.Fatal(err)
						} else {
							u.server.chanSendAllMessage <- message // отправка всем пользователям
						}
					}
				case "set_profile":
					var jprofile jsonSetProfile
					if err := json.Unmarshal([]byte(jpack.Data), &jprofile); err != nil {
						log.Panic(err)
					} else {
						log.Printf("пользователь %d меняет профиль: %#v", u.id, jprofile)
						if jprofile.Name != "" {
							if err := u.SetName(jprofile.Name); err != nil {
								u.chanSendError <- err
								break
							}
						}
						if jprofile.Color != "" {
							u.SetColorName(jprofile.Color)
						}
						u.server.chanSendAllListUsers <- u.server.users
					}

				default:
					log.Panic("неизвестная команда")
				}
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

func (uu *Users) Add(ws *websocket.Conn, server *Server) (*User, error) {

	uu.maxId++
	user := NewUser(uu.maxId, ws, server)
	//todo повторение имени

	uu.users[uu.maxId] = user

	return user, nil
}

func (uu *Users) SendAllMessage(message *Message) {

	for _, u := range uu.users {
		u.chanSendMessage <- message
	}
}

func (uu *Users) SendAllListUsers(users *Users) {

	for _, u := range uu.users {
		u.chanSendListUsers <- users
	}
}

func (uu *Users) GetList() []*User {

	users := []*User{}
	for _, user := range uu.users {
		users = append(users, user)
	}

	return users
}

func (uu *Users) Delete(id uint64) {

	delete(uu.users, id)
	// замечание: сам объект пользователя не удалиться из памяти (гарбыч коллектор не удалит его), т.к. на него есть
	// ссылки в сообщениях
	// можно сделать так:
	//uu.users[id].ws = nil
	//uu.users[id].server = nil
}

func (uu *Users) AvailableName(name string) bool {

	for _, user := range uu.users {
		if user.name == name {
			return false
		}
	}
	return true
}
