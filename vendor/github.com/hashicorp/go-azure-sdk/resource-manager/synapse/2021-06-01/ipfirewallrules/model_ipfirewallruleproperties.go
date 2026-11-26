package ipfirewallrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IPFirewallRuleProperties struct {
	EndIPAddress      *string            `json:"endIpAddress,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	StartIPAddress    *string            `json:"startIpAddress,omitempty"`
}
