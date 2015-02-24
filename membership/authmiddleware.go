package membership

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	repo UserRepo
}

func NewAuthMiddleware(repo UserRepo) *AuthMiddleware {
	return &AuthMiddleware{repo}
}

func (m *AuthMiddleware) Handler(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	err := m.CheckJWT(w, r)

	if err == nil && next != nil {
		next(w, r)
	}

}

func (m *AuthMiddleware) CheckJWT(w http.ResponseWriter, r *http.Request) error {

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		msg := "No Authorization header"
		http.Error(w, msg, http.StatusUnauthorized)
		return errors.New(msg)
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		msg := "Invalid token format: Must be Bearer {token}"
		http.Error(w, msg, http.StatusUnauthorized)
		return errors.New(msg)
	}

	userToken := authHeaderParts[1]

	var publicKey []byte
	publicKey, _ = ioutil.ReadFile("./timeapp.rsa.pub")

	parsedToken, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return err
	}

	if !parsedToken.Valid {
		//	msg := "Invalid token"
		//	http.Error(w, msg, http.StatusUnauthorized)
		//	return errors.New(msg)
	}

	return nil

}
