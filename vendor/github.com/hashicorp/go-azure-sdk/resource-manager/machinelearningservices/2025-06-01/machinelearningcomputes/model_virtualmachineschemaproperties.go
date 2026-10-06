package machinelearningcomputes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineSchemaProperties struct {
	Address                   *string                       `json:"address,omitempty"`
	AdministratorAccount      *VirtualMachineSshCredentials `json:"administratorAccount,omitempty"`
	IsNotebookInstanceCompute *bool                         `json:"isNotebookInstanceCompute,omitempty"`
	NotebookServerPort        *int64                        `json:"notebookServerPort,omitempty"`
	SshPort                   *int64                        `json:"sshPort,omitempty"`
	VirtualMachineSize        *string                       `json:"virtualMachineSize,omitempty"`
}
