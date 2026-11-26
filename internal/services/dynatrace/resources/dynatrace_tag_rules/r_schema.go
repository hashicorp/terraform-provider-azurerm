// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dynatrace_tag_rules

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/dynatrace/2023-04-27/monitors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description:  "Name of the Dynatrace tag rules.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"monitor_id": {
			Description:  "Name of the Dynatrace monitor.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: monitors.ValidateMonitorID,
		},

		"log_rule": {
			Description: "Set of rules for sending logs for the Monitor resource.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"send_azure_active_directory_logs_enabled": {
						Description: "Send Azure Active Directory logs.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"send_activity_logs_enabled": {
						Description: "Send Activity logs.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"send_subscription_logs_enabled": {
						Description: "Send Subscription logs.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},

					"filtering_tag": {
						Description: "Filtering tag for the log rule.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"action": {
									Description: "Action of the filtering tag.",
									Type:        pluginsdk.TypeString,
									Required:    true,
									ValidateFunc: validation.StringInSlice([]string{
										"Include",
										"Exclude",
									}, false),
								},

								"name": {
									Description:  "Name of the filtering tag.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"value": {
									Description:  "Value of the filtering tag.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
				},
			},
		},

		"metric_rule": {
			Description: "Set of rules for sending metrics for the Monitor resource.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
			MaxItems:    1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*schema.Schema{
					"filtering_tag": {
						Description: "Filtering tag for the metric rule.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*schema.Schema{
								"action": {
									Description: "Action of the filtering tag.",
									Type:        pluginsdk.TypeString,
									Required:    true,
									ValidateFunc: validation.StringInSlice([]string{
										"Include",
										"Exclude",
									}, false),
								},

								"name": {
									Description:  "Name of the filtering tag.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"value": {
									Description:  "Value of the filtering tag.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
					"sending_metrics_enabled": {
						Description: "If sending metrics is enabled.",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
				},
			},
		},
	}
}

func (r Resource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{}
}
