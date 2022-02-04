package accounts

type VirtualNetworkRuleProperties struct {
	SubnetId                *string                  `json:"subnetId,omitempty"`
	VirtualNetworkRuleState *VirtualNetworkRuleState `json:"virtualNetworkRuleState,omitempty"`
}
