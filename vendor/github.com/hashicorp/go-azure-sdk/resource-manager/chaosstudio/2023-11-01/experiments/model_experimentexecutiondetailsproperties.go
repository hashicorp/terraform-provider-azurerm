package experiments

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExperimentExecutionDetailsProperties struct {
	FailureReason  *string                                             `json:"failureReason,omitempty"`
	LastActionAt   *string                                             `json:"lastActionAt,omitempty"`
	RunInformation *ExperimentExecutionDetailsPropertiesRunInformation `json:"runInformation,omitempty"`
	StartedAt      *string                                             `json:"startedAt,omitempty"`
	Status         *string                                             `json:"status,omitempty"`
	StoppedAt      *string                                             `json:"stoppedAt,omitempty"`
}

func (o *ExperimentExecutionDetailsProperties) GetLastActionAtAsTime() (*time.Time, error) {
	if o.LastActionAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastActionAt, "2006-01-02T15:04:05Z07:00")
}

func (o *ExperimentExecutionDetailsProperties) SetLastActionAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastActionAt = &formatted
}

func (o *ExperimentExecutionDetailsProperties) GetStartedAtAsTime() (*time.Time, error) {
	if o.StartedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *ExperimentExecutionDetailsProperties) SetStartedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartedAt = &formatted
}

func (o *ExperimentExecutionDetailsProperties) GetStoppedAtAsTime() (*time.Time, error) {
	if o.StoppedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StoppedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *ExperimentExecutionDetailsProperties) SetStoppedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StoppedAt = &formatted
}
