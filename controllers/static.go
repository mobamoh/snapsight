package controllers

import (
	"github.com/mobamoh/snapsight/views"
	"net/http"
)

func StaticHandler(template views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		template.Execute(w, nil)
	}
}
