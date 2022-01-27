package rules

type Rule struct {
	Id         *string         `json:"id,omitempty"`
	Name       *string         `json:"name,omitempty"`
	Properties *RuleProperties `json:"properties,omitempty"`
	SystemData *SystemData     `json:"systemData,omitempty"`
	Type       *string         `json:"type,omitempty"`
}
