package migration

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/parse"
)

func UserAssignedIdentityV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		Version: 0,
		Type:    userAssignedIdentityV0Schema().CoreConfigSchema().ImpliedType(),
		Upgrade: userAssignedIdentityUpgradeV0ToV1,
	}
}

func userAssignedIdentityV0Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 128),
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"tags": tags.Schema(),

			"principal_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func userAssignedIdentityUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	log.Printf("[DEBUG] Migrating IDs to correct casing for User Assigned Idenitity")
	name := rawState["name"].(string)
	resourceGroup := rawState["resource_group_name"].(string)
	id := parse.NewUserAssignedIdentityID(subscriptionId, resourceGroup, name)

	rawState["id"] = id.ID()
	return rawState, nil
}
