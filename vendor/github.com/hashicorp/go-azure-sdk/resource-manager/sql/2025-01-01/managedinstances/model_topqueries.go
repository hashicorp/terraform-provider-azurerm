package managedinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TopQueries struct {
	AggregationFunction *string                      `json:"aggregationFunction,omitempty"`
	EndTime             *string                      `json:"endTime,omitempty"`
	IntervalType        *QueryTimeGrainType          `json:"intervalType,omitempty"`
	NumberOfQueries     *int64                       `json:"numberOfQueries,omitempty"`
	ObservationMetric   *string                      `json:"observationMetric,omitempty"`
	Queries             *[]QueryStatisticsProperties `json:"queries,omitempty"`
	StartTime           *string                      `json:"startTime,omitempty"`
}
