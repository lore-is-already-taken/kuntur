// Package bio renders the biography page ("/biografia") of the Kuntur site.
package bio

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"

	"kuntur/app/view"
)

// Get returns an http.HandlerFunc that renders the biography page.
func Get(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view.View(w, r, tmpl, getGroupInfo())
	}
}

func getGroupInfo() groupInfo {
	resp, err := http.Get("http://127.0.0.1:8000/bio")
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
