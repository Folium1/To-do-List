package handler

import (
	"net/http"
	"text/template"

	controller "todo/controllers/taskController"
	usercontroller "todo/controllers/userController"
	db "todo/db/tasks"
	"todo/db/users"
)

var (
	task           = db.NewService()
	taskController = controller.New(task)
	user           = users.New()
	userController = usercontroller.New(user)
	templ          *template.Template
)

func StartServer() {
	var err error
	templ, err = template.ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", signUpHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/todo", main)

	http.HandleFunc("/create/", createTask)
	http.HandleFunc("/update/", updateData)
	http.HandleFunc("/delete/", deleteTask)

	err = http.ListenAndServe(":9090", nil)
	if err != nil {
		panic(err)
	}
}
