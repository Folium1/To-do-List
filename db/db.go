package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// DbConnect establishes a connection to the database.
func DbConnect(table string) (*sql.DB, error) {
	// Load environment variables from .env file.
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	dataSourceName := os.Getenv("DB_SOURCE") + table
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return db, err
	}
	return db, nil
}

// DbTableInit creates the tasks table in the database if it does not already exist.
func DbTableInit() error {
	db, err := DbConnect("")
	if err != nil {
		log.Printf("couldn't connect to db, err: %v", err)
		return err
	}
	_, err = db.Query("CREATE TABLE IF NOT EXISTS tasks (task_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,description VARCHAR(50),deadline DATETIME,is_done TINYINT(1));")
	if err != nil {
		log.Printf("Couldn't create table")
		return err
	}
	return nil
}
