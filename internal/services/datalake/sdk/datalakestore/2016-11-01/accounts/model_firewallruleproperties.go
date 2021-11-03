package accounts

type FirewallRuleProperties struct {
	EndIpAddress   *string `json:"endIpAddress,omitempty"`
	StartIpAddress *string `json:"startIpAddress,omitempty"`
}
