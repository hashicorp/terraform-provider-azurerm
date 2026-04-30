package managedinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type QueryMetricInterval struct {
	ExecutionCount    *int64                   `json:"executionCount,omitempty"`
	IntervalStartTime *string                  `json:"intervalStartTime,omitempty"`
	IntervalType      *QueryTimeGrainType      `json:"intervalType,omitempty"`
	Metrics           *[]QueryMetricProperties `json:"metrics,omitempty"`
}
