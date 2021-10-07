package migration

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2021-07-01-preview/insights"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ pluginsdk.StateUpgrade = AutoscaleSettingUpgradeV0ToV1{}

type AutoscaleSettingUpgradeV0ToV1 struct{}

func (AutoscaleSettingUpgradeV0ToV1) Schema() map[string]*pluginsdk.Schema {
	return autoscaleSettingSchemaForV0AndV1()
}

func (AutoscaleSettingUpgradeV0ToV1) UpgradeFunc() pluginsdk.StateUpgraderFunc {
	return func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
		// old
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/microsoft.insights/autoscalesettings/{settingName}
		// new:
		// 	/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/autoscaleSettings/{settingName}
		oldId, err := azure.ParseAzureResourceID(rawState["id"].(string))
		if err != nil {
			return rawState, err
		}

		settingName := ""
		for key, value := range oldId.Path {
			if strings.EqualFold(key, "autoscaleSettings") {
				settingName = value
				break
			}
		}

		if settingName == "" {
			return rawState, fmt.Errorf("couldn't find the `autoscaleSettings` segment in the old resource id %q", oldId)
		}

		newId := parse.NewAutoscaleSettingID(oldId.SubscriptionID, oldId.ResourceGroup, settingName)

		log.Printf("[DEBUG] Updating ID from %q to %q", oldId, newId.ID())

		rawState["id"] = newId.ID()

		return rawState, nil
	}
}

func autoscaleSettingSchemaForV0AndV1() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"location": azure.SchemaLocation(),

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
												ValidateFunc: validate.ISO8601Duration,
											},
											"statistic": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.StringInSlice([]string{
													string(insights.MetricStatisticTypeAverage),
													string(insights.MetricStatisticTypeMax),
													string(insights.MetricStatisticTypeMin),
													string(insights.MetricStatisticTypeSum),
												}, true),
												DiffSuppressFunc: suppress.CaseDifference,
											},
											"time_window": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validate.ISO8601Duration,
											},
											"time_aggregation": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.StringInSlice([]string{
													string(insights.TimeAggregationTypeAverage),
													string(insights.TimeAggregationTypeCount),
													string(insights.TimeAggregationTypeMaximum),
													string(insights.TimeAggregationTypeMinimum),
													string(insights.TimeAggregationTypeTotal),
													string(insights.TimeAggregationTypeLast),
												}, true),
												DiffSuppressFunc: suppress.CaseDifference,
											},
											"operator": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.StringInSlice([]string{
													string(insights.ComparisonOperationTypeEquals),
													string(insights.ComparisonOperationTypeGreaterThan),
													string(insights.ComparisonOperationTypeGreaterThanOrEqual),
													string(insights.ComparisonOperationTypeLessThan),
													string(insights.ComparisonOperationTypeLessThanOrEqual),
													string(insights.ComparisonOperationTypeNotEquals),
												}, true),
												DiffSuppressFunc: suppress.CaseDifference,
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
															Type:     pluginsdk.TypeString,
															Required: true,
															ValidateFunc: validation.StringInSlice([]string{
																string(insights.ScaleRuleMetricDimensionOperationTypeEquals),
																string(insights.ScaleRuleMetricDimensionOperationTypeNotEquals),
															}, false),
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
													string(insights.ScaleDirectionDecrease),
													string(insights.ScaleDirectionIncrease),
												}, true),
												DiffSuppressFunc: suppress.CaseDifference,
											},
											"type": {
												Type:     pluginsdk.TypeString,
												Required: true,
												ValidateFunc: validation.StringInSlice([]string{
													string(insights.ScaleTypeChangeCount),
													string(insights.ScaleTypeExactCount),
													string(insights.ScaleTypePercentChangeCount),
												}, true),
												DiffSuppressFunc: suppress.CaseDifference,
											},
											"value": {
												Type:         pluginsdk.TypeInt,
												Required:     true,
												ValidateFunc: validation.IntAtLeast(0),
											},
											"cooldown": {
												Type:         pluginsdk.TypeString,
												Required:     true,
												ValidateFunc: validate.ISO8601Duration,
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
										}, true),
										DiffSuppressFunc: suppress.CaseDifference,
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
										Type: pluginsdk.TypeString,
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

		"tags": tags.Schema(),
	}
}

func validateAutoScaleSettingsTimeZone() pluginsdk.SchemaValidateFunc {
	// from https://docs.microsoft.com/en-us/rest/api/monitor/autoscalesettings/createorupdate#timewindow
	timeZones := []string{
		"Dateline Standard Time",
		"UTC-11",
		"Hawaiian Standard Time",
		"Alaskan Standard Time",
		"Pacific Standard Time (Mexico)",
		"Pacific Standard Time",
		"US Mountain Standard Time",
		"Mountain Standard Time (Mexico)",
		"Mountain Standard Time",
		"Central America Standard Time",
		"Central Standard Time",
		"Central Standard Time (Mexico)",
		"Canada Central Standard Time",
		"SA Pacific Standard Time",
		"Eastern Standard Time",
		"US Eastern Standard Time",
		"Venezuela Standard Time",
		"Paraguay Standard Time",
		"Atlantic Standard Time",
		"Central Brazilian Standard Time",
		"SA Western Standard Time",
		"Pacific SA Standard Time",
		"Newfoundland Standard Time",
		"E. South America Standard Time",
		"Argentina Standard Time",
		"SA Eastern Standard Time",
		"Greenland Standard Time",
		"Montevideo Standard Time",
		"Bahia Standard Time",
		"UTC-02",
		"Mid-Atlantic Standard Time",
		"Azores Standard Time",
		"Cape Verde Standard Time",
		"Morocco Standard Time",
		"UTC",
		"GMT Standard Time",
		"Greenwich Standard Time",
		"W. Europe Standard Time",
		"Central Europe Standard Time",
		"Romance Standard Time",
		"Central European Standard Time",
		"W. Central Africa Standard Time",
		"Namibia Standard Time",
		"Jordan Standard Time",
		"GTB Standard Time",
		"Middle East Standard Time",
		"Egypt Standard Time",
		"Syria Standard Time",
		"E. Europe Standard Time",
		"South Africa Standard Time",
		"FLE Standard Time",
		"Turkey Standard Time",
		"Israel Standard Time",
		"Kaliningrad Standard Time",
		"Libya Standard Time",
		"Arabic Standard Time",
		"Arab Standard Time",
		"Belarus Standard Time",
		"Russian Standard Time",
		"E. Africa Standard Time",
		"Iran Standard Time",
		"Arabian Standard Time",
		"Azerbaijan Standard Time",
		"Russia Time Zone 3",
		"Mauritius Standard Time",
		"Georgian Standard Time",
		"Caucasus Standard Time",
		"Afghanistan Standard Time",
		"West Asia Standard Time",
		"Ekaterinburg Standard Time",
		"Pakistan Standard Time",
		"India Standard Time",
		"Sri Lanka Standard Time",
		"Nepal Standard Time",
		"Central Asia Standard Time",
		"Bangladesh Standard Time",
		"N. Central Asia Standard Time",
		"Myanmar Standard Time",
		"SE Asia Standard Time",
		"North Asia Standard Time",
		"China Standard Time",
		"North Asia East Standard Time",
		"Singapore Standard Time",
		"W. Australia Standard Time",
		"Taipei Standard Time",
		"Ulaanbaatar Standard Time",
		"Tokyo Standard Time",
		"Korea Standard Time",
		"Yakutsk Standard Time",
		"Cen. Australia Standard Time",
		"AUS Central Standard Time",
		"E. Australia Standard Time",
		"AUS Eastern Standard Time",
		"West Pacific Standard Time",
		"Tasmania Standard Time",
		"Magadan Standard Time",
		"Vladivostok Standard Time",
		"Russia Time Zone 10",
		"Central Pacific Standard Time",
		"Russia Time Zone 11",
		"New Zealand Standard Time",
		"UTC+12",
		"Fiji Standard Time",
		"Kamchatka Standard Time",
		"Tonga Standard Time",
		"Samoa Standard Time",
		"Line Islands Standard Time",
	}
	return validation.StringInSlice(timeZones, false)
}
