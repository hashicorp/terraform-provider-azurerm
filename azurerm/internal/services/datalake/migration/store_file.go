package migration

import (
	"context"
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = StoreFileV0ToV1{}

type StoreFileV0ToV1 struct{}

func (StoreFileV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"account_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"remote_file_path": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"local_file_path": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
	}
}

func (StoreFileV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		client := meta.(*clients.Client).Datalake.StoreFilesClient

		storageAccountName := rawState["account_name"].(string)
		filePath := rawState["remote_file_path"].(string)
		newID := fmt.Sprintf("%s.%s%s", storageAccountName, client.AdlsFileSystemDNSSuffix, filePath)
		rawState["id"] = newID
		return rawState, nil
	}
}
