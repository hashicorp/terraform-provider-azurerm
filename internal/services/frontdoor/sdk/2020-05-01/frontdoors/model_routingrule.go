package frontdoors

type RoutingRule struct {
	Id         *string                `json:"id,omitempty"`
	Name       *string                `json:"name,omitempty"`
	Properties *RoutingRuleProperties `json:"properties,omitempty"`
	Type       *string                `json:"type,omitempty"`
}
