package networkrulesets

type NetworkRuleSet struct {
	Id         *string                   `json:"id,omitempty"`
	Location   *string                   `json:"location,omitempty"`
	Name       *string                   `json:"name,omitempty"`
	Properties *NetworkRuleSetProperties `json:"properties,omitempty"`
	SystemData *SystemData               `json:"systemData,omitempty"`
	Type       *string                   `json:"type,omitempty"`
}
