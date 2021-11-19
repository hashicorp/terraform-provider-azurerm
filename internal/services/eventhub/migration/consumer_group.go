package migration

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ConsumerGroupsV0ToV1{}

type ConsumerGroupsV0ToV1 struct{}

func (ConsumerGroupsV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)

		newId := strings.Replace(rawState["id"].(string), "/consumergroups/", "/consumerGroups/", 1)

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId

		return rawState, nil
	}
}

func (ConsumerGroupsV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return consumerGroupsSchemaForV0AndV1()
}

func consumerGroupsSchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
		},

		"namespace_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
		},

		"eventhub_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
		},

		"resource_group_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
		},

		"user_metadata": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
		},
	}
}