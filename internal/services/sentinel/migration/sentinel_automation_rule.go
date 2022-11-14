package migration

import (
	"context"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = SentinelAutomationRuleV0ToV1{}

type SentinelAutomationRuleV0ToV1 struct{}

func (SentinelAutomationRuleV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"action_incident": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"classification": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"classification_comment": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"labels": {
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Optional: true,
					Type:     pluginsdk.TypeList,
				},
				"order": {
					Required: true,
					Type:     pluginsdk.TypeInt,
				},
				"owner_id": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"severity": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
				"status": {
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"action_playbook": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"logic_app_id": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"order": {
					Required: true,
					Type:     pluginsdk.TypeInt,
				},
				"tenant_id": {
					Computed: true,
					Optional: true,
					Type:     pluginsdk.TypeString,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"condition": {
			Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
				"operator": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"property": {
					Required: true,
					Type:     pluginsdk.TypeString,
				},
				"values": {
					Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
					Required: true,
					Type:     pluginsdk.TypeList,
				},
			}},
			Optional: true,
			Type:     pluginsdk.TypeList,
		},
		"display_name": {
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"enabled": {
			Optional: true,
			Type:     pluginsdk.TypeBool,
		},
		"expiration": {
			Optional: true,
			Type:     pluginsdk.TypeString,
		},
		"log_analytics_workspace_id": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"name": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"order": {
			Required: true,
			Type:     pluginsdk.TypeInt,
		},
	}
}

func (SentinelAutomationRuleV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// Changing the property name "condition" to "condition_property"
		if v, ok := rawState["condition"]; ok {
			rawState["condition_property"] = v
			delete(rawState, "condition")
		}
		return rawState, nil
	}
}
