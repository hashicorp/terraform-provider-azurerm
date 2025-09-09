package firewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VwanConfiguration struct {
	IPOfTrustSubnetForUdr     *IPAddress      `json:"ipOfTrustSubnetForUdr,omitempty"`
	NetworkVirtualApplianceId *string         `json:"networkVirtualApplianceId,omitempty"`
	TrustSubnet               *IPAddressSpace `json:"trustSubnet,omitempty"`
	UnTrustSubnet             *IPAddressSpace `json:"unTrustSubnet,omitempty"`
	VHub                      IPAddressSpace  `json:"vHub"`
}
