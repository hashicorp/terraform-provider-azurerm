package webapps

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TriggeredJobRun struct {
	Duration   *string                `json:"duration,omitempty"`
	EndTime    *string                `json:"end_time,omitempty"`
	ErrorUrl   *string                `json:"error_url,omitempty"`
	JobName    *string                `json:"job_name,omitempty"`
	OutputUrl  *string                `json:"output_url,omitempty"`
	StartTime  *string                `json:"start_time,omitempty"`
	Status     *TriggeredWebJobStatus `json:"status,omitempty"`
	Trigger    *string                `json:"trigger,omitempty"`
	Url        *string                `json:"url,omitempty"`
	WebJobId   *string                `json:"web_job_id,omitempty"`
	WebJobName *string                `json:"web_job_name,omitempty"`
}

func (o *TriggeredJobRun) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *TriggeredJobRun) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *TriggeredJobRun) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *TriggeredJobRun) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
