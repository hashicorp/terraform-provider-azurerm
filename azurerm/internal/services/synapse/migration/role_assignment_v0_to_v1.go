package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = RoleAssignmentV0ToV1{}

type RoleAssignmentV0ToV1 struct{}

func (RoleAssignmentV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*schema.Schema{
		"synapse_workspace_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"principal_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"role_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}

func (RoleAssignmentV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		log.Printf("[DEBUG] Migrating for Synapse Role Assignment")

		name := rawState["role_name"].(string)
		switch name {
		case "Workspace Admin":
			name = "Synapse Administrator"
		case "Apache Spark Admin":
			name = "Synapse Spark Administrator"
		case "Sql Admin":
			name = "Synapse SQL Administrator"
		}
		rawState["role_name"] = name

		return rawState, nil
	}
}
