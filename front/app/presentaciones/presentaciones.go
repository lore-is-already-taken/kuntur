// Package presentaciones renders the presentations page ("/presentaciones") of the Kuntur site.
package presentaciones

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"

	"kuntur/app/view"
)

// Get returns an http.HandlerFunc that renders the presentations page. The
// apiBaseURL is the backend's base URL; the shows endpoint path is appended
// internally so the upstream location is not hardcoded in this package.
func Get(tmpl *template.Template, apiBaseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view.View(w, r, tmpl, getPresentations(apiBaseURL))
	}
}

// getPresentations fetches the list of shows from the backend and maps it
// into the template view model. On any network, decode, or shape error the
// function logs via slog and returns an empty slice so the page still renders
// without a hard failure.
func getPresentations(apiBaseURL string) []presentation {
	resp, err := http.Get(apiBaseURL + "/presentaciones/")
	if err != nil {
		slog.Error("presentaciones: backend unreachable", "err", err)
		return nil
	}
	defer resp.Body.Close()
	slog.Info("presentaciones: backend response", "status", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		slog.Error("presentaciones: backend returned non-2xx", "status", resp.StatusCode)
		return nil
	}

	var raw []showResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		slog.Error("presentaciones: failed to decode response", "err", err)
		return nil
	}

	out := make([]presentation, len(raw))
	for i, r := range raw {
		out[i] = toPresentation(r)
	}
	return out
}
