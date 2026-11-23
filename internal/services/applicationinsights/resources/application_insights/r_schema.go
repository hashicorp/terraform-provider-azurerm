// Copyright IBM Corp. 2014, 2026
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
			Description: "The name of the Application Insights component.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"application_type": {
			Description: "The type of Application Insights to create. Please note these values are case sensitive; unmatched values are treated as _ASP.NET_ by Azure.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
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
			Description: "The ID of a Log Analytics Workspace resource.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// NOTE: O+C A Log Analytics Workspace will be attached to the Application Insight by default, which should be computed=true
			Computed:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"retention_in_days": {
			Description: "The retention period in days.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
			Default:     90,
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
			Description:  "The percentage of the data produced by the monitored application that is sampled for Application Insights telemetry.",
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Default:      100,
			ValidateFunc: validation.FloatBetween(0, 100),
		},

		"ip_masking_enabled": {
			Description: "By default the real client IP is masked as `0.0.0.0` in the logs. Set this argument to `false` to disable masking and log the real client IP.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"tags": commonschema.Tags(),

		"daily_data_cap_in_gb": {
			Description:  "The Application Insights component daily data volume cap in GB.",
			Type:         pluginsdk.TypeFloat,
			Optional:     true,
			Default:      100,
			ValidateFunc: validation.FloatAtLeast(0),
		},

		"daily_data_cap_notifications_enabled": {
			Description: "Whether a notification email will be sent when the daily data volume cap is met.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"app_id": {
			Description: "The App ID associated with this Application Insights component.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"instrumentation_key": {
			Description: "The Instrumentation Key for this Application Insights component.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"connection_string": {
			Description: "The Connection String for this Application Insights component.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Sensitive:   true,
		},

		"local_authentication_enabled": {
			Description: "Whether Non-Azure AD based Auth is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"internet_ingestion_enabled": {
			Description: "Should the Application Insights component support ingestion over the Public Internet?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"internet_query_enabled": {
			Description: "Should the Application Insights component support querying over the Public Internet?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},
		"force_customer_storage_for_profiler": {
			Description: "Should the Application Insights component force users to create their own Storage Account for profiling?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},
	}
}
