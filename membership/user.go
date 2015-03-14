package membership

import (
	"time"
)

type User struct {
	Id       int
	Email    string
	Created  time.Time
	Password string
}

type UserRepo interface {
	GetByEmail(email string) (*User, error)
	Add(*User) (int, error)
}

func NewUser(email string) *User {
	return &User{
		Email:   email,
		Created: time.Now(),
	}
}
