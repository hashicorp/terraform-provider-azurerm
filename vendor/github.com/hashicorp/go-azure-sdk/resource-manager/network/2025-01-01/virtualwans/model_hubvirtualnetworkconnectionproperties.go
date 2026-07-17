package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HubVirtualNetworkConnectionProperties struct {
	AllowHubToRemoteVnetTransit         *bool                 `json:"allowHubToRemoteVnetTransit,omitempty"`
	AllowRemoteVnetToUseHubVnetGateways *bool                 `json:"allowRemoteVnetToUseHubVnetGateways,omitempty"`
	EnableInternetSecurity              *bool                 `json:"enableInternetSecurity,omitempty"`
	ProvisioningState                   *ProvisioningState    `json:"provisioningState,omitempty"`
	RemoteVirtualNetwork                *SubResource          `json:"remoteVirtualNetwork,omitempty"`
	RoutingConfiguration                *RoutingConfiguration `json:"routingConfiguration,omitempty"`
}
