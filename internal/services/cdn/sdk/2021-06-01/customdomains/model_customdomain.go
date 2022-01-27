package customdomains

type CustomDomain struct {
	Id         *string                 `json:"id,omitempty"`
	Name       *string                 `json:"name,omitempty"`
	Properties *CustomDomainProperties `json:"properties,omitempty"`
	SystemData *SystemData             `json:"systemData,omitempty"`
	Type       *string                 `json:"type,omitempty"`
}
