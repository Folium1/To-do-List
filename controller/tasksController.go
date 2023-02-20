package controller

import (
	"todo/db"
)

type taskController struct {
	t db.TaskService
}

func New(db db.TaskService) TaskController {
	return &taskController{t: db}
}

type TaskController interface {
	Create(newTask db.TaskCreateDTO) error
	Tasks() ([]db.Task, error)
	DeleteTask(taskId string) error
	ChangeData(newTaskData db.Task) error
}

func (c *taskController) Tasks() ([]db.Task, error) {
	tasks, err := c.t.Tasks()
	if err != nil {
		return []db.Task{}, err
	}
	return tasks, nil
}

func (c *taskController) Create(newTask db.TaskCreateDTO) error {
	err := c.t.Create(newTask)
	if err != nil {
		return err
	}
	return nil
}

func (c *taskController) DeleteTask(taskId string) error {
	return c.t.DeleteTask(taskId)
}

func (c *taskController) ChangeData(newTaskData db.Task) error {
	task, err := c.t.GetTask(newTaskData.Id)
	if err != nil {
		return err
	}
	if task.Deadline != newTaskData.Deadline {
		task.Deadline = newTaskData.Deadline
	}
	if task.Description != newTaskData.Description {
		task.Description = newTaskData.Description
	}
	if task.Done != newTaskData.Done {
		task.Done = newTaskData.Done
	}
	if task.Done {
		err := c.t.DeleteTask(newTaskData.Id)
		if err != nil {
			return err
		}
		return nil
	}
	err = c.t.SaveTask(task)
	if err != nil {
		return err
	}
	return nil
}
