package replicationprotecteditems

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CurrentScenarioDetails struct {
	JobId        *string `json:"jobId,omitempty"`
	ScenarioName *string `json:"scenarioName,omitempty"`
	StartTime    *string `json:"startTime,omitempty"`
}

func (o *CurrentScenarioDetails) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *CurrentScenarioDetails) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
