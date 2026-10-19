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

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.VirtualMachinesClient

			var config VirtualMachineImplicitDataDiskFromSourceResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			virtualMachineId, err := virtualmachines.ParseVirtualMachineID(config.VirtualMachineId)
			if err != nil {
				return err
			}

			locks.ByName(virtualMachineId.VirtualMachineName, helpers.VirtualMachineResourceName)
			defer locks.UnlockByName(virtualMachineId.VirtualMachineName, helpers.VirtualMachineResourceName)

			id := parse.NewDataDiskID(subscriptionId, virtualMachineId.ResourceGroupName, virtualMachineId.VirtualMachineName, config.Name)

			resp, err := client.Get(ctx, *virtualMachineId, virtualmachines.DefaultGetOperationOptions())
			if err != nil {
				return fmt.Errorf("checking for the presence of an existing %s: %+v", *virtualMachineId, err)
			}

			caching := string(virtualmachines.CachingTypesNone)
			if config.Caching != "" {
				caching = config.Caching
			}

			expandedDisk := virtualmachines.DataDisk{
				Name:         pointer.To(config.Name),
				Caching:      pointer.ToEnum[virtualmachines.CachingTypes](caching),
				CreateOption: virtualmachines.DiskCreateOptionTypes(config.CreateOption),
				DiskSizeGB:   pointer.To(config.DiskSizeGb),
				Lun:          config.Lun,
				SourceResource: &virtualmachines.ApiEntityReference{
					Id: pointer.To(config.SourceResourceId),
				},
				WriteAcceleratorEnabled: pointer.To(config.WriteAcceleratorEnabled),
			}

			if model := resp.Model; model != nil {
				if model.Properties != nil {
					// there are ways to provision a VM without a StorageProfile and/or DataDisks
					if model.Properties.StorageProfile == nil {
						model.Properties.StorageProfile = &virtualmachines.StorageProfile{}
					}

					if model.Properties.StorageProfile.DataDisks == nil {
						model.Properties.StorageProfile.DataDisks = pointer.To(make([]virtualmachines.DataDisk, 0))
					}

					existingIndex := -1
					disks := *model.Properties.StorageProfile.DataDisks
					for i, disk := range disks {
						if pointer.From(disk.Name) == config.Name {
							existingIndex = i
							break
						}
					}

					if existingIndex != -1 {
						return metadata.ResourceRequiresImport(r.ResourceType(), id)
					}

					disks = append(disks, expandedDisk)
					model.Properties.StorageProfile.DataDisks = &disks
					// fixes #24145
					model.Properties.ApplicationProfile = nil
					// fixes #2485
					model.Identity = nil
					// fixes #1600
					model.Resources = nil

					// TODO: implement `CallbackThenPoll`, requires migrating to an ID that implements `resourceids.ResourceId`
					if err = client.CreateOrUpdateThenPoll(ctx, *virtualMachineId, *model, virtualmachines.DefaultCreateOrUpdateOperationOptions()); err != nil {
						return fmt.Errorf("creating %s: %+v", id, err)
					}
				}
			}

			metadata.SetID(id)
			return nil
		},
	}
}
