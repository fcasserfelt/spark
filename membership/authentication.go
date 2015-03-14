package membership

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"time"
)

var (
	privateKey []byte
)

func init() {
	privateKey, _ = ioutil.ReadFile("./timeapp.rsa")
}

type AuthenticationManager struct {
	repo UserRepo
}

func NewAuthenticationManager(repo UserRepo) *AuthenticationManager {
	return &AuthenticationManager{repo}
}

func (a AuthenticationManager) Login(email string, password string) (*User, string, error) {

	if len(email) == 0 {
		return nil, "", fmt.Errorf("Email required")
	}

	if len(password) == 0 {
		return nil, "", fmt.Errorf("Password required")
	}

	user, err := a.repo.GetByEmail(email)

	if err != nil {
		return nil, "", err
	}

	if user == nil {
		return nil, "", fmt.Errorf("User not found")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, "", fmt.Errorf("Invalid password")
	}

	tokenString, err := GenerateToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("Could not generate token %v", err)
	}

	return user, tokenString, nil
}

func GenerateToken(user *User) (string, error) {
	token := jwt.New(jwt.GetSigningMethod("RS256"))

	token.Claims["ID"] = user.Id
	token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	tokenString, err := token.SignedString(privateKey)

	if err != nil {
		return "", fmt.Errorf("", err)
	}

	return tokenString, nil

}
