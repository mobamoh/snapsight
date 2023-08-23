package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
	"github.com/mobamoh/snapsight/controllers"
	"github.com/mobamoh/snapsight/migrations"
	"github.com/mobamoh/snapsight/models"
	"github.com/mobamoh/snapsight/templates"
	"github.com/mobamoh/snapsight/views"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {

	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	// Setup DB
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()
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
	pwResetService := models.PasswordResetService{
		DB: db,
	}
	emailService := models.NewEmailService(cfg.SMTP)

	userCtrl := controllers.Users{
		UserService:          &userService,
		SessionService:       &sessionService,
		PasswordResetService: &pwResetService,
		EmailService:         emailService,
	}

	// Setup Middleware
	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
	)

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
	userCtrl.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "forgot-pw.gohtml"))
	userCtrl.Templates.CheckYourEmail = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "check-your-email.gohtml"))
	userCtrl.Templates.ResetPassword = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "reset-pw.gohtml"))

	router.Get("/signup", userCtrl.GetSignUp)
	router.Post("/signup", userCtrl.PostSignUp)

	router.Get("/signin", userCtrl.GetSignIn)
	router.Post("/signin", userCtrl.PostSignIn)

	router.Post("/signout", userCtrl.PostSignOut)

	router.Get("/forgot-pw", userCtrl.GetForgotPassword)
	router.Post("/forgot-pw", userCtrl.PostForgotPassword)

	router.Get("/reset-pw", userCtrl.GetResetPassword)
	router.Post("/reset-pw", userCtrl.PostResetPassword)

	router.Route("/users/me", func(r chi.Router) {
		r.Use(userMw.RequireUser)
		r.Get("/", userCtrl.CurrentUser)
	})

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Oops!! Page not found", http.StatusNotFound)
	})

	// Starting Server
	fmt.Printf("Server listening on %s...\n", cfg.Server.Address)
	if err := http.ListenAndServe(cfg.Server.Address, csrfMw(router)); err != nil {
		log.Fatal(err)
	}

}

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	// TODO: Read the PSQL values from an ENV variable
	cfg.PSQL = models.DefaultPostgresConfig()

	// TODO: SMTP
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	// TODO: Read the CSRF values from an ENV variable
	cfg.CSRF.Key = "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"
	cfg.CSRF.Secure = false

	// TODO: Read the server values from an ENV variable
	cfg.Server.Address = ":1313"

	return cfg, nil
}
