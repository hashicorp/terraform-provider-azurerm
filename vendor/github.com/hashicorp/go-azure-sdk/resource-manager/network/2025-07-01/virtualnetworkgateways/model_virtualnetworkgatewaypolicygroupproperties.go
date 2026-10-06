package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkGatewayPolicyGroupProperties struct {
	IsDefault                         bool                                     `json:"isDefault"`
	PolicyMembers                     []VirtualNetworkGatewayPolicyGroupMember `json:"policyMembers"`
	Priority                          int64                                    `json:"priority"`
	ProvisioningState                 *ProvisioningState                       `json:"provisioningState,omitempty"`
	VngClientConnectionConfigurations *[]SubResource                           `json:"vngClientConnectionConfigurations,omitempty"`
}
