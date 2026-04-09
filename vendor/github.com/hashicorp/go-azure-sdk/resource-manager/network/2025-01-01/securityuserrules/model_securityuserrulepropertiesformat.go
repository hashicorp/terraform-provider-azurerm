package securityuserrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityUserRulePropertiesFormat struct {
	Description           *string                            `json:"description,omitempty"`
	DestinationPortRanges *[]string                          `json:"destinationPortRanges,omitempty"`
	Destinations          *[]AddressPrefixItem               `json:"destinations,omitempty"`
	Direction             SecurityConfigurationRuleDirection `json:"direction"`
	Protocol              SecurityConfigurationRuleProtocol  `json:"protocol"`
	ProvisioningState     *ProvisioningState                 `json:"provisioningState,omitempty"`
	ResourceGuid          *string                            `json:"resourceGuid,omitempty"`
	SourcePortRanges      *[]string                          `json:"sourcePortRanges,omitempty"`
	Sources               *[]AddressPrefixItem               `json:"sources,omitempty"`
}
