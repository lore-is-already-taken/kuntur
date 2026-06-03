// Package server implements the HTTP adapter for the Kuntur website.
//
// The router is built with the standard library's http.ServeMux, which since
// Go 1.22 supports method+pattern matching (e.g. "GET /{$}"). Static assets
// are served from an embedded filesystem owned by internal/web, so this
// package does not touch the real filesystem at runtime.
package server

import (
	"log/slog"
	"net/http"

	"kuntur/internal/router"
	"kuntur/internal/web"
)

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
	mux.HandleFunc("GET /{$}", router.HomeHandler(web.Templates()))

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
