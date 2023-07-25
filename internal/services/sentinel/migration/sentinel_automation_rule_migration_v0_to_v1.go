// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SentinelAutomationRuleV0ToV1 struct{}

func (s SentinelAutomationRuleV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"log_analytics_workspace_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"order": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"expiration": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"condition": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"property": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"operator": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"values": {
						Type:     pluginsdk.TypeList,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"action_incident": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"order": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"status": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"classification": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"classification_comment": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"labels": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"owner_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},

					"severity": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},

		"action_playbook": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"order": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"logic_app_id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"tenant_id": {
						Type: pluginsdk.TypeString,
						// We'll use the current tenant id if this property is absent.
						Optional: true,
						Computed: true,
					},
				},
			},
		},
	}
}

func (s SentinelAutomationRuleV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := parse.AutomationRuleIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
