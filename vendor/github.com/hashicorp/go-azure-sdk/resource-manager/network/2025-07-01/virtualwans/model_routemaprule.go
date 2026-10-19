package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RouteMapRule struct {
	Actions           *[]Action    `json:"actions,omitempty"`
	MatchCriteria     *[]Criterion `json:"matchCriteria,omitempty"`
	Name              *string      `json:"name,omitempty"`
	NextStepIfMatched *NextStep    `json:"nextStepIfMatched,omitempty"`
}
