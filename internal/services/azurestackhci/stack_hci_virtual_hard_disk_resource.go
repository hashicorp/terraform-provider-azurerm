// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package azurestackhci

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/storagecontainers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/virtualharddisks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource           = StackHCIVirtualHardDiskResource{}
	_ sdk.ResourceWithUpdate = StackHCIVirtualHardDiskResource{}
)

type StackHCIVirtualHardDiskResource struct{}

func (StackHCIVirtualHardDiskResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return virtualharddisks.ValidateVirtualHardDiskID
}

func (StackHCIVirtualHardDiskResource) ResourceType() string {
	return "azurerm_stack_hci_virtual_hard_disk"
}

func (StackHCIVirtualHardDiskResource) ModelObject() interface{} {
	return &StackHCIVirtualHardDiskResourceModel{}
}

type StackHCIVirtualHardDiskResourceModel struct {
	Name                  string                 `tfschema:"name"`
	ResourceGroupName     string                 `tfschema:"resource_group_name"`
	Location              string                 `tfschema:"location"`
	CustomLocationId      string                 `tfschema:"custom_location_id"`
	BlockSizeInBytes      int64                  `tfschema:"block_size_in_bytes"`
	DiskFileFormat        string                 `tfschema:"disk_file_format"`
	DiskSizeInGB          int64                  `tfschema:"disk_size_in_gb"`
	DynamicEnabled        bool                   `tfschema:"dynamic_enabled"`
	HypervGeneration      string                 `tfschema:"hyperv_generation"`
	LogicalSectorInBytes  int64                  `tfschema:"logical_sector_in_bytes"`
	PhysicalSectorInBytes int64                  `tfschema:"physical_sector_in_bytes"`
	StoragePathId         string                 `tfschema:"storage_path_id"`
	Tags                  map[string]interface{} `tfschema:"tags"`
}

func (StackHCIVirtualHardDiskResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9][-._a-zA-Z0-9]{0,62}[a-zA-Z0-9]$`),
				"name must be between 2 and 64 characters and can only contain alphanumberic characters, hyphen, dot and underline",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"custom_location_id": commonschema.ResourceIDReferenceRequiredForceNew(&customlocations.CustomLocationId{}),

		"disk_size_in_gb": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"block_size_in_bytes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"disk_file_format": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualharddisks.DiskFileFormatVhd),
				string(virtualharddisks.DiskFileFormatVhdx),
			}, false),
		},

		"dynamic_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			ForceNew: true,
			Default:  false,
		},

		"hyperv_generation": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(virtualharddisks.HyperVGenerationVOne),
				string(virtualharddisks.HyperVGenerationVTwo),
			}, false),
		},

		"logical_sector_in_bytes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"physical_sector_in_bytes": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntAtLeast(1),
		},

		"storage_path_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: storagecontainers.ValidateStorageContainerID,
		},

		"tags": commonschema.Tags(),
	}
}

func (StackHCIVirtualHardDiskResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r StackHCIVirtualHardDiskResource) Create() sdk.ResourceFunc {
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

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
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
				payload.Properties.DiskFileFormat = pointer.To(virtualharddisks.DiskFileFormat(config.DiskFileFormat))
			}

			if config.HypervGeneration != "" {
				payload.Properties.HyperVGeneration = pointer.To(virtualharddisks.HyperVGeneration(config.HypervGeneration))
			}

			if config.LogicalSectorInBytes != 0 {
				payload.Properties.LogicalSectorBytes = pointer.To(config.LogicalSectorInBytes)
			}

			if config.PhysicalSectorInBytes != 0 {
				payload.Properties.PhysicalSectorBytes = pointer.To(config.PhysicalSectorInBytes)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
				return fmt.Errorf("performing create %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r StackHCIVirtualHardDiskResource) Read() sdk.ResourceFunc {
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
					schema.DiskFileFormat = string(pointer.From(props.DiskFileFormat))
					schema.DiskSizeInGB = pointer.From(props.DiskSizeGB)
					schema.DynamicEnabled = pointer.From(props.Dynamic)
					schema.HypervGeneration = string(pointer.From(props.HyperVGeneration))
					schema.LogicalSectorInBytes = pointer.From(props.LogicalSectorBytes)
					schema.PhysicalSectorInBytes = pointer.From(props.PhysicalSectorBytes)
					schema.StoragePathId = pointer.From(props.ContainerId)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r StackHCIVirtualHardDiskResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.VirtualHardDisks

			id, err := virtualharddisks.ParseVirtualHardDiskID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model StackHCIVirtualHardDiskResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := virtualharddisks.VirtualHardDisksUpdateRequest{}
			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = tags.Expand(model.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r StackHCIVirtualHardDiskResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AzureStackHCI.VirtualHardDisks

			id, err := virtualharddisks.ParseVirtualHardDiskID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}
