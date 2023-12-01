package systemcentervirtualmachinemanager

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/inventoryitems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/vmmservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SystemCenterVirtualMachineManagerServerVirtualMachineTemplateInventoryItemModel struct {
	Name                                      string `tfschema:"name"`
	SystemCenterVirtualMachineManagerServerId string `tfschema:"system_center_virtual_machine_manager_server_id"`
}

var _ sdk.Resource = SystemCenterVirtualMachineManagerServerVirtualMachineTemplateInventoryItemResource{}
var _ sdk.ResourceWithUpdate = SystemCenterVirtualMachineManagerServerVirtualMachineTemplateInventoryItemResource{}

type SystemCenterVirtualMachineManagerServerVirtualMachineTemplateInventoryItemResource struct{}

func (r SystemCenterVirtualMachineManagerServerVirtualMachineTemplateInventoryItemResource) ModelObject() interface{} {
	return &SystemCenterVirtualMachineManagerServerVirtualMachineTemplateInventoryItemModel{}
}

func (r SystemCenterVirtualMachineManagerServerVirtualMachineTemplateInventoryItemResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return inventoryitems.ValidateInventoryItemID
}

func (r SystemCenterVirtualMachineManagerServerVirtualMachineTemplateInventoryItemResource) ResourceType() string {
	return "azurerm_system_center_virtual_machine_manager_server_virtual_machine_template_inventory_item"
}

func (r SystemCenterVirtualMachineManagerServerVirtualMachineTemplateInventoryItemResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"system_center_virtual_machine_manager_server_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r SystemCenterVirtualMachineManagerServerVirtualMachineTemplateInventoryItemResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r SystemCenterVirtualMachineManagerServerVirtualMachineTemplateInventoryItemResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.SystemCenterVirtualMachineManager.InventoryItems

			var model SystemCenterVirtualMachineManagerServerVirtualMachineTemplateInventoryItemModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			scvmmServerId, err := vmmservers.ParseVMmServerID(model.SystemCenterVirtualMachineManagerServerId)
			if err != nil {
				return err
			}

			id := inventoryitems.NewInventoryItemID(subscriptionId, scvmmServerId.ResourceGroupName, scvmmServerId.VmmServerName, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := inventoryitems.InventoryItem{
				Properties: inventoryitems.VirtualMachineTemplateInventoryItem{},
			}

			if _, err := client.Create(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r SystemCenterVirtualMachineManagerServerVirtualMachineTemplateInventoryItemResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.InventoryItems

			id, err := inventoryitems.ParseInventoryItemID(metadata.ResourceData.Id())
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

			state := SystemCenterVirtualMachineManagerServerVirtualMachineTemplateInventoryItemModel{}
			if model := resp.Model; model != nil {
				state.Name = id.InventoryItemName
				state.SystemCenterVirtualMachineManagerServerId = vmmservers.NewVMmServerID(id.SubscriptionId, id.ResourceGroupName, id.VmmServerName).ID()
			}

			return metadata.Encode(&state)
		},
	}
}

func (r SystemCenterVirtualMachineManagerServerVirtualMachineTemplateInventoryItemResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.InventoryItems

			id, err := inventoryitems.ParseInventoryItemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SystemCenterVirtualMachineManagerServerVirtualMachineTemplateInventoryItemModel
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

			if _, err := client.Create(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r SystemCenterVirtualMachineManagerServerVirtualMachineTemplateInventoryItemResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.InventoryItems

			id, err := inventoryitems.ParseInventoryItemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}
