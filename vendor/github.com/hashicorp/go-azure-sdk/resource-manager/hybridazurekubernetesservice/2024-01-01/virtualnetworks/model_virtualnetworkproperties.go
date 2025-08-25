package virtualnetworks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkProperties struct {
	DnsServers        *[]string                                  `json:"dnsServers,omitempty"`
	Gateway           *string                                    `json:"gateway,omitempty"`
	IPAddressPrefix   *string                                    `json:"ipAddressPrefix,omitempty"`
	InfraVnetProfile  *VirtualNetworkPropertiesInfraVnetProfile  `json:"infraVnetProfile,omitempty"`
	ProvisioningState *ProvisioningState                         `json:"provisioningState,omitempty"`
	Status            *VirtualNetworkPropertiesStatus            `json:"status,omitempty"`
	VMipPool          *[]VirtualNetworkPropertiesVMipPoolInlined `json:"vmipPool,omitempty"`
	VipPool           *[]VirtualNetworkPropertiesVipPoolInlined  `json:"vipPool,omitempty"`
	VlanID            *int64                                     `json:"vlanID,omitempty"`
}
