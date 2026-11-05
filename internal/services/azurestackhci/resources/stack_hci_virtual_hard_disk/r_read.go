// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package stack_hci_virtual_hard_disk

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/virtualharddisks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.VirtualHardDisks

			id, err := virtualharddisks.ParseVirtualHardDiskID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			schema := StackHCIVirtualHardDiskResourceModel{
				Name:              id.VirtualHardDiskName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				schema.Location = location.Normalize(model.Location)
				schema.Tags = tags.Flatten(model.Tags)

				if model.ExtendedLocation != nil && model.ExtendedLocation.Name != nil {
					customLocationId, err := customlocations.ParseCustomLocationIDInsensitively(*model.ExtendedLocation.Name)
					if err != nil {
						return err
					}

					schema.CustomLocationId = customLocationId.ID()
				}

				if props := model.Properties; props != nil {
					schema.BlockSizeInBytes = pointer.From(props.BlockSizeBytes)
					schema.DiskFileFormat = pointer.FromEnum(props.DiskFileFormat)
					schema.DiskSizeInGB = pointer.From(props.DiskSizeGB)
					schema.DynamicEnabled = pointer.From(props.Dynamic)
					schema.HypervGeneration = pointer.FromEnum(props.HyperVGeneration)
					schema.LogicalSectorInBytes = pointer.From(props.LogicalSectorBytes)
					schema.PhysicalSectorInBytes = pointer.From(props.PhysicalSectorBytes)
					schema.StoragePathId = pointer.From(props.ContainerId)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}
