// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package log_analytics_data_export_rule

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/dataexport"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func logAnalyticsDataExportRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:      "The name of the Log Analytics Data Export Rule.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			ForceNew:         true,
			ValidateFunc:     validate.LogAnalyticsDataExportName,
			DiffSuppressFunc: suppress.CaseDifference,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"workspace_resource_id": {
			Description:  "The resource ID of the workspace.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: dataexport.ValidateWorkspaceID,
		},

		"destination_resource_id": {
			Description:  "The destination resource ID. If the destination is an Event Hub namespace, an Event Hub would be created for each table automatically.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"enabled": {
			Description: "Whether this Log Analytics Data Export Rule is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"table_names": {
			Description: "A list of table names to export to the destination resource, for example: `[\"Heartbeat\", \"SecurityEvent\"]`.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			MinItems:    1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.NoZeroValues,
			},
		},

		"export_rule_id": {
			Description: "The ID of the created Data Export Rule.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}
