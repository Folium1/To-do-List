package db


type TaskCreateDTO struct {
	Description string `json:"description"`
	Deadline string `json:"deadline"`
}