package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	router := chi.NewRouter()
	router.Get("/", homeHandler)
	router.Get("/contact", contactHandler)
	router.Get("/faq", faqHandler)
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "This page is maybe disappeared in a black hole!! 🕳️")
	})

	router.Get("/param/{id}", func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "id")
		fmt.Fprintf(w, "Hello %s", param)
	})

	fmt.Println("server listening at :1313...")
	if err := http.ListenAndServe(":1313", router); err != nil {
		log.Fatal(err)
	}

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Hello Server</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, "<h1>Contact info</h1><p>contact me at mo@mail.com</p>")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>FAQ</h1><p1>Some frequently asked questions and their answers</p1>")
}
