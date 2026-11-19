// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_watchlist_item

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/watchlistitems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/watchlists"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,

		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.WatchlistItemsClient
			id, err := watchlistitems.ParseWatchlistItemID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			watchlistId := watchlists.NewWatchlistID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.WatchlistAlias)

			var properties map[string]any
			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					if itemsKV := props.ItemsKeyValue; itemsKV != nil {
						properties = itemsKV.(map[string]any)
					}
				}
			}

			watchListItem := WatchlistItemModel{
				WatchlistID: watchlistId.ID(),
				Name:        id.WatchlistItemId,
				Properties:  properties,
			}
			return metadata.Encode(&watchListItem)
		},
	}
}
