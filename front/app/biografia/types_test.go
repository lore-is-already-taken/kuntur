package biografia

import (
	"encoding/json"
	"testing"
)

func TestToInitials(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"two words", "Lorenzo Saavedra", "LS"},
		{"single word", "Victoria", "V"},
		{"more than two words uses first two", "Alan Marchant Mamani", "AM"},
		{"lowercase is uppercased", "dafne yufla", "DY"},
		{"multibyte first letter", "Ángel Pérez", "ÁP"},
		{"empty name", "", ""},
		{"only spaces", "   ", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toInitials(tt.in); got != tt.want {
				t.Fatalf("toInitials(%q): got %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestToMember_MapsPhotoAndInitials(t *testing.T) {
	r := memberResponse{
		Name:        "Victoria Flores",
		Description: "Voz principal.",
		Instrument: []memberInstrument{
			{Type: "voice", Name: "voz"},
			{Type: "string", Name: "charango"},
		},
		Photo: "/static/img/integrantes/victoria-flores.webp",
	}
	m := toMember(r)
	if m.Photo != r.Photo {
		t.Fatalf("photo: got %q, want %q", m.Photo, r.Photo)
	}
	if m.Initials != "VF" {
		t.Fatalf("initials: got %q, want VF", m.Initials)
	}
	if m.Instruments != "voz · charango" {
		t.Fatalf("instruments: got %q, want \"voz · charango\"", m.Instruments)
	}
}

func TestToMember_NullPhotoDecodesToEmpty(t *testing.T) {
	// The backend sends "photo": null for members without a portrait; the
	// decode must leave the zero value so the template falls back to the
	// initials placeholder. Decodes real JSON so the null-into-string
	// path is exercised, not just the struct zero value.
	payload := `[{"id":"1","name":"Alan Marchant","description":"d",
		"instrument":[{"type":"string","name":"charango"}],"photo":null}]`
	var raw []memberResponse
	if err := json.Unmarshal([]byte(payload), &raw); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	m := toMember(raw[0])
	if m.Photo != "" {
		t.Fatalf("photo: got %q, want empty", m.Photo)
	}
	if m.Initials != "AM" {
		t.Fatalf("initials: got %q, want AM", m.Initials)
	}
}

func TestToMember_WhitespacePhotoFallsBackToPlaceholder(t *testing.T) {
	// A whitespace-only photo would be truthy in the template's
	// {{ if .Photo }} and render a broken <img>; the mapping trims it.
	m := toMember(memberResponse{Name: "Alan Marchant", Photo: "   "})
	if m.Photo != "" {
		t.Fatalf("photo: got %q, want empty", m.Photo)
	}
}
