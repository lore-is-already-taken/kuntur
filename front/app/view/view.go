// Package view provides the shared template rendering helper used by every
// feature package. The intent is to keep a single place that knows how to
// turn a parsed *template.Template plus a data value into an HTTP response.
package view

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"
)

// View executes tmpl against data and writes the result to w with a 200 OK
// status and a text/html content type. Template execution errors and response
// write errors are logged via slog; the response status is set to 500 on
// template errors so the client never sees a half-written body.
func View(w http.ResponseWriter, r *http.Request, tmpl *template.Template, data any) {
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "layout", data); err != nil {
		slog.Error("template execution failed", "err", err, "path", r.URL.Path)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if _, err := buf.WriteTo(w); err != nil {
		slog.Error("response write failed", "err", err, "path", r.URL.Path)
	}
}
