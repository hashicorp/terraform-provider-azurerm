package migration

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func TableV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		// this should have been applied from pre-0.12 migration system; backporting just in-case
		Type:    tableSchemaV0AndV1().CoreConfigSchema().ImpliedType(),
		Upgrade: tableStateUpgradeV0ToV1,
		Version: 0,
	}
}

func TableV1ToV2() schema.StateUpgrader {
	return schema.StateUpgrader{
		// this should have been applied from pre-0.12 migration system; backporting just in-case
		Type:    tableSchemaV0AndV1().CoreConfigSchema().ImpliedType(),
		Upgrade: tableStateUpgradeV1ToV2,
		Version: 1,
	}
}

// the schema schema was used for both V0 and V1
func tableSchemaV0AndV1() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"acl": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"access_policy": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start": {
										Type:     schema.TypeString,
										Required: true,
									},
									"expiry": {
										Type:     schema.TypeString,
										Required: true,
									},
									"permissions": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func tableStateUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	tableName := rawState["name"].(string)
	accountName := rawState["storage_account_name"].(string)
	environment := meta.(*clients.Client).Account.Environment

	id := rawState["id"].(string)
	newResourceID := fmt.Sprintf("https://%s.table.%s/%s", accountName, environment.StorageEndpointSuffix, tableName)
	log.Printf("[DEBUG] Updating ID from %q to %q", id, newResourceID)

	rawState["id"] = newResourceID
	return rawState, nil
}

func tableStateUpgradeV1ToV2(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	tableName := rawState["name"].(string)
	accountName := rawState["storage_account_name"].(string)
	environment := meta.(*clients.Client).Account.Environment

	id := rawState["id"].(string)
	newResourceID := fmt.Sprintf("https://%s.table.%s/Tables('%s')", accountName, environment.StorageEndpointSuffix, tableName)
	log.Printf("[DEBUG] Updating ID from %q to %q", id, newResourceID)

	rawState["id"] = newResourceID
	return rawState, nil
}
