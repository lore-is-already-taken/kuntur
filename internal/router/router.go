// Package router this is used to redirect all views
package router

import (
	"html/template"
	"net/http"
)

// PageData es ideal para el Home. Si otras páginas usan datos distintos,
// puedes crear structs específicos como HeroData, AboutData, etc.
type PageData struct {
	Title   string
	Group   string
	Tagline string
}

func HomeHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Title:   "Kuntur",
			Group:   "Kuntur",
			Tagline: "Aguanten las cabras, somos poderosas",
		}
		render(w, r, tmpl, "index.html", data)
	}
}

func HeroHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, r, tmpl, "hero.html", nil)
	}
}

