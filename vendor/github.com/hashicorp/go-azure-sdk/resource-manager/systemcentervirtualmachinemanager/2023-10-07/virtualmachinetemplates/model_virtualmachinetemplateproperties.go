package virtualmachinetemplates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineTemplateProperties struct {
	ComputerName         *string               `json:"computerName,omitempty"`
	CpuCount             *int64                `json:"cpuCount,omitempty"`
	Disks                *[]VirtualDisk        `json:"disks,omitempty"`
	DynamicMemoryEnabled *DynamicMemoryEnabled `json:"dynamicMemoryEnabled,omitempty"`
	DynamicMemoryMaxMB   *int64                `json:"dynamicMemoryMaxMB,omitempty"`
	DynamicMemoryMinMB   *int64                `json:"dynamicMemoryMinMB,omitempty"`
	Generation           *int64                `json:"generation,omitempty"`
	InventoryItemId      *string               `json:"inventoryItemId,omitempty"`
	IsCustomizable       *IsCustomizable       `json:"isCustomizable,omitempty"`
	IsHighlyAvailable    *IsHighlyAvailable    `json:"isHighlyAvailable,omitempty"`
	LimitCPUForMigration *LimitCPUForMigration `json:"limitCpuForMigration,omitempty"`
	MemoryMB             *int64                `json:"memoryMB,omitempty"`
	NetworkInterfaces    *[]NetworkInterface   `json:"networkInterfaces,omitempty"`
	OsName               *string               `json:"osName,omitempty"`
	OsType               *OsType               `json:"osType,omitempty"`
	ProvisioningState    *ProvisioningState    `json:"provisioningState,omitempty"`
	Uuid                 *string               `json:"uuid,omitempty"`
	VMmServerId          *string               `json:"vmmServerId,omitempty"`
}
