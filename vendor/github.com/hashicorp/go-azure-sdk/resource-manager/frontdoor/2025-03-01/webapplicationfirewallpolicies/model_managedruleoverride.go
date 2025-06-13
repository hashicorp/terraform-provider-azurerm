package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedRuleOverride struct {
	Action       *ActionType              `json:"action,omitempty"`
	EnabledState *ManagedRuleEnabledState `json:"enabledState,omitempty"`
	Exclusions   *[]ManagedRuleExclusion  `json:"exclusions,omitempty"`
	RuleId       string                   `json:"ruleId"`
}
