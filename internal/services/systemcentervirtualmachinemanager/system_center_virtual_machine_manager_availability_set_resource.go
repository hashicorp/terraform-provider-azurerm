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
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/availabilitysets"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/vmmservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SystemCenterVirtualMachineManagerAvailabilitySetModel struct {
	Name                                      string            `tfschema:"name"`
	ResourceGroupName                         string            `tfschema:"resource_group_name"`
	Location                                  string            `tfschema:"location"`
	CustomLocationId                          string            `tfschema:"custom_location_id"`
	SystemCenterVirtualMachineManagerServerId string            `tfschema:"system_center_virtual_machine_manager_server_id"`
	Tags                                      map[string]string `tfschema:"tags"`
}

var _ sdk.Resource = SystemCenterVirtualMachineManagerAvailabilitySetResource{}
var _ sdk.ResourceWithUpdate = SystemCenterVirtualMachineManagerAvailabilitySetResource{}

type SystemCenterVirtualMachineManagerAvailabilitySetResource struct{}

func (r SystemCenterVirtualMachineManagerAvailabilitySetResource) ModelObject() interface{} {
	return &SystemCenterVirtualMachineManagerAvailabilitySetModel{}
}

func (r SystemCenterVirtualMachineManagerAvailabilitySetResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return availabilitysets.ValidateAvailabilitySetID
}

func (r SystemCenterVirtualMachineManagerAvailabilitySetResource) ResourceType() string {
	return "azurerm_system_center_virtual_machine_manager_availability_set"
}

func (r SystemCenterVirtualMachineManagerAvailabilitySetResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SystemCenterVirtualMachineManagerAvailabilitySetName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"custom_location_id": commonschema.ResourceIDReferenceRequiredForceNew(&customlocations.CustomLocationId{}),

		"system_center_virtual_machine_manager_server_id": commonschema.ResourceIDReferenceRequiredForceNew(&vmmservers.VMmServerId{}),

		"tags": commonschema.Tags(),
	}
}

func (r SystemCenterVirtualMachineManagerAvailabilitySetResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r SystemCenterVirtualMachineManagerAvailabilitySetResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.SystemCenterVirtualMachineManager.AvailabilitySets

			var model SystemCenterVirtualMachineManagerAvailabilitySetModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			scvmmServerId, err := vmmservers.ParseVMmServerID(model.SystemCenterVirtualMachineManagerServerId)
			if err != nil {
				return err
			}

			id := availabilitysets.NewAvailabilitySetID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := availabilitysets.AvailabilitySet{
				Location: location.Normalize(model.Location),
				ExtendedLocation: availabilitysets.ExtendedLocation{
					Type: utils.String("customLocation"),
					Name: utils.String(model.CustomLocationId),
				},
				Properties: &availabilitysets.AvailabilitySetProperties{
					AvailabilitySetName: utils.String(id.AvailabilitySetName),
					VMmServerId:         utils.String(scvmmServerId.ID()),
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

func (r SystemCenterVirtualMachineManagerAvailabilitySetResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.AvailabilitySets

			id, err := availabilitysets.ParseAvailabilitySetID(metadata.ResourceData.Id())
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

			state := SystemCenterVirtualMachineManagerAvailabilitySetModel{}
			if model := resp.Model; model != nil {
				state.Name = id.AvailabilitySetName
				state.ResourceGroupName = id.ResourceGroupName
				state.Location = location.Normalize(model.Location)
				state.CustomLocationId = pointer.From(model.ExtendedLocation.Name)
				state.Tags = pointer.From(model.Tags)

				scvmmServerId, err := vmmservers.ParseVMmServerIDInsensitively(pointer.From(model.Properties.VMmServerId))
				if err != nil {
					return err
				}
				state.SystemCenterVirtualMachineManagerServerId = scvmmServerId.ID()
			}

			return metadata.Encode(&state)
		},
	}
}

func (r SystemCenterVirtualMachineManagerAvailabilitySetResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.AvailabilitySets

			id, err := availabilitysets.ParseAvailabilitySetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SystemCenterVirtualMachineManagerAvailabilitySetModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := availabilitysets.AvailabilitySetTagsUpdate{}

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

func (r SystemCenterVirtualMachineManagerAvailabilitySetResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.AvailabilitySets

			id, err := availabilitysets.ParseAvailabilitySetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			opts := availabilitysets.DefaultDeleteOperationOptions()
			opts.Force = pointer.To(availabilitysets.ForceDeleteTrue)
			if err := client.DeleteThenPoll(ctx, *id, opts); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
