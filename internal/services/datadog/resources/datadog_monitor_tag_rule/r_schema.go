// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package datadog_monitor_tag_rule

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/monitorsresource"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func datadogMonitorTagRuleSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"datadog_monitor_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: monitorsresource.ValidateMonitorID,
		},

		"name": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "default",
		},

		"log": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"aad_log_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"subscription_log_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"resource_log_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
					},

					"filter": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"value": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"action": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(rules.PossibleValuesForTagAction(), false),
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
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"value": {
									Type:     pluginsdk.TypeString,
									Required: true,
								},
								"action": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(rules.PossibleValuesForTagAction(), false),
								},
							},
						},
					},
				},
			},
		},
	}
}
