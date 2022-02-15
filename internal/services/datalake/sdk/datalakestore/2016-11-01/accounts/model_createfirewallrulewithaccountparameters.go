package accounts

type CreateFirewallRuleWithAccountParameters struct {
	Name       string                               `json:"name"`
	Properties CreateOrUpdateFirewallRuleProperties `json:"properties"`
}
