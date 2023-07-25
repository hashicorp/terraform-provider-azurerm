// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SentinelAlertRuleTemplateV0ToV1 struct{}

func (s SentinelAlertRuleTemplateV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"log_analytics_workspace_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"scheduled_template": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"description": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"tactics": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"severity": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"query": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"query_frequency": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"query_period": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"trigger_operator": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"trigger_threshold": {
						Type:     pluginsdk.TypeInt,
						Computed: true,
					},
				},
			},
		},

		"security_incident_template": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"description": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"product_filter": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},

		"nrt_template": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"description": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"tactics": {
						Type:     pluginsdk.TypeList,
						Computed: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"severity": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
					"query": {
						Type:     pluginsdk.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

func (s SentinelAlertRuleTemplateV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := parse.SentinelAlertRuleTemplateIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
