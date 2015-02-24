package membership

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"testing"
	"time"
)

type stubAuthRepo struct{}

func (userRepo stubAuthRepo) Add(*User) (id int) {
	return 1
}

func (userRepo stubAuthRepo) GetByEmail(email string) *User {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("secret"), 10)
	if err != nil {
		panic(err)
	}

	if email == "existing@example.com" {
		user := User{1, "existing@example.com", time.Now(), string(hashedPassword)}
		return &user
	}

	return nil
}

var stub stubAuthRepo

func TestNoEmail(t *testing.T) {
	a := NewAuthenticationManager(stub)

	_, _, err := a.Login("", "secret")

	if err == nil {
		t.Errorf("expected login to fail")
	}
}

func TestNoPassword(t *testing.T) {
	a := NewAuthenticationManager(stub)

	_, _, err := a.Login("existing@example.com", "")

	if err == nil {
		t.Errorf("expected login to fail")
	}
}

func TestNoUser(t *testing.T) {
	a := NewAuthenticationManager(stub)

	_, _, err := a.Login("non-existing@example.com", "secret")

	if err == nil {
		t.Errorf("expected login to fail")
	}
}

func TestPasswordMissMatch(t *testing.T) {
	a := NewAuthenticationManager(stub)

	_, _, err := a.Login("existing@example.com", "incorrec")

	if err == nil {
		t.Errorf("expected login to fail")
	}
}

func TestSuccess(t *testing.T) {
	a := NewAuthenticationManager(stub)

	user, _, err := a.Login("existing@example.com", "secret")

	if err != nil {
		t.Errorf("expected login to succeed %v", err)
		return
	}

	if user.Email != "existing@example.com" {
		t.Errorf("expected correct user got %v", user)
	}
}

func TestToken(t *testing.T) {
	var publicKey []byte
	publicKey, _ = ioutil.ReadFile("./timeapp.rsa.pub")

	var stub stubAuthRepo
	a := NewAuthenticationManager(stub)

	_, myToken, _ := a.Login("existing@example.com", "secret")

	if len(myToken) == 0 {
		t.Error("Expected a token got %s", myToken)
		return
	}

	token, err := jwt.Parse(myToken, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		t.Errorf("Could not parse token", err)
		return
	}

	if !token.Valid {
		t.Errorf("Token not valid")
	}

	if int(token.Claims["ID"].(float64)) != 1 {
		t.Errorf("Expected user id 1 got %v", int(token.Claims["ID"].(float64)))
	}

}
