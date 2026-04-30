package firewallrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerFirewallRuleProperties struct {
	EndIPAddress   *string `json:"endIpAddress,omitempty"`
	StartIPAddress *string `json:"startIpAddress,omitempty"`
}
