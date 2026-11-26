// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_workspace

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/monitor/2023-04-03/azuremonitorworkspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model WorkspaceResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.Monitor.WorkspacesClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := azuremonitorworkspaces.NewAccountID(subscriptionId, model.ResourceGroupName, model.Name)
			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			publicNetworkAccess := azuremonitorworkspaces.PublicNetworkAccessEnabled
			if !model.PublicNetworkAccessEnabled {
				publicNetworkAccess = azuremonitorworkspaces.PublicNetworkAccessDisabled
			}

			properties := azuremonitorworkspaces.AzureMonitorWorkspaceResource{
				Location: location.Normalize(model.Location),
				Properties: &azuremonitorworkspaces.AzureMonitorWorkspace{
					PublicNetworkAccess: pointer.To(publicNetworkAccess),
				},
				Tags: &model.Tags,
			}

			if _, err := client.Create(ctx, id, properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id)
		},
	}
}
