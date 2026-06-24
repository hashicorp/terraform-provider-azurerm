package jobs

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudJob struct {
	AllowTaskPreemption         *bool                    `json:"allowTaskPreemption,omitempty"`
	CommonEnvironmentSettings   *[]EnvironmentSetting    `json:"commonEnvironmentSettings,omitempty"`
	Constraints                 *JobConstraints          `json:"constraints,omitempty"`
	CreationTime                *string                  `json:"creationTime,omitempty"`
	DisplayName                 *string                  `json:"displayName,omitempty"`
	ETag                        *string                  `json:"eTag,omitempty"`
	ExecutionInfo               *JobExecutionInformation `json:"executionInfo,omitempty"`
	Id                          *string                  `json:"id,omitempty"`
	JobManagerTask              *JobManagerTask          `json:"jobManagerTask,omitempty"`
	JobPreparationTask          *JobPreparationTask      `json:"jobPreparationTask,omitempty"`
	JobReleaseTask              *JobReleaseTask          `json:"jobReleaseTask,omitempty"`
	LastModified                *string                  `json:"lastModified,omitempty"`
	MaxParallelTasks            *int64                   `json:"maxParallelTasks,omitempty"`
	Metadata                    *[]MetadataItem          `json:"metadata,omitempty"`
	NetworkConfiguration        *JobNetworkConfiguration `json:"networkConfiguration,omitempty"`
	OnAllTasksComplete          *OnAllTasksComplete      `json:"onAllTasksComplete,omitempty"`
	OnTaskFailure               *OnTaskFailure           `json:"onTaskFailure,omitempty"`
	PoolInfo                    *PoolInformation         `json:"poolInfo,omitempty"`
	PreviousState               *JobState                `json:"previousState,omitempty"`
	PreviousStateTransitionTime *string                  `json:"previousStateTransitionTime,omitempty"`
	Priority                    *int64                   `json:"priority,omitempty"`
	State                       *JobState                `json:"state,omitempty"`
	StateTransitionTime         *string                  `json:"stateTransitionTime,omitempty"`
	Stats                       *JobStatistics           `json:"stats,omitempty"`
	Url                         *string                  `json:"url,omitempty"`
	UsesTaskDependencies        *bool                    `json:"usesTaskDependencies,omitempty"`
}

func (o *CloudJob) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *CloudJob) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *CloudJob) GetLastModifiedAsTime() (*time.Time, error) {
	if o.LastModified == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModified, "2006-01-02T15:04:05Z07:00")
}

func (o *CloudJob) SetLastModifiedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModified = &formatted
}

func (o *CloudJob) GetPreviousStateTransitionTimeAsTime() (*time.Time, error) {
	if o.PreviousStateTransitionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.PreviousStateTransitionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *CloudJob) SetPreviousStateTransitionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.PreviousStateTransitionTime = &formatted
}

func (o *CloudJob) GetStateTransitionTimeAsTime() (*time.Time, error) {
	if o.StateTransitionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StateTransitionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *CloudJob) SetStateTransitionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StateTransitionTime = &formatted
}
