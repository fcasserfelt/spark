package data

import (
	"github.com/fcasserfelt/spark/membership"
)

type DbRepo struct {
}

type DbUserRepo DbRepo

func NewDbUserRepo() *DbUserRepo {
	dbUserRepo := new(DbUserRepo)
	return dbUserRepo
}

func (userRepo DbUserRepo) Add(*membership.User) (id int) {
	return 1
}

func (userRepo DbUserRepo) GetByEmail(email string) *membership.User {
	if email == "existing@example.com" {
		return membership.NewUser("existing@example.com")
	}
	return nil

}
