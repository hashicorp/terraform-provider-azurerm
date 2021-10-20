package videoanalyzer

type AccessPolicyEntity struct {
	Id         *string                 `json:"id,omitempty"`
	Name       *string                 `json:"name,omitempty"`
	Properties *AccessPolicyProperties `json:"properties,omitempty"`
	SystemData *SystemData             `json:"systemData,omitempty"`
	Type       *string                 `json:"type,omitempty"`
}
