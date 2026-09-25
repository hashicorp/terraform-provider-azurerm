package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EffectiveNetworkSecurityRule struct {
	Access                           *SecurityRuleAccess            `json:"access,omitempty"`
	DestinationAddressPrefix         *string                        `json:"destinationAddressPrefix,omitempty"`
	DestinationAddressPrefixes       *[]string                      `json:"destinationAddressPrefixes,omitempty"`
	DestinationPortRange             *string                        `json:"destinationPortRange,omitempty"`
	DestinationPortRanges            *[]string                      `json:"destinationPortRanges,omitempty"`
	Direction                        *SecurityRuleDirection         `json:"direction,omitempty"`
	ExpandedDestinationAddressPrefix *[]string                      `json:"expandedDestinationAddressPrefix,omitempty"`
	ExpandedSourceAddressPrefix      *[]string                      `json:"expandedSourceAddressPrefix,omitempty"`
	Name                             *string                        `json:"name,omitempty"`
	Priority                         *int64                         `json:"priority,omitempty"`
	Protocol                         *EffectiveSecurityRuleProtocol `json:"protocol,omitempty"`
	SourceAddressPrefix              *string                        `json:"sourceAddressPrefix,omitempty"`
	SourceAddressPrefixes            *[]string                      `json:"sourceAddressPrefixes,omitempty"`
	SourcePortRange                  *string                        `json:"sourcePortRange,omitempty"`
	SourcePortRanges                 *[]string                      `json:"sourcePortRanges,omitempty"`
}
