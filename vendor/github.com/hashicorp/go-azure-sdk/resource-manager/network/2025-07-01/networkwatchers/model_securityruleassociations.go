package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SecurityRuleAssociations struct {
	DefaultSecurityRules        *[]SecurityRule                 `json:"defaultSecurityRules,omitempty"`
	EffectiveSecurityRules      *[]EffectiveNetworkSecurityRule `json:"effectiveSecurityRules,omitempty"`
	NetworkInterfaceAssociation *NetworkInterfaceAssociation    `json:"networkInterfaceAssociation,omitempty"`
	SubnetAssociation           *SubnetAssociation              `json:"subnetAssociation,omitempty"`
}
