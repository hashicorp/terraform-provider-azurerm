// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_datasource_windows_performance_counter

import (
	"math"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func logAnalyticsDatasourceWindowsPerformanceCounterSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The Name which should be used for this Log Analytics Windows Performance Counter DataSource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"workspace_name": {
			Description:      "The name of the Log Analytics Workspace where the Log Analytics Windows Performance Counter DataSource should exist.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateFunc:     validate.LogAnalyticsWorkspaceName,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"counter_name": {
			Description:  "The friendly name of the performance counter.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"instance_name": {
			Description:  "The name of the Virtual Machine instance to which the Windows Performance Counter DataSource to be applied. Specify a `*` will apply to all instances.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"interval_seconds": {
			Description:  "The time of sample interval in seconds. Supports values between 10 and 2147483647.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(10, math.MaxInt32),
		},

		"object_name": {
			Description:  "The object name of the Log Analytics Windows Performance Counter DataSource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}
