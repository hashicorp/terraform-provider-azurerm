package networksecuritygroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkInterfacePropertiesFormat struct {
	AuxiliaryMode               *NetworkInterfaceAuxiliaryMode      `json:"auxiliaryMode,omitempty"`
	AuxiliarySku                *NetworkInterfaceAuxiliarySku       `json:"auxiliarySku,omitempty"`
	DisableTcpStateTracking     *bool                               `json:"disableTcpStateTracking,omitempty"`
	DnsSettings                 *NetworkInterfaceDnsSettings        `json:"dnsSettings,omitempty"`
	DscpConfiguration           *SubResource                        `json:"dscpConfiguration,omitempty"`
	EnableAcceleratedNetworking *bool                               `json:"enableAcceleratedNetworking,omitempty"`
	EnableIPForwarding          *bool                               `json:"enableIPForwarding,omitempty"`
	HostedWorkloads             *[]string                           `json:"hostedWorkloads,omitempty"`
	IPConfigurations            *[]NetworkInterfaceIPConfiguration  `json:"ipConfigurations,omitempty"`
	MacAddress                  *string                             `json:"macAddress,omitempty"`
	MigrationPhase              *NetworkInterfaceMigrationPhase     `json:"migrationPhase,omitempty"`
	NetworkSecurityGroup        *NetworkSecurityGroup               `json:"networkSecurityGroup,omitempty"`
	NicType                     *NetworkInterfaceNicType            `json:"nicType,omitempty"`
	Primary                     *bool                               `json:"primary,omitempty"`
	PrivateEndpoint             *PrivateEndpoint                    `json:"privateEndpoint,omitempty"`
	PrivateLinkService          *PrivateLinkService                 `json:"privateLinkService,omitempty"`
	ProvisioningState           *ProvisioningState                  `json:"provisioningState,omitempty"`
	ResourceGuid                *string                             `json:"resourceGuid,omitempty"`
	TapConfigurations           *[]NetworkInterfaceTapConfiguration `json:"tapConfigurations,omitempty"`
	VirtualMachine              *SubResource                        `json:"virtualMachine,omitempty"`
	VnetEncryptionSupported     *bool                               `json:"vnetEncryptionSupported,omitempty"`
	WorkloadType                *string                             `json:"workloadType,omitempty"`
}
