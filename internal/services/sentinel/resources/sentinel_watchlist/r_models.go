// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_watchlist

type WatchlistModel struct {
	Name                    string   `tfschema:"name"`
	LogAnalyticsWorkspaceId string   `tfschema:"log_analytics_workspace_id"`
	DisplayName             string   `tfschema:"display_name"`
	Description             string   `tfschema:"description"`
	Labels                  []string `tfschema:"labels"`
	DefaultDuration         string   `tfschema:"default_duration"`
	ItemSearchKey           string   `tfschema:"item_search_key"`
}

func (r Resource) ModelObject() any {
	return &WatchlistModel{}
}
