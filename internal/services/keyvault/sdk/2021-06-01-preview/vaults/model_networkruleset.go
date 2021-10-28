package vaults

type NetworkRuleSet struct {
	Bypass              *NetworkRuleBypassOptions `json:"bypass,omitempty"`
	DefaultAction       *NetworkRuleAction        `json:"defaultAction,omitempty"`
	IpRules             *[]IPRule                 `json:"ipRules,omitempty"`
	VirtualNetworkRules *[]VirtualNetworkRule     `json:"virtualNetworkRules,omitempty"`
}
