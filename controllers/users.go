package controllers

import (
	"fmt"
	"github.com/mobamoh/snapsight/context"
	"github.com/mobamoh/snapsight/models"
	"log"
	"net/http"
)

type Users struct {
	Templates struct {
		SignUp Template
		SignIn Template
	}
	UserService    *models.UserService
	SessionService *models.SessionService
}

func (u Users) GetSignUp(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignUp.Execute(w, r, data)
}

func (u Users) PostSignUp(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	pwd := r.FormValue("password")

	user, err := u.UserService.Create(email, pwd)
	if err != nil {
		log.Printf("signin: %v", err)
		http.Error(w, "something went wrong!", http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		// TODO: Long term, we should show a warning about not being able to sign the user in.
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) GetSignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}
	data.Email = r.FormValue("email")
	u.Templates.SignIn.Execute(w, r, data)
}

func (u Users) PostSignIn(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	pwd := r.FormValue("password")

	user, err := u.UserService.Authenticate(email, pwd)
	if err != nil {
		log.Printf("signin: %v", err)
		http.Error(w, "something went wrong!", http.StatusInternalServerError)
		return
	}

	session, err := u.SessionService.Create(user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusInternalServerError)
		return
	}
	setCookie(w, CookieSession, session.Token)
	http.Redirect(w, r, "/users/me", http.StatusFound)
}

func (u Users) PostSignOut(w http.ResponseWriter, r *http.Request) {

	cookie, err := readCookie(r, CookieSession)
	if err != nil {
		http.Redirect(w, r, "/signin", http.StatusFound)
		return
	}
	if err = u.SessionService.Delete(cookie); err != nil {
		log.Printf("%s", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	deleteCookie(w, CookieSession)
	http.Redirect(w, r, "/signin", http.StatusFound)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	fmt.Fprintf(w, "Current user: %s\n", user.Email)
}

type UserMiddleware struct {
	SessionService *models.SessionService
}

func (umw UserMiddleware) SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := readCookie(r, CookieSession)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		user, err := umw.SessionService.User(token)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (umw UserMiddleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
