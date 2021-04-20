package migration

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func ApiPropertyV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		Type:    apiPropertySchemaForV0().CoreConfigSchema().ImpliedType(),
		Upgrade: apiPropertyUpgradeV0ToV1,
		Version: 0,
	}
}

func apiPropertySchemaForV0() *schema.Resource {
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

			"value": {
				Type:     schema.TypeString,
				Required: true,
			},

			"secret": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func apiPropertyUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	oldId := rawState["id"].(string)
	newId := strings.Replace(rawState["id"].(string), "/properties/", "/namedValues/", 1)

	log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

	rawState["id"] = newId

	return rawState, nil
}
