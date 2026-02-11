package networkwatchers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EvaluatedNetworkSecurityGroup struct {
	AppliedTo              *string                                 `json:"appliedTo,omitempty"`
	MatchedRule            *MatchedRule                            `json:"matchedRule,omitempty"`
	NetworkSecurityGroupId *string                                 `json:"networkSecurityGroupId,omitempty"`
	RulesEvaluationResult  *[]NetworkSecurityRulesEvaluationResult `json:"rulesEvaluationResult,omitempty"`
}
