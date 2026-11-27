// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package system_center_virtual_machine_manager_server

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/inventoryitems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07/vmmservers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	systemcentervirtualmachinemanagerClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := systemcentervirtualmachinemanagerClient.From(metadata.Client).VMmServers

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

			return r.flatten(metadata, id, resp.Model)
		},
	}
}

func systemCenterVirtualMachineManagerServerStateRefreshFunc(ctx context.Context, metadata sdk.ResourceMetaData, id vmmservers.VMmServerId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		client := systemcentervirtualmachinemanagerClient.From(metadata.Client).InventoryItems
		scvmmServerId := inventoryitems.NewVMmServerID(id.SubscriptionId, id.ResourceGroupName, id.VmmServerName)
		checkTimes := 10
		lastInventoryItemCount := 0

		for i := range checkTimes {
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
