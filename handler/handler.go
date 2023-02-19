package handler

import (
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
		http.Error(w, err.Error(), 404)
	}

	err = templ.ExecuteTemplate(w, "index.html", tasks)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {
	err := taskController.Create(r)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func updateData(w http.ResponseWriter, r *http.Request) {
	err := taskController.ChangeData(r)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func InitServer() {
	var err error
	templ, err = template.ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}

	router.HandleFunc("/", main).Methods("GET")
	router.HandleFunc("/create/", createTask).Methods("POST")
	router.HandleFunc("/update/{taskId}", updateData).Methods("PATCH")

	http.ListenAndServe(":8080", router)
}
