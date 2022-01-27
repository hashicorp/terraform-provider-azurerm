package secrets

type Secret struct {
	Id         *string           `json:"id,omitempty"`
	Name       *string           `json:"name,omitempty"`
	Properties *SecretProperties `json:"properties,omitempty"`
	SystemData *SystemData       `json:"systemData,omitempty"`
	Type       *string           `json:"type,omitempty"`
}
