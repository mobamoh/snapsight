package controllers

import (
	"net/http"
)

func StaticHandler(template Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		template.Execute(w, r, nil)
	}
}

func FaqHandler(template Template) http.HandlerFunc {
	faq := []struct {
		Question string
		Answer   string
	}{
		{
			Question: "What payment methods do you accept?",
			Answer:   "We accept credit cards, PayPal, and bank transfers.",
		},
		{
			Question: "How do I reset my password?",
			Answer:   "Click \"Forgot Password\" on the login page and follow the instructions.",
		},
		{
			Question: "How long does shipping take?",
			Answer:   "Shipping times vary, 3-7 business days domestically, 7-14 internationally.",
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		template.Execute(w, r, faq)
	}
}
