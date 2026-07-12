package biografia

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// bioResponse mirrors the JSON contract from the /biografia API endpoint.
type bioResponse struct {
	Resume string `json:"resume"`
}

// memberResponse mirrors the JSON contract from the /integrantes API endpoint
// (back/app/types/integrantes.py::MemberResponse).
type memberResponse struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Instrument  []memberInstrument `json:"instrument"`
	Photo       string             `json:"photo"`
}

type memberInstrument struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// groupInfo is the view model passed to the template.
type groupInfo struct {
	Paragraphs []string
	Members    []member
}

// member is the per-integrante view model. Instruments is intentionally
// pre-formatted (e.g. "charango · quena") so the template stays free of
// logic. Initials feeds the photo placeholder when Photo is empty.
type member struct {
	Name        string
	Description string
	Instruments string
	Photo       string
	Initials    string
}

func toGroupInfo(r bioResponse) groupInfo {
	return groupInfo{
		Paragraphs: toParagraphs(r.Resume),
	}
}

// toParagraphs splits the plain-text resume on line breaks so the template
// can wrap each paragraph in its own element instead of relying on
// white-space tricks over a single block.
func toParagraphs(resume string) []string {
	var out []string
	for _, p := range strings.Split(resume, "\n") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func toMember(r memberResponse) member {
	names := make([]string, len(r.Instrument))
	for i, ins := range r.Instrument {
		names[i] = ins.Name
	}
	return member{
		Name:        r.Name,
		Description: r.Description,
		Instruments: strings.Join(names, " · "),
		// Trimmed so a whitespace-only value falls through to the
		// initials placeholder instead of rendering a broken <img>.
		Photo:    strings.TrimSpace(r.Photo),
		Initials: toInitials(r.Name),
	}
}

// toInitials builds the placeholder monogram from the first letter of the
// first two words of the name ("Lorenzo Saavedra" → "LS"). Single-word
// names yield a single letter; an empty name yields an empty string.
func toInitials(name string) string {
	words := strings.Fields(name)
	if len(words) > 2 {
		words = words[:2]
	}
	var b strings.Builder
	for _, w := range words {
		r, _ := utf8.DecodeRuneInString(w)
		b.WriteRune(unicode.ToUpper(r))
	}
	return b.String()
}
