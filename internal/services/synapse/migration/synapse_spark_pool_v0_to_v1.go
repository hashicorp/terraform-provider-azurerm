// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package migration

import (
	"context"

	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/bigdatapools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type SynapseSparkPoolV0ToV1 struct{}

var _ pluginsdk.StateUpgrade = &SynapseSparkPoolV0ToV1{}

func (s SynapseSparkPoolV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"synapse_workspace_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"node_size_family": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"node_size": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"cache_size": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"compute_isolation_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"dynamic_executor_allocation_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"min_executors": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"max_executors": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"node_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Computed: true,
		},

		"auto_scale": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"min_node_count": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},

					"max_node_count": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
				},
			},
		},

		"auto_pause": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"delay_in_minutes": {
						Type:     pluginsdk.TypeInt,
						Required: true,
					},
				},
			},
		},

		"session_level_packages_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"spark_config": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"content": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"filename": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"spark_events_folder": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"spark_log_folder": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"library_requirement": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"content": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},

					"filename": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},

		"spark_version": {
			Type:     pluginsdk.TypeString,
			Required: true,
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

func (s SynapseSparkPoolV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		oldId := rawState["id"].(string)
		newId, err := bigdatapools.ParseBigDataPoolIDInsensitively(oldId)
		if err != nil {
			return nil, err
		}
		rawState["id"] = newId.ID()
		return rawState, nil
	}
}
