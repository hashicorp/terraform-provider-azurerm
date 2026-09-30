package supercomputers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SupercomputerProperties struct {
	CustomerManagedKeys            *CustomerManagedKeys     `json:"customerManagedKeys,omitempty"`
	DiskEncryptionSetId            *string                  `json:"diskEncryptionSetId,omitempty"`
	Identities                     SupercomputerIdentities  `json:"identities"`
	LogAnalyticsClusterId          *string                  `json:"logAnalyticsClusterId,omitempty"`
	ManagedOnBehalfOfConfiguration *WithMoboBrokerResources `json:"managedOnBehalfOfConfiguration,omitempty"`
	ManagedResourceGroup           *string                  `json:"managedResourceGroup,omitempty"`
	ManagementSubnetId             *string                  `json:"managementSubnetId,omitempty"`
	OutboundType                   *NetworkEgressType       `json:"outboundType,omitempty"`
	ProvisioningState              *ProvisioningState       `json:"provisioningState,omitempty"`
	SubnetId                       string                   `json:"subnetId"`
	SystemSku                      *SystemSku               `json:"systemSku,omitempty"`
}
