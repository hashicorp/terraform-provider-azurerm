package ipfirewallrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplaceAllIPFirewallRulesRequest struct {
	IPFirewallRules *map[string]IPFirewallRuleProperties `json:"ipFirewallRules,omitempty"`
}
