package handler

import (
	"errors"
	"fmt"
	"net/http"
	"text/template"

	"todo/controller"
	"todo/db"

	"github.com/gorilla/mux"
)

var (
	task           = db.NewService()
	taskController = controller.New(task)
	router         = mux.NewRouter()
	templ          *template.Template
)

func main(w http.ResponseWriter, r *http.Request) {
	tasks, err := taskController.Tasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	err = templ.ExecuteTemplate(w, "index.html", tasks)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask db.TaskCreateDTO
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
	err := taskController.Create(newTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func updateData(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		method := r.FormValue("_method")
		if method == "PATCH" {

		}
	}
	id := r.PostFormValue("id")
	if id == "" {
		err := errors.New("wrond task id handler")
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	var newTaskData db.UpdateTaskDTO
	newTaskData.Id = id
	newTaskData.Deadline = r.FormValue("deadline")
	newTaskData.Description = r.FormValue("description")

	taskController.ChangeData(newTaskData)
	fmt.Println(newTaskData)
	err := taskController.ChangeData(newTaskData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	err := taskController.DeleteTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func InitServer() {
	var err error
	templ, err = template.ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}

	router.HandleFunc("/", main).Methods("GET")
	router.HandleFunc("/create/", createTask).Methods("POST")
	router.HandleFunc("/update/{taskId}", updateData).Methods("POST")
	router.HandleFunc("/delete/{id}/", deleteTask).Methods("GET")
	http.ListenAndServe(":8080", router)
}
