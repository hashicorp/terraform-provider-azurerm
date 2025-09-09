package prerules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RuleCounterReset struct {
	FirewallName  *string `json:"firewallName,omitempty"`
	Priority      *string `json:"priority,omitempty"`
	RuleListName  *string `json:"ruleListName,omitempty"`
	RuleName      *string `json:"ruleName,omitempty"`
	RuleStackName *string `json:"ruleStackName,omitempty"`
}
