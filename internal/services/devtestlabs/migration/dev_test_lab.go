package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = DevTestLabUpgradeV0ToV1{}

type DevTestLabUpgradeV0ToV1 struct{}

func (DevTestLabUpgradeV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return devTestLabSchemaForV0AndV1()
}

func (DevTestLabUpgradeV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.devtestlab/labs/{labName}
		// new:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{labName}
		oldId, err := azure.ParseAzureResourceID(rawState["id"].(string))
		if err != nil {
			return rawState, err
		}

		labName, err := oldId.PopSegment("labs")
		if err != nil {
			return rawState, err
		}

		newId := parse.NewDevTestLabID(oldId.SubscriptionID, oldId.ResourceGroup, labName)

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId.ID())

		rawState["id"] = newId.ID()

		return rawState, nil
	}
}

func devTestLabSchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"location": azure.SchemaLocation(),

		"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

		"storage_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"tags": tags.Schema(),

		"artifacts_storage_account_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"default_storage_account_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"default_premium_storage_account_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"key_vault_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"premium_data_disk_storage_account_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"unique_identifier": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}
