package registries

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistryPropertiesUpdateParameters struct {
	AdminUserEnabled         *bool                     `json:"adminUserEnabled,omitempty"`
	AnonymousPullEnabled     *bool                     `json:"anonymousPullEnabled,omitempty"`
	DataEndpointEnabled      *bool                     `json:"dataEndpointEnabled,omitempty"`
	Encryption               *EncryptionProperty       `json:"encryption,omitempty"`
	NetworkRuleBypassOptions *NetworkRuleBypassOptions `json:"networkRuleBypassOptions,omitempty"`
	NetworkRuleSet           *NetworkRuleSet           `json:"networkRuleSet,omitempty"`
	Policies                 *Policies                 `json:"policies,omitempty"`
	PublicNetworkAccess      *PublicNetworkAccess      `json:"publicNetworkAccess,omitempty"`
}
