// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StreamAnalyticsReferenceInputMsSqlV0ToV1 struct{}

func (s StreamAnalyticsReferenceInputMsSqlV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"server": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"database": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"username": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"password": {
			Type:      pluginsdk.TypeString,
			Required:  true,
			Sensitive: true,
		},

		"refresh_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"refresh_interval_duration": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"full_snapshot_query": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"delta_snapshot_query": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"table": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (s StreamAnalyticsReferenceInputMsSqlV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := inputs.ParseInputIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
