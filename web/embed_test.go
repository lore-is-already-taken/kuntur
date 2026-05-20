package web_test

import (
	"io/fs"
	"testing"

	"github.com/lore-is-already-taken/kuntur/web"
)

// TestTemplatesFS_containsIndex verifies that the public sub-FS interface
// (post fs.Sub) exposes "index.html" at its root — the same shape callers receive.
func TestTemplatesFS_containsIndex(t *testing.T) {
	sub, err := fs.Sub(web.TemplatesFS, "templates")
	if err != nil {
		t.Fatalf("fs.Sub(TemplatesFS, templates): %v", err)
	}
	_, err = fs.Stat(sub, "index.html")
	if err != nil {
		t.Fatalf("TemplatesFS sub does not contain index.html: %v", err)
	}
}

// TestStaticFS_containsExpectedAssets verifies that the static embed exposes
// the real assets added in PR #2 via the public sub-FS interface.
func TestStaticFS_containsExpectedAssets(t *testing.T) {
	sub, err := fs.Sub(web.StaticFS, "static")
	if err != nil {
		t.Fatalf("fs.Sub(StaticFS, static): %v", err)
	}
	assets := []string{
		"css/style.css",
		"js/counter.js",
		"svg/vite.svg",
		"svg/typescript.svg",
	}
	for _, asset := range assets {
		t.Run(asset, func(t *testing.T) {
			_, err := fs.Stat(sub, asset)
			if err != nil {
				t.Errorf("StaticFS sub does not contain %s: %v", asset, err)
			}
		})
	}
}
