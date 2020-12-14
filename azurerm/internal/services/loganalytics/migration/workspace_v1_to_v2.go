package migration

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
)

func WorkspaceV1ToV2() schema.StateUpgrader {
	// V1 to V2 is the same as v1 to v0 - to workaround a historical issue where `resource_group` was
	// used in place of `resource_group_name` - ergo using the same schema is fine.
	return schema.StateUpgrader{
		Version: 1,
		Type:    workspaceV0V1Schema().CoreConfigSchema().ImpliedType(),
		Upgrade: workspaceUpgradeV1ToV2,
	}
}

func workspaceUpgradeV1ToV2(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	log.Printf("[DEBUG] Migrating IDs to correct casing for Log Analytics Workspace")
	name := rawState["name"].(string)
	resourceGroup := rawState["resource_group_name"].(string)
	id := parse.NewLogAnalyticsWorkspaceID(subscriptionId, resourceGroup, name)

	rawState["id"] = id.ID()
	return rawState, nil
}
