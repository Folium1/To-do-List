package users

import (
	"fmt"
	"log"
	"todo/db"
)

type Service interface {
	CreateUser(newUser User) (string,error)
}

func New() Service {
	return &User{}
}

// CreateUser creates new user or returns created user's id or error
func (u *User) CreateUser(newUser User) (string,error) {
	db, err := db.DbConnect("/users")
	if err != nil {
		log.Print(err)
		return "",err
	}
	query := fmt.Sprintf("INSERT INTO users(mail,name,password) VALUES(%v,%v,%v); SELECT LAST_INSERT_ID();", newUser.Mail, newUser.Name, newUser.Password)
	rows, err := db.Query(query)
	if err != nil {
		log.Print(err)
		return "",err
	}
	var userId string 
	for rows.Next() {
		rows.Scan(&userId)
	}
	return userId,nil
}
