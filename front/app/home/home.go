// Package home renders the index page ("/") of the Kuntur site.
package home

import (
	"html/template"
	"net/http"

	"kuntur/app/view"
)

// PageData is the data passed to the index template.
type PageData struct {
	Title       string
	Group       string
	Description string
}

// Get returns an http.HandlerFunc that renders the index page.
func Get(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Title:       "Kuntur Kantu",
			Group:       "Kuntur Kantu",
			Description: "Sitio oficial de Kuntur Kantu — música andina.",
		}
		view.View(w, r, tmpl, data)
	}
}
