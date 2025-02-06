package virtualmachineinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HardwareProfile struct {
	CpuCount             *int64                `json:"cpuCount,omitempty"`
	DynamicMemoryEnabled *DynamicMemoryEnabled `json:"dynamicMemoryEnabled,omitempty"`
	DynamicMemoryMaxMB   *int64                `json:"dynamicMemoryMaxMB,omitempty"`
	DynamicMemoryMinMB   *int64                `json:"dynamicMemoryMinMB,omitempty"`
	IsHighlyAvailable    *IsHighlyAvailable    `json:"isHighlyAvailable,omitempty"`
	LimitCPUForMigration *LimitCPUForMigration `json:"limitCpuForMigration,omitempty"`
	MemoryMB             *int64                `json:"memoryMB,omitempty"`
}
