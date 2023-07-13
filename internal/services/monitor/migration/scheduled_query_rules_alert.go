// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-04-16/scheduledqueryrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ScheduledQueryRulesAlertUpgradeV0ToV1{}

type ScheduledQueryRulesAlertUpgradeV0ToV1 struct{}

func (ScheduledQueryRulesAlertUpgradeV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return scheduledQueryRulesAlertSchemaForV0AndV1()
}

func (ScheduledQueryRulesAlertUpgradeV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.insights/scheduledQueryRules/rule1
		// new:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/scheduledQueryRules/rule1
		oldId := rawState["id"].(string)
		id, err := scheduledqueryrules.ParseScheduledQueryRuleIDInsensitively(oldId)
		if err != nil {
			return rawState, err
		}

		newId := id.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		return rawState, nil
	}
}

func scheduledQueryRulesAlertSchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"authorized_resource_ids": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			MaxItems: 100,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
		"action": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"action_group": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"custom_webhook_payload": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},
					"email_subject": {
						Type:     pluginsdk.TypeString,
						Optional: true,
					},
				},
			},
		},
		"data_source_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"auto_mitigation_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},
		"frequency": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},
		"query": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},
		"query_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},
		"severity": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},
		"throttling": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},
		"time_window": {
			Type:     pluginsdk.TypeInt,
			Required: true,
		},
		"trigger": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"metric_trigger": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"metric_column": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"metric_trigger_type": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"operator": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"threshold": {
									Type:     pluginsdk.TypeFloat,
									Required: true,
								},
							},
						},
					},
					"operator": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"threshold": {
						Type:     pluginsdk.TypeFloat,
						Required: true,
					},
				},
			},
		},

		"tags": tags.Schema(),
	}
}
