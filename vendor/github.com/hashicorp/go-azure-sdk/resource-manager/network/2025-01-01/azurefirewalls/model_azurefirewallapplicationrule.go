package azurefirewalls

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AzureFirewallApplicationRule struct {
	Description     *string                                 `json:"description,omitempty"`
	FqdnTags        *[]string                               `json:"fqdnTags,omitempty"`
	Name            *string                                 `json:"name,omitempty"`
	Protocols       *[]AzureFirewallApplicationRuleProtocol `json:"protocols,omitempty"`
	SourceAddresses *[]string                               `json:"sourceAddresses,omitempty"`
	SourceIPGroups  *[]string                               `json:"sourceIpGroups,omitempty"`
	TargetFqdns     *[]string                               `json:"targetFqdns,omitempty"`
}
