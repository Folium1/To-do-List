package db

import (
	"fmt"
	"log"
	"todo/db"

	_ "github.com/go-sql-driver/mysql"
)

// Returns a pointer to initiated new Task struct.
func NewService() TaskService {
	return &Task{}
}

type TaskService interface {
	Create(newTask Task) error
	Tasks() ([]Task, error)
	DeleteTask(taskId string) error
	GetTask(id string) (Task, error)
	SaveTask(task Task) error
}

// Returns a slice of all tasks
func (t *Task) Tasks() ([]Task, error) {
	db, err := db.DbConnect("/tasks")
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
	// Itarate through the query
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

// Adds a new task to the database.
func (t *Task) Create(newTask Task) error {
	db, err := db.DbConnect("/tasks")
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

// Removes a task from the database.
func (t *Task) DeleteTask(taskId string) error {
	db, err := db.DbConnect("/tasks")
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

// Returns task by it's id.
func (t *Task) GetTask(id string) (Task, error) {
	db, err := db.DbConnect("/tasks")
	if err != nil {
		return Task{}, err
	}
	query, err := db.Query(fmt.Sprintf("SELECT description,deadline,is_done FROM tasks WHERE task_id = '%v';", id))
	if err != nil {
		return Task{}, err
	}
	var task Task
	task.Id = id
	// Itarate through the query
	for query.Next() {
		err := query.Scan(&task.Description, &task.Deadline, &task.Done)
		if err != nil {
			return Task{}, err
		}
	}
	if task.Id == "" {
		panic("wrong task id")
	}
	return task, nil
}

// Saves updated task's data.
func (t *Task) SaveTask(task Task) error {
	db, err := db.DbConnect("/tasks")
	if err != nil {
		return err
	}
	defer db.Close()
	q := fmt.Sprintf("UPDATE tasks SET is_done = 0, description= '%v', deadline='%v' WHERE task_id = '%v';", task.Description, task.Deadline, task.Id)
	_, err = db.Query(q)
	if err != nil {
		return err
	}
	return nil
}
