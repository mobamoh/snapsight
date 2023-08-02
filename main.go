package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mobamoh/snapsight/controllers"
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

	fmt.Println("server listening at :1313...")
	if err := http.ListenAndServe(":1313", router); err != nil {
		log.Fatal(err)
	}

}
