package loganalytics

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
)

func WorkspaceMigrateState(v int, is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch v {
	case 0:
		log.Println("[INFO] Migrating Log Analytics Workspace State V0 to V1")
		return workspaceStateMigrateV0toV1(is, meta)
	default:
		return is, fmt.Errorf("Unexpected schema version: %d", v)
	}
}

func workspaceStateMigrateV0toV1(is *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	if is.Empty() {
		log.Println("[DEBUG] Empty InstanceState; nothing to migrate.")
		return is, nil
	}

	log.Printf("[DEBUG] Migrating IDs to correct casing for Log Analytics Workspace")
	name := is.Attributes["name"]
	resourceGroup := is.Attributes["resource_group"]
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	id := parse.NewLogAnalyticsWorkspaceID(name, resourceGroup)
	is.ID = id.ID(subscriptionId)

	return is, nil
}
