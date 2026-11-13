// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_scheduled_query_rules_alert_v2

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2023-03-15-preview/scheduledqueryrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"criteria": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"query": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"operator": {
						Type:     pluginsdk.TypeString,
						Required: true,
						ValidateFunc: validation.StringInSlice([]string{
							// see https://github.com/Azure/azure-rest-api-specs/issues/21794
							"Equal",
							string(scheduledqueryrules.ConditionOperatorGreaterThan),
							string(scheduledqueryrules.ConditionOperatorGreaterThanOrEqual),
							string(scheduledqueryrules.ConditionOperatorLessThan),
							string(scheduledqueryrules.ConditionOperatorLessThanOrEqual),
						}, false),
					},

					"time_aggregation_method": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(scheduledqueryrules.PossibleValuesForTimeAggregation(), false),
					},

					"threshold": {
						Type:     pluginsdk.TypeFloat,
						Required: true,
					},

					"dimension": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"operator": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(scheduledqueryrules.PossibleValuesForDimensionOperator(), false),
								},

								"values": {
									Type:     pluginsdk.TypeList,
									Required: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},

					"failing_periods": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"minimum_failing_periods_to_trigger_alert": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 6),
								},

								"number_of_evaluation_periods": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 6),
								},
							},
						},
					},

					"metric_measure_column": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"resource_id_column": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		// lintignore:S013
		"evaluation_frequency": {
			Type: pluginsdk.TypeString,
			// this field is required, missing this field will get an error from service
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"PT1M",
				"PT5M",
				"PT10M",
				"PT15M",
				"PT30M",
				"PT45M",
				"PT1H",
				"PT2H",
				"PT3H",
				"PT4H",
				"PT5H",
				"PT6H",
				"P1D",
			}, false),
		},

		"scopes": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			ForceNew: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: azure.ValidateResourceID,
			},
		},

		"severity": {
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(0, 4),
		},

		"window_duration": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"PT1M",
				"PT5M",
				"PT10M",
				"PT15M",
				"PT30M",
				"PT45M",
				"PT1H",
				"PT2H",
				"PT3H",
				"PT4H",
				"PT5H",
				"PT6H",
				"P1D",
				"P2D",
			}, false),
		},

		"action": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"action_groups": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: azure.ValidateResourceID,
						},
					},

					"custom_properties": {
						Type:     pluginsdk.TypeMap,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
					"email_subject": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"auto_mitigation_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"workspace_alerts_storage_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"mute_actions_after_alert_duration": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"PT5M",
				"PT10M",
				"PT15M",
				"PT30M",
				"PT45M",
				"PT1H",
				"PT2H",
				"PT3H",
				"PT4H",
				"PT5H",
				"PT6H",
				"P1D",
				"P2D",
			}, false),
		},

		"query_time_range_override": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ValidateFunc: validation.StringInSlice([]string{
				"PT5M",
				"PT10M",
				"PT15M",
				"PT20M",
				"PT30M",
				"PT45M",
				"PT1H",
				"PT2H",
				"PT3H",
				"PT4H",
				"PT5H",
				"PT6H",
				"P1D",
				"P2D",
			}, false),
		},

		"skip_query_validation": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
		},

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"tags": commonschema.Tags(),

		"target_resource_types": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"created_with_api_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"is_a_legacy_log_analytics_rule": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"is_workspace_alerts_storage_configured": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},
	}
}
