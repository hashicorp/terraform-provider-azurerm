// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_sql_pool_workload_classifier

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func synapseSqlPoolWorkloadClassifierSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"workload_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SqlPoolWorkloadGroupID,
		},

		"member_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"context": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"end_time": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^\d{2}:\d{2}$`),
				"The `end_time` is of the `HH:MM` format in UTC time zone",
			),
		},

		"importance": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"low",
				"below_normal",
				"normal",
				"above_normal",
				"high",
			}, false),
		},

		"label": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"start_time": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^\d{2}:\d{2}$`),
				"The `start_time` is of the `HH:MM` format in UTC time zone",
			),
		},
	}
}
