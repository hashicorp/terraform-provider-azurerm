package firewallresources

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallResourceUpdateProperties struct {
	AssociatedRulestack      *RulestackDetails         `json:"associatedRulestack,omitempty"`
	DnsSettings              *DNSSettings              `json:"dnsSettings,omitempty"`
	FrontEndSettings         *[]FrontendSetting        `json:"frontEndSettings,omitempty"`
	IsPanoramaManaged        *BooleanEnum              `json:"isPanoramaManaged,omitempty"`
	IsStrataCloudManaged     *BooleanEnum              `json:"isStrataCloudManaged,omitempty"`
	MarketplaceDetails       *MarketplaceDetails       `json:"marketplaceDetails,omitempty"`
	NetworkProfile           *NetworkProfile           `json:"networkProfile,omitempty"`
	PanEtag                  *string                   `json:"panEtag,omitempty"`
	PanoramaConfig           *PanoramaConfig           `json:"panoramaConfig,omitempty"`
	PlanData                 *PlanData                 `json:"planData,omitempty"`
	StrataCloudManagerConfig *StrataCloudManagerConfig `json:"strataCloudManagerConfig,omitempty"`
}
