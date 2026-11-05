// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package application_insights

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func applicationInsightsSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"application_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"web",
				"other",
				"java",
				"MobileCenter",
				"phone",
				"store",
				"ios",
				"Node.JS",
			}, false),
		},

		"workspace_id": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			// NOTE: O+C A Log Analytics Workspace will be attached to the Application Insight by default, which should be computed=true
			Computed:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"retention_in_days": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  90,
			ValidateFunc: validation.IntInSlice([]int{
				30,
				60,
				90,
				120,
				180,
				270,
				365,
				550,
				730,
			}),
		},

		"sampling_percentage": {
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Default:      100,
			ValidateFunc: validation.FloatBetween(0, 100),
		},

		"ip_masking_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"tags": commonschema.Tags(),

		"daily_data_cap_in_gb": {
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Default:      100,
			ValidateFunc: validation.FloatAtLeast(0),
		},

		"daily_data_cap_notifications_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"app_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"instrumentation_key": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"connection_string": {
			Type:      pluginsdk.TypeString,
			Computed:  true,
			Sensitive: true,
		},

		"local_authentication_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"internet_ingestion_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"internet_query_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},
		"force_customer_storage_for_profiler": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},
	}
}
