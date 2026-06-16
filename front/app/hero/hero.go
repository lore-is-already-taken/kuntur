// Package hero renders the presentations page ("/presentaciones") of the Kuntur site.
package hero

import (
	"html/template"
	"net/http"

	"kuntur/app/view"
)

// Get returns an http.HandlerFunc that renders the presentations page.
func Get(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view.View(w, r, tmpl, nil)
	}
}
