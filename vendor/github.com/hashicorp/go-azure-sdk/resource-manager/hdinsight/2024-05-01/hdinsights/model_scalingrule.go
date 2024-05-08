package hdinsights

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScalingRule struct {
	ActionType      ScaleActionType `json:"actionType"`
	ComparisonRule  ComparisonRule  `json:"comparisonRule"`
	EvaluationCount int64           `json:"evaluationCount"`
	ScalingMetric   string          `json:"scalingMetric"`
}
