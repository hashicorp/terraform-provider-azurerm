package migration

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = AdvancedThreatProtectionV0ToV1{}

type AdvancedThreatProtectionV0ToV1 struct{}

func (AdvancedThreatProtectionV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*schema.Schema{
		"target_resource_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"enabled": {
			Type:     schema.TypeBool,
			Required: true,
		},
	}
}

func (AdvancedThreatProtectionV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)

		// remove the existing `/` if it's present (2.42+) which'll do nothing if it wasn't (2.38)
		newId := fmt.Sprintf("/%s", strings.TrimPrefix(oldId, "/"))

		parsedId, err := parse.AdvancedThreatProtectionID(newId)
		if err != nil {
			return nil, err
		}

		newId = parsedId.ID()

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
		rawState["id"] = newId
		return rawState, nil
	}
}
