// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2024-03-01/restorepointcollections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// VirtualMachineRestorePointCollectionResource remove this in 4.0, the resource is renamed
type VirtualMachineRestorePointCollectionResource struct{}

var _ sdk.ResourceWithUpdate = VirtualMachineRestorePointCollectionResource{}
var _ sdk.ResourceWithDeprecationReplacedBy = VirtualMachineRestorePointCollectionResource{}

func (r VirtualMachineRestorePointCollectionResource) DeprecatedInFavourOfResource() string {
	return "azurerm_virtual_machine_restore_point_collection"
}

func (r VirtualMachineRestorePointCollectionResource) ModelObject() interface{} {
	return &VirtualMachineRestorePointCollectionResourceModel{}
}

type VirtualMachineRestorePointCollectionResourceModel struct {
	Name                   string                 `tfschema:"name"`
	ResourceGroup          string                 `tfschema:"resource_group_name"`
	Location               string                 `tfschema:"location"`
	SourceVirtualMachineId string                 `tfschema:"source_virtual_machine_id"`
	Tags                   map[string]interface{} `tfschema:"tags"`
}

func (r VirtualMachineRestorePointCollectionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return restorepointcollections.ValidateRestorePointCollectionID
}

func (r VirtualMachineRestorePointCollectionResource) ResourceType() string {
	return "azurerm_virtual_machine_restore_point_collection"
}

func (r VirtualMachineRestorePointCollectionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"source_virtual_machine_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateVirtualMachineID,
		},

		"tags": commonschema.Tags(),
	}
}

func (r VirtualMachineRestorePointCollectionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r VirtualMachineRestorePointCollectionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.RestorePointCollectionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config VirtualMachineRestorePointCollectionResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := restorepointcollections.NewRestorePointCollectionID(subscriptionId, config.ResourceGroup, config.Name)

			existing, err := client.Get(ctx, id, restorepointcollections.DefaultGetOperationOptions())
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := restorepointcollections.RestorePointCollection{
				Location: location.Normalize(config.Location),
				Properties: &restorepointcollections.RestorePointCollectionProperties{
					Source: &restorepointcollections.RestorePointCollectionSourceProperties{
						Id: pointer.To(config.SourceVirtualMachineId),
					},
				},
				Tags: tags.Expand(config.Tags),
			}

			if _, err = client.CreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r VirtualMachineRestorePointCollectionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.RestorePointCollectionsClient

			schema := VirtualMachineRestorePointCollectionResourceModel{}

			id, err := restorepointcollections.ParseRestorePointCollectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id, restorepointcollections.DefaultGetOperationOptions())
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.Name = id.RestorePointCollectionName
				schema.ResourceGroup = id.ResourceGroupName

				if props := model.Properties; props != nil {
					if source := props.Source; source != nil {
						schema.SourceVirtualMachineId = pointer.From(source.Id)
						schema.Location = location.Normalize(pointer.From(source.Location))
					}
				}

				schema.Tags = tags.Flatten(model.Tags)
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r VirtualMachineRestorePointCollectionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.RestorePointCollectionsClient

			id, err := restorepointcollections.ParseRestorePointCollectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id, restorepointcollections.DefaultGetOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			payload := *existing.Model

			var config VirtualMachineRestorePointCollectionResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = tags.Expand(config.Tags)
			}

			if _, err = client.CreateOrUpdate(ctx, *id, payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r VirtualMachineRestorePointCollectionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Compute.RestorePointCollectionsClient

			id, err := restorepointcollections.ParseRestorePointCollectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
