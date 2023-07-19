// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StreamAnalyticsOutputPowerBiV0ToV1 struct{}

func (s StreamAnalyticsOutputPowerBiV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"stream_analytics_job_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"dataset": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"table": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"group_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"token_user_principal_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"token_user_display_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (s StreamAnalyticsOutputPowerBiV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := outputs.ParseOutputIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
