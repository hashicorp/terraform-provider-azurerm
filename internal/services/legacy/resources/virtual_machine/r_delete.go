// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package virtual_machine

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceVirtualMachineDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachinesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachines.ParseVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualMachineName, compute.VirtualMachineResourceName)
	defer locks.UnlockByName(id.VirtualMachineName, compute.VirtualMachineResourceName)

	virtualMachine, err := client.Get(ctx, *id, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(virtualMachine.HttpResponse) {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(virtualMachine.HttpResponse) {
		// @tombuildsstuff: sending `nil` here omits this value from being sent - which matches
		// the previous behaviour - we're only splitting this out so it's clear why
		opts := virtualmachines.DefaultDeleteOperationOptions()
		opts.ForceDeletion = nil
		if err := client.DeleteThenPoll(ctx, *id, opts); err != nil {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	// delete OS Disk if opted in
	deleteOsDisk := d.Get("delete_os_disk_on_termination").(bool)
	deleteDataDisks := d.Get("delete_data_disks_on_termination").(bool)

	if deleteOsDisk || deleteDataDisks {
		storageClient := meta.(*clients.Client).Storage

		model := virtualMachine.Model
		if model == nil {
			log.Printf("[DEBUG] deleting Disks for %s - `model` was nil, skipping", id)
			return nil
		}
		props := model.Properties
		if props == nil {
			return fmt.Errorf("deleting Disks for %s - `props` was nil", id)
		}
		storageProfile := props.StorageProfile
		if storageProfile == nil {
			return fmt.Errorf("deleting Disks for %s - `storageProfile` was nil", id)
		}

		if deleteOsDisk {
			osDisk := storageProfile.OsDisk
			if osDisk == nil {
				return fmt.Errorf("deleting OS Disk for %s - `osDisk` was nil", id)
			}
			if osDisk.Vhd == nil && osDisk.ManagedDisk == nil {
				return fmt.Errorf("unable to determine OS Disk Type to Delete it for %s", id)
			}

			if osDisk.Vhd != nil {
				if err = resourceVirtualMachineDeleteVhd(ctx, storageClient, id.SubscriptionId, osDisk.Vhd); err != nil {
					return fmt.Errorf("deleting OS Disk VHD: %+v", err)
				}
			} else if osDisk.ManagedDisk != nil {
				if err = resourceVirtualMachineDeleteManagedDisk(d, osDisk.ManagedDisk, meta); err != nil {
					return fmt.Errorf("deleting OS Managed Disk: %+v", err)
				}
			}
		}

		// delete Data disks if opted in
		if deleteDataDisks {
			dataDisks := storageProfile.DataDisks
			if dataDisks == nil {
				return fmt.Errorf("deleting Data Disks for %s: `dataDisks` was nil", id)
			}

			for _, disk := range *dataDisks {
				if disk.Vhd == nil && disk.ManagedDisk == nil {
					return fmt.Errorf("unable to determine Data Disk Type to Delete it for %s / Disk %q", id, *disk.Name)
				}

				if disk.Vhd != nil {
					if err = resourceVirtualMachineDeleteVhd(ctx, storageClient, id.SubscriptionId, disk.Vhd); err != nil {
						return fmt.Errorf("deleting Data Disk VHD: %+v", err)
					}
				} else if disk.ManagedDisk != nil {
					if err = resourceVirtualMachineDeleteManagedDisk(d, disk.ManagedDisk, meta); err != nil {
						return fmt.Errorf("deleting Data Managed Disk: %+v", err)
					}
				}
			}
		}
	}

	return nil
}
