package biografia

import "strings"

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
// pre-formatted (e.g. "charango · quena") so the template stays free of logic.
type member struct {
	Name        string
	Description string
	Instruments string
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
	}
}
