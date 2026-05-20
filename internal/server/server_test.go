package server_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/lore-is-already-taken/kuntur/internal/server"
)

// indexTemplate is a minimal Go template used across handler tests.
// Markers match the real web/templates/index.html so assertions stay meaningful.
const indexTemplate = `<!doctype html>
<html lang="en">
<head><meta charset="UTF-8" /><title>kuntur</title></head>
<body>
<div id="app"><button id="counter" type="button"></button></div>
<p class="read-the-docs">Click on the Vite and TypeScript logos to learn more</p>
</body>
</html>
`

// newTestServer builds a Server with in-memory FS fakes.
func newTestServer(t *testing.T) *server.Server {
	t.Helper()
	templates := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte(indexTemplate)},
	}
	static := fstest.MapFS{
		"style.css": &fstest.MapFile{Data: []byte(".foo{}")},
	}
	s, err := server.New(server.Config{
		Addr:      ":0",
		Templates: templates,
		Static:    static,
	})
	if err != nil {
		t.Fatalf("server.New: %v", err)
	}
	return s
}

func TestHandlers(t *testing.T) {
	s := newTestServer(t)
	handler := s.Routes()

	tests := []struct {
		name           string
		method         string
		path           string
		wantStatus     int
		wantBodyContains string
		wantContentType  string
	}{
		{
			name:             "index_returns_html_with_title",
			method:           http.MethodGet,
			path:             "/",
			wantStatus:       http.StatusOK,
			wantBodyContains: "<title>kuntur</title>",
			wantContentType:  "text/html; charset=utf-8",
		},
		{
			name:             "index_contains_counter_element",
			method:           http.MethodGet,
			path:             "/",
			wantStatus:       http.StatusOK,
			wantBodyContains: `id="counter"`,
			wantContentType:  "text/html; charset=utf-8",
		},
		{
			name:             "index_contains_read_the_docs",
			method:           http.MethodGet,
			path:             "/",
			wantStatus:       http.StatusOK,
			wantBodyContains: "read-the-docs",
			wantContentType:  "text/html; charset=utf-8",
		},
		{
			name:       "post_on_root_returns_405",
			method:     http.MethodPost,
			path:       "/",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "unknown_path_returns_404",
			method:     http.MethodGet,
			path:       "/does-not-exist",
			wantStatus: http.StatusNotFound,
		},
		{
			name:             "static_css_returns_200",
			method:           http.MethodGet,
			path:             "/static/style.css",
			wantStatus:       http.StatusOK,
			wantBodyContains: ".foo{}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatus {
				t.Errorf("status: got %d, want %d", res.StatusCode, tt.wantStatus)
			}

			if tt.wantContentType != "" {
				ct := res.Header.Get("Content-Type")
				if !strings.HasPrefix(ct, tt.wantContentType) {
					t.Errorf("Content-Type: got %q, want prefix %q", ct, tt.wantContentType)
				}
			}

			if tt.wantBodyContains != "" {
				body, _ := io.ReadAll(res.Body)
				if !strings.Contains(string(body), tt.wantBodyContains) {
					t.Errorf("body does not contain %q; got: %s", tt.wantBodyContains, body)
				}
			}
		})
	}
}

func TestNew_nilTemplates_returnsError(t *testing.T) {
	_, err := server.New(server.Config{
		Addr:      ":0",
		Templates: nil,
		Static:    fstest.MapFS{},
	})
	if err == nil {
		t.Fatal("expected error for nil Templates, got nil")
	}
}

func TestNew_nilStatic_returnsError(t *testing.T) {
	templates := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte(indexTemplate)},
	}
	_, err := server.New(server.Config{
		Addr:      ":0",
		Templates: templates,
		Static:    nil,
	})
	if err == nil {
		t.Fatal("expected error for nil Static, got nil")
	}
}

func TestNew_invalidTemplate_returnsError(t *testing.T) {
	bad := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte(`{{.Broken`)},
	}
	_, err := server.New(server.Config{
		Addr:      ":0",
		Templates: bad,
		Static:    fstest.MapFS{},
	})
	if err == nil {
		t.Fatal("expected template parse error, got nil")
	}
}

func TestStart_gracefulShutdown(t *testing.T) {
	s := newTestServer(t)

	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Start(ctx)
	}()

	// Give ListenAndServe time to bind the port.
	time.Sleep(50 * time.Millisecond)

	// Cancel context to trigger graceful shutdown.
	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Start returned unexpected error after ctx cancel: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Start did not return within 3 seconds after ctx cancel")
	}
}

// TestStart_respectsConfiguredShutdownTimeout asserts that a caller-provided
// ShutdownTimeout is wired through to http.Server.Shutdown. Uses a deliberately
// short timeout so a regression that hardcodes 5s would be caught by the outer
// 1-second test deadline.
func TestStart_respectsConfiguredShutdownTimeout(t *testing.T) {
	templates := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte(indexTemplate)},
	}
	s, err := server.New(server.Config{
		Addr:            ":0",
		Templates:       templates,
		Static:          fstest.MapFS{},
		ShutdownTimeout: 100 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("server.New: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Start(ctx)
	}()

	time.Sleep(50 * time.Millisecond)
	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Start returned unexpected error: %v", err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Start did not return within 1s with 100ms ShutdownTimeout")
	}
}
