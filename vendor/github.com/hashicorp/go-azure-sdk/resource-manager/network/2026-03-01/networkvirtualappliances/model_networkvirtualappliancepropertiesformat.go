package networkvirtualappliances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkVirtualAppliancePropertiesFormat struct {
	AdditionalNics              *[]VirtualApplianceAdditionalNicProperties             `json:"additionalNics,omitempty"`
	AddressFamily               *[]IPVersion                                           `json:"addressFamily,omitempty"`
	AddressPrefix               *string                                                `json:"addressPrefix,omitempty"`
	AddressPrefixV6             *string                                                `json:"addressPrefixV6,omitempty"`
	BootStrapConfigurationBlobs *[]string                                              `json:"bootStrapConfigurationBlobs,omitempty"`
	CloudInitConfiguration      *string                                                `json:"cloudInitConfiguration,omitempty"`
	CloudInitConfigurationBlobs *[]string                                              `json:"cloudInitConfigurationBlobs,omitempty"`
	Delegation                  *DelegationProperties                                  `json:"delegation,omitempty"`
	DeploymentType              *string                                                `json:"deploymentType,omitempty"`
	InboundSecurityRules        *[]SubResource                                         `json:"inboundSecurityRules,omitempty"`
	InternetIngressPublicIPs    *[]InternetIngressPublicIPsProperties                  `json:"internetIngressPublicIps,omitempty"`
	MigrationStatus             *NetworkVirtualApplianceMigrationStatus                `json:"migrationStatus,omitempty"`
	NetworkProfile              *NetworkVirtualAppliancePropertiesFormatNetworkProfile `json:"networkProfile,omitempty"`
	NvaInterfaceConfigurations  *[]NvaInterfaceConfigurationsProperties                `json:"nvaInterfaceConfigurations,omitempty"`
	NvaSku                      *VirtualApplianceSkuProperties                         `json:"nvaSku,omitempty"`
	PartnerManagedResource      *PartnerManagedResourceProperties                      `json:"partnerManagedResource,omitempty"`
	PrivateIPAddress            *string                                                `json:"privateIpAddress,omitempty"`
	PrivateIPAddressV6          *string                                                `json:"privateIpAddressV6,omitempty"`
	ProvisioningState           *ProvisioningState                                     `json:"provisioningState,omitempty"`
	SshPublicKey                *string                                                `json:"sshPublicKey,omitempty"`
	VirtualApplianceAsn         *int64                                                 `json:"virtualApplianceAsn,omitempty"`
	VirtualApplianceConnections *[]SubResource                                         `json:"virtualApplianceConnections,omitempty"`
	VirtualApplianceNics        *[]VirtualApplianceNicProperties                       `json:"virtualApplianceNics,omitempty"`
	VirtualApplianceSites       *[]SubResource                                         `json:"virtualApplianceSites,omitempty"`
	VirtualHub                  *SubResource                                           `json:"virtualHub,omitempty"`
}
