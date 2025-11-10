// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-07-01/skus"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
)

// The logic on this file is based on:
// Linux: https://learn.microsoft.com/en-us/azure/virtual-machines/linux/expand-disks?tabs=azure-cli%2Cubuntu#expand-without-downtime
// Windows: https://learn.microsoft.com/en-us/azure/virtual-machines/windows/expand-os-disk#expand-without-downtime
// NOTE: whilst the Windows URI says "expand OS disk" it's not supported on OS disks, this is an old document that's not been renamed

// @tombuildsstuff: this is intentionally split out into it's own file since this'll need to be reused

func determineIfDataDiskSupportsNoDowntimeResize(disk *disks.Disk, requiresDetaching bool) bool {
	if disk == nil || disk.Properties == nil {
		return false
	}

	// Only supported for data disks.
	isDataDisk := disk.Properties.OsType != nil && string(*disk.Properties.OsType) != ""
	if isDataDisk {
		return false
	}

	// Not supported for shared disks, the maxShares of which is greater than 1: https://learn.microsoft.com/azure/virtual-machines/disks-shared-enable?tabs=azure-portal#deploy-an-ultra-disk-as-a-shared-disk
	isSharedDisk := disk.Properties.MaxShares != nil && *disk.Properties.MaxShares > 1
	if isSharedDisk {
		log.Printf("[DEBUG] Disk is shared so does not support no-downtime-resize")
		return false
	}

	if requiresDetaching {
		return false
	}

	return true
}

func determineIfDataDiskRequiresDetaching(disk *disks.Disk, oldSizeGb, newSizeGb int) bool {
	if disk == nil || disk.Sku == nil {
		return false
	}

	// This limitation doesn't apply to Premium SSD v2 or Ultra Disks:
	// If a disk is 4 TiB or less, you can't expand it beyond 4 TiB without deallocating the VM.
	// If a disk is already greater than 4 TiB, you can expand it without deallocating the VM.
	diskTypeIsSupported := strings.EqualFold(string(*disk.Sku.Name), string(disks.DiskStorageAccountTypesPremiumVTwoLRS)) || strings.EqualFold(string(*disk.Sku.Name), string(disks.DiskStorageAccountTypesUltraSSDLRS))
	if !diskTypeIsSupported && oldSizeGb < 4096 && newSizeGb >= 4096 {
		return true
	}

	return false
}

func determineIfVirtualMachineSupportsNoDowntimeResize(ctx context.Context, disk *disks.Disk, virtualMachinesClient *virtualmachines.VirtualMachinesClient, skusClient *skus.SkusClient) (*bool, error) {
	if disk == nil || disk.ManagedBy == nil || disk.Sku == nil {
		return pointer.To(false), nil
	}

	virtualMachineIdRaw := disk.ManagedBy
	virtualMachineId, err := virtualmachines.ParseVirtualMachineIDInsensitively(*virtualMachineIdRaw)
	if err != nil {
		log.Printf("[DEBUG] unable to parse Virtual Machine ID %q that the Managed Disk is attached too - skipping no-downtime-resize since we can't guarantee that's available", *virtualMachineIdRaw)
		//nolint:nilerr // this is not an error as we just want to skip the check in this situation since we can't guarantee it's available
		return pointer.To(false), nil
	}

	log.Printf("[DEBUG] Retrieving %s..", *virtualMachineId)
	virtualMachine, err := virtualMachinesClient.Get(ctx, *virtualMachineId, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *virtualMachineId, err)
	}

	vmLocation := ""
	vmSku := ""
	vmDiskControllerType := ""
	if model := virtualMachine.Model; model != nil {
		vmLocation = location.Normalize(model.Location)
		if props := model.Properties; props != nil {
			if props.HardwareProfile != nil && props.HardwareProfile.VMSize != nil {
				vmSku = string(*props.HardwareProfile.VMSize)
			}

			if props.StorageProfile != nil && props.StorageProfile.DiskControllerType != nil {
				vmDiskControllerType = string(*props.StorageProfile.DiskControllerType)
			}
		}
	}

	isUltraOrPremiumV2Disk := strings.EqualFold(string(*disk.Sku.Name), string(disks.DiskStorageAccountTypesPremiumVTwoLRS)) || strings.EqualFold(string(*disk.Sku.Name), string(disks.DiskStorageAccountTypesUltraSSDLRS))
	if isUltraOrPremiumV2Disk {
		// cannot expand a VM that's using NVMe controllers for Ultra or Premium SSD v2 disks without downtime
		return pointer.To(!strings.EqualFold(vmDiskControllerType, string(virtualmachines.DiskControllerTypesNVMe))), nil
	}

	// The following limitation doesn't apply to Premium SSD v2 or Ultra Disks
	if vmLocation == "" || vmSku == "" {
		return pointer.To(false), nil
	}

	subscriptionId := commonids.NewSubscriptionID(virtualMachineId.SubscriptionId)
	opts := skus.DefaultResourceSkusListOperationOptions()
	// @tombuildsstuff: by default this API returns EVERY SKU in EVERY LOCATION meaning this will get
	// progressively larger each time - instead we filter to the current Location only.
	opts.Filter = pointer.To(fmt.Sprintf("location eq '%s'", vmLocation))
	skusResponse, err := skusClient.ResourceSkusListComplete(ctx, subscriptionId, opts)
	if err != nil {
		return nil, fmt.Errorf("retrieving information about the Resource SKUs to check if the Virtual Machine/Disk combination supports no-downtime-resizing: %+v", err)
	}

	for _, sku := range skusResponse.Items {
		if !strings.EqualFold(pointer.From(sku.ResourceType), "virtualMachines") {
			continue
		}

		if sku.Capabilities == nil {
			continue
		}

		if !strings.EqualFold(pointer.From(sku.Name), vmSku) {
			continue
		}

		for _, capability := range *sku.Capabilities {
			if capability.Name == nil || capability.Value == nil {
				continue
			}

			supportsEphemeralOSDisks := false
			supportsHyperVGen2 := false
			supportsPremiumIO := false

			// this logic is based on:
			// if (($capability.Name -eq "EphemeralOSDiskSupported" -and $capability.Value -eq "True") -or ($capability.Name -eq "PremiumIO" -and $capability.Value -eq "True") -or ($capability.Name -eq "HyperVGenerations" -and $capability.Value -match "V2"))
			if strings.EqualFold(*capability.Name, "EphemeralOSDiskSupported") && strings.EqualFold(*capability.Value, "True") {
				supportsEphemeralOSDisks = true
			}
			if strings.EqualFold(*capability.Name, "HyperVGenerations") && strings.Contains(strings.ToLower(*capability.Value), "v2") {
				supportsHyperVGen2 = true
			}
			if strings.EqualFold(*capability.Name, "PremiumIO") && strings.EqualFold(*capability.Value, "True") {
				supportsPremiumIO = true
			}

			if supportsEphemeralOSDisks || supportsPremiumIO || supportsHyperVGen2 {
				return pointer.To(true), nil
			}
		}
	}

	return pointer.To(false), nil
}
