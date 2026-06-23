// Package server implements the HTTP adapter for the Kuntur website.
//
// The router is built with the standard library's http.ServeMux, which since
// Go 1.22 supports method+pattern matching (e.g. "GET /{$}"). Static assets
// are served from an embedded filesystem owned by app/web, so this package
// does not touch the real filesystem at runtime.
package server

import (
	"log/slog"
	"net/http"
	"time"

	"kuntur/app/bio"
	"kuntur/app/config"
	"kuntur/app/contact"
	"kuntur/app/hero"
	"kuntur/app/home"
	"kuntur/app/registro"
	"kuntur/app/web"
)

const (
	contactHTTPTimeout = 10 * time.Second
)

// NewRouter returns a fully configured http.Handler.
func NewRouter(cfg config.Config) http.Handler {
	mux := http.NewServeMux()

	// Static files: http.FileServerFS serves from an fs.FS rooted at the
	// directory passed in, so /static/css/style.css maps to static/css/style.css.
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServerFS(web.Static())))

	views := web.Views()
	client := &http.Client{Timeout: contactHTTPTimeout}
	apiURL := cfg.APIBaseURL + "/contacto"

	mux.HandleFunc("GET /{$}", home.Get(views["index.html"]))
	mux.HandleFunc("GET /presentaciones", hero.Get(views["presentaciones.html"], cfg.APIBaseURL))
	mux.HandleFunc("GET /biografia", bio.Get(views["biografia.html"], cfg.APIBaseURL))
	mux.HandleFunc("GET /registro", registro.Get(views["registro.html"]))
	mux.HandleFunc("GET /contacto", contact.Get(views["contacto.html"]))
	mux.Handle("POST /contacto", contact.New(client, apiURL))

	return logMiddleware(mux)
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
