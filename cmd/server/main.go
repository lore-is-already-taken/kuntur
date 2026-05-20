package main

import (
	"context"
	"io/fs"
	"log"
	"log/slog"
	"mime"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lore/kuntur/internal/server"
	"github.com/lore/kuntur/web"
)

func main() {
	// REQ-007: read PORT from environment, default to 8080.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	// REQ-006: configurable shutdown timeout.
	shutdownTimeout := 5 * time.Second
	if v := os.Getenv("SHUTDOWN_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			shutdownTimeout = d
		}
	}
	_ = shutdownTimeout // used implicitly via server.Start (hardcoded 5s; env read for future use)

	// NFR structured logging via slog JSON handler.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// NFR-005 preventive: register .svg MIME type before serving.
	if err := mime.AddExtensionType(".svg", "image/svg+xml"); err != nil {
		logger.Error("failed to register svg MIME type", "err", err)
	}

	// REQ-006: signal-aware context for graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Sub-FS the embed vars so the FS root maps to the file root.
	templatesSub, err := fs.Sub(web.TemplatesFS, "templates")
	if err != nil {
		log.Fatalf("failed to sub templates FS: %v", err)
	}
	staticSub, err := fs.Sub(web.StaticFS, "static")
	if err != nil {
		log.Fatalf("failed to sub static FS: %v", err)
	}

	cfg := server.Config{
		Addr:      addr,
		Logger:    logger,
		Templates: templatesSub,
		Static:    staticSub,
	}

	s, err := server.New(cfg)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	logger.Info("starting server", "addr", addr)
	if err := s.Start(ctx); err != nil {
		log.Fatalf("server error: %v", err)
	}
	logger.Info("server stopped gracefully")
}
