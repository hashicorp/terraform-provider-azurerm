// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_scheduled_query_rules_alert

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func monitorScheduledQueryRulesAlertSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringDoesNotContainAny("<>*%&:\\?+/"),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"authorized_resource_ids": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			MaxItems: 100,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: azure.ValidateResourceID,
			},
		},
		"action": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"action_group": {
						Type:     pluginsdk.TypeSet,
						Required: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
					"custom_webhook_payload": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsJSON,
					},
					"email_subject": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
		"data_source_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},
		"auto_mitigation_enabled": {
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Default:       false,
			ConflictsWith: []string{"throttling"},
		},
		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 4096),
		},
		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},
		"frequency": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(5, 1440),
		},
		"query": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"query_type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  "ResultCount",
			ValidateFunc: validation.StringInSlice([]string{
				"ResultCount",
				"Number",
			}, false),
		},
		"severity": {
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 4),
		},
		"throttling": {
			Type:          pluginsdk.TypeInt,
			Optional:      true,
			ValidateFunc:  validation.IntBetween(0, 10000),
			ConflictsWith: []string{"auto_mitigation_enabled"},
		},
		"time_window": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(5, 2880),
		},
		"trigger": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"metric_trigger": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"metric_trigger_type": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										"Consecutive",
										"Total",
									}, false),
								},
								"operator": {
									Type:     pluginsdk.TypeString,
									Required: true,
									ValidateFunc: validation.StringInSlice([]string{
										"GreaterThan",
										"GreaterThanOrEqual",
										"LessThan",
										"LessThanOrEqual",
										"Equal",
									}, false),
								},
								"threshold": {
									Type:         pluginsdk.TypeFloat,
									Required:     true,
									ValidateFunc: validate.ScheduledQueryRulesAlertThreshold,
								},
								"metric_column": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
					"operator": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							"GreaterThan",
							"GreaterThanOrEqual",
							"LessThan",
							"LessThanOrEqual",
							"Equal",
						}, false),
					},
					"threshold": {
						Type:         pluginsdk.TypeFloat,
						Required:     true,
						ValidateFunc: validate.ScheduledQueryRulesAlertThreshold,
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
