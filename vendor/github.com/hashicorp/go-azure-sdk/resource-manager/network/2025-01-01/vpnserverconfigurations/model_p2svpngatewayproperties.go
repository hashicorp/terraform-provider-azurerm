package vpnserverconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type P2SVpnGatewayProperties struct {
	CustomDnsServers            *[]string                     `json:"customDnsServers,omitempty"`
	IsRoutingPreferenceInternet *bool                         `json:"isRoutingPreferenceInternet,omitempty"`
	P2SConnectionConfigurations *[]P2SConnectionConfiguration `json:"p2SConnectionConfigurations,omitempty"`
	ProvisioningState           *ProvisioningState            `json:"provisioningState,omitempty"`
	VirtualHub                  *SubResource                  `json:"virtualHub,omitempty"`
	VpnClientConnectionHealth   *VpnClientConnectionHealth    `json:"vpnClientConnectionHealth,omitempty"`
	VpnGatewayScaleUnit         *int64                        `json:"vpnGatewayScaleUnit,omitempty"`
	VpnServerConfiguration      *SubResource                  `json:"vpnServerConfiguration,omitempty"`
}
