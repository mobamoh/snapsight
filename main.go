package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func main() {
	router := chi.NewRouter()
	router.Get("/", homeHandler)
	router.Get("/contact", contactHandler)
	router.Get("/faq", faqHandler)

	fmt.Println("server listening at :1313...")
	if err := http.ListenAndServe(":1313", router); err != nil {
		log.Fatal(err)
	}

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "home.html")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "contact.html")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "faq.html")
}

func renderTemplate(w http.ResponseWriter, templateName string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	files, err := template.ParseFiles(filepath.Join("templates", templateName))
	if err != nil {
		log.Printf("error parsing template: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	err = files.Execute(w, nil)
	if err != nil {
		log.Printf("error executing template: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
