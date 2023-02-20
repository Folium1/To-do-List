package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func NewService() TaskService {
	return &Task{}
}

// Making connection to db
func dbConnect() (*sql.DB, error) {
	err := godotenv.Load()
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

type TaskService interface {
	Create(newTask TaskCreateDTO) error
	Tasks() ([]Task, error)
	DeleteTask(taskId string) error
	GetTask(id string) (Task, error)
	SaveTask(task Task) error
}

// Get all tasks
func (t *Task) Tasks() ([]Task, error) {
	db, err := dbConnect()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	defer db.Close()
	query, err := db.Query("SELECT task_id,description, deadline, is_done FROM tasks ORDER BY deadline;")
	if err != nil {
		fmt.Println(err)
		log.Print(err)
		return nil, err
	}
	var allTasks []Task
	var oneTask Task

	for query.Next() {
		err := query.Scan(&oneTask.Id, &oneTask.Description, &oneTask.Deadline, &oneTask.Done)
		if err != nil {
			log.Print(err)
			return nil, err
		}
		allTasks = append(allTasks, oneTask)
	}

	return allTasks, nil
}

// Create a new task
func (t *Task) Create(newTask TaskCreateDTO) error {
	db, err := dbConnect()
	if err != nil {
		log.Print(err)
	}
	defer db.Close()
	q := fmt.Sprintf("INSERT INTO tasks(description, deadline, is_done) VALUES('%v', '%v', false);", newTask.Description, newTask.Deadline)
	_, err = db.Query(q)
	if err != nil {
		return err
	}
	return nil
}

// DeleteTask takes task by it's id
func (t *Task) DeleteTask(taskId string) error {
	db, err := dbConnect()
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Query(fmt.Sprintf("DELETE FROM tasks WHERE task_id = '%v';", taskId))
	if err != nil {
		log.Printf("couldn't delete task: %v", err)
		return err
	}
	return nil
}

// Returns task by id or error
func (t *Task) GetTask(id string) (Task, error) {
	db, err := dbConnect()
	if err != nil {
		return Task{}, err
	}
	query, err := db.Query(fmt.Sprintf("SELECT description,deadline,is_done FROM tasks WHERE task_id = '%v';", id))
	if err != nil {
		return Task{}, err
	}
	var task Task
	task.Id = id
	for query.Next() {
		err := query.Scan(&task.Description, &task.Deadline, &task.Done)
		if err != nil {
			return Task{}, err
		}
	}
	if task.Id == "" {
		panic("wrong task id db")
	}
	return task, nil
}

// Changes task data by id
func (t *Task) SaveTask(task Task) error {
	db, err := dbConnect()
	if err != nil {
		return err
	}

	defer db.Close()
	var q string
	if task.Done {
		q = fmt.Sprintf("UPDATE tasks SET is_done = 1, description= '%v', deadline='%v' WHERE task_id = '%v';", task.Description, task.Deadline, task.Id)

	} else {
		q = fmt.Sprintf("UPDATE tasks SET is_done = 0, description= '%v', deadline='%v' WHERE task_id = '%v';", task.Description, task.Deadline, task.Id)
	}
	_, err = db.Query(q)
	if err != nil {
		return err
	}
	return nil
}
