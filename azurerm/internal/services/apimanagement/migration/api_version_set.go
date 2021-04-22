package migration

import (
	"context"
	"log"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var _ pluginsdk.StateUpgrade = ApiVersionSetV0ToV1{}

type ApiVersionSetV0ToV1 struct {
}

func (ApiVersionSetV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"versioning_scheme": {
			Type:     schema.TypeString,
			Required: true,
		},

		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},

		"version_header_name": {
			Type:     schema.TypeString,
			Optional: true,
		},

		"version_query_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}
}

func (ApiVersionSetV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId := strings.Replace(rawState["id"].(string), "/api-version-set/", "/apiVersionSets/", 1)

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId

		return rawState, nil
	}
}
