package migration

import (
	"context"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/desktopvirtualization/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ApplicationV0ToV1{}

type ApplicationV0ToV1 struct{}

func (ApplicationV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"application_group_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"friendly_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"file_path": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"command_line_setting": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"command_line_arguments": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"show_in_portal": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"icon_path": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"icon_index": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},
	}
}

func (ApplicationV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		id, err := parse.ApplicationIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}
		newId := id.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		oldApplicationGroupId := rawState["application_group_id"].(string)
		applicationGroupId, err := parse.HostPoolIDInsensitively(oldApplicationGroupId)
		if err != nil {
			return nil, err
		}
		newApplicationGroupId := applicationGroupId.ID()
		log.Printf("[DEBUG] Updating Host Pool ID from %q to %q", oldApplicationGroupId, applicationGroupId)
		rawState["application_group_id"] = newApplicationGroupId

		return rawState, nil
	}
}
