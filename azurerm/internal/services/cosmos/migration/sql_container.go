package migration

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func SqlContainerV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		Type:    sqlContainerSchemaForV0().CoreConfigSchema().ImpliedType(),
		Upgrade: sqlContainerUpgradeV0ToV1,
		Version: 0,
	}
}

func sqlContainerSchemaForV0() *schema.Resource {
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

			"database_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"partition_key_path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"throughput": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"default_ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"unique_key": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"paths": {
							Type:     schema.TypeSet,
							Required: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func sqlContainerUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	oldId := rawState["id"].(string)
	newId := strings.Replace(rawState["id"].(string), "apis/sql/databases", "sqlDatabases", 1)

	log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

	rawState["id"] = newId

	return rawState, nil
}
