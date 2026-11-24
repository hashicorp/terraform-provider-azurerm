// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor_tag_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2025-06-11/tagrules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func datadogMonitorTagRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"datadog_monitor_id": {
			Description:  "The Datadog Monitor ID which should be used for this Datadog Monitor Tag Rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: tagrules.ValidateMonitorID,
		},

		"name": {
			Description: "The name of the Tag Rules configuration.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "default",
		},

		"log": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"aad_log_enabled": {
						Description: "Whether AAD logs should be sent for the Monitor resource?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},

					"subscription_log_enabled": {
						Description: "Whether Azure subscription logs should be sent for the Monitor resource?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},

					"resource_log_enabled": {
						Description: "Whether Azure resource logs should be sent for the Monitor resource?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},

					"filter": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description: "Name of the Tag.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
								"value": {
									Description: "Value of the Tag.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
								"action": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(tagrules.PossibleValuesForTagAction(), false),
								},
							},
						},
					},
				},
			},
		},

		"metric": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"filter": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description: "Name of the Tag.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
								"value": {
									Description: "Value of the Tag.",
									Type:        pluginsdk.TypeString,
									Required:    true,
								},
								"action": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(tagrules.PossibleValuesForTagAction(), false),
								},
							},
						},
					},
				},
			},
		},
	}
}
