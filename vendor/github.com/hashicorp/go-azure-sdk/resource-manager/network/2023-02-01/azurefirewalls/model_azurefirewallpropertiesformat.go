package azurefirewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFirewallPropertiesFormat struct {
	AdditionalProperties       *map[string]string                        `json:"additionalProperties,omitempty"`
	ApplicationRuleCollections *[]AzureFirewallApplicationRuleCollection `json:"applicationRuleCollections,omitempty"`
	FirewallPolicy             *SubResource                              `json:"firewallPolicy,omitempty"`
	HubIPAddresses             *HubIPAddresses                           `json:"hubIPAddresses,omitempty"`
	IPConfigurations           *[]AzureFirewallIPConfiguration           `json:"ipConfigurations,omitempty"`
	IPGroups                   *[]AzureFirewallIPGroups                  `json:"ipGroups,omitempty"`
	ManagementIPConfiguration  *AzureFirewallIPConfiguration             `json:"managementIpConfiguration,omitempty"`
	NatRuleCollections         *[]AzureFirewallNatRuleCollection         `json:"natRuleCollections,omitempty"`
	NetworkRuleCollections     *[]AzureFirewallNetworkRuleCollection     `json:"networkRuleCollections,omitempty"`
	ProvisioningState          *ProvisioningState                        `json:"provisioningState,omitempty"`
	Sku                        *AzureFirewallSku                         `json:"sku,omitempty"`
	ThreatIntelMode            *AzureFirewallThreatIntelMode             `json:"threatIntelMode,omitempty"`
	VirtualHub                 *SubResource                              `json:"virtualHub,omitempty"`
}
