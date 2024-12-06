package pool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeploymentConfiguration struct {
	VirtualMachineConfiguration *VirtualMachineConfiguration `json:"virtualMachineConfiguration,omitempty"`
}
