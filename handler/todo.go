package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	dto "todo/dto"
	auth "todo/handler/middleware"
)

// The main function is responsible for rendering the main page.
func main(w http.ResponseWriter, r *http.Request) {
	// checking for jwt token
	if !auth.IsAuthenticated(w, r) {
		http.Redirect(w, r, "/login", 302)
	}
	userId, err := auth.GetUserIdFromCookies(r)
	if err != nil {
		log.Printf("Coudn't get userId from cookies,err: %v", err)
		http.Redirect(w, r, "/logn", 302)
	}
	// Retrieve all tasks from the task controller
	tasks, err := taskController.Tasks(userId)
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
func RedirectToMain(w http.ResponseWriter, r *http.Request) {
	r.Method = "GET"
	http.Redirect(w, r, "/todo", http.StatusMovedPermanently)
}

// createTask function is responsible for creating a new task.
func createTask(w http.ResponseWriter, r *http.Request) {
	// checking for jwt token
	if !auth.IsAuthenticated(w, r) {
		http.Redirect(w, r, "/", 302)
	}
	// Retrieve the task data from the HTTP request.
	if r.Method == "POST" {

		var newTask dto.TaskCreateDTO
		var err error
		newTask.User_id, err = auth.GetUserIdFromCookies(r)
		if err != nil {
			log.Printf("Coudn't get userId from cookies,err: %v", err)
		}
		if err != nil {
			log.Printf("Coudn't get user's id from token,err:%v", err)
			r.Method = "GET"
			http.Redirect(w, r, "/", 302)
		}
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
		err = taskController.Create(newTask)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		// Redirects the client to the main page.
		RedirectToMain(w, r)
	}

}

// updateData is responsible for updating an existing task.
func updateData(w http.ResponseWriter, r *http.Request) {
	// checking for jwt token
	if !auth.IsAuthenticated(w, r) {
		http.Redirect(w, r, "/", 302)
	}
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
			RedirectToMain(w, r)
		}
	}
}

// deleteTask function is responsible for the deletion of the task.
func deleteTask(w http.ResponseWriter, r *http.Request) {
	// checking for jwt token
	if !auth.IsAuthenticated(w, r) {
		http.Redirect(w, r, "/", 302)
	}
	path := r.URL.Path
	parts := strings.Split(path, "/")
	id := parts[2]
	err := taskController.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	// Redirects the client to the main page.
	RedirectToMain(w, r)
}
