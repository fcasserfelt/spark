package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/fcasserfelt/spark/data"
	"github.com/fcasserfelt/spark/membership"
	"github.com/gorilla/mux"
	"net/http"
)

var userRepo membership.UserRepo

func init() {
	userRepo = data.NewDbUserRepo()
}

func main() {

	n := negroni.Classic()

	authMiddleware := membership.NewAuthMiddleware(userRepo)

	router := mux.NewRouter()
	apiRoutes := mux.NewRouter()

	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/register", RegisterHandler)

	apiRoutes.HandleFunc("/secured", SecuredHandler)

	apiRoutes.HandleFunc("/secured/ping", SecuredPingHandler)

	router.PathPrefix("/secured").Handler(negroni.New(
		negroni.HandlerFunc(authMiddleware.Handler),
		negroni.Wrap(apiRoutes),
	))

	n.UseHandler(router)
	n.Run("localhost:3000")
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "I'm alive")
}

func SecuredHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "I'm secure")
}

func SecuredPingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "I'm secure ping")
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	var email, password, confirm string

	email = "new@example.com"
	password = "secret"
	confirm = "secret"

	var a *membership.Application
	var reg *membership.RegistrationManager

	a = membership.NewApplication(email, password, confirm)
	reg = membership.NewRegistrationManager(userRepo)

	user, token, err := reg.Apply(a)

	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	fmt.Fprintf(w, "Register user: %v token: %v", user, token)
}
