package ipfirewallrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPFirewallRuleInfo struct {
	Id         *string                   `json:"id,omitempty"`
	Name       *string                   `json:"name,omitempty"`
	Properties *IPFirewallRuleProperties `json:"properties,omitempty"`
	Type       *string                   `json:"type,omitempty"`
}
