package storageaccounts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkRuleSet struct {
	Bypass              *Bypass               `json:"bypass,omitempty"`
	DefaultAction       DefaultAction         `json:"defaultAction"`
	IPRules             *[]IPRule             `json:"ipRules,omitempty"`
	ResourceAccessRules *[]ResourceAccessRule `json:"resourceAccessRules,omitempty"`
	VirtualNetworkRules *[]VirtualNetworkRule `json:"virtualNetworkRules,omitempty"`
}
