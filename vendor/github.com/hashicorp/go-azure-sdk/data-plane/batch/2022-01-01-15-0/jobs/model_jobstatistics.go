package jobs

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobStatistics struct {
	KernelCPUTime     string  `json:"kernelCPUTime"`
	LastUpdateTime    string  `json:"lastUpdateTime"`
	NumFailedTasks    int64   `json:"numFailedTasks"`
	NumSucceededTasks int64   `json:"numSucceededTasks"`
	NumTaskRetries    int64   `json:"numTaskRetries"`
	ReadIOGiB         float64 `json:"readIOGiB"`
	ReadIOps          int64   `json:"readIOps"`
	StartTime         string  `json:"startTime"`
	Url               string  `json:"url"`
	UserCPUTime       string  `json:"userCPUTime"`
	WaitTime          string  `json:"waitTime"`
	WallClockTime     string  `json:"wallClockTime"`
	WriteIOGiB        float64 `json:"writeIOGiB"`
	WriteIOps         int64   `json:"writeIOps"`
}

func (o *JobStatistics) GetLastUpdateTimeAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.LastUpdateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *JobStatistics) SetLastUpdateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdateTime = formatted
}

func (o *JobStatistics) GetStartTimeAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *JobStatistics) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = formatted
}
