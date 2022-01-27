package securitypolicies

type SecurityPolicyWebApplicationFirewallAssociation struct {
	Domains         *[]ActivatedResourceReference `json:"domains,omitempty"`
	PatternsToMatch *[]string                     `json:"patternsToMatch,omitempty"`
}
