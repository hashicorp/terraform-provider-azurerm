package attestationproviders

type AttestationServiceCreationParams struct {
	Location   string                                   `json:"location"`
	Properties AttestationServiceCreationSpecificParams `json:"properties"`
	Tags       *map[string]string                       `json:"tags,omitempty"`
}
