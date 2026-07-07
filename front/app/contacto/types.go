// Package contacto owns the contact form resource for the Kuntur site,
// including the GET view (rendering the form) and the POST handler
// (submitting the form to the upstream API).
package contacto

// Payload is the JSON body sent to the upstream contact API.
type Payload struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}
