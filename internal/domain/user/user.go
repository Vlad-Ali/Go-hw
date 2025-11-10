package user

import (
	usererror "transfersservice/internal/domain/user/error"
	"transfersservice/internal/domain/user/object"
)

type User struct {
	username string
	password string
	email    string
	id       object.UserID
}

func NewUser(username string, password string, email string) *User {
	return &User{username, password, email, object.UserID{}}
}

func (u *User) Username() string {
	return u.username
}

func (u *User) Password() string {
	return u.password
}

func (u *User) Email() string {
	return u.email
}

func (u *User) ID() object.UserID {
	return u.id
}

func (u *User) SetID(id object.UserID) error {
	if u.id.IsEmpty() {
		u.id = id
		return nil
	}
	return usererror.ErrUserIDAlreadyExists
}
