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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"target_resource_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"predictive": {
			Type:     pluginsdk.TypeList,
			MaxItems: 1,
			MinItems: 1,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"scale_mode": {
						Type:     pluginsdk.TypeString,
						Required: true,
						// Disabled is not exposed, omission of this block to mean disabled
						ValidateFunc: validation.StringInSlice([]string{
							string(autoscalesettings.PredictiveAutoscalePolicyScaleModeEnabled),
							string(autoscalesettings.PredictiveAutoscalePolicyScaleModeForecastOnly),
						}, false),
					},

					"look_ahead_time": {
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
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(0, 1000),
								},
								"maximum": {
									Type:         pluginsdk.TypeInt,
									Required:     true,
									ValidateFunc: validation.IntBetween(0, 1000),
								},
								"default": {
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
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},
											"metric_resource_id": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: azure.ValidateResourceID,
											},
											"time_grain": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.ISO8601Duration,
											},
											"statistic": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.StringInSlice([]string{
													string(autoscalesettings.MetricStatisticTypeAverage),
													string(autoscalesettings.MetricStatisticTypeMax),
													string(autoscalesettings.MetricStatisticTypeMin),
													string(autoscalesettings.MetricStatisticTypeSum),
												}, false),
											},
											"time_window": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.ISO8601Duration,
											},
											"time_aggregation": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice(autoscalesettings.PossibleValuesForTimeAggregationType(), false),
											},
											"operator": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice(autoscalesettings.PossibleValuesForComparisonOperationType(), false),
											},
											"threshold": {
												Type:     pluginsdk.TypeFloat,
												Required: true,
											},

											"metric_namespace": {
												Type:         pluginsdk.TypeString,
												Optional:     true,
												ValidateFunc: validation.StringIsNotEmpty,
											},

											"divide_by_instance_count": {
												Type:     pluginsdk.TypeBool,
												Optional: true,
											},

											"dimensions": {
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
															ValidateFunc: validation.StringInSlice(autoscalesettings.PossibleValuesForScaleRuleMetricDimensionOperationType(), false),
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
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.StringInSlice([]string{
													string(autoscalesettings.ScaleDirectionDecrease),
													string(autoscalesettings.ScaleDirectionIncrease),
												}, false),
											},
											"type": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validation.StringInSlice(autoscalesettings.PossibleValuesForScaleType(), false),
											},
											"value": {
												Type:         pluginsdk.TypeInt,
												Required:     true,
												ValidateFunc: validation.IntAtLeast(0),
											},
											"cooldown": {
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
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"timezone": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      "UTC",
									ValidateFunc: validateAutoScaleSettingsTimeZone(),
								},
								"start": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.IsRFC3339Time,
								},
								"end": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.IsRFC3339Time,
								},
							},
						},
					},
					"recurrence": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"timezone": {
									Type:         pluginsdk.TypeString,
									Optional:     true,
									Default:      "UTC",
									ValidateFunc: validateAutoScaleSettingsTimeZone(),
								},
								"days": {
									Type:     pluginsdk.TypeList,
									Required: true,
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
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeInt,
										ValidateFunc: validation.IntBetween(0, 23),
									},
								},
								"minutes": {
									Type:     pluginsdk.TypeList,
									Required: true,
									MaxItems: 1,
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
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"send_to_subscription_administrator": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
								"send_to_subscription_co_administrator": {
									Type:     pluginsdk.TypeBool,
									Optional: true,
									Default:  false,
								},
								"custom_emails": {
									Type:     pluginsdk.TypeList,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type:         pluginsdk.TypeString,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
						AtLeastOneOf: []string{"notification.0.email", "notification.0.webhook"},
					},
					"webhook": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"service_uri": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ValidateFunc: validation.StringIsNotEmpty,
								},
								"properties": {
									Type:     pluginsdk.TypeMap,
									Optional: true,
									Elem: &pluginsdk.Schema{
										Type: pluginsdk.TypeString,
									},
								},
							},
						},
						AtLeastOneOf: []string{"notification.0.email", "notification.0.webhook"},
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}
