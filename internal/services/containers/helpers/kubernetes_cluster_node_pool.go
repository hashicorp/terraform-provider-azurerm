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
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"max_unavailable": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
				"drain_timeout_in_minutes": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},
				"node_soak_duration_in_minutes": {
					Type:     pluginsdk.TypeInt,
					Computed: true,
				},
				"undrainable_node_behavior": {
					Type:     pluginsdk.TypeString,
					Computed: true,
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
