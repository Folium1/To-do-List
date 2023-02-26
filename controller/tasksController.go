package controller

import (
	"errors"

	"todo/db"
	dto "todo/dto"

	"github.com/mitchellh/mapstructure"
)

type taskController struct {
	t db.TaskService // A field of type db.TaskService to interact with the database.
}

// New returns a new instance of TaskController.
// It takes a db.TaskService as an argument and sets it to the 't' field.
func New(db db.TaskService) TaskController {
	return &taskController{t: db}
}

// TaskController is an interface that defines methods to manage tasks.
type TaskController interface {
	Create(newTask dto.TaskCreateDTO) error
	Tasks() ([]dto.TasksDTO, error)
	DeleteTask(taskId string) error
	ChangeData(newTaskData dto.UpdateTaskDTO) error
}

// Adds a new task to the database.
func (c *taskController) Create(newTaskDTO dto.TaskCreateDTO) error {
	var newTask db.Task
	err := mapstructure.Decode(newTaskDTO, &newTask)
	if err != nil {
		return err
	}
	err = c.t.Create(newTask)
	if err != nil {
		return err
	}
	return nil
}

// Retrieves a list of tasks from the database.
func (c *taskController) Tasks() ([]dto.TasksDTO, error) {
	tasks, err := c.t.Tasks()
	if err != nil {
		return nil, err
	}
	tasksDTO := make([]dto.TasksDTO, len(tasks), len(tasks))
	for _, task := range tasks {
		var taskDTO dto.TasksDTO
		err = mapstructure.Decode(task, &taskDTO)
		tasksDTO = append(tasksDTO, taskDTO)
	}
	return tasksDTO, nil
}

// Deletes a task with the given Id from the database.
func (c *taskController) DeleteTask(taskId string) error {
	return c.t.DeleteTask(taskId)
}

// Updates a task with new data in the database.
func (c *taskController) ChangeData(newTaskData dto.UpdateTaskDTO) error {
	task, err := c.t.GetTask(newTaskData.Id) // Retrieve the task from the database using the Id.
	if err != nil {
		return err
	}

	// Update the task fields if they are different from the new data.
	if task.Deadline != newTaskData.Deadline {
		task.Deadline = newTaskData.Deadline
	}

	if task.Description != newTaskData.Description {
		task.Description = newTaskData.Description
	}

	if len(task.Description) < 3 {
		return errors.New("description is too short")
	}

	// Save the updated task in the database.
	err = c.t.SaveTask(task)
	if err != nil {
		return err
	}
	return nil
}
