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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"workspace_name": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			DiffSuppressFunc: suppress.CaseDifference,
			ValidateFunc:     validate.LogAnalyticsWorkspaceName,
		},

		"counter_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"instance_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"interval_seconds": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(10, math.MaxInt32),
		},

		"object_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}
