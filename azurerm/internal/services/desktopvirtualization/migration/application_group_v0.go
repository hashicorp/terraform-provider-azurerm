package migration

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/desktopvirtualization/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
)

func ApplicationGroupUpgradeV0Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"host_pool_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"friendly_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func ApplicationGroupUpgradeV0ToV1(rawState map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
	oldId := rawState["id"].(string)
	id, err := parse.ApplicationGroupIDInsensitively(oldId)
	if err != nil {
		return nil, err
	}
	newId := id.ID()
	log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
	rawState["id"] = newId

	oldHostPoolId := rawState["host_pool_id"].(string)
	hostPoolId, err := parse.HostPoolIDInsensitively(oldHostPoolId)
	if err != nil {
		return nil, err
	}
	newHostPoolId := hostPoolId.ID()
	log.Printf("[DEBUG] Updating Host Pool ID from %q to %q", oldHostPoolId, newHostPoolId)
	rawState["host_pool_id"] = newHostPoolId

	return rawState, nil
}
