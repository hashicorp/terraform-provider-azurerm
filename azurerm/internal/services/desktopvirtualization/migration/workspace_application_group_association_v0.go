package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/desktopvirtualization/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = WorkspaceApplicationGroupAssociationV0ToV1{}

type WorkspaceApplicationGroupAssociationV0ToV1 struct{}

func (WorkspaceApplicationGroupAssociationV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*schema.Schema{
		"workspace_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"application_group_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}

func (WorkspaceApplicationGroupAssociationV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)

		id, err := parse.WorkspaceApplicationGroupAssociationIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}
		newId := id.ID()

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId

		oldApplicationGroupId := rawState["application_group_id"].(string)
		newApplicationGroupId := id.ApplicationGroup.ID()
		log.Printf("[DEBUG] Updating Application Group ID from %q to %q", oldApplicationGroupId, newApplicationGroupId)
		rawState["application_group_id"] = newApplicationGroupId

		oldWorkspaceId := rawState["workspace_id"].(string)
		newWorkspaceId := id.Workspace.ID()
		log.Printf("[DEBUG] Updating Workspace ID from %q to %q", oldWorkspaceId, newWorkspaceId)
		rawState["workspace_id"] = newWorkspaceId

		return rawState, nil
	}
}
