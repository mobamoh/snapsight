package controllers

import (
	"fmt"
	"github.com/mobamoh/snapsight/models"
	"log"
	"net/http"
)

type Users struct {
	Templates struct {
		New    Template
		SignIn Template
	}
	UserService *models.UserService
}

func (u Users) GetSignUp(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
}

func (u Users) PostSignUp(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	pwd := r.FormValue("password")

	user, err := u.UserService.Create(email, pwd)
	if err != nil {
		log.Printf("signin: %v", err)
		http.Error(w, "something went wrong!", http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "user created: %+v", user)
}

func (u Users) GetSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, data)
}

func (u Users) PostSignIn(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	pwd := r.FormValue("password")

	user, err := u.UserService.Authenticate(email, pwd)
	if err != nil {
		log.Printf("signin: %v", err)
		http.Error(w, "something went wrong!", http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "user authenticated: %+v", user)
}
