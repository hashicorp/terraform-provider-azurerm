package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedRuleOverride struct {
	Action *ActionType              `json:"action,omitempty"`
	RuleId string                   `json:"ruleId"`
	State  *ManagedRuleEnabledState `json:"state,omitempty"`
}
