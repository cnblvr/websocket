package chat

type Error struct {
	code    uint16
	message string
}

func NewError(code uint16) *Error {

	return &Error{
		code:    code,
		message: Errors[code],
	}
}

var Errors = map[uint16]string{
	1: "Пользование чатом запрещено. Установи имя пользователя.",
	5: "Пользователь с таким именем уже зарегистрирован.",
}
