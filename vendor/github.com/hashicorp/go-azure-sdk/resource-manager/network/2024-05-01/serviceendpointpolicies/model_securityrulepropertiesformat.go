package serviceendpointpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityRulePropertiesFormat struct {
	Access                               SecurityRuleAccess          `json:"access"`
	Description                          *string                     `json:"description,omitempty"`
	DestinationAddressPrefix             *string                     `json:"destinationAddressPrefix,omitempty"`
	DestinationAddressPrefixes           *[]string                   `json:"destinationAddressPrefixes,omitempty"`
	DestinationApplicationSecurityGroups *[]ApplicationSecurityGroup `json:"destinationApplicationSecurityGroups,omitempty"`
	DestinationPortRange                 *string                     `json:"destinationPortRange,omitempty"`
	DestinationPortRanges                *[]string                   `json:"destinationPortRanges,omitempty"`
	Direction                            SecurityRuleDirection       `json:"direction"`
	Priority                             int64                       `json:"priority"`
	Protocol                             SecurityRuleProtocol        `json:"protocol"`
	ProvisioningState                    *ProvisioningState          `json:"provisioningState,omitempty"`
	SourceAddressPrefix                  *string                     `json:"sourceAddressPrefix,omitempty"`
	SourceAddressPrefixes                *[]string                   `json:"sourceAddressPrefixes,omitempty"`
	SourceApplicationSecurityGroups      *[]ApplicationSecurityGroup `json:"sourceApplicationSecurityGroups,omitempty"`
	SourcePortRange                      *string                     `json:"sourcePortRange,omitempty"`
	SourcePortRanges                     *[]string                   `json:"sourcePortRanges,omitempty"`
}
