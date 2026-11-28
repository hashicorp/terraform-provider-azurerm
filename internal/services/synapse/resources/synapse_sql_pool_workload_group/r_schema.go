// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_sql_pool_workload_group

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func synapseSqlPoolWorkloadGroupSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"sql_pool_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SqlPoolID,
		},

		"max_resource_percent": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(1, 100),
		},

		"min_resource_percent": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(0, 100),
		},

		"importance": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "normal",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"max_resource_percent_per_request": {
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Default:      3,
			ValidateFunc: validation.FloatBetween(0, 100),
		},

		"min_resource_percent_per_request": {
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			ValidateFunc: validation.FloatBetween(0, 100),
		},

		"query_execution_timeout_in_seconds": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntAtLeast(0),
		},
	}
}
