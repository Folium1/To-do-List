package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type TaskInterface interface {
	Create(newTask TaskCreateDTO) error 
}

// Making connection to db
func dbConnect() (*sql.DB, error) {
	dataSourceName, exists := os.LookupEnv("mysqlRoot")
	if !exists {
		panic("Couldn't connect to db")
	}
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return db, err
	}
	return db, nil
}

//Creates a new task
func (t *task) Create(newTask TaskCreateDTO) error {
	db,err := dbConnect()
	if err != nil {
		log.Print(err)
	}
	_, err = db.Query("INSERT INTO tasks(description, deadline) VALUES(%v, %v)",newTask.Description,newTask.Deadline)
	if err != nil {
		return err
	}
	return nil
}