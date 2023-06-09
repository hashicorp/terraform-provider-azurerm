package privatelinkservices

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateEndpointProperties struct {
	ApplicationSecurityGroups           *[]ApplicationSecurityGroup        `json:"applicationSecurityGroups,omitempty"`
	CustomDnsConfigs                    *[]CustomDnsConfigPropertiesFormat `json:"customDnsConfigs,omitempty"`
	CustomNetworkInterfaceName          *string                            `json:"customNetworkInterfaceName,omitempty"`
	IPConfigurations                    *[]PrivateEndpointIPConfiguration  `json:"ipConfigurations,omitempty"`
	ManualPrivateLinkServiceConnections *[]PrivateLinkServiceConnection    `json:"manualPrivateLinkServiceConnections,omitempty"`
	NetworkInterfaces                   *[]NetworkInterface                `json:"networkInterfaces,omitempty"`
	PrivateLinkServiceConnections       *[]PrivateLinkServiceConnection    `json:"privateLinkServiceConnections,omitempty"`
	ProvisioningState                   *ProvisioningState                 `json:"provisioningState,omitempty"`
	Subnet                              *Subnet                            `json:"subnet,omitempty"`
}
