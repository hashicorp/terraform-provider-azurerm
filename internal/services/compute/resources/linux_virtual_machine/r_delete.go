// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package linux_virtual_machine

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
	"github.com/hashicorp/go-azure-sdk/sdk/client/pollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/custompoller"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceLinuxVirtualMachineDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Compute.VirtualMachinesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := virtualmachines.ParseVirtualMachineID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualMachineName, helpers.VirtualMachineResourceName)
	defer locks.UnlockByName(id.VirtualMachineName, helpers.VirtualMachineResourceName)

	options := virtualmachines.DefaultGetOperationOptions()
	options.Expand = pointer.To(virtualmachines.InstanceViewTypesUserData)
	existing, err := client.Get(ctx, *id, options)
	if err != nil {
		if response.WasNotFound(existing.HttpResponse) {
			return nil
		}

		return fmt.Errorf("retrieving Linux %s: %+v", id, err)
	}

	// Force Delete is in an opt-in Preview and can only be specified (true/false) if the feature is enabled
	// as such we default this to `nil` which matches the previous behaviour (where this isn't sent) and
	// conditionally set this if required
	var forceDeletion *bool = nil
	if meta.(*clients.Client).Features.VirtualMachine.SkipShutdownAndForceDelete {
		forceDeletion = pointer.To(true)
	}
	deleteOptions := virtualmachines.DefaultDeleteOperationOptions()
	deleteOptions.ForceDeletion = forceDeletion
	if err := client.DeleteThenPoll(ctx, *id, deleteOptions); err != nil {
		return fmt.Errorf("deleting Linux %s: %+v", id, err)
	}

	deleteOSDisk := meta.(*clients.Client).Features.VirtualMachine.DeleteOSDiskOnDeletion
	if deleteOSDisk {
		disksClient := meta.(*clients.Client).Compute.DisksClient
		managedDiskId := ""
		if props := existing.Model.Properties; props != nil && props.StorageProfile != nil && props.StorageProfile.OsDisk != nil {
			if disk := props.StorageProfile.OsDisk.ManagedDisk; disk != nil && disk.Id != nil {
				managedDiskId = *disk.Id
			}
		}

		if managedDiskId != "" {
			diskId, err := commonids.ParseManagedDiskID(managedDiskId)
			if err != nil {
				return err
			}

			if err := disksClient.DeleteThenPoll(ctx, *diskId); err != nil {
				return fmt.Errorf("deleting %s for Linux %s: %+v", diskId, id, err)
			}
		} else {
			log.Printf("[DEBUG] Skipping Deleting OS Disk from Linux %s - cannot determine OS Disk ID.", id)
		}
	} else {
		log.Printf("[DEBUG] Skipping Deleting OS Disk from Linux %s", id)
		// Wait for the Disk to be no longer `Attached` so we know it's released and free to be attached to any subsequent VM etc
		disksClient := meta.(*clients.Client).Compute.DisksClient
		managedDiskId := ""
		if props := existing.Model.Properties; props != nil && props.StorageProfile != nil && props.StorageProfile.OsDisk != nil {
			if disk := props.StorageProfile.OsDisk.ManagedDisk; disk != nil && disk.Id != nil {
				managedDiskId = *disk.Id
			}
		}
		if managedDiskId != "" {
			diskId, err := commonids.ParseManagedDiskIDInsensitively(managedDiskId)
			if err != nil {
				return err
			}

			pollerType := custompoller.NewManagedDiskDetachedPoller(disksClient, *diskId, *id)
			poller := pollers.NewPoller(pollerType, 10*time.Second, pollers.DefaultNumberOfDroppedConnectionsToAllow)
			if err := poller.PollUntilDone(ctx); err != nil {
				return err
			}
		}
	}

	// Need to add a get and a state wait to avoid bug in network API where the attached disk(s) are not actually deleted
	// Service team indicated that we need to do a get after VM delete call returns to verify that the VM and all attached
	// disks have actually been deleted.

	log.Printf("[INFO] verifying Linux %s has been deleted", id)
	virtualMachine, err := client.Get(ctx, *id, virtualmachines.DefaultGetOperationOptions())
	if err != nil && !response.WasNotFound(virtualMachine.HttpResponse) {
		return fmt.Errorf("verifying Linux %s has been deleted: %+v", id, err)
	}

	if !response.WasNotFound(virtualMachine.HttpResponse) {
		log.Printf("[INFO] %s still exists, waiting on vm to be deleted", id)

		poller := custompollers.NewEventualConsistencyPoller(1, func(pollerCtx context.Context) (*http.Response, error) {
			log.Printf("[INFO] checking on state of Linux %s", id)
			resp, err := client.Get(pollerCtx, *id, virtualmachines.DefaultGetOperationOptions())
			return resp.HttpResponse, err
		}, &custompollers.EventualConsistencyPollerOptions{
			Interval:         30 * time.Second,
			TargetStatusCode: pointer.To(http.StatusNotFound),
		})
		if err := poller.PollUntilDone(ctx); err != nil {
			return fmt.Errorf("waiting for the deletion of Linux %s: %v", id, err)
		}
	}

	return nil
}
