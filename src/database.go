package src

import (
	"net/http"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/labstack/echo/v4"
)

var ErrUserNotFound = echo.NewHTTPError(http.StatusNotFound, "user not found")

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

func GetUserById(id []byte) (*User, error) {
	for _, u := range Users {
		if string(u.ID) == string(id) {
			return u, nil
		}
	}
	return nil, ErrUserNotFound
}
