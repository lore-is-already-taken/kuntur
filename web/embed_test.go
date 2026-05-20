package web_test

import (
	"io/fs"
	"testing"

	"github.com/lore-is-already-taken/kuntur/web"
)

func TestTemplatesFS_containsIndex(t *testing.T) {
	_, err := fs.Stat(web.TemplatesFS, "templates/index.html")
	if err != nil {
		t.Fatalf("TemplatesFS does not contain templates/index.html: %v", err)
	}
}

func TestStaticFS_exists(t *testing.T) {
	_, err := fs.Stat(web.StaticFS, "static")
	if err != nil {
		t.Fatalf("StaticFS does not expose static directory: %v", err)
	}
}
