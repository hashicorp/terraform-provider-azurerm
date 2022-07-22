package cognitiveservicesaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkRuleSet struct {
	DefaultAction       *NetworkRuleAction    `json:"defaultAction,omitempty"`
	IpRules             *[]IpRule             `json:"ipRules,omitempty"`
	VirtualNetworkRules *[]VirtualNetworkRule `json:"virtualNetworkRules,omitempty"`
}
