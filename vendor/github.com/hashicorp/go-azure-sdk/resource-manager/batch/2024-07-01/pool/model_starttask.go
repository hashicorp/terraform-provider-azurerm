package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StartTask struct {
	CommandLine         *string                `json:"commandLine,omitempty"`
	ContainerSettings   *TaskContainerSettings `json:"containerSettings,omitempty"`
	EnvironmentSettings *[]EnvironmentSetting  `json:"environmentSettings,omitempty"`
	MaxTaskRetryCount   *int64                 `json:"maxTaskRetryCount,omitempty"`
	ResourceFiles       *[]ResourceFile        `json:"resourceFiles,omitempty"`
	UserIdentity        *UserIdentity          `json:"userIdentity,omitempty"`
	WaitForSuccess      *bool                  `json:"waitForSuccess,omitempty"`
}
