// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_autoscale_setting

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-10-01/autoscalesettings"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func monitorAutoscaleSettingSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the AutoScale Setting.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"target_resource_id": {
			Description:  "The resource ID of the resource that the autoscale setting should be added to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"enabled": {
			Description: "Whether automatic scaling is enabled for the target resource.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			Default:     true,
		},

		"predictive": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"scale_mode": {
						Description: "The predictive scale mode.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						// Disabled is not exposed, omission of this block to mean disabled
						ValidateFunc: validation.StringInSlice([]string{
							string(autoscalesettings.PredictiveAutoscalePolicyScaleModeEnabled),
							string(autoscalesettings.PredictiveAutoscalePolicyScaleModeForecastOnly),
						}, false),
					},

					"look_ahead_time": {
						Description:  "The amount of time by which instances are launched in advance.",
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.ISO8601DurationBetween("PT1M", "PT1H"),
					},
				},
			},
		},

		"profile": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MaxItems: 20,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Description:  "The name of the profile.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"capacity": {
						Type:     pluginsdk.TypeList,
						Required: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"minimum": {
									Description:  "The minimum number of instances for this resource.",
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(0, 1000),
								},
								"maximum": {
									Description:  "The maximum number of instances for this resource.",
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(0, 1000),
								},
								"default": {
									Description:  "The number of instances that are available for scaling if metrics are not available for evaluation.",
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(0, 1000),
								},
							},
						},
					},
					"rule": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 10,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"metric_trigger": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"metric_name": {
												Description:  "The name of the metric that defines what the rule monitors, such as `Percentage CPU` for `Virtual Machine Scale Sets` and `CpuPercentage` for `App Service Plan`.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
											"metric_resource_id": {
												Description:  "The ID of the Resource which the Rule monitors.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: azure.ValidateResourceID,
											},
											"time_grain": {
												Description:  "The granularity of metrics that the rule monitors, which must be one of the pre-defined values returned from the metric definitions for the metric. This value must be between 1 minute and 12 hours and be formatted as an ISO 8601 string.", // azignore:AZS014 - TODO 6.0 move the range from the documentation to validation
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.ISO8601Duration,
											},
											"statistic": {
												Description: "How the metrics from multiple instances are combined.",
												Type:        pluginsdk.TypeString,
												Required:    true,
												ValidateFunc: validation.StringInSlice([]string{
													string(autoscalesettings.MetricStatisticTypeAverage),
													string(autoscalesettings.MetricStatisticTypeMax),
													string(autoscalesettings.MetricStatisticTypeMin),
													string(autoscalesettings.MetricStatisticTypeSum),
												}, false),
											},
											"time_window": {
												Description:  "The time range for which data is collected, which must be greater than the delay in metric collection (which varies from resource to resource). This value must be between 5 minutes and 12 hours and be formatted as an ISO 8601 string.", // azignore:AZS014 - TODO 6.0 move the range from the documentation to validation
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.ISO8601Duration,
											},
											"time_aggregation": {
												Description:  "How the data that's collected should be combined over time.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice(autoscalesettings.PossibleValuesForTimeAggregationType(), false),
											},
											"operator": {
												Description:  "The operator used to compare the metric data and threshold.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice(autoscalesettings.PossibleValuesForComparisonOperationType(), false),
											},
											"threshold": {
												Description: "The threshold of the metric that triggers the scale action.",
												Type:        pluginsdk.TypeFloat,
												Required:    true,
											},

											"metric_namespace": {
												Description:  "The namespace of the metric that defines what the rule monitors, such as `microsoft.compute/virtualmachinescalesets` for `Virtual Machine Scale Sets`.",
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"divide_by_instance_count": {
												Description: "Whether to enable metric divide by instance count.",
												Type:        pluginsdk.TypeBool,
												Optional:    true,
											},

											"dimensions": {
												Type:     pluginsdk.TypeList,
												Optional: true,
												Elem: &pluginsdk.Resource{
													Schema: map[string]*pluginsdk.Schema{
														"name": {
															Description:  "The name of the dimension.",
															Type:         pluginsdk.TypeString,
															Required:     true,
															ValidateFunc: validation.StringIsNotEmpty,
														},

														"operator": {
															Description:  "The dimension operator. `Equals` means being equal to any of the values. `NotEquals` means being not equal to any of the values.",
															Type:         pluginsdk.TypeString,
															Required:     true,
															ValidateFunc: validation.StringInSlice(autoscalesettings.PossibleValuesForScaleRuleMetricDimensionOperationType(), false),
														},

														"values": {
															Description: "A list of dimension values.",
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
										},
									},
								},
								"scale_action": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &pluginsdk.Resource{
										Schema: map[string]*pluginsdk.Schema{
											"direction": {
												Description: "The scale direction.",
												Type:        pluginsdk.TypeString,
												Required:    true,
												ValidateFunc: validation.StringInSlice([]string{
													string(autoscalesettings.ScaleDirectionDecrease),
													string(autoscalesettings.ScaleDirectionIncrease),
												}, false),
											},
											"type": {
												Description:  "The type of action that should occur.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice(autoscalesettings.PossibleValuesForScaleType(), false),
											},
											"value": {
												Description:  "The number of instances involved in the scaling action.",
												Type:         pluginsdk.TypeInt,
												Required:     true,
												ValidateFunc: validation.IntAtLeast(0),
											},
											"cooldown": {
												Description:  "The amount of time to wait since the last scaling action before this action occurs.",
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.ISO8601Duration,
											},
										},
									},
								},
							},
						},
					},
					"fixed_date": {
						Description: "This cannot be specified if a `recurrence` block is specified.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						MaxItems:    1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"timezone": {
									Description:  "The Time Zone of the `start` and `end` times. A list of [possible values can be found here](https://learn.microsoft.com/rest/api/monitor/autoscale-settings/create-or-update?view=rest-monitor-2022-10-01&tabs=HTTP#recurrentschedule).",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      "UTC",
									ValidateFunc: validateAutoScaleSettingsTimeZone(),
								},
								"start": {
									Description:  "The start date for the profile, formatted as an RFC3339 date string.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.IsRFC3339Time,
								},
								"end": {
									Description:  "The end date for the profile, formatted as an RFC3339 date string.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.IsRFC3339Time,
								},
							},
						},
					},
					"recurrence": {
						Description: "This cannot be specified if a `fixed_date` block is specified.",
						Type:        pluginsdk.TypeList,
						Optional:    true,
						MaxItems:    1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"timezone": {
									Description:  "The Time Zone used for the `hours` field. A list of [possible values can be found here](https://learn.microsoft.com/rest/api/monitor/autoscale-settings/create-or-update?view=rest-monitor-2022-10-01&tabs=HTTP#recurrentschedule).",
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      "UTC",
									ValidateFunc: validateAutoScaleSettingsTimeZone(),
								},
								"days": {
									Description: "A list of days that this profile takes effect on.",
									Type:        pluginsdk.TypeList,
									Required:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
										ValidateFunc: validation.StringInSlice([]string{
											"Monday",
											"Tuesday",
											"Wednesday",
											"Thursday",
											"Friday",
											"Saturday",
											"Sunday",
										}, false),
									},
								},
								"hours": {
									Description: "A list containing a single item, which specifies the Hour interval at which this recurrence should be triggered (in 24-hour time).",
									Type:        pluginsdk.TypeList,
									Required:    true,
									MaxItems:    1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeInt,
										ValidateFunc: validation.IntBetween(0, 23),
									},
								},
								"minutes": {
									Description: "A list containing a single item which specifies the Minute interval at which this recurrence should be triggered.",
									Type:        pluginsdk.TypeList,
									Required:    true,
									MaxItems:    1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeInt,
										ValidateFunc: validation.IntBetween(0, 59),
									},
								},
							},
						},
					},
				},
			},
		},

		"notification": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"email": {
						Type:         pluginsdk.TypeList,
						Optional:     true,
						MaxItems:     1,
						AtLeastOneOf: []string{"notification.0.email", "notification.0.webhook"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"send_to_subscription_administrator": {
									Description: "Whether email notifications should be sent to the subscription administrator.",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									Default:     false,
								},
								"send_to_subscription_co_administrator": {
									Description: "Whether email notifications should be sent to the subscription co-administrator.",
									Type:        pluginsdk.TypeBool,
									Optional:    true,
									Default:     false,
								},
								"custom_emails": {
									Description: "A list of custom email addresses to which the email notifications will be sent.",
									Type:        pluginsdk.TypeList,
									Optional:    true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
					},
					"webhook": {
						Type:         pluginsdk.TypeList,
						Optional:     true,
						AtLeastOneOf: []string{"notification.0.email", "notification.0.webhook"},
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"service_uri": {
									Description:  "The HTTPS URI which should receive scale notifications.",
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"properties": {
									Description: "A map of settings.",
									Type:        pluginsdk.TypeMap,
									Optional:    true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
