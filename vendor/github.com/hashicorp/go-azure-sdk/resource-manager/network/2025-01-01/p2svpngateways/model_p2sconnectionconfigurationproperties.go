package p2svpngateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type P2SConnectionConfigurationProperties struct {
	ConfigurationPolicyGroupAssociations         *[]SubResource                       `json:"configurationPolicyGroupAssociations,omitempty"`
	EnableInternetSecurity                       *bool                                `json:"enableInternetSecurity,omitempty"`
	PreviousConfigurationPolicyGroupAssociations *[]VpnServerConfigurationPolicyGroup `json:"previousConfigurationPolicyGroupAssociations,omitempty"`
	ProvisioningState                            *ProvisioningState                   `json:"provisioningState,omitempty"`
	RoutingConfiguration                         *RoutingConfiguration                `json:"routingConfiguration,omitempty"`
	VpnClientAddressPool                         *AddressSpace                        `json:"vpnClientAddressPool,omitempty"`
}
