package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var _ pluginsdk.StateUpgrade = RoleDefinitionV0ToV1{}

type RoleDefinitionV0ToV1 struct{}

func (RoleDefinitionV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*schema.Schema{
		"role_definition_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},

		"name": {
			Type:     schema.TypeString,
			Required: true,
		},

		"scope": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},

		//lintignore:XS003
		"permissions": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"actions": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"not_actions": {
						Type:     schema.TypeList,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					"data_actions": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Set: schema.HashString,
					},
					"not_data_actions": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Set: schema.HashString,
					},
				},
			},
		},

		"assignable_scopes": {
			Type:     schema.TypeList,
			Optional: true,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
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
