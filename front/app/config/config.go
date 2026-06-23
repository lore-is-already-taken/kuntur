// Package config loads application configuration from the environment.
//
// Configuration is built once in the composition root (cmd/server/main.go)
// and passed down via constructors — never read from os.Getenv elsewhere.
// This keeps adapters free of hidden coupling to the process environment
// and makes the wiring explicit and testable.
package config

import "os"

// Config is the application configuration loaded from the environment.
type Config struct {
	// APIBaseURL is the base URL of the backend API that the front talks to
	// (e.g. "http://127.0.0.1:8000"). Individual endpoint paths are appended
	// by the callers, so changing the environment means changing one value.
	APIBaseURL string
	// ServerAddr is the address the HTTP server listens on (e.g. ":8080").
	ServerAddr string
}

// Load reads configuration from the environment, applying defaults for
// any value that is unset or empty.
func Load() Config {
	return Config{
		APIBaseURL: envOr("API_BASE_URL", "http://127.0.0.1:8000"),
		ServerAddr: envOr("ADDR", ":8080"),
	}
}

func envOr(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

