package migration

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/maintenance/mgmt/2018-06-01-preview/maintenance"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/maintenance/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ConfigurationV0ToV1{}

type ConfigurationV0ToV1 struct{}

func (ConfigurationV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
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

		"scope": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  string(maintenance.ScopeAll),
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

func (ConfigurationV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		subscriptionId := meta.(*clients.Client).Account.SubscriptionId

		log.Printf("[DEBUG] Migrating IDs to correct casing for Maintenance Configuration")

		name := rawState["name"].(string)
		resourceGroup := rawState["resource_group_name"].(string)
		id := parse.NewMaintenanceConfigurationID(subscriptionId, resourceGroup, name)

		rawState["id"] = id.ID()

		return rawState, nil
	}
}
