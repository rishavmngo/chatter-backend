package entities

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

type User struct {
	UID          uint      `json:"id"`
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	CreatedAt    time.Time `json:"created_at"`
	LastActiveAt time.Time `json:"last_active_at"`
}

type UserController interface {
	Register(http.ResponseWriter, *http.Request)
	Login(http.ResponseWriter, *http.Request)
	GetUserByUID(http.ResponseWriter, *http.Request)
}

func (user *User) Created() {
	user.CreatedAt = time.Now()
	user.LastActiveAt = time.Now()
}

func (user *User) LastActive() {
	user.LastActiveAt = time.Now()
}

func (user *User) Validate(action string) error {

	switch strings.ToLower(action) {
	case "login":

		if user.Password == "" {
			return errors.New("Required Password")
		}
		if user.Username == "" {
			return errors.New("Required Username")
		}
		return nil
	default:
		if user.Password == "" {
			return errors.New("Required Password")
		}
		if user.Username == "" {
			return errors.New("Required Username")
		}
		return nil
	}

}
