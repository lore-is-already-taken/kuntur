// Package web bundles the static assets (HTML templates and files) used by
// the HTTP server.
//
// Keeping assets in their own package makes the embed directives predictable
// (paths are relative to this directory) and lets the server stay free of
// any knowledge about where the files live on disk. As the project grows
// you can swap this for an fs.FS backed by a CDN or object storage without
// touching internal/server.
package web

import (
	"embed"
	"html/template"
	"io/fs"
)

//go:embed all:templates
var templatesFS embed.FS

//go:embed all:static
var staticFS embed.FS

// Templates returns the parsed HTML templates, ready to be executed.
// The package panics on parse errors because they are programming mistakes
// that should be caught the first time the binary runs.
func Templates() *template.Template {
	return template.Must(template.ParseFS(templatesFS, "templates/*.html"))
}

// Static returns an fs.FS rooted at the static directory, suitable for
// http.FileServerFS. The fs.Sub call cannot fail at runtime: the embed.FS
// already contains a "static" entry, so the panic is unreachable.
func Static() fs.FS {
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic("web: static directory missing from embed.FS: " + err.Error())
	}
	return sub
}
