package migration

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func MongoCollectionV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		Type:    mongoCollectionSchemaForV0().CoreConfigSchema().ImpliedType(),
		Upgrade: mongoCollectionUpgradeV0ToV1,
		Version: 0,
	}
}

func mongoCollectionSchemaForV0() *schema.Resource {
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

			"shard_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"default_ttl_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"throughput": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func mongoCollectionUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	oldId := rawState["id"].(string)
	newId := strings.Replace(rawState["id"].(string), "apis/mongodb/databases", "mongodbDatabases", 1)

	log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

	rawState["id"] = newId

	return rawState, nil
}
