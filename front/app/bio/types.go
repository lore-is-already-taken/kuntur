package bio

// bioResponse mirrors the JSON contract from the /bio API endpoint.
type bioResponse struct {
	Resume string `json:"resume"`
}

// groupInfo is the view model passed to the template.
type groupInfo struct {
	Resume string
}

func toGroupInfo(r bioResponse) groupInfo {
	//nolint:gosimple //intentional: adapter decouples API contract from view model
	return groupInfo{
		Resume: r.Resume,
	}
}
