package storage

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

// the schema schema was used for both V0 and V1
func ResourceStorageTableStateResourceV0V1() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_account_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateArmStorageAccountName,
			},
			"resource_group_name": azure.SchemaResourceGroupName(),
			"acl": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringLenBetween(1, 64),
						},
						"access_policy": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"expiry": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"permissions": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func ResourceStorageTableStateUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	tableName := rawState["name"].(string)
	accountName := rawState["storage_account_name"].(string)
	environment := meta.(*clients.Client).Account.Environment

	id := rawState["id"].(string)
	newResourceID := fmt.Sprintf("https://%s.table.%s/%s", accountName, environment.StorageEndpointSuffix, tableName)
	log.Printf("[DEBUG] Updating ID from %q to %q", id, newResourceID)

	rawState["id"] = newResourceID
	return rawState, nil
}

func ResourceStorageTableStateUpgradeV1ToV2(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	tableName := rawState["name"].(string)
	accountName := rawState["storage_account_name"].(string)
	environment := meta.(*clients.Client).Account.Environment

	id := rawState["id"].(string)
	newResourceID := fmt.Sprintf("https://%s.table.%s/Tables('%s')", accountName, environment.StorageEndpointSuffix, tableName)
	log.Printf("[DEBUG] Updating ID from %q to %q", id, newResourceID)

	rawState["id"] = newResourceID
	return rawState, nil
}
