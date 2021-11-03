package accounts

type CreateOrUpdateFirewallRuleProperties struct {
	EndIpAddress   string `json:"endIpAddress"`
	StartIpAddress string `json:"startIpAddress"`
}
