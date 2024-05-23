package firewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallDeploymentProperties struct {
	AssociatedRulestack *RulestackDetails  `json:"associatedRulestack,omitempty"`
	DnsSettings         DNSSettings        `json:"dnsSettings"`
	FrontEndSettings    *[]FrontendSetting `json:"frontEndSettings,omitempty"`
	IsPanoramaManaged   *BooleanEnum       `json:"isPanoramaManaged,omitempty"`
	MarketplaceDetails  MarketplaceDetails `json:"marketplaceDetails"`
	NetworkProfile      NetworkProfile     `json:"networkProfile"`
	PanEtag             *string            `json:"panEtag,omitempty"`
	PanoramaConfig      *PanoramaConfig    `json:"panoramaConfig,omitempty"`
	PlanData            PlanData           `json:"planData"`
	ProvisioningState   *ProvisioningState `json:"provisioningState,omitempty"`
}
