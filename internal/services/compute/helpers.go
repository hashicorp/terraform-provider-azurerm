// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-07-03/galleryimageversions"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-07-01/virtualmachinescalesets"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
)

func expandIDsToSubResources(input []interface{}) *[]virtualmachinescalesets.SubResource {
	ids := make([]virtualmachinescalesets.SubResource, 0)

	for _, v := range input {
		ids = append(ids, virtualmachinescalesets.SubResource{
			Id: pointer.To(v.(string)),
		})
	}

	return &ids
}

func flattenSubResourcesToIDs(input *[]virtualmachinescalesets.SubResource) []interface{} {
	ids := make([]interface{}, 0)
	if input == nil {
		return ids
	}

	for _, v := range *input {
		if v.Id == nil {
			continue
		}

		ids = append(ids, *v.Id)
	}

	return ids
}

func flattenSubResourcesToStringIDs(input *[]virtualmachinescalesets.SubResource) []string {
	ids := make([]string, 0)
	if input == nil {
		return ids
	}

	for _, v := range *input {
		if v.Id == nil {
			continue
		}

		ids = append(ids, *v.Id)
	}

	return ids
}

func sortSharedImageVersions(values []galleryimageversions.GalleryImageVersion) ([]galleryimageversions.GalleryImageVersion, []error) {
	errors := make([]error, 0)
	sort.Slice(values, func(i, j int) bool {
		if values[i].Name == nil || values[j].Name == nil {
			return false
		}

		verA, err := version.NewVersion(*values[i].Name)
		if err != nil {
			errors = append(errors, err)
			return false
		}
		verA = version.Must(verA, err)

		verB, err := version.NewVersion(*values[j].Name)
		if err != nil {
			errors = append(errors, err)
			return false
		}
		verB = version.Must(verB, err)
		return verA.LessThan(verB)
	})

	if len(errors) > 0 {
		return values, errors
	}
	return values, nil
}

func resourceManagedDiskUpdateWithVmShutDown(ctx context.Context, clients *clients.Client, id *commonids.ManagedDiskId, virtualMachineId *virtualmachines.VirtualMachineId, diskUpdate disks.DiskUpdate, shouldDetach bool) error {
	expandedDisk := virtualmachines.DataDisk{}
	shouldShutDown := true
	diskClient := clients.Compute.DisksClient
	virtualMachinesClient := clients.Compute.VirtualMachinesClient

	vm, err := virtualMachinesClient.Get(ctx, *virtualMachineId, virtualmachines.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", virtualMachineId, err)
	}

	instanceView, err := virtualMachinesClient.InstanceView(ctx, *virtualMachineId)
	if err != nil {
		return fmt.Errorf("retrieving InstanceView for %s: %+v", virtualMachineId, err)
	}

	shouldTurnBackOn := virtualMachineShouldBeStarted(instanceView.Model)
	shouldDeallocate := true

	if instanceView.Model != nil && instanceView.Model.Statuses != nil {
		for _, status := range *instanceView.Model.Statuses {
			if status.Code == nil {
				continue
			}

			// could also be the provisioning state which we're not bothered with here
			state := strings.ToLower(*status.Code)
			if !strings.HasPrefix(state, "powerstate/") {
				continue
			}

			state = strings.TrimPrefix(state, "powerstate/")
			switch strings.ToLower(state) {
			case "deallocated":
				// VM already deallocated, no shutdown and deallocation needed anymore
				shouldShutDown = false
				shouldDeallocate = false
			case "deallocating":
				// VM is deallocating
				// To make sure we do not start updating before this action has finished,
				// only skip the shutdown and send another deallocation request if shouldDeallocate == true
				shouldShutDown = false
			case "stopped":
				shouldShutDown = false
			}
		}
	}

	// Detach
	if shouldDetach {
		dataDisks := make([]virtualmachines.DataDisk, 0)
		if vmModel := vm.Model; vmModel != nil && vmModel.Properties != nil && vmModel.Properties.StorageProfile != nil && vmModel.Properties.StorageProfile.DataDisks != nil {
			for _, dataDisk := range *vmModel.Properties.StorageProfile.DataDisks {
				// since this field isn't (and shouldn't be) case-sensitive; we're deliberately not using `strings.EqualFold`
				if dataDisk.Name != nil && *dataDisk.Name != id.DiskName {
					dataDisks = append(dataDisks, dataDisk)
				} else {
					if dataDisk.Caching != nil && *dataDisk.Caching != virtualmachines.CachingTypesNone {
						return fmt.Errorf("`disk_size_gb` can't be increased above 4095GB when `caching` is set to anything other than `None`")
					}
					expandedDisk = dataDisk
				}
			}

			vmModel.Properties.StorageProfile.DataDisks = &dataDisks

			// fixes #2485
			vmModel.Identity = nil
			// fixes #1600
			vmModel.Resources = nil

			if err := virtualMachinesClient.CreateOrUpdateThenPoll(ctx, *virtualMachineId, *vm.Model, virtualmachines.DefaultCreateOrUpdateOperationOptions()); err != nil {
				return fmt.Errorf("removing Disk %q from %s : %+v", id.DiskName, virtualMachineId, err)
			}
		}
	}

	// Shutdown
	if shouldShutDown {
		log.Printf("[DEBUG] Shutting Down %s", virtualMachineId)
		options := virtualmachines.DefaultPowerOffOperationOptions()
		options.SkipShutdown = pointer.To(false)
		if err := virtualMachinesClient.PowerOffThenPoll(ctx, *virtualMachineId, options); err != nil {
			return fmt.Errorf("sending Power Off to %s: %+v", virtualMachineId, err)
		}

		log.Printf("[DEBUG] Shut Down %s", virtualMachineId)
	}

	// De-allocate
	if shouldDeallocate {
		log.Printf("[DEBUG] Deallocating %s.", virtualMachineId)
		// Upgrading to 2021-07-01 exposed a new hibernate paramater to the Deallocate method
		if err := virtualMachinesClient.DeallocateThenPoll(ctx, *virtualMachineId, virtualmachines.DefaultDeallocateOperationOptions()); err != nil {
			return fmt.Errorf("deallocating to %s: %+v", virtualMachineId, err)
		}
		log.Printf("[DEBUG] Deallocated %s", virtualMachineId)
	}

	// Update Disk
	err = diskClient.UpdateThenPoll(ctx, *id, diskUpdate)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	// Reattach DataDisk
	if shouldDetach && vm.Model.Properties.StorageProfile != nil {
		disks := *vm.Model.Properties.StorageProfile.DataDisks

		expandedDisk.DiskSizeGB = diskUpdate.Properties.DiskSizeGB
		disks = append(disks, expandedDisk)

		vm.Model.Properties.StorageProfile.DataDisks = &disks

		// fixes #2485
		vm.Model.Identity = nil
		// fixes #1600
		vm.Model.Resources = nil

		// if there's too many disks we get a 409 back with:
		//   `The maximum number of data disks allowed to be attached to a VM of this size is 1.`
		// which we're intentionally not wrapping, since the errors good.
		if err := virtualMachinesClient.CreateOrUpdateThenPoll(ctx, *virtualMachineId, *vm.Model, virtualmachines.DefaultCreateOrUpdateOperationOptions()); err != nil {
			return fmt.Errorf("updating %s to reattach Disk %s: %+v", virtualMachineId, id, err)
		}
	}

	if shouldTurnBackOn && (shouldShutDown || shouldDeallocate) {
		log.Printf("[DEBUG] Starting %s", virtualMachineId)
		if err := virtualMachinesClient.StartThenPoll(ctx, *virtualMachineId); err != nil {
			return fmt.Errorf("starting %s: %+v", virtualMachineId, err)
		}
		log.Printf("[DEBUG] Started %s", virtualMachineId)
	}

	return nil
}
