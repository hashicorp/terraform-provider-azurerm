package virtualnetworkgatewayconnections

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocalNetworkGatewayPropertiesFormat struct {
	BgpSettings              *BgpSettings       `json:"bgpSettings,omitempty"`
	Fqdn                     *string            `json:"fqdn,omitempty"`
	GatewayIPAddress         *string            `json:"gatewayIpAddress,omitempty"`
	LocalNetworkAddressSpace *AddressSpace      `json:"localNetworkAddressSpace,omitempty"`
	ProvisioningState        *ProvisioningState `json:"provisioningState,omitempty"`
	ResourceGuid             *string            `json:"resourceGuid,omitempty"`
}
