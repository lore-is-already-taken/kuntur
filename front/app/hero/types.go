package hero

import "strconv"

// showResponse mirrors the JSON contract from the backend's /show endpoint
// (back/app/types/show.py::ShowResponse).
type showResponse struct {
	ID    string            `json:"id"`
	Place showResponsePlace `json:"place"`
	Fecha showResponseFecha `json:"fecha"`
}

type showResponsePlace struct {
	Name      string  `json:"name"`
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Direction *string `json:"direction"`
}

type showResponseFecha struct {
	Mes  string `json:"mes"`
	Year int    `json:"year"`
}

// presentation is the view model passed to the template. It is intentionally
// pre-formatted (Date, CityLine) so the template stays free of logic.
type presentation struct {
	ID       string
	Date     string // e.g. "Jun 2026"
	Venue    string
	CityLine string // e.g. "Santiago, Chile"
}

func toPresentation(r showResponse) presentation {
	return presentation{
		ID:       r.ID,
		Date:     r.Fecha.Mes + " " + strconv.Itoa(r.Fecha.Year),
		Venue:    r.Place.Name,
		CityLine: r.Place.City + ", " + r.Place.Country,
	}
}
