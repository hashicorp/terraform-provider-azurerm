package virtualnetworks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkPeeringPropertiesFormat struct {
	AllowForwardedTraffic            *bool                         `json:"allowForwardedTraffic,omitempty"`
	AllowGatewayTransit              *bool                         `json:"allowGatewayTransit,omitempty"`
	AllowVirtualNetworkAccess        *bool                         `json:"allowVirtualNetworkAccess,omitempty"`
	DoNotVerifyRemoteGateways        *bool                         `json:"doNotVerifyRemoteGateways,omitempty"`
	EnableOnlyIPv6Peering            *bool                         `json:"enableOnlyIPv6Peering,omitempty"`
	LocalAddressSpace                *AddressSpace                 `json:"localAddressSpace,omitempty"`
	LocalSubnetNames                 *[]string                     `json:"localSubnetNames,omitempty"`
	LocalVirtualNetworkAddressSpace  *AddressSpace                 `json:"localVirtualNetworkAddressSpace,omitempty"`
	PeerCompleteVnets                *bool                         `json:"peerCompleteVnets,omitempty"`
	PeeringState                     *VirtualNetworkPeeringState   `json:"peeringState,omitempty"`
	PeeringSyncLevel                 *VirtualNetworkPeeringLevel   `json:"peeringSyncLevel,omitempty"`
	ProvisioningState                *ProvisioningState            `json:"provisioningState,omitempty"`
	RemoteAddressSpace               *AddressSpace                 `json:"remoteAddressSpace,omitempty"`
	RemoteBgpCommunities             *VirtualNetworkBgpCommunities `json:"remoteBgpCommunities,omitempty"`
	RemoteSubnetNames                *[]string                     `json:"remoteSubnetNames,omitempty"`
	RemoteVirtualNetwork             *SubResource                  `json:"remoteVirtualNetwork,omitempty"`
	RemoteVirtualNetworkAddressSpace *AddressSpace                 `json:"remoteVirtualNetworkAddressSpace,omitempty"`
	RemoteVirtualNetworkEncryption   *VirtualNetworkEncryption     `json:"remoteVirtualNetworkEncryption,omitempty"`
	ResourceGuid                     *string                       `json:"resourceGuid,omitempty"`
	UseRemoteGateways                *bool                         `json:"useRemoteGateways,omitempty"`
}
