package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobPreparationTask struct {
	CommandLine                   string                 `json:"commandLine"`
	Constraints                   *TaskConstraints       `json:"constraints,omitempty"`
	ContainerSettings             *TaskContainerSettings `json:"containerSettings,omitempty"`
	EnvironmentSettings           *[]EnvironmentSetting  `json:"environmentSettings,omitempty"`
	Id                            *string                `json:"id,omitempty"`
	RerunOnNodeRebootAfterSuccess *bool                  `json:"rerunOnNodeRebootAfterSuccess,omitempty"`
	ResourceFiles                 *[]ResourceFile        `json:"resourceFiles,omitempty"`
	UserIdentity                  *UserIdentity          `json:"userIdentity,omitempty"`
	WaitForSuccess                *bool                  `json:"waitForSuccess,omitempty"`
}
