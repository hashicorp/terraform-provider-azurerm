package managedinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueryStatisticsProperties struct {
	DatabaseName *string                `json:"databaseName,omitempty"`
	EndTime      *string                `json:"endTime,omitempty"`
	Intervals    *[]QueryMetricInterval `json:"intervals,omitempty"`
	QueryId      *string                `json:"queryId,omitempty"`
	StartTime    *string                `json:"startTime,omitempty"`
}
