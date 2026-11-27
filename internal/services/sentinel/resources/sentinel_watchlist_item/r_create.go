// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_watchlist_item

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/watchlistitems"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-11-01/watchlists"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Sentinel.WatchlistItemsClient

			var watchListItem WatchlistItemModel
			if err := metadata.Decode(&watchListItem); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			// Generate a random UUID as the resource name if the user doesn't specify it.
			if watchListItem.Name == "" {
				watchListItem.Name = uuid.New().String()
			}

			watchlistId, err := watchlists.ParseWatchlistID(watchListItem.WatchlistID)
			if err != nil {
				return err
			}

			id := watchlistitems.NewWatchlistItemID(watchlistId.SubscriptionId, watchlistId.ResourceGroupName, watchlistId.WorkspaceName, watchlistId.WatchlistAlias, watchListItem.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			params := watchlistitems.WatchlistItem{
				Properties: &watchlistitems.WatchlistItemProperties{
					ItemsKeyValue: watchListItem.Properties,
				},
			}

			if _, err = client.CreateOrUpdate(ctx, id, params); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}
