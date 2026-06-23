// Package bio renders the biography page ("/biografia") of the Kuntur site.
package bio

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

func getGroupInfo(apiBaseURL string) groupInfo {
	resp, err := http.Get(apiBaseURL + "/bio")
	if err != nil {
		slog.Error("bio: backend unreachable", "err", err)
		return groupInfo{Resume: "layoutcomo estas mi rey precioso hermoso"}
	}
	defer resp.Body.Close()
	slog.Info("bio: backend response", "status", resp.StatusCode)

	var raw bioResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		slog.Error("bio: failed to decode response", "err", err)
		return groupInfo{Resume: "layoutcomo estas mi rey precioso hermoso"}
	}

	return toGroupInfo(raw)
}
