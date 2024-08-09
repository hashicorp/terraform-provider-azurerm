package machineruncommands

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MachineRunCommandProperties struct {
	AsyncExecution            *bool                          `json:"asyncExecution,omitempty"`
	ErrorBlobManagedIdentity  *RunCommandManagedIdentity     `json:"errorBlobManagedIdentity,omitempty"`
	ErrorBlobUri              *string                        `json:"errorBlobUri,omitempty"`
	InstanceView              *MachineRunCommandInstanceView `json:"instanceView,omitempty"`
	OutputBlobManagedIdentity *RunCommandManagedIdentity     `json:"outputBlobManagedIdentity,omitempty"`
	OutputBlobUri             *string                        `json:"outputBlobUri,omitempty"`
	Parameters                *[]RunCommandInputParameter    `json:"parameters,omitempty"`
	ProtectedParameters       *[]RunCommandInputParameter    `json:"protectedParameters,omitempty"`
	ProvisioningState         *string                        `json:"provisioningState,omitempty"`
	RunAsPassword             *string                        `json:"runAsPassword,omitempty"`
	RunAsUser                 *string                        `json:"runAsUser,omitempty"`
	Source                    *MachineRunCommandScriptSource `json:"source,omitempty"`
	TimeoutInSeconds          *int64                         `json:"timeoutInSeconds,omitempty"`
}
