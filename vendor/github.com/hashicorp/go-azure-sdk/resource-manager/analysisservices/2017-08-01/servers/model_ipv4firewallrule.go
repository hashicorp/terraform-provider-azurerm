package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPv4FirewallRule struct {
	FirewallRuleName *string `json:"firewallRuleName,omitempty"`
	RangeEnd         *string `json:"rangeEnd,omitempty"`
	RangeStart       *string `json:"rangeStart,omitempty"`
}
