// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_datasource_windows_event

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/set"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func logAnalyticsDatasourceWindowsEventSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Log Analytics Windows Event DataSource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"workspace_name": {
			Description:      "The name of the Log Analytics Workspace where the Log Analytics Windows Event DataSource should exist.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     validate.LogAnalyticsWorkspaceName,
		},

		"event_log_name": {
			Description:  "The name of the Windows Event Log to collect events from.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"event_types": {
			Description: "An array of event types applied to the specified event log.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			MinItems:    1,
			Set:         set.HashStringIgnoreCase,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				// API backend accepts event_types case-insensitively
				ValidateFunc: validation.StringInSlice([]string{"Error", "Warning", "Information"}, false),
			},
		},
	}
}
