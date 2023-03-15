package handler

import (
	"fmt"
	"log"
	"net/http"

	dto "todo/dto"
	auth "todo/handler/middleware"

	"golang.org/x/crypto/bcrypt"
)

// signing up new user
func signUpHandler(w http.ResponseWriter, r *http.Request) {
	if auth.IsAuthenticated(w, r) {
		http.Redirect(w, r, "/todo", 302)
	}
	if r.Method == "GET" {
		err := templ.ExecuteTemplate(w, "register.html", nil)
		if err != nil {
			return
		}
	}
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			return
		}
		newUser := dto.UserDTO{
			Name:     r.FormValue("name"),
			Mail:     r.FormValue("mail"),
			Password: r.FormValue("pass"),
		}

		// validating data
		if newUser.Name == "" {
			err = fmt.Errorf("Name is missing")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if newUser.Mail == "" {
			err = fmt.Errorf("Mail is missing")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if len(newUser.Password) < 8 {
			err = fmt.Errorf("Password is too short or missing")
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Hashing user's password
		hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10)
		if err != nil {
			log.Printf("Couldn't hash user's password")
		}
		newUser.Password = string(hash)
		userId, err := userController.Create(newUser)
		if err != nil {
			http.Redirect(w, r, "/sign-up", 302)
		}
		// authorizing user
		err = auth.AuthUser(w, r, userId)
		if err != nil {
			log.Printf("Couldn't authorize user(%v), err: %v", userId, err)
			http.Redirect(w, r, "/login", 302)
		}
		http.Redirect(w, r, "/todo", http.StatusFound)
	}
}


// loginHandler login's user using mail and password from input
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if auth.IsAuthenticated(w, r) {
		http.Redirect(w, r, "/todo", 302)
	}
	if r.Method == "GET" {
		err := templ.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			log.Println(err)
		}
	}
	if r.Method == "POST" {
		newUser := dto.LoginUserDTO{
			Mail:     r.FormValue("mail"),
			Password: r.FormValue("pass"),
		}
		// Validating user's data
		if newUser.Mail == "" {
			err := fmt.Errorf("missing mail")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if newUser.Password == "" {
			err := fmt.Errorf("password is too short, at least 8 elements required")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		dbUser, err := userController.GetUser(newUser)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(newUser.Password))
		if err != nil {
			log.Printf("Input password and db password not equal,err: %v", err)
			http.Error(w, "Password is not valid", http.StatusBadRequest)
		}
		// authorize user
		err = auth.AuthUser(w, r, dbUser.Id)
		if err != nil {
			log.Printf("Couldn't authorize user(%v), err: %v", dbUser.Id, err)

		}
		http.Redirect(w, r, "/todo", 302)
	}
}
