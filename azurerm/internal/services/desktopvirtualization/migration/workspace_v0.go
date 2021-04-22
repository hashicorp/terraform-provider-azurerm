package migration

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/desktopvirtualization/parse"
)

func WorkspaceV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		Type:    workspaceSchemaForV0().CoreConfigSchema().ImpliedType(),
		Upgrade: workspaceUpgradeV0ToV1,
		Version: 0,
	}
}

func workspaceSchemaForV0() *schema.Resource {
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

			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
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

func workspaceUpgradeV0ToV1(rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	oldId := rawState["id"].(string)

	id, err := parse.WorkspaceID(oldId)
	if err != nil {
		return nil, err
	}
	newId := id.ID()

	log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
	rawState["id"] = newId

	return rawState, nil
}
