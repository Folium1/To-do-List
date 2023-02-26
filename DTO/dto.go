package dto


type TaskCreateDTO struct {
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
}

type TasksDTO struct {
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
	Done        bool   `json:"done"`
}

type UpdateTaskDTO struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
}
