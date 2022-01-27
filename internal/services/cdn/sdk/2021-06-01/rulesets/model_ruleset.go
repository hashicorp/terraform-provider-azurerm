package rulesets

type RuleSet struct {
	Id         *string            `json:"id,omitempty"`
	Name       *string            `json:"name,omitempty"`
	Properties *RuleSetProperties `json:"properties,omitempty"`
	SystemData *SystemData        `json:"systemData,omitempty"`
	Type       *string            `json:"type,omitempty"`
}
