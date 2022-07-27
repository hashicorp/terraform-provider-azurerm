package scheduledqueryrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Condition struct {
	Dimensions          *[]Dimension             `json:"dimensions,omitempty"`
	FailingPeriods      *ConditionFailingPeriods `json:"failingPeriods,omitempty"`
	MetricMeasureColumn *string                  `json:"metricMeasureColumn,omitempty"`
	MetricName          *string                  `json:"metricName,omitempty"`
	Operator            *ConditionOperator       `json:"operator,omitempty"`
	Query               *string                  `json:"query,omitempty"`
	ResourceIdColumn    *string                  `json:"resourceIdColumn,omitempty"`
	Threshold           *float64                 `json:"threshold,omitempty"`
	TimeAggregation     *TimeAggregation         `json:"timeAggregation,omitempty"`
}
