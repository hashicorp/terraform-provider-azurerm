// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_aad_diagnostic_setting

import (
	"regexp"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventhub/2021-11-01/authorizationrulesnamespaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2020-08-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func monitorAadDiagnosticSettingSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.MonitorDiagnosticSettingName,
		},

		// When absent, will use the default eventhub, whilst the Diagnostic Setting API will return this property as an empty string. Therefore, it is useless to make this property as Computed.
		"eventhub_name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9]([-._a-zA-Z0-9]{0,48}[a-zA-Z0-9])?$"),
				"The event hub name can contain only letters, numbers, periods (.), hyphens (-),and underscores (_), up to 50 characters, and it must begin and end with a letter or number.",
			),
		},

		"eventhub_authorization_rule_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: authorizationrulesnamespaces.ValidateAuthorizationRuleID,
			AtLeastOneOf: []string{"eventhub_authorization_rule_id", "log_analytics_workspace_id", "storage_account_id"},
		},

		"log_analytics_workspace_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
			AtLeastOneOf: []string{"eventhub_authorization_rule_id", "log_analytics_workspace_id", "storage_account_id"},
		},

		"storage_account_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateStorageAccountID,
			AtLeastOneOf: []string{"eventhub_authorization_rule_id", "log_analytics_workspace_id", "storage_account_id"},
		},

		"enabled_log": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"category": {
						Type:     pluginsdk.TypeString,
						Required: true,
					},
				},
			},
		},
	}
}
