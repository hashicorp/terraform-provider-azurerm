package vpnserverconfigurations

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VpnServerConfigurationPolicyGroupProperties struct {
	IsDefault                   *bool                                      `json:"isDefault,omitempty"`
	P2SConnectionConfigurations *[]SubResource                             `json:"p2SConnectionConfigurations,omitempty"`
	PolicyMembers               *[]VpnServerConfigurationPolicyGroupMember `json:"policyMembers,omitempty"`
	Priority                    *int64                                     `json:"priority,omitempty"`
	ProvisioningState           *ProvisioningState                         `json:"provisioningState,omitempty"`
}
