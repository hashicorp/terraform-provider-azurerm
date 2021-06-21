package migration

import (
	"context"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = RoleAssignmentV0ToV1{}

type RoleAssignmentV0ToV1 struct{}

func (RoleAssignmentV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"synapse_workspace_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"principal_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"role_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}

func (RoleAssignmentV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		log.Printf("[DEBUG] Migrating for Synapse Role Assignment")

		name := rawState["role_name"].(string)
		rawState["role_name"] = MigrateToNewRole(name)

		return rawState, nil
	}
}

func MigrateToNewRole(roleName string) string {
	switch roleName {
	case "Workspace Admin":
		roleName = "Synapse Administrator"
	case "Apache Spark Admin":
		roleName = "Apache Spark Administrator"
	case "Sql Admin":
		roleName = "Synapse SQL Administrator"
	}
	return roleName
}
