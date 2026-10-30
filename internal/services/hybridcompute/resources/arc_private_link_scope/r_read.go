// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_private_link_scope

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridcompute/2022-11-10/privatelinkscopes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (a Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.HybridCompute.PrivateLinkScopesClient

			id, err := privatelinkscopes.ParseProviderPrivateLinkScopeID(metadata.ResourceData.Id())
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

			state := PrivateLinkScopeModel{
				Name:              id.PrivateLinkScopeName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := existing.Model; model != nil {
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				publicNetworkAccess := false
				if props := model.Properties; props != nil && props.PublicNetworkAccess != nil {
					publicNetworkAccess = *model.Properties.PublicNetworkAccess == privatelinkscopes.PublicNetworkAccessTypeEnabled
				}
				state.PublicNetworkAccessEnabled = publicNetworkAccess
			}

			return metadata.Encode(&state)
		},
	}
}
