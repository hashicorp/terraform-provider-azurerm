package accounts

type UpdateFirewallRuleWithAccountParameters struct {
	Name       string                        `json:"name"`
	Properties *UpdateFirewallRuleProperties `json:"properties,omitempty"`
}
