// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_fusion

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2023-12-01-preview/alertrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func sentinelAlertRuleFusionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"log_analytics_workspace_id": {
			Description:  "The ID of the Log Analytics Workspace this Sentinel Fusion Alert Rule belongs to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: alertrules.ValidateWorkspaceID,
		},

		"alert_rule_template_guid": {
			Description:  "The GUID of the alert rule template which is used for this Sentinel Fusion Alert Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"enabled": {
			Description: "Whether this Sentinel Fusion Alert Rule should be enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"source": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			// NOTE: O+C The API creates a source if omitted based on the `alert_rule_template_guid`
			// but overwriting this/reverting to the default can be done without issue so this can remain
			Computed: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the Fusion source signal. Refer to Fusion alert rule template for supported values.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"enabled": {
						Description: "Whether this source signal is enabled or disabled in Fusion detection.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     true,
					},
					"sub_type": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description:  "The Name of the source subtype under a given source signal in Fusion detection. Refer to Fusion alert rule template for supported values.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"enabled": {
									Description: "Whether this source subtype under source signal is enabled or disabled in Fusion detection.",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									Default:     true,
								},
								"severities_allowed": {
									Description: "A list of severities that are enabled for this source subtype consumed in Fusion detection.",
									Type:        pluginsdk.TypeSet,
									Required:    true,
									MinItems:    1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice(
											alertrules.PossibleValuesForAlertSeverity(),
											false,
										),
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
