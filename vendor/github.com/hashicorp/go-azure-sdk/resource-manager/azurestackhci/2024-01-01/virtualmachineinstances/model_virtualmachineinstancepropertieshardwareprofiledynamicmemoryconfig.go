package virtualmachineinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineInstancePropertiesHardwareProfileDynamicMemoryConfig struct {
	MaximumMemoryMB    *int64 `json:"maximumMemoryMB,omitempty"`
	MinimumMemoryMB    *int64 `json:"minimumMemoryMB,omitempty"`
	TargetMemoryBuffer *int64 `json:"targetMemoryBuffer,omitempty"`
}
