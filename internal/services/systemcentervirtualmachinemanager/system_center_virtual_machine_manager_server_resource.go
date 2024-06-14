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
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/inventoryitems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/vmmservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SystemCenterVirtualMachineManagerServerModel struct {
	Name              string            `tfschema:"name"`
	ResourceGroupName string            `tfschema:"resource_group_name"`
	Location          string            `tfschema:"location"`
	CustomLocationId  string            `tfschema:"custom_location_id"`
	Fqdn              string            `tfschema:"fqdn"`
	Username          string            `tfschema:"username"`
	Password          string            `tfschema:"password"`
	Port              int64             `tfschema:"port"`
	Tags              map[string]string `tfschema:"tags"`
}

var _ sdk.Resource = SystemCenterVirtualMachineManagerServerResource{}
var _ sdk.ResourceWithUpdate = SystemCenterVirtualMachineManagerServerResource{}

type SystemCenterVirtualMachineManagerServerResource struct{}

func (r SystemCenterVirtualMachineManagerServerResource) ModelObject() interface{} {
	return &SystemCenterVirtualMachineManagerServerModel{}
}

func (r SystemCenterVirtualMachineManagerServerResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return vmmservers.ValidateVMmServerID
}

func (r SystemCenterVirtualMachineManagerServerResource) ResourceType() string {
	return "azurerm_system_center_virtual_machine_manager_server"
}

func (r SystemCenterVirtualMachineManagerServerResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SystemCenterVirtualMachineManagerServerName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"custom_location_id": commonschema.ResourceIDReferenceRequiredForceNew(&customlocations.CustomLocationId{}),

		"fqdn": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"username": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"password": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"port": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IntBetween(1, 65535),
		},

		"tags": commonschema.Tags(),
	}
}

func (r SystemCenterVirtualMachineManagerServerResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r SystemCenterVirtualMachineManagerServerResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.SystemCenterVirtualMachineManager.VMmServers

			var model SystemCenterVirtualMachineManagerServerModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := vmmservers.NewVMmServerID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := &vmmservers.VMmServer{
				Location: location.Normalize(model.Location),
				ExtendedLocation: vmmservers.ExtendedLocation{
					Type: pointer.To("customLocation"),
					Name: pointer.To(model.CustomLocationId),
				},
				Properties: &vmmservers.VMmServerProperties{
					Credentials: &vmmservers.VMmCredential{
						Username: pointer.To(model.Username),
						Password: pointer.To(model.Password),
					},
					Fqdn: model.Fqdn,
				},
				Tags: pointer.To(model.Tags),
			}

			if v := model.Port; v != 0 {
				parameters.Properties.Port = pointer.To(v)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, id, *parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			// After System Center Virtual Machine Manager Server is created, it needs some time to sync the Inventory Items. And service team confirmed that the sync would definitely be completed within 10 minutes. In case, so we need to set a timeout of 120 minutes and check the inventory quantity continuously every minute for 10 times. If the quantity doesn't change, then we consider the sync to be complete.
			stateConf := &pluginsdk.StateChangeConf{
				Delay:        5 * time.Second,
				Pending:      []string{"SyncNotCompleted"},
				Target:       []string{"SyncCompleted"},
				Refresh:      systemCenterVirtualMachineManagerServerStateRefreshFunc(ctx, metadata, id),
				PollInterval: 1 * time.Minute,
				Timeout:      120 * time.Minute,
			}

			if _, err = stateConf.WaitForStateContext(ctx); err != nil {
				return fmt.Errorf("waiting for %s to become available: %s", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r SystemCenterVirtualMachineManagerServerResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VMmServers

			id, err := vmmservers.ParseVMmServerID(metadata.ResourceData.Id())
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

			state := SystemCenterVirtualMachineManagerServerModel{}
			if model := resp.Model; model != nil {
				state.Name = id.VmmServerName
				state.ResourceGroupName = id.ResourceGroupName
				state.Location = location.Normalize(model.Location)
				state.CustomLocationId = pointer.From(model.ExtendedLocation.Name)
				state.Fqdn = model.Properties.Fqdn
				state.Password = metadata.ResourceData.Get("password").(string)
				state.Port = pointer.From(model.Properties.Port)
				state.Tags = pointer.From(model.Tags)

				if v := model.Properties.Credentials; v != nil {
					state.Username = pointer.From(v.Username)
				}

			}

			return metadata.Encode(&state)
		},
	}
}

func (r SystemCenterVirtualMachineManagerServerResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VMmServers

			id, err := vmmservers.ParseVMmServerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SystemCenterVirtualMachineManagerServerModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := vmmservers.VMmServerTagsUpdate{
				Tags: pointer.To(model.Tags),
			}

			if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r SystemCenterVirtualMachineManagerServerResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 180 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VMmServers

			id, err := vmmservers.ParseVMmServerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			opts := vmmservers.DefaultDeleteOperationOptions()
			opts.Force = pointer.To(vmmservers.ForceDeleteTrue)
			if err := client.DeleteThenPoll(ctx, *id, opts); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func systemCenterVirtualMachineManagerServerStateRefreshFunc(ctx context.Context, metadata sdk.ResourceMetaData, id vmmservers.VMmServerId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		client := metadata.Client.SystemCenterVirtualMachineManager.InventoryItems
		scvmmServerId := inventoryitems.NewVMmServerID(id.SubscriptionId, id.ResourceGroupName, id.VmmServerName)
		checkTimes := 10
		lastInventoryItemCount := 0

		for i := 0; i < checkTimes; i++ {
			resp, err := client.ListByVMmServer(ctx, scvmmServerId)
			if err != nil {
				return nil, "", fmt.Errorf("polling for %s: %+v", id, err)
			}

			if model := resp.Model; model != nil {
				currentInventoryItemCount := len(pointer.From(model))

				if i == 0 {
					lastInventoryItemCount = currentInventoryItemCount
					continue
				}

				if currentInventoryItemCount != lastInventoryItemCount {
					return "SyncNotCompleted", "SyncNotCompleted", nil
				}

				time.Sleep(1 * time.Second) // avoid checking too quickly
			}
		}

		return "SyncCompleted", "SyncCompleted", nil
	}
}
