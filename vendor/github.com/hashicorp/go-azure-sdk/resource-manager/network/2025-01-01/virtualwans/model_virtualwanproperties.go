package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualWanProperties struct {
	AllowBranchToBranchTraffic     *bool                  `json:"allowBranchToBranchTraffic,omitempty"`
	AllowVnetToVnetTraffic         *bool                  `json:"allowVnetToVnetTraffic,omitempty"`
	DisableVpnEncryption           *bool                  `json:"disableVpnEncryption,omitempty"`
	Office365LocalBreakoutCategory *OfficeTrafficCategory `json:"office365LocalBreakoutCategory,omitempty"`
	ProvisioningState              *ProvisioningState     `json:"provisioningState,omitempty"`
	Type                           *string                `json:"type,omitempty"`
	VirtualHubs                    *[]SubResource         `json:"virtualHubs,omitempty"`
	VpnSites                       *[]SubResource         `json:"vpnSites,omitempty"`
}
