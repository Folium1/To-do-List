package users

import (
	"fmt"
	"log"
	"strconv"
	"todo/db"
)

type Service interface {
	CreateUser(newUser User) (int, error)
	GetUser(user User) (User, error)
}

func New() Service {
	return &User{}
}

// CreateUser creates new user or returns created user's id or error
func (u *User) CreateUser(newUser User) (int, error) {
	db, err := db.DbConnect()
	if err != nil {
		log.Print(err)
		return 0, err
	}
	defer db.Close()
	query := fmt.Sprintf("INSERT INTO todo.users(mail,user_name,password) VALUES('%v','%v','%v');", newUser.Mail, newUser.Name, newUser.Password)
	_, err = db.Query(query)
	if err != nil {
		log.Print(err)
		return 0, err
	}
	var userId string
	rows, err := db.Query("SELECT LAST_INSERT_ID();")
	if err != nil {
		log.Print(err)
		return 0, err
	}
	for rows.Next() {
		err := rows.Scan(&userId)
		if err != nil {
			log.Printf("Couldn't parse data from rows, err: %v", err)
			return 0, err
		}
	}
	log.Println(userId)
	intId, _ := strconv.Atoi(userId)
	return intId, nil
}

// GetUser returns user's id and password by it's mail
func (u *User) GetUser(user User) (User, error) {
	db, err := db.DbConnect()
	if err != nil {
		log.Print(err)
		return User{}, err
	}
	defer db.Close()
	query := fmt.Sprintf("SELECT user_id,password FROM todo.users WHERE mail = '%v'", user.Mail)
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return User{}, err
	}
	for rows.Next() {
		err = rows.Scan(&user.Id, &user.Password)
		if err != nil {
			log.Println(err)
			return User{}, err
		}
	}
	return user, nil
}
