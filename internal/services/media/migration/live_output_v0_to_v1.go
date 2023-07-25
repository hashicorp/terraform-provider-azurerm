// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/media/2022-08-01/liveoutputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = LiveOutputV0ToV1{}

type LiveOutputV0ToV1 struct {
}

func (LiveOutputV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"live_event_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"archive_window_duration": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"asset_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"hls_fragments_per_ts_segment": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			ForceNew: true,
		},

		"manifest_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"output_snap_time_in_seconds": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			ForceNew: true,
		},
	}
}

func (LiveOutputV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldIdRaw := rawState["id"].(string)
		oldId, err := liveoutputs.ParseLiveOutputIDInsensitively(oldIdRaw)
		if err != nil {
			return nil, err
		}

		newId := oldId.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q..", oldIdRaw, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}
