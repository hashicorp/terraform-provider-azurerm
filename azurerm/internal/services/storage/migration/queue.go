package migration

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func QueueV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		// this should have been applied from pre-0.12 migration system; backporting just in-case
		Type:    queueSchemaForV0().CoreConfigSchema().ImpliedType(),
		Upgrade: queueUpgradeV0ToV1,
		Version: 0,
	}
}

func queueSchemaForV0() *schema.Resource {
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
			"storage_account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func queueUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	environment := meta.(*clients.Client).Account.Environment

	queueName := rawState["name"]
	storageAccountName := rawState["storage_account_name"]
	newID := fmt.Sprintf("https://%s.queue.%s/%s", storageAccountName, environment.StorageEndpointSuffix, queueName)
	rawState["id"] = newID

	return rawState, nil
}
