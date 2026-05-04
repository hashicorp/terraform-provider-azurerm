package virtualrouters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualRouterPropertiesFormat struct {
	HostedGateway     *SubResource       `json:"hostedGateway,omitempty"`
	HostedSubnet      *SubResource       `json:"hostedSubnet,omitempty"`
	Peerings          *[]SubResource     `json:"peerings,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	VirtualRouterAsn  *int64             `json:"virtualRouterAsn,omitempty"`
	VirtualRouterIPs  *[]string          `json:"virtualRouterIps,omitempty"`
}
