package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkVirtualApplianceConnectionProperties struct {
	Asn                    *int64                   `json:"asn,omitempty"`
	BgpPeerAddress         *[]string                `json:"bgpPeerAddress,omitempty"`
	EnableInternetSecurity *bool                    `json:"enableInternetSecurity,omitempty"`
	Name                   *string                  `json:"name,omitempty"`
	ProvisioningState      *ProvisioningState       `json:"provisioningState,omitempty"`
	RoutingConfiguration   *RoutingConfigurationNfv `json:"routingConfiguration,omitempty"`
	TunnelIdentifier       *int64                   `json:"tunnelIdentifier,omitempty"`
}
