package managedcluster

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Subnet struct {
	EnableIPv6                        *bool                              `json:"enableIpv6,omitempty"`
	Name                              string                             `json:"name"`
	NetworkSecurityGroupId            *string                            `json:"networkSecurityGroupId,omitempty"`
	PrivateEndpointNetworkPolicies    *PrivateEndpointNetworkPolicies    `json:"privateEndpointNetworkPolicies,omitempty"`
	PrivateLinkServiceNetworkPolicies *PrivateLinkServiceNetworkPolicies `json:"privateLinkServiceNetworkPolicies,omitempty"`
}
