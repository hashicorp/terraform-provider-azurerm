package logicalnetworks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogicalNetworkProperties struct {
	DhcpOptions       *LogicalNetworkPropertiesDhcpOptions `json:"dhcpOptions,omitempty"`
	ProvisioningState *ProvisioningStateEnum               `json:"provisioningState,omitempty"`
	Status            *LogicalNetworkStatus                `json:"status,omitempty"`
	Subnets           *[]Subnet                            `json:"subnets,omitempty"`
	VMSwitchName      *string                              `json:"vmSwitchName,omitempty"`
}
