// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_implicit_data_disk_from_source

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2023-04-02/disks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachinesClient

			id, err := parse.DataDiskID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			virtualMachineId := virtualmachines.NewVirtualMachineID(id.SubscriptionId, id.ResourceGroup, id.VirtualMachineName)
			resp, err := client.Get(ctx, virtualMachineId, virtualmachines.DefaultGetOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", virtualMachineId, err)
			}

			schema := VirtualMachineImplicitDataDiskFromSourceResourceModel{
				Name:             id.Name,
				VirtualMachineId: virtualMachineId.ID(),
			}

			var disk *virtualmachines.DataDisk
			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if profile := props.StorageProfile; profile != nil {
						if dataDisks := profile.DataDisks; dataDisks != nil {
							for _, dataDisk := range *dataDisks {
								if pointer.From(dataDisk.Name) == id.Name {
									disk = &dataDisk
									break
								}
							}
						}
					}
				}
			}

			if disk == nil {
				return metadata.MarkAsGone(*id)
			}

			schema.Lun = disk.Lun
			if v := pointer.From(disk.Caching); v != virtualmachines.CachingTypesNone {
				schema.Caching = string(v)
			}

			schema.CreateOption = string(disk.CreateOption)
			schema.DiskSizeGb = pointer.From(disk.DiskSizeGB)
			if disk.SourceResource != nil {
				schema.SourceResourceId = pointer.From(disk.SourceResource.Id)
			}

			schema.WriteAcceleratorEnabled = pointer.From(disk.WriteAcceleratorEnabled)

			return metadata.Encode(&schema)
		},
	}
}

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		client := metadata.Client.Compute.VirtualMachinesClient

		id, err := parse.DataDiskID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		virtualMachineId := virtualmachines.NewVirtualMachineID(id.SubscriptionId, id.ResourceGroup, id.VirtualMachineName)
		resp, err := client.Get(ctx, virtualMachineId, virtualmachines.DefaultGetOperationOptions())
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", virtualMachineId, err)
		}

		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil {
				if profile := props.StorageProfile; profile != nil {
					if dataDisks := profile.DataDisks; dataDisks != nil {
						var disk *virtualmachines.DataDisk
						for _, dataDisk := range *dataDisks {
							if pointer.From(dataDisk.Name) == id.Name {
								disk = &dataDisk
								break
							}
						}

						if disk == nil {
							return fmt.Errorf("unable to retrieve an existing data disk %s", *id)
						}

						if disk.CreateOption != virtualmachines.DiskCreateOptionTypesCopy {
							return fmt.Errorf("the value of `create_option` for the imported `azurerm_virtual_machine_implicit_data_disk_from_source` instance must be `Copy`, whereas now is %s", disk.CreateOption)
						}
					}
				}
			}
		}

		return nil
	}
}

func resizeImplicitDataDisk(ctx context.Context, metadata sdk.ResourceMetaData, id *parse.DataDiskId) error {
	diskClient := metadata.Client.Compute.DisksClient
	skusClient := metadata.Client.Compute.SkusClient
	virtualMachinesClient := metadata.Client.Compute.VirtualMachinesClient
	shouldShutDown := false
	shouldDetach := false

	managedDiskId := commonids.NewManagedDiskID(id.SubscriptionId, id.ResourceGroup, id.Name)

	disk, err := diskClient.Get(ctx, managedDiskId)
	if err != nil {
		if response.WasNotFound(disk.HttpResponse) {
			return fmt.Errorf("%s was not found", managedDiskId)
		}

		return fmt.Errorf("checking for presence of existing %s: %+v", managedDiskId, err)
	}

	diskUpdate := disks.DiskUpdate{
		Properties: &disks.DiskUpdateProperties{},
	}

	oldSize, newSize := metadata.ResourceData.GetChange("disk_size_gb")
	canBeResizedWithoutDowntime := false
	if metadata.Client.Features.ManagedDisk.ExpandWithoutDowntime {
		shouldDetach = helpers.DetermineIfDataDiskRequiresDetaching(disk.Model, oldSize.(int), newSize.(int))
		diskSupportsNoDowntimeResize := helpers.DetermineIfDataDiskSupportsNoDowntimeResize(disk.Model, shouldDetach)

		vmSupportsNoDowntimeResize, err := helpers.DetermineIfVirtualMachineSupportsNoDowntimeResize(ctx, disk.Model, virtualMachinesClient, skusClient)
		if err != nil {
			return fmt.Errorf("determining if the Virtual Machine supports no-downtime-resizing: %+v", err)
		}

		canBeResizedWithoutDowntime = *vmSupportsNoDowntimeResize && diskSupportsNoDowntimeResize
	}

	if !canBeResizedWithoutDowntime {
		log.Printf("[INFO] The %s, or the Virtual Machine that it's attached to, doesn't support no-downtime-resizing - requiring that the VM should be shutdown", *id)
		shouldShutDown = true
	}

	diskUpdate.Properties.DiskSizeGB = pointer.To(int64(newSize.(int)))

	if shouldShutDown {
		virtualMachineId, err := virtualmachines.ParseVirtualMachineID(*disk.Model.ManagedBy)
		if err != nil {
			return err
		}

		if err = helpers.ResourceManagedDiskUpdateWithVmShutDown(ctx, metadata.Client, pointer.To(commonids.NewManagedDiskID(id.SubscriptionId, id.ResourceGroup, id.Name)), virtualMachineId, diskUpdate, shouldDetach); err != nil {
			return err
		}
	} else { // otherwise, just update it
		if err := diskClient.UpdateThenPoll(ctx, managedDiskId, diskUpdate); err != nil {
			return fmt.Errorf("updating %s: %+v", managedDiskId, err)
		}
	}

	return nil
}
