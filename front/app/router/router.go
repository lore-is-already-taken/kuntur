// Package router this is used to redirect all views
package router

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
)

api_url := "http://localhost:8000/contacto"

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
		render(w, r, tmpl, data)
	}
}

func HeroHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, r, tmpl, nil)
	}
}

func RegistroHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, r, tmpl, nil)
	}
}

func BioHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, r, tmpl, nil)
	}
}

type ContactData struct {
	Success bool
}

func ContactHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render(w, r, tmpl, ContactData{Success: r.URL.Query().Get("ok") == "1"})
	}
}

func SaveContacto() http.HandlerFunc {
	type payload struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Message string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "invalid form", http.StatusBadRequest)
			return
		}

		body, err := json.Marshal(payload{
			Name:    r.FormValue("name"),
			Email:   r.FormValue("email"),
			Message: r.FormValue("message"),
		})
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		resp, err := http.Post(api_url, "application/json", bytes.NewReader(body))
		if err != nil {
			slog.Error("backend unreachable", "err", err)
			http.Error(w, "service unavailable", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			slog.Error("backend error", "status", resp.StatusCode)
			http.Error(w, "upstream error", http.StatusBadGateway)
			return
		}

		http.Redirect(w, r, "/contacto?ok=1", http.StatusSeeOther)
	}
}
