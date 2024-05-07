package firewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPolicyPropertiesFormat struct {
	BasePolicy           *SubResource                        `json:"basePolicy,omitempty"`
	ChildPolicies        *[]SubResource                      `json:"childPolicies,omitempty"`
	DnsSettings          *DnsSettings                        `json:"dnsSettings,omitempty"`
	ExplicitProxy        *ExplicitProxy                      `json:"explicitProxy,omitempty"`
	Firewalls            *[]SubResource                      `json:"firewalls,omitempty"`
	Insights             *FirewallPolicyInsights             `json:"insights,omitempty"`
	IntrusionDetection   *FirewallPolicyIntrusionDetection   `json:"intrusionDetection,omitempty"`
	ProvisioningState    *ProvisioningState                  `json:"provisioningState,omitempty"`
	RuleCollectionGroups *[]SubResource                      `json:"ruleCollectionGroups,omitempty"`
	Sku                  *FirewallPolicySku                  `json:"sku,omitempty"`
	Snat                 *FirewallPolicySNAT                 `json:"snat,omitempty"`
	Sql                  *FirewallPolicySQL                  `json:"sql,omitempty"`
	ThreatIntelMode      *AzureFirewallThreatIntelMode       `json:"threatIntelMode,omitempty"`
	ThreatIntelWhitelist *FirewallPolicyThreatIntelWhitelist `json:"threatIntelWhitelist,omitempty"`
	TransportSecurity    *FirewallPolicyTransportSecurity    `json:"transportSecurity,omitempty"`
}
