package networkmanagereffectivesecurityadminrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DefaultAdminPropertiesFormat struct {
	Access                *SecurityConfigurationRuleAccess    `json:"access,omitempty"`
	Description           *string                             `json:"description,omitempty"`
	DestinationPortRanges *[]string                           `json:"destinationPortRanges,omitempty"`
	Destinations          *[]AddressPrefixItem                `json:"destinations,omitempty"`
	Direction             *SecurityConfigurationRuleDirection `json:"direction,omitempty"`
	Flag                  *string                             `json:"flag,omitempty"`
	Priority              *int64                              `json:"priority,omitempty"`
	Protocol              *SecurityConfigurationRuleProtocol  `json:"protocol,omitempty"`
	ProvisioningState     *ProvisioningState                  `json:"provisioningState,omitempty"`
	ResourceGuid          *string                             `json:"resourceGuid,omitempty"`
	SourcePortRanges      *[]string                           `json:"sourcePortRanges,omitempty"`
	Sources               *[]AddressPrefixItem                `json:"sources,omitempty"`
}
