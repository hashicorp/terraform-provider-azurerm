package webapplicationfirewallpolicies

type CdnWebApplicationFirewallPolicyProperties struct {
	CustomRules       *CustomRuleList      `json:"customRules,omitempty"`
	EndpointLinks     *[]CdnEndpoint       `json:"endpointLinks,omitempty"`
	ManagedRules      *ManagedRuleSetList  `json:"managedRules,omitempty"`
	PolicySettings    *PolicySettings      `json:"policySettings,omitempty"`
	ProvisioningState *ProvisioningState   `json:"provisioningState,omitempty"`
	RateLimitRules    *RateLimitRuleList   `json:"rateLimitRules,omitempty"`
	ResourceState     *PolicyResourceState `json:"resourceState,omitempty"`
}
