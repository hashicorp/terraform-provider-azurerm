package virtualnetworkrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkRuleProperties struct {
	IgnoreMissingVnetServiceEndpoint *bool                    `json:"ignoreMissingVnetServiceEndpoint,omitempty"`
	State                            *VirtualNetworkRuleState `json:"state,omitempty"`
	VirtualNetworkSubnetId           string                   `json:"virtualNetworkSubnetId"`
}
