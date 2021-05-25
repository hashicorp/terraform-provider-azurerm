package migration

import (
	"context"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = NamespaceV0ToV1{}

type NamespaceV0ToV1 struct{}

func (NamespaceV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"sku": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"capacity": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			ForceNew: true,
		},

		"default_primary_connection_string": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"default_secondary_connection_string": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"default_primary_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"default_secondary_key": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (NamespaceV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	// this should have been applied from pre-0.12 migration system; backporting just in-case
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		skuName := rawState["sku"].(string)
		if !strings.EqualFold(skuName, "Premium") {
			delete(rawState, "capacity")
		}

		return rawState, nil
	}
}
