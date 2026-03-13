// Copyright IBM Corp. 2014, 2025
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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/extendedlocation/2021-08-15/customlocations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/inventoryitems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualnetworks"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/vmmservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name system_center_virtual_machine_manager_virtual_network -service-package-name systemcentervirtualmachinemanager -properties "name,resource_group_name" -known-values "subscription_id:data.Subscriptions.Primary" -test-sequential

type SystemCenterVirtualMachineManagerVirtualNetworkModel struct {
	Name                                                   string            `tfschema:"name"`
	Location                                               string            `tfschema:"location"`
	ResourceGroupName                                      string            `tfschema:"resource_group_name"`
	CustomLocationId                                       string            `tfschema:"custom_location_id"`
	SystemCenterVirtualMachineManagerServerInventoryItemId string            `tfschema:"system_center_virtual_machine_manager_server_inventory_item_id"`
	Tags                                                   map[string]string `tfschema:"tags"`
}

var (
	_ sdk.Resource             = SystemCenterVirtualMachineManagerVirtualNetworkResource{}
	_ sdk.ResourceWithUpdate   = SystemCenterVirtualMachineManagerVirtualNetworkResource{}
	_ sdk.ResourceWithIdentity = SystemCenterVirtualMachineManagerVirtualNetworkResource{}
)

type SystemCenterVirtualMachineManagerVirtualNetworkResource struct{}

func (r SystemCenterVirtualMachineManagerVirtualNetworkResource) ModelObject() interface{} {
	return &SystemCenterVirtualMachineManagerVirtualNetworkModel{}
}

func (r SystemCenterVirtualMachineManagerVirtualNetworkResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return virtualnetworks.ValidateVirtualNetworkID
}

func (r SystemCenterVirtualMachineManagerVirtualNetworkResource) ResourceType() string {
	return "azurerm_system_center_virtual_machine_manager_virtual_network"
}

func (r SystemCenterVirtualMachineManagerVirtualNetworkResource) Identity() resourceids.ResourceId {
	return &virtualnetworks.VirtualNetworkId{}
}

func (r SystemCenterVirtualMachineManagerVirtualNetworkResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SystemCenterVirtualMachineManagerVirtualNetworkName,
		},

		"location": commonschema.Location(),

		"resource_group_name": commonschema.ResourceGroupName(),

		"custom_location_id": commonschema.ResourceIDReferenceRequiredForceNew(&customlocations.CustomLocationId{}),

		"system_center_virtual_machine_manager_server_inventory_item_id": commonschema.ResourceIDReferenceRequiredForceNew(&inventoryitems.InventoryItemId{}),

		"tags": commonschema.Tags(),
	}
}

func (r SystemCenterVirtualMachineManagerVirtualNetworkResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r SystemCenterVirtualMachineManagerVirtualNetworkResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.SystemCenterVirtualMachineManager.VirtualNetworks

			var model SystemCenterVirtualMachineManagerVirtualNetworkModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			scvmmServerInventoryItemId, err := inventoryitems.ParseInventoryItemID(model.SystemCenterVirtualMachineManagerServerInventoryItemId)
			if err != nil {
				return err
			}

			id := virtualnetworks.NewVirtualNetworkID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := virtualnetworks.VirtualNetwork{
				ExtendedLocation: virtualnetworks.ExtendedLocation{
					Type: pointer.To("customLocation"),
					Name: pointer.To(model.CustomLocationId),
				},
				Location: location.Normalize(model.Location),
				Properties: &virtualnetworks.VirtualNetworkProperties{
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
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}
			return nil
		},
	}
}

func (r SystemCenterVirtualMachineManagerVirtualNetworkResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VirtualNetworks

			id, err := virtualnetworks.ParseVirtualNetworkID(metadata.ResourceData.Id())
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

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func (r SystemCenterVirtualMachineManagerVirtualNetworkResource) flatten(metadata sdk.ResourceMetaData, id *virtualnetworks.VirtualNetworkId, model *virtualnetworks.VirtualNetwork) error {
	state := SystemCenterVirtualMachineManagerVirtualNetworkModel{
		Name:              id.VirtualNetworkName,
		ResourceGroupName: id.ResourceGroupName,
	}

	if model != nil {
		state.Location = location.Normalize(model.Location)
		state.CustomLocationId = pointer.From(model.ExtendedLocation.Name)
		state.Tags = pointer.From(model.Tags)

		if model.Properties != nil {
			state.SystemCenterVirtualMachineManagerServerInventoryItemId = pointer.From(model.Properties.InventoryItemId)
		}
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func (r SystemCenterVirtualMachineManagerVirtualNetworkResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VirtualNetworks

			id, err := virtualnetworks.ParseVirtualNetworkID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SystemCenterVirtualMachineManagerVirtualNetworkModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := virtualnetworks.VirtualNetworkTagsUpdate{}

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

func (r SystemCenterVirtualMachineManagerVirtualNetworkResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VirtualNetworks

			id, err := virtualnetworks.ParseVirtualNetworkID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, virtualnetworks.DeleteOperationOptions{Force: pointer.To(virtualnetworks.ForceDeleteTrue)}); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
