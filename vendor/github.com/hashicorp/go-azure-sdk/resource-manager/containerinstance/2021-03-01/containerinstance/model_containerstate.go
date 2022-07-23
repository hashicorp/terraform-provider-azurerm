package containerinstance

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerState struct {
	DetailStatus *string `json:"detailStatus,omitempty"`
	ExitCode     *int64  `json:"exitCode,omitempty"`
	FinishTime   *string `json:"finishTime,omitempty"`
	StartTime    *string `json:"startTime,omitempty"`
	State        *string `json:"state,omitempty"`
}

func (o *ContainerState) GetFinishTimeAsTime() (*time.Time, error) {
	if o.FinishTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.FinishTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ContainerState) SetFinishTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.FinishTime = &formatted
}

func (o *ContainerState) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ContainerState) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
