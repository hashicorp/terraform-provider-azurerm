package firewallrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FirewallRuleProperties struct {
	EndIPAddress      string             `json:"endIpAddress"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	StartIPAddress    string             `json:"startIpAddress"`
}
