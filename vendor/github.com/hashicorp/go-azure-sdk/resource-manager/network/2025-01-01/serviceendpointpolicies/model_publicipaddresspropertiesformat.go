package serviceendpointpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PublicIPAddressPropertiesFormat struct {
	DdosSettings             *DdosSettings                  `json:"ddosSettings,omitempty"`
	DeleteOption             *DeleteOptions                 `json:"deleteOption,omitempty"`
	DnsSettings              *PublicIPAddressDnsSettings    `json:"dnsSettings,omitempty"`
	IPAddress                *string                        `json:"ipAddress,omitempty"`
	IPConfiguration          *IPConfiguration               `json:"ipConfiguration,omitempty"`
	IPTags                   *[]IPTag                       `json:"ipTags,omitempty"`
	IdleTimeoutInMinutes     *int64                         `json:"idleTimeoutInMinutes,omitempty"`
	LinkedPublicIPAddress    *PublicIPAddress               `json:"linkedPublicIPAddress,omitempty"`
	MigrationPhase           *PublicIPAddressMigrationPhase `json:"migrationPhase,omitempty"`
	NatGateway               *NatGateway                    `json:"natGateway,omitempty"`
	ProvisioningState        *ProvisioningState             `json:"provisioningState,omitempty"`
	PublicIPAddressVersion   *IPVersion                     `json:"publicIPAddressVersion,omitempty"`
	PublicIPAllocationMethod *IPAllocationMethod            `json:"publicIPAllocationMethod,omitempty"`
	PublicIPPrefix           *SubResource                   `json:"publicIPPrefix,omitempty"`
	ResourceGuid             *string                        `json:"resourceGuid,omitempty"`
	ServicePublicIPAddress   *PublicIPAddress               `json:"servicePublicIPAddress,omitempty"`
}
