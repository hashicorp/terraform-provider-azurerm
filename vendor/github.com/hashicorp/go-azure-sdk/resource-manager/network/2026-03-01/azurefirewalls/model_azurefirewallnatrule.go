package azurefirewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFirewallNatRule struct {
	Description          *string                             `json:"description,omitempty"`
	DestinationAddresses *[]string                           `json:"destinationAddresses,omitempty"`
	DestinationPorts     *[]string                           `json:"destinationPorts,omitempty"`
	Name                 *string                             `json:"name,omitempty"`
	Protocols            *[]AzureFirewallNetworkRuleProtocol `json:"protocols,omitempty"`
	SourceAddresses      *[]string                           `json:"sourceAddresses,omitempty"`
	SourceIPGroups       *[]string                           `json:"sourceIpGroups,omitempty"`
	TranslatedAddress    *string                             `json:"translatedAddress,omitempty"`
	TranslatedFqdn       *string                             `json:"translatedFqdn,omitempty"`
	TranslatedPort       *string                             `json:"translatedPort,omitempty"`
}
