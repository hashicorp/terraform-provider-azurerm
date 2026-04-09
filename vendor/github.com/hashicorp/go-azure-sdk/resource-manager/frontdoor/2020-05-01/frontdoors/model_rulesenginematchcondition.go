package frontdoors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RulesEngineMatchCondition struct {
	NegateCondition          *bool                    `json:"negateCondition,omitempty"`
	RulesEngineMatchValue    []string                 `json:"rulesEngineMatchValue"`
	RulesEngineMatchVariable RulesEngineMatchVariable `json:"rulesEngineMatchVariable"`
	RulesEngineOperator      RulesEngineOperator      `json:"rulesEngineOperator"`
	Selector                 *string                  `json:"selector,omitempty"`
	Transforms               *[]Transform             `json:"transforms,omitempty"`
}
