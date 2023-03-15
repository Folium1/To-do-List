package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// DbConnect establishes a connection to the database.
func DbConnect() (*sql.DB, error) {
	// Load environment variables from .env file.
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	dataSourceName := os.Getenv("DB_SOURCE")
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return db, err
	}
	return db, nil
}

// DbTableInit creates the tasks table in the database if it does not already exist.
func DbTableInit() error {
	db, err := DbConnect()
	if err != nil {
		log.Printf("couldn't connect to db, err: %v", err)
		return err
	}
	
	// creating todo schema
	_, err = db.Query("CREATE SCHEMA IF NOT EXISTS todo;")
	if err != nil {
		log.Printf("Couldn't create todo schema")
		return err
	}

	// creating users table
	_, err = db.Query("CREATE TABLE IF NOT EXISTS todo.users(user_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,user_name VARCHAR(50),mail VARCHAR(50) NOT NULL UNIQUE,password VARCHAR(200));")
	if err != nil {
		log.Printf("Couldn't create users table")
		return err
	}

	// creating idx_user_id if it is not created
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM information_schema.statistics WHERE table_name ='users' AND index_name = 'idx_user_id'").Scan(&count)
	if err != nil {
		log.Printf("Error checking if index exists: %v", err)
		return err
	}
	if count == 0 {
		_, err = db.Query("CREATE INDEX idx_user_id ON todo.users(user_id);")
		if err != nil {
			log.Printf("Couldn't create index")
			return err
		}
	}
	// creting tasks table
	_, err = db.Query("CREATE TABLE IF NOT EXISTS todo.tasks(task_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,user_id INT,description VARCHAR(50),deadline DATETIME,is_done TINYINT(1),CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id));")
	if err != nil {
		log.Printf("Couldn't create tasks table")
		return err
	}

	return nil
}
