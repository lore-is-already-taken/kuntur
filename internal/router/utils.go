package router

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"
)

func render(w http.ResponseWriter, r *http.Request, tmpl *template.Template, data any) {
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "layout", data); err != nil {
		slog.Error("template execution failed", "err", err, "path", r.URL.Path)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK) // Asegura el estatus 200 explícito antes de escribir el cuerpo

	if _, err := buf.WriteTo(w); err != nil {
		slog.Error("response write failed", "err", err, "path", r.URL.Path)
	}
}
