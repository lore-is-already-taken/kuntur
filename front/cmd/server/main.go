// Command server runs the Kuntur music group website.
//
// It is the composition root of the application: it reads configuration from
// the environment, wires the HTTP router, starts an http.Server with safe
// timeouts, and shuts it down gracefully when it receives SIGINT/SIGTERM.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kuntur/app/server"
)

const (
	defaultAddr            = ":8080"
	defaultShutdownTimeout = 10 * time.Second
	readHeaderTimeout      = 5 * time.Second
	readTimeout            = 10 * time.Second
	writeTimeout           = 10 * time.Second
	idleTimeout            = 60 * time.Second
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	srv := &http.Server{
		Addr:              envOr("ADDR", defaultAddr),
		Handler:           server.NewRouter(),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	if err := run(srv); err != nil {
		logger.Error("server stopped with error", "err", err)
		os.Exit(1)
	}
}

// run blocks until the server exits or receives a termination signal, then
// performs a graceful shutdown bounded by a timeout.
func run(srv *http.Server) error {
	logger := slog.Default()

	errCh := make(chan error, 1)
	go func() {
		logger.Info("listening", "addr", srv.Addr)
		// http.ErrServerClosed is the expected return value after Shutdown.
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case sig := <-stop:
		logger.Info("shutdown signal received", "signal", sig.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	logger.Info("server stopped cleanly")
	return nil
}

func envOr(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}
