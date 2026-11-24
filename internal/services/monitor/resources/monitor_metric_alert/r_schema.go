// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_metric_alert

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-02-02/componentsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2022-06-15/webtestsapis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-03-01/metricalerts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func monitorMetricAlertSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Metric Alert.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"scopes": {
			Description: "A set of strings of resource IDs at which the metric criteria should be applied.",
			Type:        pluginsdk.TypeSet,
			Required:    true,
			MinItems:    1,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: azure.ValidateResourceID,
			},
			Set: pluginsdk.HashString,
		},

		"target_resource_type": {
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Computed:    true, // azignore:AZS007 - pre-existing violation
			Description: "The resource type (e.g. `Microsoft.Compute/virtualMachines`) of the target resource.",
		},

		"target_resource_location": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true, // azignore:AZS007 - pre-existing violation
			StateFunc:        location.StateFunc,
			DiffSuppressFunc: location.DiffSuppressFunc,
			Description:      "The location of the target resource.",
		},

		// static criteria
		"criteria": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MinItems:     1,
			ExactlyOneOf: []string{"criteria", "dynamic_criteria", "application_insights_web_test_location_availability_criteria"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"metric_namespace": {
						Description:  "One of the metric namespaces to be monitored.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"metric_name": {
						Description:  "One of the metric names to be monitored.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"aggregation": {
						Description: "The statistic that runs over the metric values.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"Average",
							"Count",
							"Minimum",
							"Maximum",
							"Total",
						}, false),
					},
					"dimension": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description:  "One of the dimension names.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"operator": {
									Description: "The dimension operator.",
									Type:        pluginsdk.TypeString,
									Required:    true,
									ValidateFunc: validation.StringInSlice([]string{
										"Include",
										"Exclude",
										"StartsWith",
									}, false),
								},
								"values": {
									Description: "The list of dimension values.",
									Type:        pluginsdk.TypeList,
									Required:    true,
									MinItems:    1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
					"operator": {
						Description:  "The criteria operator.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(metricalerts.PossibleValuesForOperator(), false),
					},
					"threshold": {
						Description: "The criteria threshold value that activates the alert.",
						Type:        pluginsdk.TypeFloat,
						Required:    true,
					},
					"skip_metric_validation": {
						Description: "Skip the metric validation to allow creating an alert rule on a custom metric that isn't yet emitted?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
						Default:     false,
					},
				},
			},
		},

		// lintignore: S018
		"dynamic_criteria": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			// Currently, it allows to define only one dynamic criteria in one metric alert.
			MaxItems:     1,
			ExactlyOneOf: []string{"criteria", "dynamic_criteria", "application_insights_web_test_location_availability_criteria"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"metric_namespace": {
						Description:  "One of the metric namespaces to be monitored.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"metric_name": {
						Description:  "One of the metric names to be monitored.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"aggregation": {
						Description: "The statistic that runs over the metric values.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							"Average",
							"Count",
							"Minimum",
							"Maximum",
							"Total",
						}, false),
					},
					"dimension": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"name": {
									Description:  "One of the dimension names.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"operator": {
									Description: "The dimension operator.",
									Type:        pluginsdk.TypeString,
									Required:    true,
									ValidateFunc: validation.StringInSlice([]string{
										"Include",
										"Exclude",
										"StartsWith",
									}, false),
								},
								"values": {
									Description: "The list of dimension values.",
									Type:        pluginsdk.TypeList,
									Required:    true,
									MinItems:    1,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
					"operator": {
						Description:  "The criteria operator.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(metricalerts.PossibleValuesForDynamicThresholdOperator(), false),
					},
					"alert_sensitivity": {
						Description:  "The extent of deviation required to trigger an alert.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringInSlice(metricalerts.PossibleValuesForDynamicThresholdSensitivity(), false),
					},

					"evaluation_total_count": {
						Description:  "The number of aggregated lookback points. The lookback time window is calculated based on the aggregation granularity (`window_size`) and the selected number of aggregated points.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(1),
						Default:      4,
					},

					"evaluation_failure_count": {
						Description:  "The number of violations to trigger an alert.",
						Type:         pluginsdk.TypeInt,
						Optional:     true,
						ValidateFunc: validation.IntAtLeast(1),
						Default:      4,
					},

					"ignore_data_before": {
						Description:  "The [ISO8601](https://en.wikipedia.org/wiki/ISO_8601) date from which to start learning the metric historical data and calculate the dynamic thresholds.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsRFC3339Time,
					},
					"skip_metric_validation": {
						Description: "Skip the metric validation to allow creating an alert rule on a custom metric that isn't yet emitted?",
						Type:        pluginsdk.TypeBool,
						Optional:    true,
					},
				},
			},
		},

		"application_insights_web_test_location_availability_criteria": {
			Type:         pluginsdk.TypeList,
			Optional:     true,
			MinItems:     1,
			MaxItems:     1,
			ExactlyOneOf: []string{"criteria", "dynamic_criteria", "application_insights_web_test_location_availability_criteria"},
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"web_test_id": {
						Description:  "The ID of the Application Insights Web Test.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: webtestsapis.ValidateWebTestID,
					},
					"component_id": {
						Description:  "The ID of the Application Insights Resource.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: componentsapis.ValidateComponentID,
					},
					"failed_location_count": {
						Description:  "The number of failed locations.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(1),
					},
				},
			},
		},

		"action": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"action_group_id": {
						Description:  "The ID of the Action Group can be sourced from [the `azurerm_monitor_action_group` resource](./monitor_action_group.html).",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: azure.ValidateResourceID,
					},
					"webhook_properties": {
						Description: "The map of custom string properties to include with the post operation. These data are appended to the webhook payload.",
						Type:        pluginsdk.TypeMap,
						Optional:    true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
			Set: resourceMonitorMetricAlertActionHash,
		},

		"auto_mitigate": {
			Description: "Should the alerts in this Metric Alert be auto resolved?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"description": {
			Description: "The description of this Metric Alert.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"enabled": {
			Description: "Should this Metric Alert be enabled?",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"frequency": {
			Description: "The evaluation frequency of this Metric Alert, represented in ISO 8601 duration format.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "PT1M",
			ValidateFunc: validation.StringInSlice([]string{
				"PT1M",
				"PT5M",
				"PT15M",
				"PT30M",
				"PT1H",
			}, false),
		},

		"severity": {
			Description:  "The severity of this Metric Alert.",
			Type:         pluginsdk.TypeInt,
			Optional:     true,
			Default:      3,
			ValidateFunc: validation.IntBetween(0, 4),
		},

		"window_size": {
			Description: "The period of time that is used to monitor alert activity, represented in ISO 8601 duration format. This value must be greater than `frequency`.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			Default:     "PT5M",
			ValidateFunc: validation.StringInSlice([]string{
				"PT1M",
				"PT5M",
				"PT15M",
				"PT30M",
				"PT1H",
				"PT6H",
				"PT12H",
				"P1D",
			}, false),
		},

		"tags": commonschema.Tags(),
	}
}
