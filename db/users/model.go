package users


type User struct {
	Id int `json:"user_id"`
	Name string `json:"name"`
	Mail string `json:"mail"`
	Password string `json:"pass"`
}