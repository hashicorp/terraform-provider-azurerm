// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2020-03-01/streamingjobs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StreamAnalyticsJobV0ToV1 struct{}

func (s StreamAnalyticsJobV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"stream_analytics_cluster_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"compatibility_level": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"data_locale": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"events_late_arrival_max_delay_in_seconds": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"events_out_of_order_max_delay_in_seconds": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"events_out_of_order_policy": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"output_error_policy": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"streaming_units": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"content_storage_policy": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"job_storage_account": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"authentication_mode": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"account_name": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"account_key": {
						Type:      pluginsdk.TypeString,
						Required:  true,
						Sensitive: true,
					},
				},
			},
		},

		"transformation_query": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"identity": {
			Type:     schema.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:     schema.TypeString,
						Required: true,
					},
					"principal_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"tenant_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},

		"job_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func (s StreamAnalyticsJobV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := streamingjobs.ParseStreamingJobIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
