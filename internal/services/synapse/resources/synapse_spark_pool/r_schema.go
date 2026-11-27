// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_spark_pool

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/bigdatapools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func synapseSparkPoolSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SparkPoolName,
		},

		"synapse_workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.AsGeneratedID(workspaces.ParseWorkspaceIDInsensitively),
		},

		"node_size_family": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(bigdatapools.PossibleValuesForNodeSizeFamily(), false),
		},

		"node_size": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(bigdatapools.PossibleValuesForNodeSize(), false),
		},

		"cache_size": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
		},

		"compute_isolation_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"dynamic_executor_allocation_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"min_executors": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 200),
		},

		"max_executors": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 200),
		},

		"node_count": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			// NOTE: O+C There is a bug in the API where this gets set when auto_scale is enabled resulting in a diff
			Computed:     true,
			ValidateFunc: validation.IntBetween(3, 200),
			ExactlyOneOf: []string{"node_count", "auto_scale"},
		},

		"auto_scale": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MaxItems:     1,
			ExactlyOneOf: []string{"node_count", "auto_scale"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"min_node_count": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(3, 200),
					},

					"max_node_count": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(3, 200),
					},
				},
			},
		},

		"auto_pause": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"delay_in_minutes": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntBetween(5, 10080),
					},
				},
			},
		},

		"session_level_packages_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"spark_config": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"content": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"filename": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"spark_events_folder": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "/events",
		},

		"spark_log_folder": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "/logs",
		},

		"library_requirement": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
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
			ValidateFunc: validation.StringInSlice([]string{
				"3.4",
				"3.5",
			}, false),
		},

		"tags": commonschema.Tags(),
	}
}
