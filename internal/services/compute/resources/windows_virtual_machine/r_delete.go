// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package windows_virtual_machine

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceWindowsVirtualMachineDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachinesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachines.ParseVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualMachineName, helpers.VirtualMachineResourceName)
	defer locks.UnlockByName(id.VirtualMachineName, helpers.VirtualMachineResourceName)

	existing, err := client.Get(ctx, *id, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return nil
		}

		return fmt.Errorf("retrieving Windows %s: %+v", id, err)
	}

	// Force Delete is in an opt-in Preview and can only be specified (true/false) if the feature is enabled
	// as such we default this to `nil` which matches the previous behaviour (where this isn't sent) and
	// conditionally set this if required
	options := virtualmachines.DefaultDeleteOperationOptions()
	if meta.(*clients.Client).Features.VirtualMachine.SkipShutdownAndForceDelete {
		options.ForceDeletion = pointer.To(true)
	}
	if err := client.DeleteThenPoll(ctx, *id, options); err != nil {
		return fmt.Errorf("deleting Windows %s: %+v", id, err)
	}

	deleteOSDisk := meta.(*clients.Client).Features.VirtualMachine.DeleteOSDiskOnDeletion
	if deleteOSDisk {
		disksClient := meta.(*clients.Client).Compute.DisksClient
		managedDiskId := ""
		if model := existing.Model; model != nil {
			if props := model.Properties; props.StorageProfile != nil && props.StorageProfile.OsDisk != nil {
				if disk := props.StorageProfile.OsDisk.ManagedDisk; disk != nil && disk.Id != nil {
					managedDiskId = *disk.Id
				}
			}
		}

		if managedDiskId != "" {
			diskId, err := commonids.ParseManagedDiskIDInsensitively(managedDiskId)
			if err != nil {
				return err
			}

			diskDeleteFuture, err := disksClient.Delete(ctx, *diskId)
			if err != nil {
				if !response.WasNotFound(diskDeleteFuture.HttpResponse) {
					return fmt.Errorf("deleting OS Disk %q (Resource Group %q) for Windows %s: %+v", diskId.DiskName, diskId.ResourceGroupName, id, err)
				}
			}
			if !response.WasNotFound(diskDeleteFuture.HttpResponse) {
				if err := diskDeleteFuture.Poller.PollUntilDone(ctx); err != nil {
					return fmt.Errorf("OS Disk %s for Windows %s: %+v", diskId, id, err)
				}
			}
		} else {
			log.Printf("[DEBUG] Skipping Deleting OS Disk from Windows %s - cannot determine OS Disk ID.", id)
		}
	} else {
		log.Printf("[DEBUG] Skipping Deleting OS Disk from Windows %s", id)
	}

	// Need to add a get and a state wait to avoid bug in network API where the attached disk(s) are not actually deleted
	// Service team indicated that we need to do a get after VM delete call returns to verify that the VM and all attached
	// disks have actually been deleted.

	log.Printf("[INFO] verifying Windows %s has been deleted", id)
	virtualMachine, err := client.Get(ctx, *id, virtualmachines.DefaultGetOperationOptions())
	if err != nil && !response.WasNotFound(virtualMachine.HttpResponse) {
		return fmt.Errorf("verifying Windows %s has been deleted: %+v", id, err)
	}

	if !response.WasNotFound(virtualMachine.HttpResponse) {
		log.Printf("[INFO] Windows %s still exists, waiting on vm to be deleted", id)

		poller := custompollers.NewEventualConsistencyPoller(1, func(pollerCtx context.Context) (*http.Response, error) {
			resp, err := client.Get(pollerCtx, *id, virtualmachines.DefaultGetOperationOptions())
			return resp.HttpResponse, err
		}, &custompollers.EventualConsistencyPollerOptions{
			Interval:         30 * time.Second,
			TargetStatusCode: pointer.To(http.StatusNotFound),
		})
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("waiting for the deletion of Windows %s: %v", id, err)
		}
	}

	return nil
}
