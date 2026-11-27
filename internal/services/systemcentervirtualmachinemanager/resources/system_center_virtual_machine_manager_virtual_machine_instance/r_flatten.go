// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package system_center_virtual_machine_manager_virtual_machine_instance

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualmachineinstances"
)

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceInfrastructureProfile(input *virtualmachineinstances.InfrastructureProfile) []Infrastructure {
	result := make([]Infrastructure, 0)
	if input == nil {
		return result
	}

	return append(result, Infrastructure{
		CheckpointType:                                          pointer.From(input.CheckpointType),
		SystemCenterVirtualMachineManagerCloudId:                pointer.From(input.CloudId),
		SystemCenterVirtualMachineManagerInventoryItemId:        pointer.From(input.InventoryItemId),
		SystemCenterVirtualMachineManagerTemplateId:             pointer.From(input.TemplateId),
		SystemCenterVirtualMachineManagerVirtualMachineServerId: pointer.From(input.VMmServerId),
	})
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceHardwareProfile(input *virtualmachineinstances.HardwareProfile) []Hardware {
	result := make([]Hardware, 0)
	if input == nil {
		return result
	}

	return append(result, Hardware{
		CpuCount:                    pointer.From(input.CpuCount),
		DynamicMemoryMaxInMb:        pointer.From(input.DynamicMemoryMaxMB),
		DynamicMemoryMinInMb:        pointer.From(input.DynamicMemoryMinMB),
		LimitCpuForMigrationEnabled: pointer.From(input.LimitCPUForMigration) == virtualmachineinstances.LimitCPUForMigrationTrue,
		MemoryInMb:                  pointer.From(input.MemoryMB),
	})
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceNetworkInterfaces(input *[]virtualmachineinstances.NetworkInterface) []NetworkInterface {
	result := make([]NetworkInterface, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		result = append(result, NetworkInterface{
			Name:             pointer.From(v.Name),
			VirtualNetworkId: pointer.From(v.VirtualNetworkId),
			Ipv4AddressType:  pointer.FromEnum(v.IPv4AddressType),
			Ipv6AddressType:  pointer.FromEnum(v.IPv6AddressType),
			MacAddressType:   pointer.FromEnum(v.MacAddressType),
		})
	}

	return result
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceOSProfile(input *virtualmachineinstances.OsProfileForVMInstance, adminPassword string) []OperatingSystem {
	result := make([]OperatingSystem, 0)
	if input == nil {
		return result
	}

	return append(result, OperatingSystem{
		ComputerName:  pointer.From(input.ComputerName),
		AdminPassword: adminPassword,
	})
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceStorageDisks(input *[]virtualmachineinstances.VirtualDisk) []StorageDisk {
	result := make([]StorageDisk, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		// service team confirmed that diskSizeGB in the request payload and maxDiskSizeGB in the API response both represent the maximum size of the disk but diskSizeGB in the API response represents the actual disk size
		storageDisk := StorageDisk{
			Bus:            pointer.From(v.Bus),
			BusType:        pointer.From(v.BusType),
			DiskSizeGB:     pointer.From(v.MaxDiskSizeGB),
			Lun:            pointer.From(v.Lun),
			Name:           pointer.From(v.Name),
			TemplateDiskId: pointer.From(v.TemplateDiskId),
			VhdType:        pointer.From(v.VhdType),
		}

		if storageQoSPolicy := v.StorageQoSPolicy; storageQoSPolicy != nil {
			storageDisk.StorageQoSPolicyName = pointer.From(storageQoSPolicy.Name)
		}

		result = append(result, storageDisk)
	}

	return result
}

func flattenSystemCenterVirtualMachineManagerVirtualMachineInstanceAvailabilitySets(input *[]virtualmachineinstances.AvailabilitySetListItem) []string {
	result := make([]string, 0)
	if input == nil {
		return result
	}

	for _, v := range *input {
		result = append(result, pointer.From(v.Id))
	}

	return result
}
