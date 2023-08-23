// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/inputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StreamAnalyticsReferenceInputBlobV0ToV1 struct{}

func (s StreamAnalyticsReferenceInputBlobV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"date_format": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"path_pattern": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"storage_account_key": {
			Type:      pluginsdk.TypeString,
			Required:  true,
			Sensitive: true,
		},

		"storage_account_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"storage_container_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"time_format": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"serialization": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"field_delimiter": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"encoding": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"authentication_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
	}
}

func (s StreamAnalyticsReferenceInputBlobV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
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
