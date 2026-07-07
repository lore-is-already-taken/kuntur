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
	Resume  string
	Members []member
}

// member is the per-integrante view model. Instruments is intentionally
// pre-formatted (e.g. "charango · quena") so the template stays free of logic.
type member struct {
	Name        string
	Description string
	Instruments string
}

func toGroupInfo(r bioResponse) groupInfo {
	//nolint:gosimple //intentional: adapter decouples API contract from view model
	return groupInfo{
		Resume: r.Resume,
	}
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
