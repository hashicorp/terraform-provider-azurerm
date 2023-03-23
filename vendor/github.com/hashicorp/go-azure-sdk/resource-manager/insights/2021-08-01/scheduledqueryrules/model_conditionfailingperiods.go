package scheduledqueryrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConditionFailingPeriods struct {
	MinFailingPeriodsToAlert  *int64 `json:"minFailingPeriodsToAlert,omitempty"`
	NumberOfEvaluationPeriods *int64 `json:"numberOfEvaluationPeriods,omitempty"`
}
