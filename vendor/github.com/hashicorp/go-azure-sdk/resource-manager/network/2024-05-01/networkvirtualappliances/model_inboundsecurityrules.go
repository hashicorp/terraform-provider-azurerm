package networkvirtualappliances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InboundSecurityRules struct {
	AppliesOn             *[]string                     `json:"appliesOn,omitempty"`
	DestinationPortRange  *int64                        `json:"destinationPortRange,omitempty"`
	DestinationPortRanges *[]string                     `json:"destinationPortRanges,omitempty"`
	Name                  *string                       `json:"name,omitempty"`
	Protocol              *InboundSecurityRulesProtocol `json:"protocol,omitempty"`
	SourceAddressPrefix   *string                       `json:"sourceAddressPrefix,omitempty"`
}
