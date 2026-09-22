package vaults

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkRuleSet struct {
	Bypass              *NetworkRuleBypassOptions `json:"bypass,omitempty"`
	DefaultAction       *NetworkRuleAction        `json:"defaultAction,omitempty"`
	IPRules             *[]IPRule                 `json:"ipRules,omitempty"`
	VirtualNetworkRules *[]VirtualNetworkRule     `json:"virtualNetworkRules,omitempty"`
}
