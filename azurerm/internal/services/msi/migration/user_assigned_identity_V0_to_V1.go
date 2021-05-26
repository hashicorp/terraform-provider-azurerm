package migration

import (
	"context"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = UserAssignedIdentityV0ToV1{}

type UserAssignedIdentityV0ToV1 struct{}

func (UserAssignedIdentityV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"tags": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"principal_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"client_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (UserAssignedIdentityV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		id, err := parse.UserAssignedIdentityID(oldId)
		if err != nil {
			return rawState, err
		}

		newId := id.ID()
		log.Printf("Updating `id` from %q to %q", oldId, newId)
		rawState["id"] = newId
		return rawState, nil
	}
}
