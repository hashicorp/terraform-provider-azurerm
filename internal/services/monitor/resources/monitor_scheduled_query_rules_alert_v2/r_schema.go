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
			Description:  "The name which should be used for this Monitor Scheduled Query Rule.",
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
						Description:  "The query to run on logs. The results returned by this query are used to populate the alert.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"operator": {
						Description: "The criteria operator.",
						Type:        pluginsdk.TypeString,
						Required:    true,
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
						Description:  "The type of aggregation to apply to the data points in aggregation granularity.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(scheduledqueryrules.PossibleValuesForTimeAggregation(), false),
					},

					"threshold": {
						Description: "The criteria threshold value that activates the alert.",
						Type:        pluginsdk.TypeFloat,
						Required:    true,
					},

					"dimension": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description:  "Name of the dimension.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},

								"operator": {
									Description:  "Operator for dimension values.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringInSlice(scheduledqueryrules.PossibleValuesForDimensionOperator(), false),
								},

								"values": {
									Description: "List of dimension values. Use a wildcard `*` to collect all.",
									Type:        pluginsdk.TypeList,
									Required:    true,
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
									Description:  "The number of violations to trigger an alert.",
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 6),
								},

								"number_of_evaluation_periods": {
									Description:  "The number of aggregated look-back points. The look-back time window is calculated based on the aggregation granularity `window_duration` and the selected number of aggregated points.",
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(1, 6),
								},
							},
						},
					},

					"metric_measure_column": {
						Description:  "The column containing the metric measure number.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"resource_id_column": {
						Description:  "The column containing the resource ID. The content of the column must be a URI formatted as resource ID.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		// lintignore:S013
		"evaluation_frequency": {
			Description: "How often the scheduled query rule is evaluated, represented in ISO 8601 duration format.",
			Type:        pluginsdk.TypeString,
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
			Description: "The list of resource IDs that this scheduled query rule is scoped to. Currently, the API supports exactly 1 resource ID in the scopes list.",
			Type:        pluginsdk.TypeList,
			Required:    true,
			ForceNew:    true,
			MinItems:    1,
			MaxItems:    1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: azure.ValidateResourceID,
			},
		},

		"severity": {
			Description:  "Severity of the alert. Severity decreases as the value increases.",
			Type:         pluginsdk.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(0, 4),
		},

		"window_duration": {
			Description: "The period of time in ISO 8601 duration format on which the Scheduled Query Rule will be executed (bin size).",
			Type:        pluginsdk.TypeString,
			Required:    true,
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
						Description: "List of Action Group resource IDs to invoke when the alert fires.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: azure.ValidateResourceID,
						},
					},

					"custom_properties": {
						Description: "The properties of an alert payload.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
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

		"auto_mitigation_enabled": {
			Description: "The flag that indicates whether the alert should be automatically resolved.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"workspace_alerts_storage_enabled": {
			Description: "The flag which indicates whether this scheduled query rule check if storage is configured.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     false,
		},

		"description": {
			Description: "The description of the scheduled query rule.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"display_name": {
			Description:  "The display name of the alert rule.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"enabled": {
			Description: "The flag which indicates whether this scheduled query rule is enabled.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"mute_actions_after_alert_duration": {
			Description: "Mute actions for the chosen period of time in ISO 8601 duration format after the alert is fired.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
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
			Description: "Set this if the alert evaluation period is different from the query time range. If not specified, the value is `window_duration`*`number_of_evaluation_periods`.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
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
			Description: "The flag which indicates whether the provided query should be validated.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},

		"identity": commonschema.SystemOrUserAssignedIdentityOptional(),

		"tags": commonschema.Tags(),

		"target_resource_types": {
			Description: "List of resource type of the target resource(s) on which the alert is created/updated. For example if the scope is a Resource Group and targetResourceTypes is `Microsoft.Compute/virtualMachines`, then a different alert will be fired for each Virtual Machine in the Resource Group which meet the alert criteria.",
			Type:        pluginsdk.TypeList,
			Optional:    true,
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
			Description: "The api-version used when creating this alert rule.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"is_a_legacy_log_analytics_rule": {
			Description: "True if this alert rule is a legacy Log Analytic Rule.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},

		"is_workspace_alerts_storage_configured": {
			Description: "The flag indicates whether this Scheduled Query Rule has been configured to be stored in the customer's storage.",
			Type:        pluginsdk.TypeBool,
			Computed:    true,
		},
	}
}
