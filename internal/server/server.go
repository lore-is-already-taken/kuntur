package server

import (
	"context"
	"errors"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"time"
)

// Config holds the dependencies and settings for the HTTP server.
type Config struct {
	Addr      string       // e.g. ":8080"
	Logger    *slog.Logger // if nil, slog.Default() is used
	Templates fs.FS        // expected to contain "index.html" at the root of the sub-FS
	Static    fs.FS        // expected to be the result of fs.Sub(web.StaticFS, "static")
}

// Server is the HTTP server for kuntur.
type Server struct {
	cfg     Config
	handler http.Handler
	tmpl    *template.Template
}

// New constructs a Server from cfg. Returns an error if Templates or Static are nil,
// or if template parsing fails.
func New(cfg Config) (*Server, error) {
	if cfg.Templates == nil {
		return nil, errors.New("server: Config.Templates must not be nil")
	}
	if cfg.Static == nil {
		return nil, errors.New("server: Config.Static must not be nil")
	}
	if cfg.Logger == nil {
		cfg.Logger = slog.Default()
	}

	t, err := template.ParseFS(cfg.Templates, "*.html")
	if err != nil {
		return nil, err
	}

	s := &Server{cfg: cfg, tmpl: t}
	s.handler = s.loggingMiddleware(s.routes())
	return s, nil
}

// Routes returns the HTTP handler for the server. Useful for testing without
// starting a real TCP listener.
func (s *Server) Routes() http.Handler {
	return s.handler
}

// routes builds and returns the internal ServeMux (unwrapped, before middleware).
func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()

	// GET /{$} — exact match on "/" only (Go 1.22+ syntax)
	mux.HandleFunc("GET /{$}", s.handleIndex)

	// GET /static/ — serve embedded static assets
	staticHandler := http.StripPrefix("/static/", http.FileServer(http.FS(s.cfg.Static)))
	mux.Handle("GET /static/", staticHandler)

	return mux
}

// handleIndex renders the index template.
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
		s.cfg.Logger.Error("template execution failed", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// statusRecorder wraps http.ResponseWriter to capture the written status code.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

// Header returns the header map from the underlying ResponseWriter.
func (rec *statusRecorder) Header() http.Header {
	return rec.ResponseWriter.Header()
}

// loggingMiddleware wraps next and logs method, path, status and duration.
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		s.cfg.Logger.Info("http",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

// Start runs the HTTP server until ctx is cancelled or ListenAndServe returns
// a non-ErrServerClosed error. On ctx cancellation it initiates a graceful
// shutdown with a 5-second timeout.
func (s *Server) Start(ctx context.Context) error {
	httpSrv := &http.Server{
		Addr:    s.cfg.Addr,
		Handler: s.handler,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- httpSrv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return httpSrv.Shutdown(shutdownCtx)
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
