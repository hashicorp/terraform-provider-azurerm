package testjobstream

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobStreamProperties struct {
	JobStreamId *string                 `json:"jobStreamId,omitempty"`
	StreamText  *string                 `json:"streamText,omitempty"`
	StreamType  *JobStreamType          `json:"streamType,omitempty"`
	Summary     *string                 `json:"summary,omitempty"`
	Time        *string                 `json:"time,omitempty"`
	Value       *map[string]interface{} `json:"value,omitempty"`
}

func (o *JobStreamProperties) GetTimeAsTime() (*time.Time, error) {
	if o.Time == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Time, "2006-01-02T15:04:05Z07:00")
}

func (o *JobStreamProperties) SetTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Time = &formatted
}
