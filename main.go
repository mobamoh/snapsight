package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/mobamoh/snapsight/controllers"
	"github.com/mobamoh/snapsight/migrations"
	"github.com/mobamoh/snapsight/models"
	"github.com/mobamoh/snapsight/templates"
	"github.com/mobamoh/snapsight/views"
	"log"
	"net/http"
)

func main() {
	// Setup DB
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	//defer func(db *sql.DB) {
	//	err := db.Close()
	//	if err != nil {
	//		panic(err)
	//	}
	//}(db)
	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("DB connected...")
	if err = models.MigrateFS(db, migrations.FS, "."); err != nil {
		panic(err)
	}

	// Setup App Services
	userService := models.UserService{
		DB: db,
	}
	sessionService := models.SessionService{
		DB: db,
	}
	userCtrl := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}

	// Setup Middleware
	csrfKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	csrfMw := csrf.Protect([]byte(csrfKey), csrf.Secure(false)) // TODO: change for prod
	userMw := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	// Setup Router
	router := chi.NewRouter()
	router.Use(csrfMw)
	router.Use(userMw.SetUser)

	router.Get("/", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "layout.gohtml", "home.gohtml"))))
	router.Get("/contact", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "layout.gohtml", "contact.gohtml"))))
	router.Get("/faq", controllers.FaqHandler(views.Must(views.ParseFS(templates.FS, "layout.gohtml", "faq.gohtml"))))

	userCtrl.Templates.SignUp = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "signup.gohtml"))
	userCtrl.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "signin.gohtml"))
	router.Get("/signup", userCtrl.GetSignUp)
	router.Post("/signup", userCtrl.PostSignUp)
	router.Get("/signin", userCtrl.GetSignIn)
	router.Post("/signin", userCtrl.PostSignIn)
	router.Post("/signout", userCtrl.PostSignOut)

	router.Route("/users/me", func(r chi.Router) {
		r.Use(userMw.RequireUser)
		r.Get("/", userCtrl.CurrentUser)
	})

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Oops!! Page not found", http.StatusNotFound)
	})

	// Starting Server
	fmt.Println("server listening at :1313...")
	if err := http.ListenAndServe(":1313", csrfMw(router)); err != nil {
		log.Fatal(err)
	}

}
