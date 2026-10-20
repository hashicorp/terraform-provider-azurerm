// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_implicit_data_disk_from_source

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachinesClient

			var config VirtualMachineImplicitDataDiskFromSourceResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := parse.DataDiskID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.VirtualMachineName, helpers.VirtualMachineResourceName)
			defer locks.UnlockByName(id.VirtualMachineName, helpers.VirtualMachineResourceName)

			virtualMachineId := virtualmachines.NewVirtualMachineID(id.SubscriptionId, id.ResourceGroup, id.VirtualMachineName)
			resp, err := client.Get(ctx, virtualMachineId, virtualmachines.DefaultGetOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", virtualMachineId, err)
			}

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if profile := props.StorageProfile; profile != nil {
						if dataDisks := profile.DataDisks; dataDisks != nil {
							existingIndex := -1
							disks := *dataDisks
							for i, disk := range disks {
								if pointer.From(disk.Name) == config.Name {
									existingIndex = i
									break
								}
							}

							if existingIndex == -1 {
								return fmt.Errorf("unable to retrieve the data disk %s ", *id)
							}

							needToUpdate := false
							expandedDisk := &disks[existingIndex]
							if metadata.ResourceData.HasChange("caching") {
								caching := string(virtualmachines.CachingTypesNone)
								if config.Caching != "" {
									caching = config.Caching
								}

								expandedDisk.Caching = pointer.ToEnum[virtualmachines.CachingTypes](caching)
								needToUpdate = true
							}

							if metadata.ResourceData.HasChange("write_accelerator_enabled") {
								expandedDisk.WriteAcceleratorEnabled = pointer.To(config.WriteAcceleratorEnabled)
								needToUpdate = true
							}

							if needToUpdate {
								profile.DataDisks = &disks
								// fixes #24145
								model.Properties.ApplicationProfile = nil
								// fixes #2485
								model.Identity = nil
								// fixes #1600
								model.Resources = nil

								if err = client.CreateOrUpdateThenPoll(ctx, virtualMachineId, *model, virtualmachines.DefaultCreateOrUpdateOperationOptions()); err != nil {
									return fmt.Errorf("updating %s: %+v", id, err)
								}
							}

							// VirtualMachineClient does not support expanding implicit data disk size, need to use DisksClient
							if metadata.ResourceData.HasChange("disk_size_gb") {
								if err = resizeImplicitDataDisk(ctx, metadata, id); err != nil {
									return err
								}
							}
						}
					}
				}
			}

			return nil
		},
	}
}
