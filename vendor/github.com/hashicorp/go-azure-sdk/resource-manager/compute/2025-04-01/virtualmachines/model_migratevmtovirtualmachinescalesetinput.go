package virtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MigrateVMToVirtualMachineScaleSetInput struct {
	TargetFaultDomain *int64  `json:"targetFaultDomain,omitempty"`
	TargetVMSize      *string `json:"targetVMSize,omitempty"`
	TargetZone        *string `json:"targetZone,omitempty"`
}
