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
			Description:  "The name of the scheduled query rule.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringDoesNotContainAny("<>*%&:\\?+/"),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"authorized_resource_ids": {
			Description: "List of Resource IDs referred into query.",
			Type:        pluginsdk.TypeSet,
			Optional:    true,
			MaxItems:    100,
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
						Description: "List of action group reference resource IDs.",
						Type:        pluginsdk.TypeSet,
						Required:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
					"custom_webhook_payload": {
						Description:  "Custom payload to be sent for all webhook payloads in alerting action.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsJSON,
					},
					"email_subject": {
						Description:  "Custom subject override for all email IDs in Azure action group.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},
		"data_source_id": {
			Description:  "The resource URI over which log search query is to be run.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},
		"auto_mitigation_enabled": {
			Description:   "Should the alerts in this Metric Alert be auto resolved?",
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Default:       false,
			ConflictsWith: []string{"throttling"},
		},
		"description": {
			Description:  "The description of the scheduled query rule.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 4096),
		},
		"enabled": {
			Description: "Whether this scheduled query rule is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},
		"frequency": {
			Description:  "Frequency (in minutes) at which rule condition should be evaluated.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(5, 1440),
		},
		"query": {
			Description:  "Log search query.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"query_type": {
			Description: "The type of query results. If set to `ResultCount`, `query` must include an `AggregatedValue` column of a numeric type, for example, `Heartbeat | summarize AggregatedValue = count() by bin(TimeGenerated, 5m)`.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "ResultCount",
			ValidateFunc: validation.StringInSlice([]string{
				"ResultCount",
				"Number",
			}, false),
		},
		"severity": {
			Description:  "Severity of the alert.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			ValidateFunc: validation.IntBetween(0, 4),
		},
		"throttling": {
			Description:   "Time (in minutes) for which Alerts should be throttled or suppressed.",
			Type:          pluginsdk.TypeInt,
			Optional:      true,
			ValidateFunc:  validation.IntBetween(0, 10000),
			ConflictsWith: []string{"auto_mitigation_enabled"},
		},
		"time_window": {
			Description:  "Time window for which data needs to be fetched for query (must be greater than or equal to `frequency`).",
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
						Description: "Trigger condition for metric query rule.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						MaxItems:    1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"metric_trigger_type": {
									Description: "Metric Trigger Type.",
									Type:        pluginsdk.TypeString,
									Required:    true,
									ValidateFunc: validation.StringInSlice([]string{
										"Consecutive",
										"Total",
									}, false),
								},
								"operator": {
									Description: "Evaluation operation for rule.",
									Type:        pluginsdk.TypeString,
									Required:    true,
									ValidateFunc: validation.StringInSlice([]string{
										"GreaterThan",
										"GreaterThanOrEqual",
										"LessThan",
										"LessThanOrEqual",
										"Equal",
									}, false),
								},
								"threshold": {
									Description:  "The threshold of the metric trigger.",
									Type:         pluginsdk.TypeFloat,
									Required:     true,
									ValidateFunc: validate.ScheduledQueryRulesAlertThreshold,
								},
								"metric_column": {
									Description:  "Evaluation of metric on a particular column.",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
							},
						},
					},
					"operator": {
						Description: "Evaluation operation for rule.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"GreaterThan",
							"GreaterThanOrEqual",
							"LessThan",
							"LessThanOrEqual",
							"Equal",
						}, false),
					},
					"threshold": {
						Description:  "Result or count threshold based on which rule should be triggered.",
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
