package frontdoors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RulesEngineRule struct {
	Action                  RulesEngineAction            `json:"action"`
	MatchConditions         *[]RulesEngineMatchCondition `json:"matchConditions,omitempty"`
	MatchProcessingBehavior *MatchProcessingBehavior     `json:"matchProcessingBehavior,omitempty"`
	Name                    string                       `json:"name"`
	Priority                int64                        `json:"priority"`
}
