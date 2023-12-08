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
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-07-01/skus"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
)

// The logic on this file is based on:
// Linux: https://learn.microsoft.com/en-us/azure/virtual-machines/linux/expand-disks?tabs=azure-cli%2Cubuntu#expand-without-downtime
// Windows: https://learn.microsoft.com/en-us/azure/virtual-machines/windows/expand-os-disk#expand-without-downtime
// NOTE: whilst the Windows URI says "expand OS disk" it's not supported on OS disks, this is an old document that's not been renamed

// @tombuildsstuff: this is intentionally split out into it's own file since this'll need to be reused

func determineIfDataDiskSupportsNoDowntimeResize(disk *disks.Disk, oldSizeGb, newSizeGb int) *bool {
	if disk == nil || disk.Properties == nil || disk.Sku == nil {
		return pointer.To(false)
	}

	// Only supported for data disks.
	isDataDisk := disk.Properties.OsType != nil && string(*disk.Properties.OsType) != ""
	if isDataDisk {
		return pointer.To(false)
	}

	// Not supported for shared disks.
	isSharedDisk := disk.Properties.MaxShares != nil && *disk.Properties.MaxShares >= 0
	if isSharedDisk {
		log.Printf("[DEBUG] Disk is shared so does not support no-downtime-resize")
		return pointer.To(false)
	}

	// If a disk is 4 TiB or less, you can't expand it beyond 4 TiB without deallocating the VM.
	// If a disk is already greater than 4 TiB, you can expand it without deallocating the VM.
	if oldSizeGb < 4096 && newSizeGb >= 4096 {
		return pointer.To(false)
	}

	// Not supported for Ultra disks or Premium SSD v2 disks.
	diskTypeIsSupported := false
	if disk.Sku.Name != nil {
		for _, supportedDiskType := range []disks.DiskStorageAccountTypes{
			disks.DiskStorageAccountTypesPremiumLRS,
			disks.DiskStorageAccountTypesPremiumZRS,
			disks.DiskStorageAccountTypesStandardSSDLRS,
			disks.DiskStorageAccountTypesStandardSSDZRS,
		} {
			if strings.EqualFold(string(*disk.Sku.Name), string(supportedDiskType)) {
				diskTypeIsSupported = true
			}
		}
	}
	return pointer.To(diskTypeIsSupported)
}

func determineIfVirtualMachineSkuSupportsNoDowntimeResize(ctx context.Context, virtualMachineIdRaw *string, virtualMachinesClient *virtualmachines.VirtualMachinesClient, skusClient *skus.SkusClient) (*bool, error) {
	if virtualMachineIdRaw == nil {
		return pointer.To(false), nil
	}

	virtualMachineId, err := virtualmachines.ParseVirtualMachineIDInsensitively(*virtualMachineIdRaw)
	if err != nil {
		log.Printf("[DEBUG] unable to parse Virtual Machine ID %q that the Managed Disk is attached too - skipping no-downtime-resize since we can't guarantee that's available", *virtualMachineIdRaw)
		return pointer.To(false), nil
	}

	log.Printf("[DEBUG] Retrieving %s..", *virtualMachineId)
	virtualMachine, err := virtualMachinesClient.Get(ctx, *virtualMachineId, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *virtualMachineId, err)
	}

	vmSku := ""
	if model := virtualMachine.Model; model != nil && model.Properties != nil && model.Properties.HardwareProfile != nil && model.Properties.HardwareProfile.VMSize != nil {
		vmSku = string(*model.Properties.HardwareProfile.VMSize)
	}
	if vmSku == "" {
		return pointer.To(false), nil
	}

	subscriptionId := commonids.NewSubscriptionID(virtualMachineId.SubscriptionId)
	skusResponse, err := skusClient.ResourceSkusListComplete(ctx, subscriptionId, skus.DefaultResourceSkusListOperationOptions())
	if err != nil {
		return nil, fmt.Errorf("retrieving information about the Resource SKUs to check if the Virtual Machine/Disk combination supports no-downtime-resizing: %+v", err)
	}

	supportsEphemeralOSDisks := false
	supportsHyperVGen2 := false
	supportsPremiumIO := false
	for _, sku := range skusResponse.Items {
		if sku.ResourceType == nil || !strings.EqualFold(*sku.ResourceType, "virtualMachines") {
			continue
		}
		if sku.Capabilities == nil {
			continue
		}

		for _, capability := range *sku.Capabilities {
			if capability.Name == nil || capability.Value == nil {
				continue
			}

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
		}
	}
	result := supportsEphemeralOSDisks || supportsPremiumIO || supportsHyperVGen2
	return pointer.To(result), nil
}
