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
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.VirtualHardDisks

			var config StackHCIVirtualHardDiskResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			subscriptionId := metadata.Client.Account.SubscriptionId
			id := virtualharddisks.NewVirtualHardDiskID(subscriptionId, config.ResourceGroupName, config.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			payload := virtualharddisks.VirtualHardDisks{
				Name:     pointer.To(config.Name),
				Location: location.Normalize(config.Location),
				Tags:     tags.Expand(config.Tags),
				ExtendedLocation: &virtualharddisks.ExtendedLocation{
					Name: pointer.To(config.CustomLocationId),
					Type: pointer.To(virtualharddisks.ExtendedLocationTypesCustomLocation),
				},
				Properties: &virtualharddisks.VirtualHardDiskProperties{
					Dynamic:    pointer.To(config.DynamicEnabled),
					DiskSizeGB: pointer.To(config.DiskSizeInGB),
				},
			}

			if config.BlockSizeInBytes != 0 {
				payload.Properties.BlockSizeBytes = pointer.To(config.BlockSizeInBytes)
			}

			if config.StoragePathId != "" {
				payload.Properties.ContainerId = pointer.To(config.StoragePathId)
			}

			if config.DiskFileFormat != "" {
				payload.Properties.DiskFileFormat = pointer.ToEnum[virtualharddisks.DiskFileFormat](config.DiskFileFormat)
			}

			if config.HypervGeneration != "" {
				payload.Properties.HyperVGeneration = pointer.ToEnum[virtualharddisks.HyperVGeneration](config.HypervGeneration)
			}

			if config.LogicalSectorInBytes != 0 {
				payload.Properties.LogicalSectorBytes = pointer.To(config.LogicalSectorInBytes)
			}

			if config.PhysicalSectorInBytes != 0 {
				payload.Properties.PhysicalSectorBytes = pointer.To(config.PhysicalSectorInBytes)
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, payload, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("performing create %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}
