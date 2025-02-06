package networkvirtualappliances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkVirtualAppliancePropertiesFormat struct {
	AdditionalNics              *[]VirtualApplianceAdditionalNicProperties `json:"additionalNics,omitempty"`
	AddressPrefix               *string                                    `json:"addressPrefix,omitempty"`
	BootStrapConfigurationBlobs *[]string                                  `json:"bootStrapConfigurationBlobs,omitempty"`
	CloudInitConfiguration      *string                                    `json:"cloudInitConfiguration,omitempty"`
	CloudInitConfigurationBlobs *[]string                                  `json:"cloudInitConfigurationBlobs,omitempty"`
	Delegation                  *DelegationProperties                      `json:"delegation,omitempty"`
	DeploymentType              *string                                    `json:"deploymentType,omitempty"`
	InboundSecurityRules        *[]SubResource                             `json:"inboundSecurityRules,omitempty"`
	InternetIngressPublicIPs    *[]InternetIngressPublicIPsProperties      `json:"internetIngressPublicIps,omitempty"`
	NvaSku                      *VirtualApplianceSkuProperties             `json:"nvaSku,omitempty"`
	PartnerManagedResource      *PartnerManagedResourceProperties          `json:"partnerManagedResource,omitempty"`
	ProvisioningState           *ProvisioningState                         `json:"provisioningState,omitempty"`
	SshPublicKey                *string                                    `json:"sshPublicKey,omitempty"`
	VirtualApplianceAsn         *int64                                     `json:"virtualApplianceAsn,omitempty"`
	VirtualApplianceConnections *[]SubResource                             `json:"virtualApplianceConnections,omitempty"`
	VirtualApplianceNics        *[]VirtualApplianceNicProperties           `json:"virtualApplianceNics,omitempty"`
	VirtualApplianceSites       *[]SubResource                             `json:"virtualApplianceSites,omitempty"`
	VirtualHub                  *SubResource                               `json:"virtualHub,omitempty"`
}
