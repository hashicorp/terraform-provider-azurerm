package scheduledqueryrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TriggerCondition struct {
	MetricTrigger     *LogMetricTrigger   `json:"metricTrigger,omitempty"`
	Threshold         float64             `json:"threshold"`
	ThresholdOperator ConditionalOperator `json:"thresholdOperator"`
}
