package migration

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iotcentral/parse"
)

func ApplicationV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		Type:    applicationSchemaForV0().CoreConfigSchema().ImpliedType(),
		Upgrade: applicationUpgradeV0ToV1,
		Version: 0,
	}
}

func applicationSchemaForV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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

			"sub_domain": {
				Type:     schema.TypeString,
				Required: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"sku": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ST1",
			},

			"template": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func applicationUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	oldId := rawState["id"].(string)
	id, err := parse.ApplicationIDInsensitively(oldId)
	if err != nil {
		return rawState, err
	}

	newId := id.ID()
	log.Printf("Updating `id` from %q to %q", oldId, newId)
	rawState["id"] = newId
	return rawState, nil
}
