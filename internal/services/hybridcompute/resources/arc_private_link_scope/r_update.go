// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_private_link_scope

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/privatelinkscopes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (a Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.HybridCompute.PrivateLinkScopesClient

			id, err := privatelinkscopes.ParseProviderPrivateLinkScopeID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model PrivateLinkScopeModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := existing.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("public_network_access_enabled") {
				publicNetwork := privatelinkscopes.PublicNetworkAccessTypeDisabled
				if model.PublicNetworkAccessEnabled {
					publicNetwork = privatelinkscopes.PublicNetworkAccessTypeEnabled
				}
				properties.Properties.PublicNetworkAccess = &publicNetwork
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
