package firewallrules

type FirewallRuleProperties struct {
	EndIpAddress   string `json:"endIpAddress"`
	StartIpAddress string `json:"startIpAddress"`
}
