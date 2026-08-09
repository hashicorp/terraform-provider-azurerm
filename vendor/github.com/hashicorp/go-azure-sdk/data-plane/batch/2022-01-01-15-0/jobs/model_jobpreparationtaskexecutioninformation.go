package jobs

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobPreparationTaskExecutionInformation struct {
	ContainerInfo        *TaskContainerExecutionInformation `json:"containerInfo,omitempty"`
	EndTime              *string                            `json:"endTime,omitempty"`
	ExitCode             *int64                             `json:"exitCode,omitempty"`
	FailureInfo          *TaskFailureInformation            `json:"failureInfo,omitempty"`
	LastRetryTime        *string                            `json:"lastRetryTime,omitempty"`
	Result               *TaskExecutionResult               `json:"result,omitempty"`
	RetryCount           int64                              `json:"retryCount"`
	StartTime            string                             `json:"startTime"`
	State                JobPreparationTaskState            `json:"state"`
	TaskRootDirectory    *string                            `json:"taskRootDirectory,omitempty"`
	TaskRootDirectoryURL *string                            `json:"taskRootDirectoryUrl,omitempty"`
}

func (o *JobPreparationTaskExecutionInformation) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *JobPreparationTaskExecutionInformation) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *JobPreparationTaskExecutionInformation) GetLastRetryTimeAsTime() (*time.Time, error) {
	if o.LastRetryTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastRetryTime, "2006-01-02T15:04:05Z07:00")
}

func (o *JobPreparationTaskExecutionInformation) SetLastRetryTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastRetryTime = &formatted
}

func (o *JobPreparationTaskExecutionInformation) GetStartTimeAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *JobPreparationTaskExecutionInformation) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = formatted
}
