package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MatchCondition struct {
	MatchValue      []string         `json:"matchValue"`
	MatchVariable   MatchVariable    `json:"matchVariable"`
	NegateCondition *bool            `json:"negateCondition,omitempty"`
	Operator        Operator         `json:"operator"`
	Selector        *string          `json:"selector,omitempty"`
	Transforms      *[]TransformType `json:"transforms,omitempty"`
}
