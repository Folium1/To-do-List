package usercontroller

import (
	"log"
	parser "todo/controllers"
	db "todo/db/users"
	"todo/dto"
)

type UserController struct {
	db db.Service
}

type ControllerService interface {
	Create(newUser dto.UserDTO) (int, error)
	GetUser(user dto.LoginUserDTO) (dto.LoginUserDTO, error)
}

func New(userService db.Service) ControllerService {
	return &UserController{userService}
}

func (u *UserController) Create(newUser dto.UserDTO) (int, error) {
	var dbUser db.User
	err := parser.UsersDTOtoDB(newUser, &dbUser)
	if err != nil {
		log.Printf("Coudn't parse from dto to db data,err: %v", err)
		return 0, err
	}
	userId, err := u.db.CreateUser(dbUser)
	if err != nil {
		log.Printf("Couldn't create user, err: %v", err)
		return 0, err
	}
	return userId, nil
}

// GetUser returns user data
func (u *UserController) GetUser(user dto.LoginUserDTO) (dto.LoginUserDTO, error) {
	var dbUser db.User
	err := parser.UsersDTOtoDB(user, &dbUser)
	if err != nil {
		log.Printf("Couldn't parse data, err = %v", err)
		return dto.LoginUserDTO{}, err
	}

	gotUser, err := u.db.GetUser(dbUser)
	if err != nil {
		log.Printf("Couldn't get user from db,err: %v", err)
		return dto.LoginUserDTO{}, err
	}

	var gotUserDTO dto.LoginUserDTO
	err = parser.UsersDBtoDTO(gotUser, &gotUserDTO)
	if err != nil {
		log.Printf("Couldn't parse data from db to dto, err: %v", err)
		return dto.LoginUserDTO{}, err
	}
	return gotUserDTO, nil
}
