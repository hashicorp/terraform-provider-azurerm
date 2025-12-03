package firewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallPolicyDraftProperties struct {
	BasePolicy           *SubResource                        `json:"basePolicy,omitempty"`
	DnsSettings          *DnsSettings                        `json:"dnsSettings,omitempty"`
	ExplicitProxy        *ExplicitProxy                      `json:"explicitProxy,omitempty"`
	Insights             *FirewallPolicyInsights             `json:"insights,omitempty"`
	IntrusionDetection   *FirewallPolicyIntrusionDetection   `json:"intrusionDetection,omitempty"`
	Snat                 *FirewallPolicySNAT                 `json:"snat,omitempty"`
	Sql                  *FirewallPolicySQL                  `json:"sql,omitempty"`
	ThreatIntelMode      *AzureFirewallThreatIntelMode       `json:"threatIntelMode,omitempty"`
	ThreatIntelWhitelist *FirewallPolicyThreatIntelWhitelist `json:"threatIntelWhitelist,omitempty"`
}
