package views

import (
	"bytes"
	"fmt"
	"github.com/gorilla/csrf"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
)

type Template struct {
	htmlTmpl *template.Template
}

func ParseFS(fs fs.FS, pattern ...string) (Template, error) {
	tmpl := template.New(pattern[0])
	tmpl.Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", fmt.Errorf("implement csrfField method")
		},
	})

	tmpl, err := tmpl.ParseFS(fs, pattern...)
	if err != nil {
		return Template{}, fmt.Errorf("error parsing template: %w", err)
	}
	return Template{htmlTmpl: tmpl}, nil
}
func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {
	tmpl, err := t.htmlTmpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
		return
	}
	tmpl = tmpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
	})
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		log.Printf("error executing template: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}
