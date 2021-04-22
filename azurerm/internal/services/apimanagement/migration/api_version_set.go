package migration

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func ApiVersionSetV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		Type:    apiVersionSetSchemaForV0().CoreConfigSchema().ImpliedType(),
		Upgrade: apiVersionSetUpgradeV0ToV1,
		Version: 0,
	}
}

func apiVersionSetSchemaForV0() *schema.Resource {
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

			"api_management_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"versioning_scheme": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"version_header_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"version_query_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func apiVersionSetUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	oldId := rawState["id"].(string)
	newId := strings.Replace(rawState["id"].(string), "/api-version-set/", "/apiVersionSets/", 1)

	log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

	rawState["id"] = newId

	return rawState, nil
}
