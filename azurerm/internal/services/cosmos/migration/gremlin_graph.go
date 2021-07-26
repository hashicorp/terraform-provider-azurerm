package migration

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

var _ pluginsdk.StateUpgrade = GremlinGraphV0ToV1{}

type GremlinGraphV0ToV1 struct{}

func (GremlinGraphV0ToV1) Schema() map[string]*pluginsdk.Schema {
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

		"account_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"database_name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},

		"throughput": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},

		"partition_key_path": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"index_policy": {
			Type:     schema.TypeList,
			Required: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"automatic": {
						Type:     schema.TypeBool,
						Optional: true,
						Default:  true,
					},

					"indexing_mode": {
						Type:     schema.TypeString,
						Required: true,
					},

					"included_paths": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Set: schema.HashString,
					},

					"excluded_paths": {
						Type:     schema.TypeSet,
						Optional: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
						Set: schema.HashString,
					},
				},
			},
		},

		"conflict_resolution_policy": {
			Type:     schema.TypeList,
			Required: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"mode": {
						Type:     schema.TypeString,
						Required: true,
					},

					"conflict_resolution_path": {
						Type:     schema.TypeString,
						Optional: true,
					},

					"conflict_resolution_procedure": {
						Type:     schema.TypeString,
						Optional: true,
					},
				},
			},
		},

		"unique_key": {
			Type:     schema.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"paths": {
						Type:     schema.TypeSet,
						Required: true,
						ForceNew: true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func (GremlinGraphV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId := strings.Replace(rawState["id"].(string), "apis/gremlin/databases", "gremlinDatabases", 1)

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId)

		rawState["id"] = newId

		return rawState, nil
	}
}
