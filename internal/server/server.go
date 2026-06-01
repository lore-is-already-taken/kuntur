// Package server implements the HTTP adapter for the Kuntur website.
//
// The router is built with the standard library's http.ServeMux, which since
// Go 1.22 supports method+pattern matching (e.g. "GET /{$}"). Static assets
// are served from an embedded filesystem owned by internal/web, so this
// package does not touch the real filesystem at runtime.
package server

import (
	"html/template"
	"log/slog"
	"net/http"

	"kuntur/internal/web"
)

// PageData is the data passed to the index template. Keeping it as an
// explicit struct (instead of a map) catches typos at compile time and makes
// it obvious what the template can render.
type PageData struct {
	Title   string
	Group   string
	Tagline string
}

// NewRouter returns a fully configured http.Handler.
func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// Static files: http.FileServerFS serves from an fs.FS rooted at the
	// directory passed in, so /static/css/style.css maps to static/css/style.css.
	// http.FileServerFS takes an fs.FS directly. http.StripPrefix lets the
	// prefix be removed before the file server looks up the path.
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(web.Static())))

	// "GET /{$}" matches exactly "/" and nothing else — prevents accidental
	// shadowing of the static handler by paths like /foo.
	mux.HandleFunc("GET /{$}", homeHandler(web.Templates()))

	return logMiddleware(mux)
}

func homeHandler(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Title:   "Kuntur",
			Group:   "Kuntur",
			Tagline: "Andean music, alive.",
		}
		// Set Content-Type before writing the body so http.Error below can
		// still emit a proper text/plain error if the template blows up.
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
			slog.Error("template execution failed", "err", err, "path", r.URL.Path)
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
	}
}

// logMiddleware emits one structured log line per request. It is intentionally
// minimal: real production servers should add request IDs, response status,
// duration, and the user-agent. Use the slog handler from stdlib or a wrapper
// from your router of choice.
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info(
			"request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote", r.RemoteAddr,
		)
		next.ServeHTTP(w, r)
	})
}
