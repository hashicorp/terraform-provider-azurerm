package updateruns

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateStatus struct {
	CompletedTime *string      `json:"completedTime,omitempty"`
	Error         *ErrorDetail `json:"error,omitempty"`
	StartTime     *string      `json:"startTime,omitempty"`
	State         *UpdateState `json:"state,omitempty"`
}

func (o *UpdateStatus) GetCompletedTimeAsTime() (*time.Time, error) {
	if o.CompletedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CompletedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *UpdateStatus) SetCompletedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CompletedTime = &formatted
}

func (o *UpdateStatus) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *UpdateStatus) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
