package migration

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TableV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		Type:    tableSchemaForV0().CoreConfigSchema().ImpliedType(),
		Upgrade: tableUpgradeV0ToV1,
		Version: 0,
	}
}

func tableSchemaForV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"throughput": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func tableUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	oldId := rawState["id"].(string)
	newId := strings.Replace(rawState["id"].(string), "apis/table/tables", "tables", 1)

	log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

	rawState["id"] = newId

	return rawState, nil
}
