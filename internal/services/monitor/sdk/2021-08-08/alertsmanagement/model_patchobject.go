package alertsmanagement

type PatchObject struct {
	Properties *PatchProperties   `json:"properties,omitempty"`
	Tags       *map[string]string `json:"tags,omitempty"`
}
