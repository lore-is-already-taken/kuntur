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
// apiBaseURL is the backend's base URL; the biografia and integrantes
// endpoint paths are appended internally so the upstream locations are not
// hardcoded in this package.
func Get(tmpl *template.Template, apiBaseURL string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info := getGroupInfo(apiBaseURL)
		info.Members = getMembers(apiBaseURL)
		view.View(w, r, tmpl, info)
	}
}

// fallbackResume is shown when the backend is unreachable or returns an
// unreadable payload, so the page never renders empty.
const fallbackResume = "La biografía no está disponible en este momento. Inténtalo de nuevo más tarde."

func getGroupInfo(apiBaseURL string) groupInfo {
	resp, err := http.Get(apiBaseURL + "/biografia/")
	if err != nil {
		slog.Error("biografia: backend unreachable", "err", err)
		return groupInfo{Resume: fallbackResume}
	}
	defer resp.Body.Close()
	slog.Info("biografia: backend response", "status", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		slog.Error("biografia: backend returned non-2xx", "status", resp.StatusCode)
		return groupInfo{Resume: fallbackResume}
	}

	var raw bioResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		slog.Error("biografia: failed to decode response", "err", err)
		return groupInfo{Resume: fallbackResume}
	}

	return toGroupInfo(raw)
}

// getMembers fetches the band members. A page without members is still a
// valid page (the template has an empty state), so every failure path
// degrades to nil instead of an error.
func getMembers(apiBaseURL string) []member {
	resp, err := http.Get(apiBaseURL + "/integrantes/")
	if err != nil {
		slog.Error("biografia: integrantes unreachable", "err", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("biografia: integrantes returned non-2xx", "status", resp.StatusCode)
		return nil
	}

	var raw []memberResponse
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		slog.Error("biografia: failed to decode integrantes", "err", err)
		return nil
	}

	out := make([]member, len(raw))
	for i, r := range raw {
		out[i] = toMember(r)
	}
	return out
}
