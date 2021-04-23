package migration

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = StandardEnvironmentAccessPolicyV0ToV1{}

type StandardEnvironmentAccessPolicyV0ToV1 struct{}

func (StandardEnvironmentAccessPolicyV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"time_series_insights_environment_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"principal_object_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},

		"roles": {
			Type:     schema.TypeSet,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func (StandardEnvironmentAccessPolicyV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		log.Println("[DEBUG] Migrating ResourceType from v0 to v1 format")
		oldId := rawState["id"].(string)
		newId := strings.Replace(oldId, "/accesspolicies/", "/accessPolicies/", 1)

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId

		return rawState, nil
	}
}
