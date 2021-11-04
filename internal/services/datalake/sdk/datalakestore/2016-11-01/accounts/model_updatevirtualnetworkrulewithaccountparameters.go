package accounts

type UpdateVirtualNetworkRuleWithAccountParameters struct {
	Name       string                              `json:"name"`
	Properties *UpdateVirtualNetworkRuleProperties `json:"properties,omitempty"`
}
