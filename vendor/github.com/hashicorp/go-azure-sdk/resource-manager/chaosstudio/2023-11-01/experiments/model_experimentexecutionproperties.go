package experiments

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExperimentExecutionProperties struct {
	StartedAt *string `json:"startedAt,omitempty"`
	Status    *string `json:"status,omitempty"`
	StoppedAt *string `json:"stoppedAt,omitempty"`
}

func (o *ExperimentExecutionProperties) GetStartedAtAsTime() (*time.Time, error) {
	if o.StartedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *ExperimentExecutionProperties) SetStartedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartedAt = &formatted
}

func (o *ExperimentExecutionProperties) GetStoppedAtAsTime() (*time.Time, error) {
	if o.StoppedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StoppedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *ExperimentExecutionProperties) SetStoppedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StoppedAt = &formatted
}
