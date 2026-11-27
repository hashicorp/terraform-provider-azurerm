// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_ms_security_incident

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func sentinelAlertRuleMsSecurityIncidentSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name which should be used for this Sentinel MS Security Incident Alert Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"log_analytics_workspace_id": {
			Description:  "The ID of the Log Analytics Workspace this Sentinel MS Security Incident Alert Rule belongs to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: alertrules.ValidateWorkspaceID,
		},

		"display_name": {
			Description:  "The friendly name of this Sentinel MS Security Incident Alert Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"product_filter": {
			Description:  "The Microsoft Security Service from where the alert will be generated.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForMicrosoftSecurityProductName(), false),
		},

		"severity_filter": {
			Description: "Only create incidents from alerts when alert severity level is contained in this list.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			MinItems:    1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice(alertrules.PossibleValuesForAlertSeverity(), false),
			},
		},

		"alert_rule_template_guid": {
			Description:  "The GUID of the alert rule template which is used to create this Sentinel Scheduled Alert Rule.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"description": {
			Description:  "The description of this Sentinel MS Security Incident Alert Rule.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"enabled": {
			Description: "Whether this Sentinel MS Security Incident Alert Rule should be enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"display_name_filter": {
			Description: "Only create incidents when the alert display name contain text from this list, leave empty to apply no filter.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			MinItems:    1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"display_name_exclude_filter": {
			Description: "Only create incidents when the alert display name doesn't contain text from this list.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			MinItems:    1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}
