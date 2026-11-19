// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_watchlist_item

type WatchlistItemModel struct {
	Name        string         `tfschema:"name"`
	WatchlistID string         `tfschema:"watchlist_id"`
	Properties  map[string]any `tfschema:"properties"`
}

func (r Resource) ModelObject() any {
	return &WatchlistItemModel{}
}
