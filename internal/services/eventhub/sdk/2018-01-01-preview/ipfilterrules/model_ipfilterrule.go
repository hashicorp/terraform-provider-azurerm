package ipfilterrules

type IpFilterRule struct {
	Id         *string                 `json:"id,omitempty"`
	Name       *string                 `json:"name,omitempty"`
	Properties *IpFilterRuleProperties `json:"properties,omitempty"`
	Type       *string                 `json:"type,omitempty"`
}
