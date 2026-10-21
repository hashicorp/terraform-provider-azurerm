// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_manager

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkmanagers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.NetworkManager.NetworkManagers
			id, err := networkmanagers.ParseNetworkManagerID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}
			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *id)
			}
			if resp.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: model properties was nil", *id)
			}

			properties := resp.Model.Properties
			var description string
			var scope []ManagerScopeModel
			var scopeAccesses []string
			if properties.Description != nil {
				description = *properties.Description
			}
			scope = flattenNetworkManagerScope(properties.NetworkManagerScopes)
			scopeAccesses = flattenNetworkManagerScopeAccesses(properties.NetworkManagerScopeAccesses)

			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
				return err
			}

			return metadata.Encode(&ManagerModel{
				CrossTenantScopes: flattenNetworkManagerCrossTenantScopes(properties.NetworkManagerScopes.CrossTenantScopes),
				Description:       description,
				Location:          location.NormalizeNilable(resp.Model.Location),
				Name:              id.NetworkManagerName,
				ResourceGroupName: id.ResourceGroupName,
				ScopeAccesses:     scopeAccesses,
				Scope:             scope,
				Tags:              pluginsdk.FlattenPtrMapStringString(resp.Model.Tags),
			})
		},
	}
}
