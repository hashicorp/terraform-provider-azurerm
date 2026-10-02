package fleets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VMAttributes struct {
	AcceleratorCount          *VMAttributeMinMaxInteger  `json:"acceleratorCount,omitempty"`
	AcceleratorManufacturers  *[]AcceleratorManufacturer `json:"acceleratorManufacturers,omitempty"`
	AcceleratorSupport        *VMAttributeSupport        `json:"acceleratorSupport,omitempty"`
	AcceleratorTypes          *[]AcceleratorType         `json:"acceleratorTypes,omitempty"`
	ArchitectureTypes         *[]ArchitectureType        `json:"architectureTypes,omitempty"`
	BurstableSupport          *VMAttributeSupport        `json:"burstableSupport,omitempty"`
	CpuManufacturers          *[]CPUManufacturer         `json:"cpuManufacturers,omitempty"`
	DataDiskCount             *VMAttributeMinMaxInteger  `json:"dataDiskCount,omitempty"`
	ExcludedVMSizes           *[]string                  `json:"excludedVMSizes,omitempty"`
	LocalStorageDiskTypes     *[]LocalStorageDiskType    `json:"localStorageDiskTypes,omitempty"`
	LocalStorageInGiB         *VMAttributeMinMaxDouble   `json:"localStorageInGiB,omitempty"`
	LocalStorageSupport       *VMAttributeSupport        `json:"localStorageSupport,omitempty"`
	MemoryInGiB               VMAttributeMinMaxDouble    `json:"memoryInGiB"`
	MemoryInGiBPerVCPU        *VMAttributeMinMaxDouble   `json:"memoryInGiBPerVCpu,omitempty"`
	NetworkBandwidthInMbps    *VMAttributeMinMaxDouble   `json:"networkBandwidthInMbps,omitempty"`
	NetworkInterfaceCount     *VMAttributeMinMaxInteger  `json:"networkInterfaceCount,omitempty"`
	RdmaNetworkInterfaceCount *VMAttributeMinMaxInteger  `json:"rdmaNetworkInterfaceCount,omitempty"`
	RdmaSupport               *VMAttributeSupport        `json:"rdmaSupport,omitempty"`
	VCPUCount                 VMAttributeMinMaxInteger   `json:"vCpuCount"`
	VMCategories              *[]VMCategory              `json:"vmCategories,omitempty"`
}
