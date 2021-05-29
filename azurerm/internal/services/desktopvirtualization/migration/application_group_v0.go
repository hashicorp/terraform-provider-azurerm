package migration

import (
	"context"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/desktopvirtualization/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ApplicationGroupV0ToV1{}

type ApplicationGroupV0ToV1 struct{}

func (ApplicationGroupV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"host_pool_id": {
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

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (ApplicationGroupV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		id, err := parse.ApplicationGroupIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}
		newId := id.ID()
		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		oldHostPoolId := rawState["host_pool_id"].(string)
		hostPoolId, err := parse.HostPoolIDInsensitively(oldHostPoolId)
		if err != nil {
			return nil, err
		}
		newHostPoolId := hostPoolId.ID()
		log.Printf("[DEBUG] Updating Host Pool ID from %q to %q", oldHostPoolId, newHostPoolId)
		rawState["host_pool_id"] = newHostPoolId

		return rawState, nil
	}
}
