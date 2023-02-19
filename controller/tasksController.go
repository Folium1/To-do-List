package controller

import (
	"errors"
	"fmt"
	"net/http"

	"todo/db"
)

type taskController struct {
	t db.TaskService
}

func New(db db.TaskService) TaskController {
	return &taskController{t: db}
}

type TaskController interface {
	Create(params *http.Request) error
	Tasks() ([]db.Task, error)
	DeleteTask(taskId string) error
	IsDone(taskId string) error
}

func (c *taskController) Tasks() ([]db.Task, error) {
	tasks, err := c.t.Tasks()
	if err != nil {
		return []db.Task{}, err
	}
	return tasks, nil
}

func (c *taskController) Create(params *http.Request) error {
	newTask := db.TaskCreateDTO{}
	newTask.Description = params.PostFormValue("description")
	if newTask.Description == "" {
		err := fmt.Sprintf("wrong description: %v", newTask.Description)
		return errors.New(err)
	}
	deadline := params.PostFormValue("date")
	if deadline == "" {
		err := fmt.Sprintf("wrong deadline: %v", deadline)
		return errors.New(err)
	}
	newTask.Deadline = deadline
	err := c.t.Create(newTask)
	if err != nil {
		return err
	}
	return nil
}

func (c *taskController) DeleteTask(taskId string) error {
	return c.t.DeleteTask(taskId)
}

func (c *taskController) IsDone(taskId string) error {
	return c.t.IsDone(taskId)
}
