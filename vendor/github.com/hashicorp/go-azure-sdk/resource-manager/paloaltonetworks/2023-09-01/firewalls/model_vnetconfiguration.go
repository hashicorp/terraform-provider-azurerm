package firewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VnetConfiguration struct {
	IPOfTrustSubnetForUdr *IPAddress     `json:"ipOfTrustSubnetForUdr,omitempty"`
	TrustSubnet           IPAddressSpace `json:"trustSubnet"`
	UnTrustSubnet         IPAddressSpace `json:"unTrustSubnet"`
	Vnet                  IPAddressSpace `json:"vnet"`
}
