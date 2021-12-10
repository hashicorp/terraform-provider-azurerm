package networkrulesets

type NetworkRuleSet struct {
	Id         *string                   `json:"id,omitempty"`
	Name       *string                   `json:"name,omitempty"`
	Properties *NetworkRuleSetProperties `json:"properties,omitempty"`
	Type       *string                   `json:"type,omitempty"`
}
