package publicipaddresses

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkTapPropertiesFormat struct {
	DestinationLoadBalancerFrontEndIPConfiguration *FrontendIPConfiguration            `json:"destinationLoadBalancerFrontEndIPConfiguration,omitempty"`
	DestinationNetworkInterfaceIPConfiguration     *NetworkInterfaceIPConfiguration    `json:"destinationNetworkInterfaceIPConfiguration,omitempty"`
	DestinationPort                                *int64                              `json:"destinationPort,omitempty"`
	NetworkInterfaceTapConfigurations              *[]NetworkInterfaceTapConfiguration `json:"networkInterfaceTapConfigurations,omitempty"`
	ProvisioningState                              *ProvisioningState                  `json:"provisioningState,omitempty"`
	ResourceGuid                                   *string                             `json:"resourceGuid,omitempty"`
}
