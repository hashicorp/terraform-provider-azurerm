package hybridconnections

type HybridConnection struct {
	Id         *string                     `json:"id,omitempty"`
	Name       *string                     `json:"name,omitempty"`
	Properties *HybridConnectionProperties `json:"properties,omitempty"`
	Type       *string                     `json:"type,omitempty"`
}
