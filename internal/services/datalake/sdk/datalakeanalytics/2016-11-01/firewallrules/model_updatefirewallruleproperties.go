package firewallrules

type UpdateFirewallRuleProperties struct {
	EndIpAddress   *string `json:"endIpAddress,omitempty"`
	StartIpAddress *string `json:"startIpAddress,omitempty"`
}
