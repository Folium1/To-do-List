package controller

import (
	"errors"
	"log"

	parser "todo/controllers"
	db "todo/db/tasks"
	dto "todo/dto"
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
	Tasks(userId string) ([]dto.TasksDTO, error)
	DeleteTask(taskId string) error
	ChangeData(newTaskData dto.UpdateTaskDTO) error
}

// Adds a new task to the database.
func (c *taskController) Create(newTaskDTO dto.TaskCreateDTO) error {
	var newTask db.Task
	err := parser.TaskDTOtoDB(newTaskDTO, &newTask)
	if err != nil {
		log.Printf("Coudn't parse from dto to db data,err: %v", err)
		return err
	}
	err = c.t.Create(newTask)
	if err != nil {
		return err
	}
	return nil
}

// Retrieves a list of tasks from the database.
func (c *taskController) Tasks(userId string) ([]dto.TasksDTO, error) {
	tasks, err := c.t.Tasks(userId)
	if err != nil {
		return nil, err
	}
	tasksDTO := make([]dto.TasksDTO, 0, len(tasks))
	for _, task := range tasks {
		// Parse data from db.Task struct into dto.TasksDTO struct
		var taskDTO dto.TasksDTO
		err = parser.TaskDBtoDTO(task, &taskDTO)
		tasksDTO = append(tasksDTO, taskDTO)
	}
	return tasksDTO, nil
}

// DeleteTask delets a task with the given Id from the database.
func (c *taskController) DeleteTask(taskId string) error {
	return c.t.DeleteTask(taskId)
}

// Updates a task with new data in the database.
func (c *taskController) ChangeData(newTaskData dto.UpdateTaskDTO) error {
	task, err := c.t.GetTask(newTaskData.Id) // Retrieves the task from the database using the Id.
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
