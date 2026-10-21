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

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.ScopeConnections

			id, err := scopeconnections.ParseScopeConnectionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}
			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model properties was nil", *id)
			}

			properties := existing.Model.Properties

			state := ManagerScopeConnectionModel{
				Name:             id.ScopeConnectionName,
				NetworkManagerId: scopeconnections.NewNetworkManagerID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName).ID(),
			}

			if properties.ConnectionState != nil {
				state.ConnectionState = string(*properties.ConnectionState)
			}

			if properties.Description != nil {
				state.Description = *properties.Description
			}

			if properties.ResourceId != nil {
				state.ResourceId = *properties.ResourceId
			}

			if properties.TenantId != nil {
				state.TenantId = *properties.TenantId
			}

			return metadata.Encode(&state)
		},
	}
}
