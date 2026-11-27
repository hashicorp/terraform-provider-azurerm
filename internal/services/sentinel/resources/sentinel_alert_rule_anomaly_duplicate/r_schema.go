// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package sentinel_alert_rule_anomaly_duplicate

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2023-09-01/workspaces"
	"github.com/hashicorp/go-azure-sdk/resource-manager/securityinsights/2022-10-01-preview/securitymlanalyticssettings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"display_name": {
			Description:  "The Display Name of the built-in Anomaly Alert Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"built_in_rule_id": {
			Description:  "The ID of the built-in Anomaly Alert Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: securitymlanalyticssettings.ValidateSecurityMLAnalyticsSettingID,
		},

		"log_analytics_workspace_id": {
			Description:  "The ID of the Log Analytics Workspace.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: workspaces.ValidateWorkspaceID,
		},

		"enabled": {
			Description: "Whether the Duplicated Anomaly Alert Rule should be enabled.",
			Type:        pluginsdk.TypeBool,
			Required:    true,
		},

		"mode": {
			Description:  "Mode of the Duplicated Anomaly Alert Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(securitymlanalyticssettings.PossibleValuesForSettingsStatus(), false),
		},

		"multi_select_observation": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Description:  "The name of the multi select observation.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"values": {
						Description: "A list of values of the multi select observation.",
						Type:        pluginsdk.TypeList,
						Required:    true,
						Elem: &schema.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"description": {
						Description: "The description of the multi select observation.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"supported_values": {
						Description: "A list of supported values of the multi select observation.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &schema.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"single_select_observation": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Description:  "The name of the single select observation.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"description": {
						Description: "The description of the single select observation.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"supported_values": {
						Description: "A list of supported values of the single select observation.",
						Type:        pluginsdk.TypeList,
						Computed:    true,
						Elem: &schema.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"value": {
						Description:  "The value of the multi select observation.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"prioritized_exclude_observation": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Description: "The name of the prioritized exclude observation.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},

					"description": {
						Description: "The description of the prioritized exclude observation.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"prioritize": {
						Description:  "The prioritized value per `description`.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"exclude": {
						Description:  "The excluded value per `description`.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"threshold_observation": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Computed: true, // azignore:AZS007 - pre-existing violation
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Description:  "The name of the threshold observation.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"description": {
						Description: "The description of the threshold observation.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"max": {
						Description: "The max value of the threshold observation.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"min": {
						Description: "The min value of the threshold observation.",
						Type:        pluginsdk.TypeString,
						Computed:    true,
					},

					"value": {
						Description: "The value of the threshold observation.",
						Type:        pluginsdk.TypeString,
						Required:    true,
					},
				},
			},
		},
	}
}

func (r Resource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

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
			Description: "The frequency the Anomaly Alert Rule will be run, such as `P1D`.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"is_default_settings": {
			Description: "Whether the current settings of the Anomaly Alert Rule equal default settings.",
			Type:        pluginsdk.TypeBool,
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
	}
}
