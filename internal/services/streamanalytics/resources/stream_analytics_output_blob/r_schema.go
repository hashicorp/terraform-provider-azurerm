// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_output_blob

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/outputs"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func streamAnalyticsOutputBlobSchema() map[string]*pluginsdk.Schema {
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

		"resource_group_name": commonschema.ResourceGroupName(),

		"date_format": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"path_pattern": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"storage_account_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"storage_container_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"time_format": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"serialization": helpers.SchemaStreamAnalyticsOutputSerialization(),

		"authentication_mode": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(outputs.AuthenticationModeConnectionString),
			ValidateFunc: validation.StringInSlice([]string{
				string(outputs.AuthenticationModeConnectionString),
				string(outputs.AuthenticationModeMsi),
			}, false),
		},

		"batch_max_wait_time": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validate.BatchMaxWaitTime,
		},
		"batch_min_rows": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 1000000),
		},

		"storage_account_key": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"blob_write_mode": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(outputs.BlobWriteModeAppend),
			ValidateFunc: validation.StringInSlice(outputs.PossibleValuesForBlobWriteMode(), false),
		},
	}
}
