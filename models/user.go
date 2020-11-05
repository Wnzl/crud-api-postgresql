package models

import (
	"errors"
	"net/http"
)

type User struct {
	ID    int    `json:"id"`
	First string `json:"first"`
	Last  string `json:"last"`
}

type UserStorage interface {
	Get(id int) (*User, error)
	GetAll() ([]*User, error)
	Store(user *User) (id int, err error)
	Update(id int, userData *User) (*User, error)
	Delete(id int) error
	UserExist(user *User) (bool, error)
}

// validating request payload
func (u *User) Bind(r *http.Request) error {
	if u.First == "" {
		return errors.New("missing required \"first\" fields")
	}

	if u.Last == "" {
		return errors.New("missing required \"last\" fields")
	}

	return nil
}
