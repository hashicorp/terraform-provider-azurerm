// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package systemcentervirtualmachinemanager

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/clouds"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/inventoryitems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/vmmservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SystemCenterVirtualMachineManagerCloudModel struct {
	Name                                                   string            `tfschema:"name"`
	Location                                               string            `tfschema:"location"`
	ResourceGroupName                                      string            `tfschema:"resource_group_name"`
	CustomLocationId                                       string            `tfschema:"custom_location_id"`
	SystemCenterVirtualMachineManagerServerInventoryItemId string            `tfschema:"system_center_virtual_machine_manager_server_inventory_item_id"`
	Tags                                                   map[string]string `tfschema:"tags"`
}

var _ sdk.Resource = SystemCenterVirtualMachineManagerCloudResource{}
var _ sdk.ResourceWithUpdate = SystemCenterVirtualMachineManagerCloudResource{}

type SystemCenterVirtualMachineManagerCloudResource struct{}

func (r SystemCenterVirtualMachineManagerCloudResource) ModelObject() interface{} {
	return &SystemCenterVirtualMachineManagerCloudModel{}
}

func (r SystemCenterVirtualMachineManagerCloudResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return clouds.ValidateCloudID
}

func (r SystemCenterVirtualMachineManagerCloudResource) ResourceType() string {
	return "azurerm_system_center_virtual_machine_manager_cloud"
}

func (r SystemCenterVirtualMachineManagerCloudResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SystemCenterVirtualMachineManagerCloudName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"custom_location_id": commonschema.ResourceIDReferenceRequiredForceNew(&customlocations.CustomLocationId{}),

		"system_center_virtual_machine_manager_server_inventory_item_id": commonschema.ResourceIDReferenceRequiredForceNew(&inventoryitems.InventoryItemId{}),

		"tags": commonschema.Tags(),
	}
}

func (r SystemCenterVirtualMachineManagerCloudResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r SystemCenterVirtualMachineManagerCloudResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.SystemCenterVirtualMachineManager.Clouds

			var model SystemCenterVirtualMachineManagerCloudModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			scvmmServerInventoryItemId, err := inventoryitems.ParseInventoryItemID(model.SystemCenterVirtualMachineManagerServerInventoryItemId)
			if err != nil {
				return err
			}

			id := clouds.NewCloudID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := clouds.Cloud{
				Location: location.Normalize(model.Location),
				ExtendedLocation: clouds.ExtendedLocation{
					Type: pointer.To("customLocation"),
					Name: pointer.To(model.CustomLocationId),
				},
				Properties: &clouds.CloudProperties{
					InventoryItemId: pointer.To(scvmmServerInventoryItemId.ID()),
					Uuid:            pointer.To(scvmmServerInventoryItemId.InventoryItemName),
					VMmServerId:     pointer.To(vmmservers.NewVMmServerID(scvmmServerInventoryItemId.SubscriptionId, scvmmServerInventoryItemId.ResourceGroupName, scvmmServerInventoryItemId.VmmServerName).ID()),
				},
				Tags: pointer.To(model.Tags),
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r SystemCenterVirtualMachineManagerCloudResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.Clouds

			id, err := clouds.ParseCloudID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := SystemCenterVirtualMachineManagerCloudModel{
				Name:              id.CloudName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.CustomLocationId = pointer.From(model.ExtendedLocation.Name)
				state.SystemCenterVirtualMachineManagerServerInventoryItemId = pointer.From(model.Properties.InventoryItemId)
				state.Tags = pointer.From(model.Tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r SystemCenterVirtualMachineManagerCloudResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.Clouds

			id, err := clouds.ParseCloudID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SystemCenterVirtualMachineManagerCloudModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := clouds.CloudTagsUpdate{}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = pointer.To(model.Tags)
			}

			if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r SystemCenterVirtualMachineManagerCloudResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.Clouds

			id, err := clouds.ParseCloudID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, clouds.DeleteOperationOptions{Force: pointer.To(clouds.ForceDeleteTrue)}); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
