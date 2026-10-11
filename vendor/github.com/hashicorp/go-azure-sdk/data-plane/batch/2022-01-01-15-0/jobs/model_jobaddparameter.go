package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobAddParameter struct {
	AllowTaskPreemption       *bool                    `json:"allowTaskPreemption,omitempty"`
	CommonEnvironmentSettings *[]EnvironmentSetting    `json:"commonEnvironmentSettings,omitempty"`
	Constraints               *JobConstraints          `json:"constraints,omitempty"`
	DisplayName               *string                  `json:"displayName,omitempty"`
	Id                        string                   `json:"id"`
	JobManagerTask            *JobManagerTask          `json:"jobManagerTask,omitempty"`
	JobPreparationTask        *JobPreparationTask      `json:"jobPreparationTask,omitempty"`
	JobReleaseTask            *JobReleaseTask          `json:"jobReleaseTask,omitempty"`
	MaxParallelTasks          *int64                   `json:"maxParallelTasks,omitempty"`
	Metadata                  *[]MetadataItem          `json:"metadata,omitempty"`
	NetworkConfiguration      *JobNetworkConfiguration `json:"networkConfiguration,omitempty"`
	OnAllTasksComplete        *OnAllTasksComplete      `json:"onAllTasksComplete,omitempty"`
	OnTaskFailure             *OnTaskFailure           `json:"onTaskFailure,omitempty"`
	PoolInfo                  PoolInformation          `json:"poolInfo"`
	Priority                  *int64                   `json:"priority,omitempty"`
	UsesTaskDependencies      *bool                    `json:"usesTaskDependencies,omitempty"`
}
