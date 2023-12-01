package systemcentervirtualmachinemanager

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/inventoryitems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/vmmservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SystemCenterVirtualMachineManagerServerVirtualMachineInventoryItemModel struct {
	Name                                      string   `tfschema:"name"`
	SystemCenterVirtualMachineManagerServerId string   `tfschema:"system_center_virtual_machine_manager_server_id"`
	IPAddressess                              []string `tfschema:"ip_addresses"`
	CloudInventoryItemId                      string   `tfschema:"cloud_inventory_item_id"`
}

var _ sdk.Resource = SystemCenterVirtualMachineManagerServerVirtualMachineInventoryItemResource{}
var _ sdk.ResourceWithUpdate = SystemCenterVirtualMachineManagerServerVirtualMachineInventoryItemResource{}

type SystemCenterVirtualMachineManagerServerVirtualMachineInventoryItemResource struct{}

func (r SystemCenterVirtualMachineManagerServerVirtualMachineInventoryItemResource) ModelObject() interface{} {
	return &SystemCenterVirtualMachineManagerServerVirtualMachineInventoryItemModel{}
}

func (r SystemCenterVirtualMachineManagerServerVirtualMachineInventoryItemResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return inventoryitems.ValidateInventoryItemID
}

func (r SystemCenterVirtualMachineManagerServerVirtualMachineInventoryItemResource) ResourceType() string {
	return "azurerm_system_center_virtual_machine_manager_server_virtual_machine_inventory_item"
}

func (r SystemCenterVirtualMachineManagerServerVirtualMachineInventoryItemResource) Arguments() map[string]*pluginsdk.Schema {
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

		"ip_addresses": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.IsIPv4Address,
			},
		},

		"cloud_inventory_item_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r SystemCenterVirtualMachineManagerServerVirtualMachineInventoryItemResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r SystemCenterVirtualMachineManagerServerVirtualMachineInventoryItemResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			subscriptionId := metadata.Client.Account.SubscriptionId
			client := metadata.Client.SystemCenterVirtualMachineManager.InventoryItems

			var model SystemCenterVirtualMachineManagerServerVirtualMachineInventoryItemModel
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
				Properties: inventoryitems.VirtualMachineInventoryItem{
					IPAddresses: pointer.To(model.IPAddressess),
				},
			}

			if v := model.CloudInventoryItemId; v != "" {
				cloudInventoryItemId, err := inventoryitems.ParseInventoryItemID(model.CloudInventoryItemId)
				if err != nil {
					return err
				}

				cloudInventoryItem := pointer.To(parameters.Properties.(inventoryitems.VirtualMachineInventoryItem))
				cloudInventoryItem.Cloud = &inventoryitems.InventoryItemDetails{
					InventoryItemId:   utils.String(cloudInventoryItemId.ID()),
					InventoryItemName: utils.String(cloudInventoryItemId.InventoryItemName),
				}
			}

			if _, err := client.Create(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r SystemCenterVirtualMachineManagerServerVirtualMachineInventoryItemResource) Read() sdk.ResourceFunc {
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

			state := SystemCenterVirtualMachineManagerServerVirtualMachineInventoryItemModel{}
			if model := resp.Model; model != nil {
				state.Name = id.InventoryItemName
				state.SystemCenterVirtualMachineManagerServerId = vmmservers.NewVMmServerID(id.SubscriptionId, id.ResourceGroupName, id.VmmServerName).ID()

				if props := model.Properties; props != nil {
					vmInventoryItem := props.(inventoryitems.VirtualMachineInventoryItem)
					state.IPAddressess = pointer.From(vmInventoryItem.IPAddresses)

					if v := vmInventoryItem.Cloud; v != nil {
						state.CloudInventoryItemId = pointer.From(v.InventoryItemId)
					}
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r SystemCenterVirtualMachineManagerServerVirtualMachineInventoryItemResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.InventoryItems

			id, err := inventoryitems.ParseInventoryItemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SystemCenterVirtualMachineManagerServerVirtualMachineInventoryItemModel
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

			vmInventoryItem := pointer.To(parameters.Properties.(inventoryitems.VirtualMachineInventoryItem))

			if metadata.ResourceData.HasChange("ip_addresses") {
				vmInventoryItem.IPAddresses = pointer.To(model.IPAddressess)
			}

			if metadata.ResourceData.HasChange("cloud_inventory_item_id") {
				cloudInventoryItemId, err := inventoryitems.ParseInventoryItemID(model.CloudInventoryItemId)
				if err != nil {
					return err
				}

				vmInventoryItem.Cloud = &inventoryitems.InventoryItemDetails{
					InventoryItemId:   utils.String(cloudInventoryItemId.ID()),
					InventoryItemName: utils.String(cloudInventoryItemId.InventoryItemName),
				}
			}

			if _, err := client.Create(ctx, *id, *parameters); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r SystemCenterVirtualMachineManagerServerVirtualMachineInventoryItemResource) Delete() sdk.ResourceFunc {
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
