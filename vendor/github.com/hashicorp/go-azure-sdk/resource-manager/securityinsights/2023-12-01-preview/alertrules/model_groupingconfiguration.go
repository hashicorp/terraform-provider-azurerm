package alertrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GroupingConfiguration struct {
	Enabled              bool                 `json:"enabled"`
	GroupByAlertDetails  *[]AlertDetail       `json:"groupByAlertDetails,omitempty"`
	GroupByCustomDetails *[]string            `json:"groupByCustomDetails,omitempty"`
	GroupByEntities      *[]EntityMappingType `json:"groupByEntities,omitempty"`
	LookbackDuration     string               `json:"lookbackDuration"`
	MatchingMethod       MatchingMethod       `json:"matchingMethod"`
	ReopenClosedIncident bool                 `json:"reopenClosedIncident"`
}
