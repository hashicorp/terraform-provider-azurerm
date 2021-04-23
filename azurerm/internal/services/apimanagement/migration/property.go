package migration

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ApiPropertyV0ToV1{}

type ApiPropertyV0ToV1 struct {
}

func (ApiPropertyV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"api_management_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"display_name": {
			Type:     schema.TypeString,
			Required: true,
		},

		"value": {
			Type:     schema.TypeString,
			Required: true,
		},

		"secret": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},

		"tags": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func (ApiPropertyV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId := strings.Replace(rawState["id"].(string), "/properties/", "/namedValues/", 1)

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId

		return rawState, nil
	}
}
