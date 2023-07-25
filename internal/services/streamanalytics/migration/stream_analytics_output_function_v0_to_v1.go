// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StreamAnalyticsOutputFunctionV0ToV1 struct{}

func (s StreamAnalyticsOutputFunctionV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"stream_analytics_job_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"function_app": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"function_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"api_key": {
			Type:      pluginsdk.TypeString,
			Required:  true,
			Sensitive: true,
		},

		"batch_max_in_bytes": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"batch_max_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},
	}
}

func (s StreamAnalyticsOutputFunctionV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
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
