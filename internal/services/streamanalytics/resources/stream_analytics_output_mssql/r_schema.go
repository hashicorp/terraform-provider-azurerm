// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_mssql

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func streamAnalyticsOutputMssqlSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"stream_analytics_job_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"server": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"database": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"table": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"user": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"password": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"max_batch_count": {
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Default:      10000,
			ValidateFunc: validation.FloatBetween(1, 1073741824),
		},

		"max_writer_count": {
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.FloatBetween(0, 1),
		},

		"authentication_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(outputs.AuthenticationModeConnectionString),
			ValidateFunc: validation.StringInSlice([]string{
				string(outputs.AuthenticationModeMsi),
				string(outputs.AuthenticationModeConnectionString),
			}, false),
		},
	}
}
