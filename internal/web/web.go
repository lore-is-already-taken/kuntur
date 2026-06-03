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

// Views returns a parsed template per view file. Each view is parsed into a
// fresh *template.Template seeded with a clone of the layout, so the block
// names that each view defines (e.g. "title", "content") live in their own
// namespace and never collide across views.
//
// Files in templates/ named layout.html are skipped — they are the shared
// shell consumed via {{ template "layout" . }} from each view.
//
// The package panics on parse errors because they are programming mistakes
// that should be caught the first time the binary runs.
func Views() map[string]*template.Template {
	layout := template.Must(template.ParseFS(templatesFS, "templates/layout.html"))

	entries, err := fs.ReadDir(templatesFS, "templates")
	if err != nil {
		panic("web: cannot read templates dir: " + err.Error())
	}

	views := make(map[string]*template.Template, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if name == "layout.html" {
			continue
		}

		// Clone the parsed layout into a fresh set, then layer the view on
		// top. Clone() is cheap: it copies the parse tree by reference.
		set, err := layout.Clone()
		if err != nil {
			panic("web: cannot clone layout: " + err.Error())
		}
		if _, err := set.ParseFS(templatesFS, "templates/"+name); err != nil {
			panic("web: cannot parse view " + name + ": " + err.Error())
		}
		views[name] = set
	}
	return views
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

