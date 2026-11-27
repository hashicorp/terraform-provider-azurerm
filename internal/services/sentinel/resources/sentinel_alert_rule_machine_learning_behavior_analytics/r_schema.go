// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_machine_learning_behavior_analytics

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func sentinelAlertRuleMachineLearningBehaviorAnalyticsSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"log_analytics_workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: alertrules.ValidateWorkspaceID,
		},

		"alert_rule_template_guid": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},
	}
}
