package cognitiveservicesaccounts

type NetworkRuleSet struct {
	DefaultAction       *NetworkRuleAction    `json:"defaultAction,omitempty"`
	IpRules             *[]IpRule             `json:"ipRules,omitempty"`
	VirtualNetworkRules *[]VirtualNetworkRule `json:"virtualNetworkRules,omitempty"`
}
