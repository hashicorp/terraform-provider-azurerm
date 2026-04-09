package scheduledqueryrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogMetricTrigger struct {
	MetricColumn      *string              `json:"metricColumn,omitempty"`
	MetricTriggerType *MetricTriggerType   `json:"metricTriggerType,omitempty"`
	Threshold         *float64             `json:"threshold,omitempty"`
	ThresholdOperator *ConditionalOperator `json:"thresholdOperator,omitempty"`
}
