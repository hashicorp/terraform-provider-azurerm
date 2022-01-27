package afdcustomdomains

type AFDDomain struct {
	Id         *string              `json:"id,omitempty"`
	Name       *string              `json:"name,omitempty"`
	Properties *AFDDomainProperties `json:"properties,omitempty"`
	SystemData *SystemData          `json:"systemData,omitempty"`
	Type       *string              `json:"type,omitempty"`
}
