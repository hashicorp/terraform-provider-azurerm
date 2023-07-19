// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StreamAnalyticsOutputCosmosDbV0ToV1 struct{}

func (s StreamAnalyticsOutputCosmosDbV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"cosmosdb_account_key": {
			Type:      pluginsdk.TypeString,
			Required:  true,
			Sensitive: true,
		},

		"cosmosdb_sql_database_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"container_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"document_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"partition_key": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (s StreamAnalyticsOutputCosmosDbV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
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
