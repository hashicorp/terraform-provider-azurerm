package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = TableV0ToV1{}

type TableV0ToV1 struct{}

func (TableV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return tableSchemaV0AndV1()
}

func (TableV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	// this should have been applied from pre-0.12 migration system; backporting just in-case
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		tableName := rawState["name"].(string)
		accountName := rawState["storage_account_name"].(string)
		environment := meta.(*clients.Client).Account.Environment

		id := rawState["id"].(string)
		newResourceID := fmt.Sprintf("https://%s.table.%s/%s", accountName, environment.StorageEndpointSuffix, tableName)
		log.Printf("[DEBUG] Updating ID from %q to %q", id, newResourceID)

		rawState["id"] = newResourceID
		return rawState, nil
	}
}

var _ pluginsdk.StateUpgrade = TableV1ToV2{}

type TableV1ToV2 struct{}

func (TableV1ToV2) Schema() map[string]*pluginsdk.Schema {
	return tableSchemaV0AndV1()
}

func (TableV1ToV2) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	// this should have been applied from pre-0.12 migration system; backporting just in-case
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		tableName := rawState["name"].(string)
		accountName := rawState["storage_account_name"].(string)
		environment := meta.(*clients.Client).Account.Environment

		id := rawState["id"].(string)
		newResourceID := fmt.Sprintf("https://%s.table.%s/Tables('%s')", accountName, environment.StorageEndpointSuffix, tableName)
		log.Printf("[DEBUG] Updating ID from %q to %q", id, newResourceID)

		rawState["id"] = newResourceID
		return rawState, nil
	}
}

// the schema schema was used for both V0 and V1
func tableSchemaV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"storage_account_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"acl": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"id": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
					"access_policy": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"start": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"expiry": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"permissions": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},
	}
}
