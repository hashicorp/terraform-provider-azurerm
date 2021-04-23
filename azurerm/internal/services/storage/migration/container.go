package migration

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ContainerV0ToV1{}

type ContainerV0ToV1 struct{}

func (ContainerV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*schema.Schema{
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
		"container_access_type": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  "private",
		},

		"properties": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func (ContainerV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	// this should have been applied from pre-0.12 migration system; backporting just in-case
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		environment := meta.(*clients.Client).Account.Environment

		containerName := rawState["name"]
		storageAccountName := rawState["storage_account_name"]
		newID := fmt.Sprintf("https://%s.blob.%s/%s", storageAccountName, environment.StorageEndpointSuffix, containerName)
		rawState["id"] = newID

		return rawState, nil
	}
}
