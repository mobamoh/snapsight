package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mobamoh/snapsight/controllers"
	"github.com/mobamoh/snapsight/models"
	"github.com/mobamoh/snapsight/templates"
	"github.com/mobamoh/snapsight/views"
	"log"
	"net/http"
)

func main() {
	router := chi.NewRouter()

	router.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "layout.gohtml", "home.gohtml"))))

	router.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "layout.gohtml", "contact.gohtml"))))

	router.Get("/faq", controllers.FaqHandler(views.Must(views.ParseFS(templates.FS, "layout.gohtml", "faq.gohtml"))))

	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected!")

	userService := models.UserService{
		DB: db,
	}

	userCtrl := controllers.Users{
		UserService: &userService,
	}
	userCtrl.Templates.New = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "signup.gohtml"))
	userCtrl.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "signin.gohtml"))

	router.Get("/signup", userCtrl.GetSignUp)
	router.Post("/signup", userCtrl.PostSignUp)

	router.Get("/signin", userCtrl.GetSignIn)
	router.Post("/signin", userCtrl.PostSignIn)

	fmt.Println("server listening at :1313...")
	if err := http.ListenAndServe(":1313", router); err != nil {
		log.Fatal(err)
	}

}
