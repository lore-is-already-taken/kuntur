// Package biografia renders the biography page ("/biografia") of the Kuntur site.
package biografia

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"

	"kuntur/app/view"
)

// Get returns an http.HandlerFunc that renders the biography page. The
// apiBaseURL is the backend's base URL; the bio endpoint path is appended
// internally so the upstream location is not hardcoded in this package.
func Get(tmpl *template.Template, apiBaseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view.View(w, r, tmpl, getGroupInfo(apiBaseURL))
	}
}

// fallbackResume is shown when the backend is unreachable or returns an
// unreadable payload, so the page never renders empty.
const fallbackResume = "La biografía no está disponible en este momento. Inténtalo de nuevo más tarde."

func getGroupInfo(apiBaseURL string) groupInfo {
	resp, err := http.Get(apiBaseURL + "/bio")
	if err != nil {
		slog.Error("biografia: backend unreachable", "err", err)
		return groupInfo{Resume: fallbackResume}
	}
	defer resp.Body.Close()
	slog.Info("biografia: backend response", "status", resp.StatusCode)

	var raw bioResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		slog.Error("biografia: failed to decode response", "err", err)
		return groupInfo{Resume: fallbackResume}
	}

	return toGroupInfo(raw)
}
