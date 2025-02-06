package experiments

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExperimentExecutionActionTargetDetailsProperties struct {
	Error               *ExperimentExecutionActionTargetDetailsError `json:"error,omitempty"`
	Status              *string                                      `json:"status,omitempty"`
	Target              *string                                      `json:"target,omitempty"`
	TargetCompletedTime *string                                      `json:"targetCompletedTime,omitempty"`
	TargetFailedTime    *string                                      `json:"targetFailedTime,omitempty"`
}

func (o *ExperimentExecutionActionTargetDetailsProperties) GetTargetCompletedTimeAsTime() (*time.Time, error) {
	if o.TargetCompletedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TargetCompletedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ExperimentExecutionActionTargetDetailsProperties) SetTargetCompletedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TargetCompletedTime = &formatted
}

func (o *ExperimentExecutionActionTargetDetailsProperties) GetTargetFailedTimeAsTime() (*time.Time, error) {
	if o.TargetFailedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TargetFailedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ExperimentExecutionActionTargetDetailsProperties) SetTargetFailedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TargetFailedTime = &formatted
}
