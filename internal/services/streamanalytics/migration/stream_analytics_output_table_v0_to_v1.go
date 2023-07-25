// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StreamAnalyticsOutputTableV0ToV1 struct{}

func (s StreamAnalyticsOutputTableV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"storage_account_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"storage_account_key": {
			Type:      pluginsdk.TypeString,
			Required:  true,
			Sensitive: true,
		},

		"table": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"partition_key": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"row_key": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"batch_size": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		"columns_to_remove": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (s StreamAnalyticsOutputTableV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
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
