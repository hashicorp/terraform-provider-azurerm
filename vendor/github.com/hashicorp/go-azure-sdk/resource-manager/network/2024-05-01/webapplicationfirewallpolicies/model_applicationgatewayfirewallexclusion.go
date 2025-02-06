package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayFirewallExclusion struct {
	MatchVariable         string `json:"matchVariable"`
	Selector              string `json:"selector"`
	SelectorMatchOperator string `json:"selectorMatchOperator"`
}
