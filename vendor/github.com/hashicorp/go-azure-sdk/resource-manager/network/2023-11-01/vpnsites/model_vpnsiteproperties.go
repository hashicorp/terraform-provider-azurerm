package vpnsites

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VpnSiteProperties struct {
	AddressSpace      *AddressSpace         `json:"addressSpace,omitempty"`
	BgpProperties     *BgpSettings          `json:"bgpProperties,omitempty"`
	DeviceProperties  *DeviceProperties     `json:"deviceProperties,omitempty"`
	IPAddress         *string               `json:"ipAddress,omitempty"`
	IsSecuritySite    *bool                 `json:"isSecuritySite,omitempty"`
	O365Policy        *O365PolicyProperties `json:"o365Policy,omitempty"`
	ProvisioningState *ProvisioningState    `json:"provisioningState,omitempty"`
	SiteKey           *string               `json:"siteKey,omitempty"`
	VirtualWAN        *SubResource          `json:"virtualWan,omitempty"`
	VpnSiteLinks      *[]VpnSiteLink        `json:"vpnSiteLinks,omitempty"`
}
