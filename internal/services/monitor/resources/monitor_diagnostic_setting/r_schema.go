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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.MonitorDiagnosticSettingName,
		},

		"target_resource_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.Any(
				azure.ValidateResourceID,
				commonids.ValidateManagementGroupID,
			),
		},

		"eventhub_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: eventhubValidate.ValidateEventHubName(),
		},

		"eventhub_authorization_rule_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: authorizationrulesnamespaces.ValidateAuthorizationRuleID,
			AtLeastOneOf: []string{"eventhub_authorization_rule_id", "log_analytics_workspace_id", "storage_account_id", "partner_solution_id"},
		},

		"log_analytics_workspace_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
			AtLeastOneOf: []string{"eventhub_authorization_rule_id", "log_analytics_workspace_id", "storage_account_id", "partner_solution_id"},
		},

		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
			AtLeastOneOf: []string{"eventhub_authorization_rule_id", "log_analytics_workspace_id", "storage_account_id", "partner_solution_id"},
		},

		"partner_solution_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: azure.ValidateResourceID,
			AtLeastOneOf: []string{"eventhub_authorization_rule_id", "log_analytics_workspace_id", "storage_account_id", "partner_solution_id"},
		},

		"log_analytics_destination_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			ValidateFunc: validation.StringInSlice([]string{
				"Dedicated",
				"AzureDiagnostics", // Not documented in azure API, but some resource has skew. See: https://github.com/Azure/azure-rest-api-specs/issues/9281
			}, false),
		},

		"enabled_log": {
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			AtLeastOneOf: []string{"enabled_log", "enabled_metric"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"category": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"category_group": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
			Set: resourceMonitorDiagnosticLogSettingHash,
		},

		"enabled_metric": {
			Type:         pluginsdk.TypeSet,
			Optional:     true,
			AtLeastOneOf: []string{"enabled_log", "enabled_metric"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"category": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
	}
}
