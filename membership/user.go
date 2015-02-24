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
	GetByEmail(email string) *User
	Add(*User) int
}

func NewUser(email string) *User {
	return &User{
		Email:   email,
		Created: time.Now(),
	}
}
