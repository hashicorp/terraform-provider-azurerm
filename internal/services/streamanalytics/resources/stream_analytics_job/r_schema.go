// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package stream_analytics_job

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/streamanalytics/2021-10-01-preview/streamingjobs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func streamAnalyticsJobSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"stream_analytics_cluster_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"compatibility_level": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C Best Practice from MSFT is to use the latest version 1.2, but API uses 1.0 as the default
			Computed: true,
			ValidateFunc: validation.StringInSlice([]string{
				// values found in the other API the portal uses
				string(streamingjobs.CompatibilityLevelOnePointZero),
				"1.1",
				string(streamingjobs.CompatibilityLevelOnePointTwo),
			}, false),
		},

		"data_locale": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      "en-US",
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"events_late_arrival_max_delay_in_seconds": {
			// portal allows for up to 20d 23h 59m 59s
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(-1, 1814399),
			Default:      5,
		},

		"events_out_of_order_max_delay_in_seconds": {
			// portal allows for up to 9m 59s
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 599),
			Default:      0,
		},

		"events_out_of_order_policy": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(streamingjobs.PossibleValuesForEventsOutOfOrderPolicy(), false),
			Default:      string(streamingjobs.EventsOutOfOrderPolicyAdjust),
		},

		"type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(streamingjobs.PossibleValuesForJobType(), false),
			Default:      string(streamingjobs.JobTypeCloud),
		},

		"output_error_policy": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringInSlice(streamingjobs.PossibleValuesForOutputErrorPolicy(), false),
			Default:      string(streamingjobs.OutputErrorPolicyDrop),
		},

		"streaming_units": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(1, 120),
		},

		"content_storage_policy": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Default:      string(streamingjobs.ContentStoragePolicySystemAccount),
			ValidateFunc: validation.StringInSlice(streamingjobs.PossibleValuesForContentStoragePolicy(), false),
		},

		"job_storage_account": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"authentication_mode": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Default:  string(streamingjobs.AuthenticationModeConnectionString),
						ValidateFunc: validation.StringInSlice([]string{
							string(streamingjobs.AuthenticationModeConnectionString),
							string(streamingjobs.AuthenticationModeMsi),
						}, false),
					},

					"account_name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"account_key": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Sensitive:    true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"transformation_query": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"job_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"sku_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(streamingjobs.SkuNameStandard),
			ValidateFunc: validation.StringInSlice([]string{
				string(streamingjobs.SkuNameStandard),
				"StandardV2", // missing from swagger as described here https://github.com/Azure/azure-rest-api-specs/issues/27506
			}, false),
		},

		"tags": commonschema.Tags(),
	}
}
