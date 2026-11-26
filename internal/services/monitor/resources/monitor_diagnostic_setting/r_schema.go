// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_diagnostic_setting

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/authorizationrulesnamespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	eventhubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func monitorDiagnosticSettingSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Diagnostic Setting.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.MonitorDiagnosticSettingName,
		},

		"target_resource_id": {
			Description: "The ID of an existing Resource on which to configure Diagnostic Settings.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.Any(
				azure.ValidateResourceID,
				commonids.ValidateManagementGroupID,
			),
		},

		"eventhub_name": {
			Description:  "The name of the Event Hub where Diagnostics Data should be sent.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: eventhubValidate.ValidateEventHubName(),
		},

		"eventhub_authorization_rule_id": {
			Description:  "The ID of an Event Hub Namespace Authorization Rule used to send Diagnostics Data.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: authorizationrulesnamespaces.ValidateAuthorizationRuleID,
			AtLeastOneOf: []string{"eventhub_authorization_rule_id", "log_analytics_workspace_id", "storage_account_id", "partner_solution_id"},
		},

		"log_analytics_workspace_id": {
			Description:  "The ID of a Log Analytics Workspace where Diagnostics Data should be sent.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
			AtLeastOneOf: []string{"eventhub_authorization_rule_id", "log_analytics_workspace_id", "storage_account_id", "partner_solution_id"},
		},

		"storage_account_id": {
			Description:  "The ID of the Storage Account where logs should be sent.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
			AtLeastOneOf: []string{"eventhub_authorization_rule_id", "log_analytics_workspace_id", "storage_account_id", "partner_solution_id"},
		},

		"partner_solution_id": {
			Description:  "The ID of the market partner solution where Diagnostics Data should be sent. For potential partner integrations, [click to learn more about partner integration](https://learn.microsoft.com/azure/partner-solutions/overview).",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: azure.ValidateResourceID,
			AtLeastOneOf: []string{"eventhub_authorization_rule_id", "log_analytics_workspace_id", "storage_account_id", "partner_solution_id"},
		},

		"log_analytics_destination_type": {
			Description: "When set to `Dedicated`, logs sent to a Log Analytics Workspace will go into resource specific tables, instead of the legacy `AzureDiagnostics` table.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice([]string{
				"Dedicated",
				"AzureDiagnostics", // Not documented in azure API, but some resource has skew. See: https://github.com/Azure/azure-rest-api-specs/issues/9281
			}, false),
		},

		"enabled_log": {
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			Set:          resourceMonitorDiagnosticLogSettingHash,
			AtLeastOneOf: []string{"enabled_log", "enabled_metric"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"category": {
						Description:  "The name of a Diagnostic Log Category for this Resource.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"category_group": {
						Description:  "The name of a Diagnostic Log Category Group for this Resource.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"enabled_metric": {
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			AtLeastOneOf: []string{"enabled_log", "enabled_metric"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"category": {
						Description:  "The name of a Diagnostic Metric Category for this Resource.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
}
