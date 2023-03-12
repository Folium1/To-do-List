package db

type Task struct {
	Id          string `json:"id"`
	UserId 		string `json:"user_id"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
	Done        bool   `json:"done"`
}
