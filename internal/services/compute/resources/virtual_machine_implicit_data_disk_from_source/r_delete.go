// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_implicit_data_disk_from_source

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/virtualmachines"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
)

func (r Resource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachinesClient

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
						newDisks := make([]virtualmachines.DataDisk, 0)
						var toBeDeletedDisk *virtualmachines.DataDisk
						if dataDisks := profile.DataDisks; dataDisks != nil {
							for _, dataDisk := range *dataDisks {
								if pointer.From(dataDisk.Name) != id.Name {
									newDisks = append(newDisks, dataDisk)
								} else {
									toBeDeletedDisk = pointer.To(dataDisk)
								}
							}
						}

						profile.DataDisks = &newDisks

						// fixes #24145
						model.Properties.ApplicationProfile = nil

						// fixes #2485
						model.Identity = nil
						// fixes #1600
						model.Resources = nil

						if err = client.CreateOrUpdateThenPoll(ctx, virtualMachineId, *model, virtualmachines.DefaultCreateOrUpdateOperationOptions()); err != nil {
							return fmt.Errorf("deleting %s: %+v", id, err)
						}

						// delete the data disk which was created by Azure Service when creating this resource
						detachDataDisk := metadata.Client.Features.VirtualMachine.DetachImplicitDataDiskOnDeletion
						if !detachDataDisk && toBeDeletedDisk != nil && toBeDeletedDisk.ManagedDisk != nil && toBeDeletedDisk.ManagedDisk.Id != nil {
							diskClient := metadata.Client.Compute.DisksClient
							diskId, err := commonids.ParseManagedDiskID(*toBeDeletedDisk.ManagedDisk.Id)
							if err != nil {
								return err
							}

							if err = diskClient.DeleteThenPoll(ctx, *diskId); err != nil {
								return fmt.Errorf("deleting Managed Disk %s: %+v", *diskId, err)
							}
						}
					}
				}
			}

			return nil
		},
	}
}
