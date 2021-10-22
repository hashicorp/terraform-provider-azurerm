package attestationproviders

type AttestationProviders struct {
	Id         *string            `json:"id,omitempty"`
	Location   string             `json:"location"`
	Name       *string            `json:"name,omitempty"`
	Properties *StatusResult      `json:"properties,omitempty"`
	SystemData *SystemData        `json:"systemData,omitempty"`
	Tags       *map[string]string `json:"tags,omitempty"`
	Type       *string            `json:"type,omitempty"`
}
