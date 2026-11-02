package jobs

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobManagerTask struct {
	AllowLowPriorityNode         *bool                          `json:"allowLowPriorityNode,omitempty"`
	ApplicationPackageReferences *[]ApplicationPackageReference `json:"applicationPackageReferences,omitempty"`
	AuthenticationTokenSettings  *AuthenticationTokenSettings   `json:"authenticationTokenSettings,omitempty"`
	CommandLine                  string                         `json:"commandLine"`
	Constraints                  *TaskConstraints               `json:"constraints,omitempty"`
	ContainerSettings            *TaskContainerSettings         `json:"containerSettings,omitempty"`
	DisplayName                  *string                        `json:"displayName,omitempty"`
	EnvironmentSettings          *[]EnvironmentSetting          `json:"environmentSettings,omitempty"`
	Id                           string                         `json:"id"`
	KillJobOnCompletion          *bool                          `json:"killJobOnCompletion,omitempty"`
	OutputFiles                  *[]OutputFile                  `json:"outputFiles,omitempty"`
	RequiredSlots                *int64                         `json:"requiredSlots,omitempty"`
	ResourceFiles                *[]ResourceFile                `json:"resourceFiles,omitempty"`
	RunExclusive                 *bool                          `json:"runExclusive,omitempty"`
	UserIdentity                 *UserIdentity                  `json:"userIdentity,omitempty"`
}
