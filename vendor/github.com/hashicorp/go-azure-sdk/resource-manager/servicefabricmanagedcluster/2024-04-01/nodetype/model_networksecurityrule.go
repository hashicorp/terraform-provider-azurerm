package nodetype

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkSecurityRule struct {
	Access                     Access      `json:"access"`
	Description                *string     `json:"description,omitempty"`
	DestinationAddressPrefix   *string     `json:"destinationAddressPrefix,omitempty"`
	DestinationAddressPrefixes *[]string   `json:"destinationAddressPrefixes,omitempty"`
	DestinationPortRange       *string     `json:"destinationPortRange,omitempty"`
	DestinationPortRanges      *[]string   `json:"destinationPortRanges,omitempty"`
	Direction                  Direction   `json:"direction"`
	Name                       string      `json:"name"`
	Priority                   int64       `json:"priority"`
	Protocol                   NsgProtocol `json:"protocol"`
	SourceAddressPrefix        *string     `json:"sourceAddressPrefix,omitempty"`
	SourceAddressPrefixes      *[]string   `json:"sourceAddressPrefixes,omitempty"`
	SourcePortRange            *string     `json:"sourcePortRange,omitempty"`
	SourcePortRanges           *[]string   `json:"sourcePortRanges,omitempty"`
}
