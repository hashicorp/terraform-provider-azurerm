package autoscalesettings

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MetricTrigger struct {
	Dimensions             *[]ScaleRuleMetricDimension `json:"dimensions,omitempty"`
	DividePerInstance      *bool                       `json:"dividePerInstance,omitempty"`
	MetricName             string                      `json:"metricName"`
	MetricNamespace        *string                     `json:"metricNamespace,omitempty"`
	MetricResourceLocation *string                     `json:"metricResourceLocation,omitempty"`
	MetricResourceUri      string                      `json:"metricResourceUri"`
	Operator               ComparisonOperationType     `json:"operator"`
	Statistic              MetricStatisticType         `json:"statistic"`
	Threshold              float64                     `json:"threshold"`
	TimeAggregation        TimeAggregationType         `json:"timeAggregation"`
	TimeGrain              string                      `json:"timeGrain"`
	TimeWindow             string                      `json:"timeWindow"`
}
