package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"internet_ingestion_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},

		"internet_query_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},

		"sku": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"retention_in_days": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},

		"daily_quota_gb": {
			Type:     schema.TypeFloat,
			Optional: true,
			Default:  -1.0,
		},

		"workspace_id": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"portal_url": {
			Type:     schema.TypeString,
			Computed: true,
		},

		"primary_shared_key": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"secondary_shared_key": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"tags": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}
