package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobReleaseTask struct {
	CommandLine         string                 `json:"commandLine"`
	ContainerSettings   *TaskContainerSettings `json:"containerSettings,omitempty"`
	EnvironmentSettings *[]EnvironmentSetting  `json:"environmentSettings,omitempty"`
	Id                  *string                `json:"id,omitempty"`
	MaxWallClockTime    *string                `json:"maxWallClockTime,omitempty"`
	ResourceFiles       *[]ResourceFile        `json:"resourceFiles,omitempty"`
	RetentionTime       *string                `json:"retentionTime,omitempty"`
	UserIdentity        *UserIdentity          `json:"userIdentity,omitempty"`
}
