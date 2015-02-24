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
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/register", RegisterHandler)

	n := negroni.Classic()
	n.UseHandler(router)
	n.Run("localhost:3000")
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	//user := userRepo.GetByEmail("fredrik@bitjoy.se")
	//fmt.Fprintf(w, "id: %d email:%s", user.Id, user.Email)

	fmt.Fprintf(w, "I'm alive")
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
