package migration

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
)

func StoreFileV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		Type:    storeFileSchemaForV0().CoreConfigSchema().ImpliedType(),
		Upgrade: storeFileUpgradeV0ToV1,
		Version: 0,
	}
}

func storeFileSchemaForV0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"account_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"remote_file_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"local_file_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func storeFileUpgradeV0ToV1(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	client := meta.(*clients.Client).Datalake.StoreFilesClient

	storageAccountName := rawState["account_name"].(string)
	filePath := rawState["remote_file_path"].(string)
	newID := fmt.Sprintf("%s.%s%s", storageAccountName, client.AdlsFileSystemDNSSuffix, filePath)
	rawState["id"] = newID
	return rawState, nil
}
