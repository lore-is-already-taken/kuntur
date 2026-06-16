// Package registro renders the registration page ("/registro") of the Kuntur site.
package registro

import (
	"html/template"
	"net/http"

	"kuntur/app/view"
)

// Get returns an http.HandlerFunc that renders the registration page.
func Get(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view.View(w, r, tmpl, nil)
	}
}
