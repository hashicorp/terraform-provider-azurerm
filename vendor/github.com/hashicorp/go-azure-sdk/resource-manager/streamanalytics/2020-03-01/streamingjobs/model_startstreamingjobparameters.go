package streamingjobs

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StartStreamingJobParameters struct {
	OutputStartMode *OutputStartMode `json:"outputStartMode,omitempty"`
	OutputStartTime *string          `json:"outputStartTime,omitempty"`
}

func (o *StartStreamingJobParameters) GetOutputStartTimeAsTime() (*time.Time, error) {
	if o.OutputStartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.OutputStartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *StartStreamingJobParameters) SetOutputStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.OutputStartTime = &formatted
}
