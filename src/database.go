package src

import (
	"errors"

	"github.com/go-webauthn/webauthn/webauthn"
)

var ErrUserNotFound = errors.New("user not found")

var Sessions = map[string]*webauthn.SessionData{}
var Users = map[string]*User{}

func InsertSession(id string, session *webauthn.SessionData) {
	Sessions[id] = session
}

func GetSession(id string) (*webauthn.SessionData, error) {
	s, ok := Sessions[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return s, nil
}

func InsertUser(id string, user *User) {
	Users[id] = user
}

func GetUser(name string) (*User, error) {
	u, ok := Users[name]
	if !ok {
		return nil, ErrUserNotFound
	}
	return u, nil
}
