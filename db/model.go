package db



type task struct {
	Id string `json:"id"`
	Description string `json:"description"`
	Deadline string `json:"deadline"`
	Done  bool `json:"done"`
}