// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_workspace

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/monitor/2023-04-03/azuremonitorworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Monitor.WorkspacesClient

			id, err := azuremonitorworkspaces.ParseAccountID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model WorkspaceResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}

			if metadata.ResourceData.HasChange("tags") {
				resp.Model.Tags = pointer.To(model.Tags)
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				if properties.Properties == nil {
					properties.Properties = &azuremonitorworkspaces.AzureMonitorWorkspace{}
				}

				publicNetworkAccess := azuremonitorworkspaces.PublicNetworkAccessEnabled
				if !model.PublicNetworkAccessEnabled {
					publicNetworkAccess = azuremonitorworkspaces.PublicNetworkAccessDisabled
				}

				resp.Model.Properties.PublicNetworkAccess = pointer.To(publicNetworkAccess)
			}

			if _, err := client.Create(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
