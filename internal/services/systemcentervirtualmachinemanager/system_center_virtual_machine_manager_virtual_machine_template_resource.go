package systemcentervirtualmachinemanager

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/virtualmachinetemplates"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SystemCenterVirtualMachineManagerVirtualMachineTemplateModel struct {
	Name                                                   string            `tfschema:"name"`
	Location                                               string            `tfschema:"location"`
	ResourceGroupName                                      string            `tfschema:"resource_group_name"`
	CustomLocationId                                       string            `tfschema:"custom_location_id"`
	SystemCenterVirtualMachineManagerServerId              string            `tfschema:"system_center_virtual_machine_manager_server_id"`
	SystemCenterVirtualMachineManagerServerInventoryItemId string            `tfschema:"system_center_virtual_machine_manager_server_inventory_item_id"`
	Uuid                                                   string            `tfschema:"uuid"`
	Tags                                                   map[string]string `tfschema:"tags"`
}

var _ sdk.Resource = SystemCenterVirtualMachineManagerVirtualMachineTemplateResource{}
var _ sdk.ResourceWithUpdate = SystemCenterVirtualMachineManagerVirtualMachineTemplateResource{}

type SystemCenterVirtualMachineManagerVirtualMachineTemplateResource struct{}

func (r SystemCenterVirtualMachineManagerVirtualMachineTemplateResource) ModelObject() interface{} {
	return &SystemCenterVirtualMachineManagerVirtualMachineTemplateModel{}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineTemplateResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return virtualmachinetemplates.ValidateVirtualMachineTemplateID
}

func (r SystemCenterVirtualMachineManagerVirtualMachineTemplateResource) ResourceType() string {
	return "azurerm_system_center_virtual_machine_manager_virtual_machine_template"
}

func (r SystemCenterVirtualMachineManagerVirtualMachineTemplateResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SystemCenterVirtualMachineManagerVirtualMachineTemplateName,
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

func (r SystemCenterVirtualMachineManagerVirtualMachineTemplateResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r SystemCenterVirtualMachineManagerVirtualMachineTemplateResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.SystemCenterVirtualMachineManager.VirtualMachineTemplates

			var model SystemCenterVirtualMachineManagerVirtualMachineTemplateModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := virtualmachinetemplates.NewVirtualMachineTemplateID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := virtualmachinetemplates.VirtualMachineTemplate{
				ExtendedLocation: virtualmachinetemplates.ExtendedLocation{
					Type: utils.String("customLocation"),
					Name: utils.String(model.CustomLocationId),
				},
				Location: location.Normalize(model.Location),
				Properties: virtualmachinetemplates.VirtualMachineTemplateProperties{
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

func (r SystemCenterVirtualMachineManagerVirtualMachineTemplateResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VirtualMachineTemplates

			id, err := virtualmachinetemplates.ParseVirtualMachineTemplateID(metadata.ResourceData.Id())
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

			state := SystemCenterVirtualMachineManagerVirtualMachineTemplateModel{}
			if model := resp.Model; model != nil {
				state.Name = id.VirtualMachineTemplateName
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

func (r SystemCenterVirtualMachineManagerVirtualMachineTemplateResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VirtualMachineTemplates

			id, err := virtualmachinetemplates.ParseVirtualMachineTemplateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SystemCenterVirtualMachineManagerVirtualMachineTemplateModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving existing %s: %+v", *id, err)
			}

			parameters := existing.Model
			if parameters == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			if metadata.ResourceData.HasChange("custom_location_id") {
				parameters.ExtendedLocation = virtualmachinetemplates.ExtendedLocation{
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

func (r SystemCenterVirtualMachineManagerVirtualMachineTemplateResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.VirtualMachineTemplates

			id, err := virtualmachinetemplates.ParseVirtualMachineTemplateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, virtualmachinetemplates.DeleteOperationOptions{Force: pointer.To(virtualmachinetemplates.ForceTrue)}); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
