package systemcentervirtualmachinemanager

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/clouds"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type SystemCenterVirtualMachineManagerCloudModel struct {
	Name                                                   string            `tfschema:"name"`
	Location                                               string            `tfschema:"location"`
	ResourceGroupName                                      string            `tfschema:"resource_group_name"`
	CustomLocationId                                       string            `tfschema:"custom_location_id"`
	SystemCenterVirtualMachineManagerServerId              string            `tfschema:"system_center_virtual_machine_manager_server_id"`
	SystemCenterVirtualMachineManagerServerInventoryItemId string            `tfschema:"system_center_virtual_machine_manager_server_inventory_item_id"`
	Uuid                                                   string            `tfschema:"uuid"`
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

		"custom_location_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"system_center_virtual_machine_manager_server_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"system_center_virtual_machine_manager_server_inventory_item_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"uuid": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.IsUUID,
		},

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
				ExtendedLocation: clouds.ExtendedLocation{
					Type: utils.String("customLocation"),
					Name: utils.String(model.CustomLocationId),
				},
				Location: location.Normalize(model.Location),
				Properties: clouds.CloudProperties{
					InventoryItemId: utils.String(model.SystemCenterVirtualMachineManagerServerInventoryItemId),
					Uuid:            utils.String(model.Uuid),
					VMmServerId:     utils.String(model.SystemCenterVirtualMachineManagerServerId),
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

			state := SystemCenterVirtualMachineManagerCloudModel{}
			if model := resp.Model; model != nil {
				state.Name = id.CloudName
				state.Location = location.Normalize(model.Location)
				state.ResourceGroupName = id.ResourceGroupName
				state.CustomLocationId = pointer.From(model.ExtendedLocation.Name)
				state.Tags = pointer.From(model.Tags)

				if v := model.Properties.Uuid; v != nil {
					state.Uuid = pointer.From(v)
				}

				if v := model.Properties.InventoryItemId; v != nil {
					state.SystemCenterVirtualMachineManagerServerInventoryItemId = pointer.From(v)
				}

				if v := model.Properties.VMmServerId; v != nil {
					state.SystemCenterVirtualMachineManagerServerId = pointer.From(v)
				}
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

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			parameters := existing.Model
			if parameters == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			if metadata.ResourceData.HasChange("custom_location_id") {
				parameters.ExtendedLocation = clouds.ExtendedLocation{
					Type: utils.String("customLocation"),
					Name: utils.String(model.CustomLocationId),
				}
			}

			if metadata.ResourceData.HasChange("system_center_virtual_machine_manager_server_id") {
				parameters.Properties.VMmServerId = pointer.To(model.SystemCenterVirtualMachineManagerServerId)
			}

			if metadata.ResourceData.HasChange("system_center_virtual_machine_manager_server_inventory_item_id") {
				parameters.Properties.InventoryItemId = pointer.To(model.SystemCenterVirtualMachineManagerServerInventoryItemId)
			}

			if metadata.ResourceData.HasChange("uuid") {
				parameters.Properties.Uuid = pointer.To(model.Uuid)
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Tags = pointer.To(model.Tags)
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *parameters); err != nil {
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

			if err := client.DeleteThenPoll(ctx, *id, clouds.DeleteOperationOptions{Force: pointer.To(clouds.ForceTrue)}); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
