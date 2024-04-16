package managedhsms

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MHSMNetworkRuleSet struct {
	Bypass              *NetworkRuleBypassOptions `json:"bypass,omitempty"`
	DefaultAction       *NetworkRuleAction        `json:"defaultAction,omitempty"`
	IPRules             *[]MHSMIPRule             `json:"ipRules,omitempty"`
	VirtualNetworkRules *[]MHSMVirtualNetworkRule `json:"virtualNetworkRules,omitempty"`
}
