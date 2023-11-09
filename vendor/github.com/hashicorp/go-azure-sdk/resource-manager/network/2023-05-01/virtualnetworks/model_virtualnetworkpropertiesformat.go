package virtualnetworks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkPropertiesFormat struct {
	AddressSpace           *AddressSpace                 `json:"addressSpace,omitempty"`
	BgpCommunities         *VirtualNetworkBgpCommunities `json:"bgpCommunities,omitempty"`
	DdosProtectionPlan     *SubResource                  `json:"ddosProtectionPlan,omitempty"`
	DhcpOptions            *DhcpOptions                  `json:"dhcpOptions,omitempty"`
	EnableDdosProtection   *bool                         `json:"enableDdosProtection,omitempty"`
	EnableVMProtection     *bool                         `json:"enableVmProtection,omitempty"`
	Encryption             *VirtualNetworkEncryption     `json:"encryption,omitempty"`
	FlowLogs               *[]FlowLog                    `json:"flowLogs,omitempty"`
	FlowTimeoutInMinutes   *int64                        `json:"flowTimeoutInMinutes,omitempty"`
	IPAllocations          *[]SubResource                `json:"ipAllocations,omitempty"`
	ProvisioningState      *ProvisioningState            `json:"provisioningState,omitempty"`
	ResourceGuid           *string                       `json:"resourceGuid,omitempty"`
	Subnets                *[]Subnet                     `json:"subnets,omitempty"`
	VirtualNetworkPeerings *[]VirtualNetworkPeering      `json:"virtualNetworkPeerings,omitempty"`
}
