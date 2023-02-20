package db

type TaskCreateDTO struct {
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
}

type TasksDTO struct {
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
	Done        bool   `json:"done"`
}
