package controller

import (
	// "encoding/json"
	"encoding/json"
	tasks "todo/db/tasks"
	users "todo/db/users"
	"todo/dto"
)

// Parses dto data to db.User struct
func UsersDTOtoDB[userDataDTO dto.LoginUserDTO | dto.UserDTO](dtoData userDataDTO, user *users.User) error {
	data, err := json.Marshal(dtoData)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &user)
	if err != nil {
		return err
	}
	return nil
}

// Parses db.User's data to dto structs
func UsersDBtoDTO[userDataDTO dto.LoginUserDTO | dto.UserDTO](user users.User, userDto *userDataDTO) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &userDto)
	if err != nil {
		return err
	}
	return nil
}

func TaskDTOtoDB[T dto.TaskCreateDTO | dto.TasksDTO | dto.UpdateTaskDTO](dtoData T, task *tasks.Task) error {
	data, err := json.Marshal(dtoData)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &task)
	if err != nil {
		return err
	}
	return nil
}

func TaskDBtoDTO[T dto.TaskCreateDTO | dto.TasksDTO | dto.UpdateTaskDTO](task tasks.Task, taskDTO *T) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &taskDTO)
	if err != nil {
		return err
	}
	return nil
}
