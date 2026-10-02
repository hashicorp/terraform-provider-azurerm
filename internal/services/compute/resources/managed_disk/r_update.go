// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package managed_disk

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceManagedDiskUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.DisksClient
	virtualMachinesClient := meta.(*clients.Client).Compute.VirtualMachinesClient
	skusClient := meta.(*clients.Client).Compute.SkusClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	maxShares := d.Get("max_shares").(int)
	storageAccountType := d.Get("storage_account_type").(string)
	diskSizeGB := d.Get("disk_size_gb").(int)
	onDemandBurstingEnabled := d.Get("on_demand_bursting_enabled").(bool)
	shouldShutDown := false
	shouldDetach := false

	id, err := commonids.ParseManagedDiskID(d.Id())
	if err != nil {
		return err
	}

	disk, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(disk.HttpResponse) {
			return fmt.Errorf("managed disk %q (Resource Group %q) was not found", name, resourceGroup)
		}

		return fmt.Errorf("making Read request on Azure Managed Disk %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	diskUpdate := disks.DiskUpdate{
		Properties: &disks.DiskUpdateProperties{},
	}

	if d.HasChange("max_shares") {
		diskUpdate.Properties.MaxShares = pointer.To(int64(maxShares))
		var skuName disks.DiskStorageAccountTypes
		for _, v := range disks.PossibleValuesForDiskStorageAccountTypes() {
			if strings.EqualFold(storageAccountType, v) {
				skuName = disks.DiskStorageAccountTypes(v)
			}
		}
		diskUpdate.Sku = &disks.DiskSku{
			Name: &skuName,
		}
	}

	if d.HasChange("tier") {
		if storageAccountType != string(disks.DiskStorageAccountTypesPremiumZRS) && storageAccountType != string(disks.DiskStorageAccountTypesPremiumLRS) {
			return fmt.Errorf("`tier` can only be specified when `storage_account_type` is set to `Premium_LRS` or `Premium_ZRS`")
		}
		shouldShutDown = true
		tier := d.Get("tier").(string)
		diskUpdate.Properties.Tier = &tier
	}

	if d.HasChange("tags") {
		t := d.Get("tags").(map[string]any)
		diskUpdate.Tags = tags.Expand(t)
	}

	if d.HasChange("storage_account_type") {
		shouldShutDown = true
		var skuName disks.DiskStorageAccountTypes
		for _, v := range disks.PossibleValuesForDiskStorageAccountTypes() {
			if strings.EqualFold(storageAccountType, v) {
				skuName = disks.DiskStorageAccountTypes(v)
			}
		}
		diskUpdate.Sku = &disks.DiskSku{
			Name: &skuName,
		}
	}

	if strings.EqualFold(storageAccountType, string(disks.DiskStorageAccountTypesUltraSSDLRS)) || storageAccountType == string(disks.DiskStorageAccountTypesPremiumVTwoLRS) {
		if d.HasChange("disk_iops_read_write") {
			v := d.Get("disk_iops_read_write")
			diskUpdate.Properties.DiskIOPSReadWrite = pointer.To(int64(v.(int)))
		}

		if d.HasChange("disk_mbps_read_write") {
			v := d.Get("disk_mbps_read_write")
			diskUpdate.Properties.DiskMBpsReadWrite = pointer.To(int64(v.(int)))
		}

		if d.HasChange("disk_iops_read_only") {
			if maxShares == 0 {
				return fmt.Errorf("[ERROR] disk_iops_read_only is only available for UltraSSD disks with shared disk enabled")
			}

			v := d.Get("disk_iops_read_only")
			diskUpdate.Properties.DiskIOPSReadOnly = pointer.To(int64(v.(int)))
		}

		if d.HasChange("disk_mbps_read_only") {
			if maxShares == 0 {
				return fmt.Errorf("[ERROR] disk_mbps_read_only is only available for UltraSSD disks with shared disk enabled")
			}

			v := d.Get("disk_mbps_read_only")
			diskUpdate.Properties.DiskMBpsReadOnly = pointer.To(int64(v.(int)))
		}
	} else if d.HasChanges("disk_iops_read_write", "disk_mbps_read_write", "disk_iops_read_only", "disk_mbps_read_only") {
		return fmt.Errorf("[ERROR] disk_iops_read_write, disk_mbps_read_write, disk_iops_read_only and disk_mbps_read_only are only available for UltraSSD disks and PremiumV2 disks")
	}

	if d.HasChange("optimized_frequent_attach_enabled") {
		diskUpdate.Properties.OptimizedForFrequentAttach = pointer.To(d.Get("optimized_frequent_attach_enabled").(bool))
	}

	if d.HasChange("os_type") {
		operatingSystemType := disks.OperatingSystemTypes(d.Get("os_type").(string))
		diskUpdate.Properties.OsType = &operatingSystemType
	}

	if d.HasChange("disk_size_gb") {
		if oldSize, newSize := d.GetChange("disk_size_gb"); newSize.(int) > oldSize.(int) {
			canBeResizedWithoutDowntime := false
			if meta.(*clients.Client).Features.ManagedDisk.ExpandWithoutDowntime {
				shouldDetach = helpers.DetermineIfDataDiskRequiresDetaching(disk.Model, oldSize.(int), newSize.(int))
				diskSupportsNoDowntimeResize := helpers.DetermineIfDataDiskSupportsNoDowntimeResize(disk.Model, shouldDetach)

				vmSupportsNoDowntimeResize, err := helpers.DetermineIfVirtualMachineSupportsNoDowntimeResize(ctx, disk.Model, virtualMachinesClient, skusClient)
				if err != nil {
					return fmt.Errorf("determining if the Virtual Machine the Disk is attached to supports no-downtime-resize: %+v", err)
				}

				canBeResizedWithoutDowntime = *vmSupportsNoDowntimeResize && diskSupportsNoDowntimeResize
			}
			if !canBeResizedWithoutDowntime {
				log.Printf("[INFO] The %s, or the Virtual Machine that it's attached to, doesn't support no-downtime-resizing - requiring that the VM should be shutdown", *id)
				shouldShutDown = true
			}
			diskUpdate.Properties.DiskSizeGB = pointer.To(int64(newSize.(int)))
		} else {
			return fmt.Errorf("- New size must be greater than original size. Shrinking disks is not supported on Azure")
		}
	}

	if d.HasChange("encryption_settings") {
		diskUpdate.Properties.EncryptionSettingsCollection = helpers.ExpandManagedDiskEncryptionSettings(d.Get("encryption_settings").([]any))
	}

	if d.HasChange("disk_encryption_set_id") {
		shouldShutDown = true
		if diskEncryptionSetId := d.Get("disk_encryption_set_id").(string); diskEncryptionSetId != "" {
			encryptionType, err := helpers.RetrieveDiskEncryptionSetEncryptionType(ctx, meta.(*clients.Client).Compute.DiskEncryptionSetsClient, diskEncryptionSetId)
			if err != nil {
				return err
			}

			diskUpdate.Properties.Encryption = &disks.Encryption{
				Type:                encryptionType,
				DiskEncryptionSetId: pointer.To(diskEncryptionSetId),
			}
		} else {
			return fmt.Errorf("once a customer-managed key is used, you can’t change the selection back to a platform-managed key")
		}
	}

	if d.HasChange("network_access_policy") {
		diskUpdate.Properties.NetworkAccessPolicy = pointer.ToEnum[disks.NetworkAccessPolicy](d.Get("network_access_policy").(string))
	}

	if diskAccessID := d.Get("disk_access_id").(string); d.HasChange("disk_access_id") {
		networkAccessPolicy := disks.NetworkAccessPolicy(d.Get("network_access_policy").(string))
		switch {
		case networkAccessPolicy == disks.NetworkAccessPolicyAllowPrivate:
			diskUpdate.Properties.DiskAccessId = pointer.To(diskAccessID)
		case diskAccessID != "" && networkAccessPolicy != disks.NetworkAccessPolicyAllowPrivate:
			return fmt.Errorf("[ERROR] disk_access_id is only available when network_access_policy is set to AllowPrivate")
		default:
			diskUpdate.Properties.DiskAccessId = nil
		}
	}

	if d.HasChange("public_network_access_enabled") {
		if d.Get("public_network_access_enabled").(bool) {
			diskUpdate.Properties.PublicNetworkAccess = pointer.To(disks.PublicNetworkAccessEnabled)
		} else {
			diskUpdate.Properties.PublicNetworkAccess = pointer.To(disks.PublicNetworkAccessDisabled)
		}
	}

	if onDemandBurstingEnabled {
		switch storageAccountType {
		case string(disks.DiskStorageAccountTypesPremiumLRS):
		case string(disks.DiskStorageAccountTypesPremiumZRS):
		default:
			return fmt.Errorf("`on_demand_bursting_enabled` can only be set to true when `storage_account_type` is set to `Premium_LRS` or `Premium_ZRS`")
		}

		if diskSizeGB != 0 && diskSizeGB <= 512 {
			return fmt.Errorf("`on_demand_bursting_enabled` can only be set to true when `disk_size_gb` is larger than 512GB")
		}
	}

	if d.HasChange("on_demand_bursting_enabled") {
		shouldShutDown = true
		diskUpdate.Properties.BurstingEnabled = pointer.To(onDemandBurstingEnabled)
	}

	// whilst we need to shut this down, if we're not attached to anything there's no point
	if shouldShutDown && disk.Model.ManagedBy == nil {
		shouldShutDown = false
	}

	// if we are attached to a VM we bring down the VM as necessary for the operations which are not allowed while it's online
	if shouldShutDown {
		virtualMachineId, err := virtualmachines.ParseVirtualMachineID(*disk.Model.ManagedBy)
		if err != nil {
			return fmt.Errorf("parsing VMID %q for disk attachment: %+v", *disk.Model.ManagedBy, err)
		}
		// check instanceView State

		locks.ByName(virtualMachineId.VirtualMachineName, helpers.VirtualMachineResourceName)
		defer locks.UnlockByName(virtualMachineId.VirtualMachineName, helpers.VirtualMachineResourceName)

		if err = helpers.ResourceManagedDiskUpdateWithVmShutDown(ctx, meta.(*clients.Client), id, virtualMachineId, diskUpdate, shouldDetach); err != nil {
			return err
		}
	} else { // otherwise, just update it
		if err := client.UpdateThenPoll(ctx, *id, diskUpdate); err != nil {
			return fmt.Errorf("expanding managed disk %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return resourceManagedDiskRead(d, meta)
}
