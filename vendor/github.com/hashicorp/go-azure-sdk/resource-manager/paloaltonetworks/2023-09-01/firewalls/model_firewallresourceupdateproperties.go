package firewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallResourceUpdateProperties struct {
	AssociatedRulestack *RulestackDetails   `json:"associatedRulestack,omitempty"`
	DnsSettings         *DNSSettings        `json:"dnsSettings,omitempty"`
	FrontEndSettings    *[]FrontendSetting  `json:"frontEndSettings,omitempty"`
	IsPanoramaManaged   *BooleanEnum        `json:"isPanoramaManaged,omitempty"`
	MarketplaceDetails  *MarketplaceDetails `json:"marketplaceDetails,omitempty"`
	NetworkProfile      *NetworkProfile     `json:"networkProfile,omitempty"`
	PanEtag             *string             `json:"panEtag,omitempty"`
	PanoramaConfig      *PanoramaConfig     `json:"panoramaConfig,omitempty"`
	PlanData            *PlanData           `json:"planData,omitempty"`
}
