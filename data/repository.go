package data

import (
	"database/sql"
	"github.com/fcasserfelt/spark/membership"
	_ "github.com/lib/pq"
	"log"
)

type DbRepo struct {
	db *sql.DB
}

type DbUserRepo DbRepo

func NewDbUserRepo(db *sql.DB) *DbUserRepo {
	dbUserRepo := new(DbUserRepo)
	dbUserRepo.db = db
	return dbUserRepo
}

func (repo DbUserRepo) Add(u *membership.User) (id int) {

	var userid int
	err := repo.db.QueryRow(`INSERT INTO users(email, password)
	VALUES($1, $2) RETURNING id`, u.Email, u.Password).Scan(&userid)
	if err != nil {
		log.Printf("Sql error: %v", err)
		return 0
	}
	return userid
}

func (repo *DbUserRepo) GetByEmail(email string) *membership.User {

	var user membership.User
	err := repo.db.QueryRow("select id, email, password from users where email = $1", email).Scan(&user.Id, &user.Email, &user.Password)
	if err != nil {
		log.Printf("Sql error: %v", err)
		return nil
	}
	return &user
	/*
		fmt.Println(name)

		row, err := repo.db.Query("SELECT id, email FROM users WHERE email = $1", email)
		if err != nil {
			fmt.Printf("error: %v", err)
		}

		for rows.Next() {
			err := rows.Scan(&id, &name)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(id, name)
		}

		var id int
		row.Next()
		row.Scan(&id, &email)
		user := membership.User{Id: id, Email: email}
		return &user
	*/
}
