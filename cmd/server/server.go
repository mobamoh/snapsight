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
	"net/http"
	"os"
	"strconv"
)

func main() {

	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	err = run(cfg)
	if err != nil {
		panic(err)
	}
}
func run(cfg config) error {
	// Setup DB
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		return err
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		return err
	}
	fmt.Println("DB connected...")
	if err = models.MigrateFS(db, migrations.FS, "."); err != nil {
		return err
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
	galleryService := models.GalleryService{
		DB: db,
	}

	userCtrl := controllers.Users{
		UserService:          &userService,
		SessionService:       &sessionService,
		PasswordResetService: &pwResetService,
		EmailService:         emailService,
	}
	galleryCtrl := controllers.Galleries{
		GalleryService: &galleryService,
	}

	// Setup Middleware
	csrfMw := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
		csrf.Path("/"),
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

	galleryCtrl.Templates.New = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "galleries/new.gohtml"))
	galleryCtrl.Templates.Edit = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "galleries/edit.gohtml"))
	galleryCtrl.Templates.Index = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "galleries/index.gohtml"))
	galleryCtrl.Templates.Show = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "galleries/show.gohtml"))

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

	router.Route("/galleries", func(subR chi.Router) {
		subR.Get("/{id}", galleryCtrl.Show)
		subR.Get("/{id}/images/{filename}", galleryCtrl.Image)
		subR.Group(func(r chi.Router) {
			r.Use(userMw.RequireUser)
			r.Get("/", galleryCtrl.Index)
			r.Get("/new", galleryCtrl.New)
			r.Post("/", galleryCtrl.Create)
			r.Get("/{id}/edit", galleryCtrl.Edit)
			r.Post("/{id}", galleryCtrl.Update)
			r.Post("/{id}/delete", galleryCtrl.Delete)
			r.Post("/{id}/images/{filename}/delete", galleryCtrl.DeleteImage)
			r.Post("/{id}/images", galleryCtrl.UploadImage)
		})
	})

	assetsHandler := http.FileServer(http.Dir("assets"))
	router.Get("/assets/*", http.StripPrefix("/assets", assetsHandler).ServeHTTP)

	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Oops!! Page not found", http.StatusNotFound)
	})

	// Starting Server
	fmt.Printf("Server listening on %s...\n", cfg.Server.Address)
	return http.ListenAndServe(cfg.Server.Address, csrfMw(router))
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
	cfg.PSQL = models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Database: os.Getenv("PSQL_DATABASE"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	}
	if cfg.PSQL.Host == "" && cfg.PSQL.Port == "" {
		return cfg, fmt.Errorf("no plsql config provided")
	}

	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	cfg.CSRF.Secure = os.Getenv("CSRF_SECURE") == "true"

	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")

	return cfg, nil
}
