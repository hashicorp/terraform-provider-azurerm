package networkinterfaces

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EffectiveNetworkSecurityGroup struct {
	Association            *EffectiveNetworkSecurityGroupAssociation `json:"association,omitempty"`
	EffectiveSecurityRules *[]EffectiveNetworkSecurityRule           `json:"effectiveSecurityRules,omitempty"`
	NetworkSecurityGroup   *SubResource                              `json:"networkSecurityGroup,omitempty"`
	TagMap                 *map[string][]string                      `json:"tagMap,omitempty"`
}
