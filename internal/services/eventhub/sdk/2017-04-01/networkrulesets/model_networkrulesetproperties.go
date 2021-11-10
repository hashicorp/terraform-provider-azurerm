package networkrulesets

type NetworkRuleSetProperties struct {
	DefaultAction       *DefaultAction                  `json:"defaultAction,omitempty"`
	IpRules             *[]NWRuleSetIpRules             `json:"ipRules,omitempty"`
	VirtualNetworkRules *[]NWRuleSetVirtualNetworkRules `json:"virtualNetworkRules,omitempty"`
}
