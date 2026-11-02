package vnetpeering

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkPeeringPropertiesFormat struct {
	AllowForwardedTraffic     *bool                                                          `json:"allowForwardedTraffic,omitempty"`
	AllowGatewayTransit       *bool                                                          `json:"allowGatewayTransit,omitempty"`
	AllowVirtualNetworkAccess *bool                                                          `json:"allowVirtualNetworkAccess,omitempty"`
	DatabricksAddressSpace    *AddressSpace                                                  `json:"databricksAddressSpace,omitempty"`
	DatabricksVirtualNetwork  *VirtualNetworkPeeringPropertiesFormatDatabricksVirtualNetwork `json:"databricksVirtualNetwork,omitempty"`
	PeeringState              *PeeringState                                                  `json:"peeringState,omitempty"`
	ProvisioningState         *PeeringProvisioningState                                      `json:"provisioningState,omitempty"`
	RemoteAddressSpace        *AddressSpace                                                  `json:"remoteAddressSpace,omitempty"`
	RemoteVirtualNetwork      VirtualNetworkPeeringPropertiesFormatRemoteVirtualNetwork      `json:"remoteVirtualNetwork"`
	UseRemoteGateways         *bool                                                          `json:"useRemoteGateways,omitempty"`
}
