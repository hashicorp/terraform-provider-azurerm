package migration

import (
	"context"
	"log"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var _ pluginsdk.StateUpgrade = CdnProfileV0ToV1{}

type CdnProfileV0ToV1 struct{}

func (CdnProfileV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*schema.Schema{
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

		"sku": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"tags": {
			Type:     schema.TypeMap,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func (CdnProfileV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old
		// 	/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}
		// new:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Cdn/profiles/{profileName}
		// summary:
		// resourcegroups -> resourceGroups
		oldId := rawState["id"].(string)
		oldParsedId, err := azure.ParseAzureResourceID(oldId)
		if err != nil {
			return rawState, err
		}

		resourceGroup := oldParsedId.ResourceGroup
		name, err := oldParsedId.PopSegment("profiles")
		if err != nil {
			return rawState, err
		}

		newId := parse.NewProfileID(oldParsedId.SubscriptionID, resourceGroup, name)
		newIdStr := newId.ID()

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newIdStr)

		rawState["id"] = newIdStr

		return rawState, nil
	}
}
