package managedhsms

type MHSMNetworkRuleSet struct {
	Bypass              *NetworkRuleBypassOptions `json:"bypass,omitempty"`
	DefaultAction       *NetworkRuleAction        `json:"defaultAction,omitempty"`
	IpRules             *[]MHSMIPRule             `json:"ipRules,omitempty"`
	VirtualNetworkRules *[]MHSMVirtualNetworkRule `json:"virtualNetworkRules,omitempty"`
}
