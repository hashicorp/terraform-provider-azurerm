// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func UpgradeSettingsForDataSourceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Computed: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"max_surge": {
					Description: "The maximum number of percentage of nodes that will be added to the Node Pool size during an upgrade.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},
				"max_unavailable": {
					Description: "The maximum number or percentage of nodes which can be unavailable during the upgrade.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},
				"drain_timeout_in_minutes": {
					Description: "The amount of time in minutes to wait on eviction of pods and graceful termination per node. This eviction wait time honours waiting on pod disruption budgets. If this time is exceeded, the upgrade fails.",
					Type:        pluginsdk.TypeInt,
					Computed:    true,
				},
				"node_soak_duration_in_minutes": {
					Description: "The amount of time in minutes to wait after draining a node and before reimaging it and moving on to next node.",
					Type:        pluginsdk.TypeInt,
					Computed:    true,
				},
				"undrainable_node_behavior": {
					Description: "The action when a node is undrainable during upgrade.",
					Type:        pluginsdk.TypeString,
					Computed:    true,
				},
			},
		},
	}
}

func ExpandNodeLabels(input map[string]any) *map[string]string {
	result := make(map[string]string)
	for k, v := range input {
		result[k] = v.(string)
	}
	return &result
}
