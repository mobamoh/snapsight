package views

import (
	"bytes"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/mobamoh/snapsight/context"
	"github.com/mobamoh/snapsight/errors"
	"github.com/mobamoh/snapsight/models"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"
)

type Template struct {
	htmlTmpl *template.Template
}

// This is used to determine if an error provides the Public method.
type public interface {
	Public() string
}

func ParseFS(fs fs.FS, pattern ...string) (Template, error) {
	tmpl := template.New(path.Base(pattern[0]))
	tmpl.Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", fmt.Errorf("implement csrfField method")
		},
		"currentUser": func() (template.HTML, error) {
			return "", fmt.Errorf("implement currentUser method")
		},
		"errors": func() []string {
			return nil
		},
	})

	tmpl, err := tmpl.ParseFS(fs, pattern...)
	if err != nil {
		return Template{}, fmt.Errorf("error parsing template: %w", err)
	}
	return Template{htmlTmpl: tmpl}, nil
}
func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {
	tmpl, err := t.htmlTmpl.Clone()
	if err != nil {
		log.Printf("cloning template: %v", err)
		http.Error(w, "There was an error rendering the page.", http.StatusInternalServerError)
		return
	}

	errMsgs := errMessages(errs...)

	tmpl = tmpl.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrf.TemplateField(r)
		},
		"currentUser": func() *models.User {
			return context.User(r.Context())
		},
		"errors": func() []string {
			// return the pre-processed err messages inside the closure.
			return errMsgs
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

func errMessages(errs ...error) []string {
	var msgs []string
	for _, err := range errs {
		var pubErr public
		if errors.As(err, &pubErr) {
			msgs = append(msgs, pubErr.Public())
		} else {
			fmt.Println(err)
			msgs = append(msgs, "Something went wrong.")
		}
	}
	return msgs
}
