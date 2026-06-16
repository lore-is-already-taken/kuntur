// Package contact owns the contact form resource for the Kuntur site,
// including the GET view (rendering the form) and the POST handler
// (submitting the form to the upstream API).
package contact

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"kuntur/app/view"
)

// PageData is the data passed to the contact form template.
type PageData struct {
	Success bool
}

// Get returns an http.HandlerFunc that renders the contact form. The optional
// ?ok=1 query parameter toggles the success message.
func Get(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view.View(w, r, tmpl, PageData{Success: r.URL.Query().Get("ok") == "1"})
	}
}

// New returns an http.Handler that processes a POST submission of the contact
// form. The client and apiURL are injected so tests can swap them and so the
// caller (server.NewRouter) controls the timeout policy.
func New(client *http.Client, apiURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		submit(client, apiURL, w, r)
	})
}

// submit is the private worker for the POST handler. It validates the
// Content-Type, parses the form, marshals the payload, POSTs to apiURL via
// the injected client, and translates upstream errors into HTTP 502.
func submit(client *http.Client, apiURL string, w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/x-www-form-urlencoded") && !strings.HasPrefix(ct, "multipart/form-data") {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	body, err := json.Marshal(Payload{
		Name:    r.FormValue("name"),
		Email:   r.FormValue("email"),
		Message: r.FormValue("message"),
	})
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	resp, err := client.Post(apiURL, "application/json", bytes.NewReader(body))
	if err != nil {
		slog.Error("backend unreachable", "err", err, "path", r.URL.Path)
		http.Error(w, "service unavailable", http.StatusBadGateway)
		return
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode >= 400 {
		slog.Error("backend error", "status", resp.StatusCode, "path", r.URL.Path)
		http.Error(w, "upstream error", http.StatusBadGateway)
		return
	}

	http.Redirect(w, r, "/contacto?ok=1", http.StatusSeeOther)
}
