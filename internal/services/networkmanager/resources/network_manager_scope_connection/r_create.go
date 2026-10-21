// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager_scope_connection

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/scopeconnections"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ManagerScopeConnectionModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.NetworkManager.ScopeConnections
			networkManagerId, err := scopeconnections.ParseNetworkManagerID(model.NetworkManagerId)
			if err != nil {
				return err
			}

			id := scopeconnections.NewScopeConnectionID(networkManagerId.SubscriptionId, networkManagerId.ResourceGroupName, networkManagerId.NetworkManagerName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			scopeConnection := scopeconnections.ScopeConnection{
				Properties: &scopeconnections.ScopeConnectionProperties{},
			}

			if model.Description != "" {
				scopeConnection.Properties.Description = &model.Description
			}

			if model.ResourceId != "" {
				scopeConnection.Properties.ResourceId = &model.ResourceId
			}

			if model.TenantId != "" {
				scopeConnection.Properties.TenantId = &model.TenantId
			}

			if _, err := client.CreateOrUpdate(ctx, id, scopeConnection); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
