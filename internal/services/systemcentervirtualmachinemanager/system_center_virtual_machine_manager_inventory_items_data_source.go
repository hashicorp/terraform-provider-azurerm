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
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/inventoryitems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/vmmservers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type SystemCenterVirtualMachineManagerInventoryItemsDataSource struct{}

var _ sdk.DataSource = SystemCenterVirtualMachineManagerInventoryItemsDataSource{}

type SystemCenterVirtualMachineManagerInventoryItemsDataSourceModel struct {
	SystemCenterVirtualMachineManagerServerId string          `tfschema:"system_center_virtual_machine_manager_server_id"`
	InventoryItems                            []InventoryItem `tfschema:"inventory_items"`
	InventoryType                             string          `tfschema:"inventory_type"`
}

type InventoryItem struct {
	id   string `tfschema:"id"`
	name string `tfschema:"name"`
	Uuid string `tfschema:"uuid"`
}

func (l SystemCenterVirtualMachineManagerInventoryItemsDataSource) ResourceType() string {
	return "azurerm_system_center_virtual_machine_manager_inventory_items"
}

func (l SystemCenterVirtualMachineManagerInventoryItemsDataSource) ModelObject() interface{} {
	return &SystemCenterVirtualMachineManagerInventoryItemsDataSourceModel{}
}

func (l SystemCenterVirtualMachineManagerInventoryItemsDataSource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"inventory_type": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(inventoryitems.PossibleValuesForInventoryType(), false),
		},

		"system_center_virtual_machine_manager_server_id": commonschema.ResourceIDReferenceRequired(&vmmservers.VMmServerId{}),
	}
}

func (l SystemCenterVirtualMachineManagerInventoryItemsDataSource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"inventory_items": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"name": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},

					"uuid": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (l SystemCenterVirtualMachineManagerInventoryItemsDataSource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.SystemCenterVirtualMachineManager.InventoryItems

			var state SystemCenterVirtualMachineManagerInventoryItemsDataSourceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			scvmmServerId, err := inventoryitems.ParseVMmServerID(state.SystemCenterVirtualMachineManagerServerId)
			if err != nil {
				return err
			}

			resp, err := client.ListByVMmServerComplete(ctx, *scvmmServerId)
			if err != nil {
				if response.WasNotFound(resp.LatestHttpResponse) {
					return fmt.Errorf("%s was not found", scvmmServerId)
				}
				return fmt.Errorf("reading %s: %+v", scvmmServerId, err)
			}

			if model := resp.Items; model != nil {
				inventoryItems, err := flattenInventoryItems(model, state.InventoryType)
				if err != nil {
					return err
				}
				if len(pointer.From(inventoryItems)) == 0 {
					return fmt.Errorf("no inventory items were found for %s", scvmmServerId)
				}
				state.InventoryItems = pointer.From(inventoryItems)
			}

			metadata.ResourceData.SetId(scvmmServerId.ID())

			return metadata.Encode(&state)
		},
	}
}

func flattenInventoryItems(input []inventoryitems.InventoryItem, inventoryType string) (*[]InventoryItem, error) {
	results := make([]InventoryItem, 0)

	for _, item := range input {
		if props := item.Properties; props != nil {
			inventoryItem := InventoryItem{}

			if v, ok := props.(inventoryitems.CloudInventoryItem); ok && inventoryType == string(inventoryitems.InventoryTypeCloud) {
				// Service API indicates that the static segment `inventoryItems` in the resource ID of the Inventory Item should start with lowercase. See more details from https://github.com/Azure/azure-rest-api-specs/blob/92c409d93f895a30d51603b2fda78a49b3a2cd60/specification/scvmm/resource-manager/Microsoft.ScVmm/stable/2023-10-07/scvmm.json#L1785
				// But the static segment `InventoryItems` in the resource ID of the Inventory Item returned by the API starts with uppercase. So all instances of setting the inventory item ID must use ParseInventoryItemIDInsensitively() in Read() to normalize the resource ID
				scvmmServerInventoryItemId, err := inventoryitems.ParseInventoryItemIDInsensitively(pointer.From(item.Id))
				if err != nil {
					return nil, err
				}
				inventoryItem.id = scvmmServerInventoryItemId.ID()

				inventoryItem.name = pointer.From(v.InventoryItemName)
				inventoryItem.Uuid = pointer.From(v.Uuid)
				results = append(results, inventoryItem)
			} else if v, ok := props.(inventoryitems.VirtualMachineInventoryItem); ok && inventoryType == string(inventoryitems.InventoryTypeVirtualMachine) {
				scvmmServerInventoryItemId, err := inventoryitems.ParseInventoryItemIDInsensitively(pointer.From(item.Id))
				if err != nil {
					return nil, err
				}
				inventoryItem.id = scvmmServerInventoryItemId.ID()

				inventoryItem.name = pointer.From(v.InventoryItemName)
				inventoryItem.Uuid = pointer.From(v.Uuid)
				results = append(results, inventoryItem)
			} else if v, ok := props.(inventoryitems.VirtualMachineTemplateInventoryItem); ok && inventoryType == string(inventoryitems.InventoryTypeVirtualMachineTemplate) {
				scvmmServerInventoryItemId, err := inventoryitems.ParseInventoryItemIDInsensitively(pointer.From(item.Id))
				if err != nil {
					return nil, err
				}
				inventoryItem.id = scvmmServerInventoryItemId.ID()

				inventoryItem.name = pointer.From(v.InventoryItemName)
				inventoryItem.Uuid = pointer.From(v.Uuid)
				results = append(results, inventoryItem)
			} else if v, ok := props.(inventoryitems.VirtualNetworkInventoryItem); ok && inventoryType == string(inventoryitems.InventoryTypeVirtualNetwork) {
				scvmmServerInventoryItemId, err := inventoryitems.ParseInventoryItemIDInsensitively(pointer.From(item.Id))
				if err != nil {
					return nil, err
				}
				inventoryItem.id = scvmmServerInventoryItemId.ID()

				inventoryItem.name = pointer.From(v.InventoryItemName)
				inventoryItem.Uuid = pointer.From(v.Uuid)
				results = append(results, inventoryItem)
			}
		}
	}

	return &results, nil
}
