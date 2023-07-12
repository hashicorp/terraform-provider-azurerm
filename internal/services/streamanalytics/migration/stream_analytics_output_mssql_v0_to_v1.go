// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StreamAnalyticsOutputMsSqlV0ToV1 struct{}

func (s StreamAnalyticsOutputMsSqlV0ToV1) Schema() map[string]*pluginsdk.Schema {
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
			ForceNew: true,
		},

		"database": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"table": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"user": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"password": {
			Type:      pluginsdk.TypeString,
			Required:  true,
			Sensitive: true,
		},

		"max_batch_count": {
			Type:     pluginsdk.TypeFloat,
			Optional: true,
		},

		"max_writer_count": {
			Type:     pluginsdk.TypeFloat,
			Optional: true,
		},

		"authentication_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (s StreamAnalyticsOutputMsSqlV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
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
