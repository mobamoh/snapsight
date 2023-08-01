package views

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type Template struct {
	htmlTmpl *template.Template
}

func ParseTemplate(templateName string) (Template, error) {
	tmpl, err := template.ParseFiles(filepath.Join("templates", templateName))
	if err != nil {
		return Template{}, fmt.Errorf("error parsing template: %w", err)
	}
	return Template{htmlTmpl: tmpl}, nil
}
func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := t.htmlTmpl.Execute(w, data)
	if err != nil {
		log.Printf("error executing template: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}
