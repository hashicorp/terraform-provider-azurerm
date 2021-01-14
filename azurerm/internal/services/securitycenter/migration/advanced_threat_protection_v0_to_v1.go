package migration

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/parse"
)

func AdvancedThreatProtectionV0ToV1() schema.StateUpgrader {
	return schema.StateUpgrader{
		Version: 0,
		Type:    advancedThreatProtectionV0Schema().CoreConfigSchema().ImpliedType(),
		Upgrade: advancedThreadProtectionV0toV1Upgrade,
	}
}

func advancedThreatProtectionV0Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"target_resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func advancedThreadProtectionV0toV1Upgrade(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	oldId := rawState["id"].(string)

	// remove the existing `/` if it's present (2.42+) which'll do nothing if it wasn't (2.38)
	newId := strings.TrimPrefix(oldId, "/")
	newId = fmt.Sprintf("/%s", oldId)

	parsedId, err := parse.AdvancedThreatProtectionID(newId)
	if err != nil {
		return nil, err
	}

	newId = parsedId.ID()

	log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)
	rawState["id"] = newId
	return rawState, nil
}
