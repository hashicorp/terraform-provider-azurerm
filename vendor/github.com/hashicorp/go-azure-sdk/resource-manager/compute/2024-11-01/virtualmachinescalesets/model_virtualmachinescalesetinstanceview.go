package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineScaleSetInstanceView struct {
	Extensions            *[]VirtualMachineScaleSetVMExtensionsSummary       `json:"extensions,omitempty"`
	OrchestrationServices *[]OrchestrationServiceSummary                     `json:"orchestrationServices,omitempty"`
	Statuses              *[]InstanceViewStatus                              `json:"statuses,omitempty"`
	VirtualMachine        *VirtualMachineScaleSetInstanceViewStatusesSummary `json:"virtualMachine,omitempty"`
}
