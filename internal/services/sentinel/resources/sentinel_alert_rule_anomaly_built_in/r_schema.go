// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_anomaly_built_in

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2023-09-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/securitymlanalyticssettings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description:  "The Name of the built-in Anomaly Alert Rule.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"name", "display_name"},
		},

		"display_name": {
			Description:  "The Display Name of the built-in Anomaly Alert Rule.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true, // azignore:AZS007 - pre-existing violation
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ExactlyOneOf: []string{"name", "display_name"},
		},

		"log_analytics_workspace_id": {
			Description:  "The ID of the Log Analytics Workspace.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"enabled": {
			Description: "Whether the Built-in Anomaly Alert Rule should be enabled.",
			Type:        pluginsdk.TypeBool,
			Required:    true,
		},

		"mode": {
			Description:  "Mode of the Built-in Anomaly Alert Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(securitymlanalyticssettings.PossibleValuesForSettingsStatus(), false),
		},
	}
}

func (r Resource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"anomaly_version": {
			Description: "The anomaly version of the Anomaly Alert Rule.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"anomaly_settings_version": {
			Description: "The version of the Anomaly Security ML Analytics Settings.",
			Type:        pluginsdk.TypeInt,
			Computed:    true,
		},

		"description": {
			Description: "The description of the Anomaly Alert Rule.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"frequency": {
			Description: "The frequency the Anomaly Alert Rule will be run.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"required_data_connector": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"connector_id": {
						Description: "The ID of the required Data Connector.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"data_types": {
						Description: "A list of data types of the required Data Connector.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &schema.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"settings_definition_id": {
			Description: "The ID of the anomaly settings definition ID.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"tactics": {
			Description: "A list of categories of attacks by which to classify the rule.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &schema.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"techniques": {
			Description: "A list of techniques of attacks by which to classify the rule.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &schema.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"multi_select_observation": helpers.AnomalyRuleMultiSelectSchema(),

		"single_select_observation": helpers.AnomalyRuleSingleSelectSchema(),

		"prioritized_exclude_observation": helpers.AnomalyRulePrioritySchema(),

		"threshold_observation": helpers.AnomalyRuleThresholdSchema(),
	}
}
