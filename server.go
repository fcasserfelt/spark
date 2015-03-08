package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/fcasserfelt/spark/data"
	"github.com/fcasserfelt/spark/membership"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"os"
)

var userRepo membership.UserRepo

func init() {

	var DB_ENV_POSTGRES_USER = os.Getenv("DB_ENV_POSTGRES_USER")
	var DB_ENV_POSTGRES_PASSWORD = os.Getenv("DB_ENV_POSTGRES_PASSWORD")
	var DB_PORT_5432_TCP_ADDR = os.Getenv("DB_PORT_5432_TCP_ADDR")

	var err error
	var db *sql.DB

	s := fmt.Sprintf("user=%s password=%s dbname=spark sslmode=disable host=%s", DB_ENV_POSTGRES_USER, DB_ENV_POSTGRES_PASSWORD, DB_PORT_5432_TCP_ADDR)
	fmt.Println(s)

	//	db, err = sql.Open("postgres", "user=%s password=%s dbname=spark sslmode=disable host=%s", DB_ENV_POSTGRES_USER, DB_ENV_POSTGRES_PASSWORD, DB_PORT)

	db, err = sql.Open("postgres", s)
	if err != nil {
		panic(err)
	}
	userRepo = data.NewDbUserRepo(db)

}

func main() {

	n := negroni.Classic()

	authMiddleware := membership.NewAuthMiddleware(userRepo)

	router := mux.NewRouter()
	securedRoutes := mux.NewRouter()

	router.HandleFunc("/", HomeHandler)
	router.HandleFunc("/register", RegisterHandler)
	router.HandleFunc("/token", TokenHandler)

	securedRoutes.HandleFunc("/secured", SecuredHandler)
	securedRoutes.HandleFunc("/secured/ping", SecuredPingHandler)

	router.PathPrefix("/secured").Handler(negroni.New(
		negroni.HandlerFunc(authMiddleware.Handler),
		negroni.Wrap(securedRoutes),
	))

	n.UseHandler(router)
	n.Run(":3000")
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

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	var email, password string

	email = r.PostFormValue("email")
	password = r.PostFormValue("password")

	var a = membership.NewAuthenticationManager(userRepo)
	_, token, err := a.Login(email, password)

	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	fmt.Fprint(w, token)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	var a membership.Application
	err = json.Unmarshal(body, &a)

	var reg *membership.RegistrationManager

	reg = membership.NewRegistrationManager(userRepo)

	user, token, err := reg.Apply(&a)

	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	fmt.Fprintf(w, "Register user: %v token: %v", user, token)
}
