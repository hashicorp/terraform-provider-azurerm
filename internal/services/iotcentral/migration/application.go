package migration

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iotcentral/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = ApplicationV0ToV1{}

type ApplicationV0ToV1 struct{}

func (ApplicationV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"location": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"sub_domain": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"sku": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "ST1",
		},

		"template": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true,
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
		
		"public_network_access" : {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"network_rule_sets" : {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"private_endpoint_connections" : { 
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
			// Elem: &pluginsdk.Schema{
			// 	Type: pluginsdk.TypeMap,
			// 	Elem: &pluginsdk.Schema{
			// 		Type: pluginsdk.TypeString,
			// 	},
			// },
		},
	}
}

func (ApplicationV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		id, err := parse.ApplicationIDInsensitively(oldId)
		if err != nil {
			return rawState, err
		}

		newId := id.ID()
		log.Printf("Updating `id` from %q to %q", oldId, newId)
		rawState["id"] = newId
		return rawState, nil
	}
}
