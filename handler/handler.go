package handler

import (
	"errors"
	"fmt"
	"strings"

	"net/http"
	"text/template"

	"todo/controller"
	"todo/db"
	dto "todo/DTO"
)

var (
	task           = db.NewService()
	taskController = controller.New(task)
	templ          *template.Template
)

// The main function is responsible for rendering the main page.
func main(w http.ResponseWriter, r *http.Request) {
	// Retrieve all tasks from the task controller
	tasks, err := taskController.Tasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	// Pass the tasks to the template for rendering
	err = templ.ExecuteTemplate(w, "index.html", tasks)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// Redirects to main page.
func redirectToMain(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo", http.StatusMovedPermanently)
}

// The createTask function is responsible for creating a new task.
func createTask(w http.ResponseWriter, r *http.Request) {
	// Retrieve the task data from the HTTP request.
	if r.Method == "POST" {
		var newTask dto.TaskCreateDTO
		newTask.Description = r.PostFormValue("description")
		if newTask.Description == "" {
			err := fmt.Sprintf("wrong description: %v", newTask.Description)
			http.Error(w, err, http.StatusBadRequest)
		}
		deadline := r.PostFormValue("date")
		if deadline == "" {
			err := fmt.Sprintf("wrong deadline: %v", deadline)
			http.Error(w, err, http.StatusBadRequest)
		}
		newTask.Deadline = deadline

		// Pass the new task data to the task controller for creation.
		err := taskController.Create(newTask)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		// Redirects the client to the main page.
		redirectToMain(w, r)
	}

}

// The updateData function is responsible for updating an existing task.
func updateData(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		method := r.FormValue("_method")
		if method == "PATCH" {
			// Retrieve the task ID and updated data from the HTTP request.
			id := r.PostFormValue("id")
			if id == "" {
				err := errors.New("wrond task id")
				http.Error(w, err.Error(), http.StatusNotFound)
			}
			var newTaskData dto.UpdateTaskDTO
			newTaskData.Id = id
			newTaskData.Deadline = r.FormValue("deadline")
			newTaskData.Description = r.FormValue("description")
			err := taskController.ChangeData(newTaskData)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			// Redirects the client to the main page.
			redirectToMain(w, r)
		}
	}
}

// The deleteTask function is responsible for the deletion of the task.
func deleteTask(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	id := parts[2]
	err := taskController.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	// Redirects the client to the main page.
	redirectToMain(w, r)
}

// Starts server.
func InitServer() {
	var err error
	templ, err = template.ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/todo", main)
	http.HandleFunc("/create/", createTask)
	http.HandleFunc("/update/", updateData)
	http.HandleFunc("/delete/", deleteTask)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
