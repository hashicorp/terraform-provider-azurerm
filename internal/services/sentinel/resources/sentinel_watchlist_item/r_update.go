// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_watchlist_item

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/watchlistitems"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.WatchlistItemsClient

			var model WatchlistItemModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			watchlistId, err := watchlistitems.ParseWatchlistID(model.WatchlistID)
			if err != nil {
				return err
			}
			id := watchlistitems.NewWatchlistItemID(watchlistId.SubscriptionId, watchlistId.ResourceGroupName, watchlistId.WorkspaceName, watchlistId.WatchlistAlias, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			update := watchlistitems.WatchlistItem{}

			if model := existing.Model; model != nil {
				update.Properties = model.Properties
			}
			if metadata.ResourceData.HasChange("properties") {
				if update.Properties == nil {
					update.Properties = &watchlistitems.WatchlistItemProperties{}
				}
				update.Properties.ItemsKeyValue = model.Properties
			}

			if _, err = client.CreateOrUpdate(ctx, id, update); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
