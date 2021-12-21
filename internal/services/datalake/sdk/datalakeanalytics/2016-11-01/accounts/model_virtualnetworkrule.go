package accounts

type VirtualNetworkRule struct {
	Id         *string                       `json:"id,omitempty"`
	Name       *string                       `json:"name,omitempty"`
	Properties *VirtualNetworkRuleProperties `json:"properties,omitempty"`
	Type       *string                       `json:"type,omitempty"`
}
