package standbyvirtualmachinepools

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StandbyVirtualMachinePoolResourceUpdateProperties struct {
	AttachedVirtualMachineScaleSetId *string                                     `json:"attachedVirtualMachineScaleSetId,omitempty"`
	ElasticityProfile                *StandbyVirtualMachinePoolElasticityProfile `json:"elasticityProfile,omitempty"`
	VirtualMachineState              *VirtualMachineState                        `json:"virtualMachineState,omitempty"`
}
