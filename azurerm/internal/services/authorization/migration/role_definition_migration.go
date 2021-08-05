package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = RoleDefinitionV0ToV1{}

type RoleDefinitionV0ToV1 struct{}

func (RoleDefinitionV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"role_definition_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"scope": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		//lintignore:XS003
		"permissions": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"actions": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"not_actions": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"data_actions": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Set: pluginsdk.HashString,
					},
					"not_data_actions": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
						Set: pluginsdk.HashString,
					},
				},
			},
		},

		"assignable_scopes": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (RoleDefinitionV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		log.Println("[DEBUG] Migrating ID from v0 to v1 format")

		oldID := rawState["id"].(string)
		if oldID == "" {
			return nil, fmt.Errorf("failed to migrate state: old ID empty")
		}
		scope := rawState["scope"].(string)
		if scope == "" {
			return nil, fmt.Errorf("failed to migrate state: scope missing")
		}

		newID := fmt.Sprintf("%s|%s", oldID, scope)

		rawState["id"] = newID

		return rawState, nil
	}
}
