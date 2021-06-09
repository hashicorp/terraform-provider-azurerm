package migration

import (
	"context"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = WorkspaceV0ToV1{}

type WorkspaceV0ToV1 struct{}

func (WorkspaceV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return workspaceSchemaForV0AndV1()
}

func (WorkspaceV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		subscriptionId := meta.(*clients.Client).Account.SubscriptionId

		log.Printf("[DEBUG] Migrating IDs to correct casing for Log Analytics Workspace")
		name := rawState["name"].(string)
		resourceGroup := rawState["resource_group_name"].(string)
		id := parse.NewLogAnalyticsWorkspaceID(subscriptionId, resourceGroup, name)

		rawState["id"] = id.ID()
		return rawState, nil
	}
}

func workspaceSchemaForV0AndV1() map[string]*pluginsdk.Schema {
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

		"internet_ingestion_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"internet_query_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"sku": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"retention_in_days": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"daily_quota_gb": {
			Type:     pluginsdk.TypeFloat,
			Optional: true,
			Default:  -1.0,
		},

		"workspace_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"portal_url": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"primary_shared_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_shared_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
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
