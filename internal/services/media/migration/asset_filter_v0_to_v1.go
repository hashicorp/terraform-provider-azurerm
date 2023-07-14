// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/assetsandassetfilters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = AssetFilterV0ToV1{}

type AssetFilterV0ToV1 struct {
}

func (AssetFilterV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"asset_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"first_quality_bitrate": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"presentation_time_range": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"end_in_units": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						AtLeastOneOf: []string{
							"presentation_time_range.0.end_in_units", "presentation_time_range.0.force_end", "presentation_time_range.0.live_backoff_in_units",
							"presentation_time_range.0.presentation_window_in_units", "presentation_time_range.0.start_in_units", "presentation_time_range.0.unit_timescale_in_miliseconds",
						},
					},

					"force_end": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						AtLeastOneOf: []string{
							"presentation_time_range.0.end_in_units", "presentation_time_range.0.force_end", "presentation_time_range.0.live_backoff_in_units",
							"presentation_time_range.0.presentation_window_in_units", "presentation_time_range.0.start_in_units", "presentation_time_range.0.unit_timescale_in_miliseconds",
						},
					},

					"live_backoff_in_units": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						AtLeastOneOf: []string{
							"presentation_time_range.0.end_in_units", "presentation_time_range.0.force_end", "presentation_time_range.0.live_backoff_in_units",
							"presentation_time_range.0.presentation_window_in_units", "presentation_time_range.0.start_in_units", "presentation_time_range.0.unit_timescale_in_miliseconds",
						},
					},

					"presentation_window_in_units": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						AtLeastOneOf: []string{
							"presentation_time_range.0.end_in_units", "presentation_time_range.0.force_end", "presentation_time_range.0.live_backoff_in_units",
							"presentation_time_range.0.presentation_window_in_units", "presentation_time_range.0.start_in_units", "presentation_time_range.0.unit_timescale_in_miliseconds",
						},
					},

					"start_in_units": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						AtLeastOneOf: []string{
							"presentation_time_range.0.end_in_units", "presentation_time_range.0.force_end", "presentation_time_range.0.live_backoff_in_units",
							"presentation_time_range.0.presentation_window_in_units", "presentation_time_range.0.start_in_units", "presentation_time_range.0.unit_timescale_in_miliseconds",
						},
					},

					"unit_timescale_in_miliseconds": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						AtLeastOneOf: []string{
							"presentation_time_range.0.end_in_units", "presentation_time_range.0.force_end", "presentation_time_range.0.live_backoff_in_units",
							"presentation_time_range.0.presentation_window_in_units", "presentation_time_range.0.start_in_units", "presentation_time_range.0.unit_timescale_in_miliseconds",
						},
					},
				},
			},
		},

		"track_selection": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					//lintignore:XS003
					"condition": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"operation": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"property": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},

								"value": {
									Type:     pluginsdk.TypeString,
									Optional: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (AssetFilterV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldIdRaw := rawState["id"].(string)
		oldId, err := assetsandassetfilters.ParseAssetFilterIDInsensitively(oldIdRaw)
		if err != nil {
			return nil, err
		}

		newId := oldId.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q..", oldIdRaw, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
