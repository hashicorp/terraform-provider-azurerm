// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2022-10-01/autoscalesettings"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMonitorAutoScaleSetting() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorAutoScaleSettingCreateUpdate,
		Read:   resourceMonitorAutoScaleSettingRead,
		Update: resourceMonitorAutoScaleSettingCreateUpdate,
		Delete: resourceMonitorAutoScaleSettingDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := autoscalesettings.ParseAutoScaleSettingID(id)
			return err
		}),

		SchemaVersion: 2,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AutoscaleSettingUpgradeV0ToV1{},
			1: migration.AutoscaleSettingUpgradeV1ToV2{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
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
							ValidateFunc: validate.ISO8601DurationBetween("PT1M", "PT1H"),
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
													ValidateFunc: validate.ISO8601Duration,
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
													ValidateFunc: validate.ISO8601Duration,
												},
												"time_aggregation": {
													Type:     pluginsdk.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(autoscalesettings.TimeAggregationTypeAverage),
														string(autoscalesettings.TimeAggregationTypeCount),
														string(autoscalesettings.TimeAggregationTypeMaximum),
														string(autoscalesettings.TimeAggregationTypeMinimum),
														string(autoscalesettings.TimeAggregationTypeTotal),
														string(autoscalesettings.TimeAggregationTypeLast),
													}, false),
												},
												"operator": {
													Type:     pluginsdk.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(autoscalesettings.ComparisonOperationTypeEquals),
														string(autoscalesettings.ComparisonOperationTypeGreaterThan),
														string(autoscalesettings.ComparisonOperationTypeGreaterThanOrEqual),
														string(autoscalesettings.ComparisonOperationTypeLessThan),
														string(autoscalesettings.ComparisonOperationTypeLessThanOrEqual),
														string(autoscalesettings.ComparisonOperationTypeNotEquals),
													}, false),
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
																	string(autoscalesettings.ScaleRuleMetricDimensionOperationTypeEquals),
																	string(autoscalesettings.ScaleRuleMetricDimensionOperationTypeNotEquals),
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
														string(autoscalesettings.ScaleDirectionDecrease),
														string(autoscalesettings.ScaleDirectionIncrease),
													}, false),
												},
												"type": {
													Type:     pluginsdk.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(autoscalesettings.ScaleTypeChangeCount),
														string(autoscalesettings.ScaleTypeExactCount),
														string(autoscalesettings.ScaleTypePercentChangeCount),
														string(autoscalesettings.ScaleTypeServiceAllowedNextValue),
													}, false),
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

			"tags": tags.Schema(),
		},
	}
}

func resourceMonitorAutoScaleSettingCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.AutoscaleSettingsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := autoscalesettings.NewAutoScaleSettingID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Monitor %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_monitor_autoscale_setting", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	enabled := d.Get("enabled").(bool)
	targetResourceId := d.Get("target_resource_id").(string)

	notificationsRaw := d.Get("notification").([]interface{})
	notifications := expandAzureRmMonitorAutoScaleSettingNotifications(notificationsRaw)

	profilesRaw := d.Get("profile").([]interface{})
	profiles, err := expandAzureRmMonitorAutoScaleSettingProfile(profilesRaw)
	if err != nil {
		return fmt.Errorf("expanding `profile`: %+v", err)
	}

	t := d.Get("tags").(map[string]interface{})

	parameters := autoscalesettings.AutoscaleSettingResource{
		Location: location,
		Properties: autoscalesettings.AutoscaleSetting{
			Enabled:                   &enabled,
			Profiles:                  profiles,
			PredictiveAutoscalePolicy: expandAzureRmMonitorAutoScaleSettingPredictive(d.Get("predictive").([]interface{})),
			Notifications:             notifications,
			TargetResourceUri:         &targetResourceId,
		},
		Tags: utils.ExpandPtrMapStringString(t),
	}

	if _, err = client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating Monitor %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMonitorAutoScaleSettingRead(d, meta)
}

func resourceMonitorAutoScaleSettingRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.AutoscaleSettingsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := autoscalesettings.ParseAutoScaleSettingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Monitor %s: %+v", *id, err)
	}

	d.Set("name", id.AutoScaleSettingName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		props := model.Properties

		d.Set("location", azure.NormalizeLocation(model.Location))
		d.Set("enabled", props.Enabled)
		d.Set("target_resource_id", props.TargetResourceUri)

		profile, err := flattenAzureRmMonitorAutoScaleSettingProfile(props.Profiles)
		if err != nil {
			return fmt.Errorf("flattening `profile` of %s: %+v", *id, err)
		}
		if err = d.Set("profile", profile); err != nil {
			return fmt.Errorf("setting `profile` of %s: %+v", *id, err)
		}

		if err = d.Set("predictive", flattenAzureRmMonitorAutoScaleSettingPredictive(props.PredictiveAutoscalePolicy)); err != nil {
			return fmt.Errorf("setting `predictive_scale_mode` of %s: %+v", *id, err)
		}

		notifications := flattenAzureRmMonitorAutoScaleSettingNotification(props.Notifications)
		if err = d.Set("notification", notifications); err != nil {
			return fmt.Errorf("setting `notification` of %s: %+v", *id, err)
		}

		// Return a new tag map filtered by the specified tag names.
		tagMap := tags.Filter(model.Tags, "$type")

		if err = d.Set("tags", utils.FlattenPtrMapStringString(tagMap)); err != nil {
			return err
		}

	}
	return nil
}

func resourceMonitorAutoScaleSettingDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.AutoscaleSettingsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := autoscalesettings.ParseAutoScaleSettingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting Monitor %s: %+v", *id, err)
		}
	}

	return nil
}

func expandAzureRmMonitorAutoScaleSettingProfile(input []interface{}) ([]autoscalesettings.AutoscaleProfile, error) {
	results := make([]autoscalesettings.AutoscaleProfile, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		name := raw["name"].(string)

		// this is Required, so we don't need to check for optionals here
		capacitiesRaw := raw["capacity"].([]interface{})
		capacityRaw := capacitiesRaw[0].(map[string]interface{})
		capacity := autoscalesettings.ScaleCapacity{
			Minimum: strconv.Itoa(capacityRaw["minimum"].(int)),
			Maximum: strconv.Itoa(capacityRaw["maximum"].(int)),
			Default: strconv.Itoa(capacityRaw["default"].(int)),
		}

		recurrencesRaw := raw["recurrence"].([]interface{})
		recurrence := expandAzureRmMonitorAutoScaleSettingRecurrence(recurrencesRaw)

		rulesRaw := raw["rule"].([]interface{})
		rules := expandAzureRmMonitorAutoScaleSettingRule(rulesRaw)

		fixedDatesRaw := raw["fixed_date"].([]interface{})
		fixedDate, err := expandAzureRmMonitorAutoScaleSettingFixedDate(fixedDatesRaw)
		if err != nil {
			return nil, fmt.Errorf("expanding `fixed_date`: %+v", err)
		}

		result := autoscalesettings.AutoscaleProfile{
			Name:       name,
			Capacity:   capacity,
			FixedDate:  fixedDate,
			Recurrence: recurrence,
			Rules:      rules,
		}
		results = append(results, result)
	}

	return results, nil
}

func expandAzureRmMonitorAutoScaleSettingPredictive(input []interface{}) *autoscalesettings.PredictiveAutoscalePolicy {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	predictive := autoscalesettings.PredictiveAutoscalePolicy{
		ScaleMode: autoscalesettings.PredictiveAutoscalePolicyScaleMode(raw["scale_mode"].(string)),
	}

	if lookAheadTime := raw["look_ahead_time"].(string); lookAheadTime != "" {
		predictive.ScaleLookAheadTime = pointer.To(lookAheadTime)
	}

	return &predictive
}

func expandAzureRmMonitorAutoScaleSettingRule(input []interface{}) []autoscalesettings.ScaleRule {
	rules := make([]autoscalesettings.ScaleRule, 0)

	for _, v := range input {
		ruleRaw := v.(map[string]interface{})

		triggersRaw := ruleRaw["metric_trigger"].([]interface{})
		triggerRaw := triggersRaw[0].(map[string]interface{})
		metricTrigger := autoscalesettings.MetricTrigger{
			MetricName:        triggerRaw["metric_name"].(string),
			MetricNamespace:   utils.String(triggerRaw["metric_namespace"].(string)),
			MetricResourceUri: triggerRaw["metric_resource_id"].(string),
			TimeGrain:         triggerRaw["time_grain"].(string),
			Statistic:         autoscalesettings.MetricStatisticType(triggerRaw["statistic"].(string)),
			TimeWindow:        triggerRaw["time_window"].(string),
			TimeAggregation:   autoscalesettings.TimeAggregationType(triggerRaw["time_aggregation"].(string)),
			Operator:          autoscalesettings.ComparisonOperationType(triggerRaw["operator"].(string)),
			Threshold:         triggerRaw["threshold"].(float64),
			Dimensions:        expandAzureRmMonitorAutoScaleSettingRuleDimensions(triggerRaw["dimensions"].([]interface{})),
			DividePerInstance: utils.Bool(triggerRaw["divide_by_instance_count"].(bool)),
		}

		actionsRaw := ruleRaw["scale_action"].([]interface{})
		actionRaw := actionsRaw[0].(map[string]interface{})
		scaleAction := autoscalesettings.ScaleAction{
			Direction: autoscalesettings.ScaleDirection(actionRaw["direction"].(string)),
			Type:      autoscalesettings.ScaleType(actionRaw["type"].(string)),
			Value:     utils.String(strconv.Itoa(actionRaw["value"].(int))),
			Cooldown:  actionRaw["cooldown"].(string),
		}

		rule := autoscalesettings.ScaleRule{
			MetricTrigger: metricTrigger,
			ScaleAction:   scaleAction,
		}

		rules = append(rules, rule)
	}

	return rules
}

func expandAzureRmMonitorAutoScaleSettingFixedDate(input []interface{}) (*autoscalesettings.TimeWindow, error) {
	if len(input) == 0 {
		return nil, nil
	}

	raw := input[0].(map[string]interface{})

	startString := raw["start"].(string)
	startTime, err := date.ParseTime(time.RFC3339, startString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse `start` time %q as an RFC3339 date: %+v", startString, err)
	}
	endString := raw["end"].(string)
	endTime, err := date.ParseTime(time.RFC3339, endString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse `end` time %q as an RFC3339 date: %+v", endString, err)
	}

	timeZone := raw["timezone"].(string)
	timeWindow := autoscalesettings.TimeWindow{
		TimeZone: utils.String(timeZone),
	}

	timeWindow.SetStartAsTime(startTime)
	timeWindow.SetEndAsTime(endTime)

	return &timeWindow, nil
}

func expandAzureRmMonitorAutoScaleSettingRecurrence(input []interface{}) *autoscalesettings.Recurrence {
	if len(input) == 0 {
		return nil
	}

	recurrenceRaw := input[0].(map[string]interface{})

	timeZone := recurrenceRaw["timezone"].(string)
	days := make([]string, 0)
	for _, dayItem := range recurrenceRaw["days"].([]interface{}) {
		days = append(days, dayItem.(string))
	}

	hours := make([]int64, 0)
	for _, hourItem := range recurrenceRaw["hours"].([]interface{}) {
		hours = append(hours, int64(hourItem.(int)))
	}

	minutes := make([]int64, 0)
	for _, minuteItem := range recurrenceRaw["minutes"].([]interface{}) {
		minutes = append(minutes, int64(minuteItem.(int)))
	}

	return &autoscalesettings.Recurrence{
		// API docs say this has to be `Week`.
		Frequency: autoscalesettings.RecurrenceFrequencyWeek,
		Schedule: autoscalesettings.RecurrentSchedule{
			TimeZone: timeZone,
			Days:     days,
			Hours:    hours,
			Minutes:  minutes,
		},
	}
}

func expandAzureRmMonitorAutoScaleSettingNotifications(input []interface{}) *[]autoscalesettings.AutoscaleNotification {
	notifications := make([]autoscalesettings.AutoscaleNotification, 0)

	for _, v := range input {
		notificationRaw := v.(map[string]interface{})

		configsRaw := notificationRaw["webhook"].([]interface{})
		webhooks := expandAzureRmMonitorAutoScaleSettingNotificationWebhook(configsRaw)

		notification := autoscalesettings.AutoscaleNotification{
			Operation: "scale",
			WebHooks:  webhooks,
		}

		emailsRaw := notificationRaw["email"].([]interface{})
		if len(emailsRaw) > 0 && emailsRaw[0] != nil {
			notification.Email = expandAzureRmMonitorAutoScaleSettingNotificationEmail(emailsRaw[0].(map[string]interface{}))
		}

		notifications = append(notifications, notification)
	}

	return &notifications
}

func expandAzureRmMonitorAutoScaleSettingNotificationEmail(input map[string]interface{}) *autoscalesettings.EmailNotification {
	customEmails := make([]string, 0)
	if v, ok := input["custom_emails"]; ok {
		for _, item := range v.([]interface{}) {
			customEmails = append(customEmails, item.(string))
		}
	}

	email := autoscalesettings.EmailNotification{
		CustomEmails:                       &customEmails,
		SendToSubscriptionAdministrator:    utils.Bool(input["send_to_subscription_administrator"].(bool)),
		SendToSubscriptionCoAdministrators: utils.Bool(input["send_to_subscription_co_administrator"].(bool)),
	}

	return &email
}

func expandAzureRmMonitorAutoScaleSettingNotificationWebhook(input []interface{}) *[]autoscalesettings.WebhookNotification {
	webhooks := make([]autoscalesettings.WebhookNotification, 0)

	for _, v := range input {
		if v == nil {
			continue
		}
		webhookRaw := v.(map[string]interface{})

		webhook := autoscalesettings.WebhookNotification{
			ServiceUri: utils.String(webhookRaw["service_uri"].(string)),
		}

		if props, ok := webhookRaw["properties"]; ok {
			properties := make(map[string]string)
			for key, value := range props.(map[string]interface{}) {
				properties[key] = value.(string)
			}

			webhook.Properties = &properties
		}

		webhooks = append(webhooks, webhook)
	}

	return &webhooks
}

func expandAzureRmMonitorAutoScaleSettingRuleDimensions(input []interface{}) *[]autoscalesettings.ScaleRuleMetricDimension {
	dimensions := make([]autoscalesettings.ScaleRuleMetricDimension, 0)

	for _, v := range input {
		if v == nil {
			continue
		}
		dimensionRaw := v.(map[string]interface{})

		dimension := autoscalesettings.ScaleRuleMetricDimension{
			DimensionName: dimensionRaw["name"].(string),
			Operator:      autoscalesettings.ScaleRuleMetricDimensionOperationType(dimensionRaw["operator"].(string)),
			Values:        expandStringValues(dimensionRaw["values"].([]interface{})),
		}

		dimensions = append(dimensions, dimension)
	}

	return &dimensions
}

func flattenAzureRmMonitorAutoScaleSettingProfile(profiles []autoscalesettings.AutoscaleProfile) ([]interface{}, error) {
	if profiles == nil {
		return []interface{}{}, nil
	}

	results := make([]interface{}, 0)
	for _, profile := range profiles {
		result := make(map[string]interface{})

		result["name"] = profile.Name

		capacity, err := flattenAzureRmMonitorAutoScaleSettingCapacity(profile.Capacity)
		if err != nil {
			return nil, fmt.Errorf("flattening `capacity`: %+v", err)
		}
		result["capacity"] = capacity

		result["fixed_date"] = flattenAzureRmMonitorAutoScaleSettingFixedDate(profile.FixedDate)
		result["recurrence"] = flattenAzureRmMonitorAutoScaleSettingRecurrence(profile.Recurrence)

		rule, err := flattenAzureRmMonitorAutoScaleSettingRules(profile.Rules)
		if err != nil {
			return nil, fmt.Errorf("flattening Rule: %s", err)
		}
		result["rule"] = rule

		results = append(results, result)
	}
	return results, nil
}

func flattenAzureRmMonitorAutoScaleSettingPredictive(input *autoscalesettings.PredictiveAutoscalePolicy) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	// omit the block if disabled
	if input.ScaleMode == autoscalesettings.PredictiveAutoscalePolicyScaleModeDisabled {
		return []interface{}{}
	}

	result := map[string]interface{}{
		"look_ahead_time": pointer.From(input.ScaleLookAheadTime),
		"scale_mode":      string(input.ScaleMode),
	}

	return []interface{}{result}
}

func flattenAzureRmMonitorAutoScaleSettingCapacity(input autoscalesettings.ScaleCapacity) ([]interface{}, error) {

	result := make(map[string]interface{})

	min, err := strconv.Atoi(input.Minimum)
	if err != nil {
		return nil, fmt.Errorf("converting Minimum Scale Capacity %q to an int: %+v", input.Minimum, err)
	}
	result["minimum"] = min

	max, err := strconv.Atoi(input.Maximum)
	if err != nil {
		return nil, fmt.Errorf("converting Maximum Scale Capacity %q to an int: %+v", input.Maximum, err)
	}
	result["maximum"] = max

	defaultCapacity, err := strconv.Atoi(input.Default)
	if err != nil {
		return nil, fmt.Errorf("converting Default Scale Capacity %q to an int: %+v", input.Default, err)
	}
	result["default"] = defaultCapacity

	return []interface{}{result}, nil
}

func flattenAzureRmMonitorAutoScaleSettingRules(input []autoscalesettings.ScaleRule) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	results := make([]interface{}, 0)
	for _, rule := range input {
		result := make(map[string]interface{})

		metricTriggers := make([]interface{}, 0)
		var metricNamespace string
		var dividePerInstance bool

		if v := rule.MetricTrigger.MetricNamespace; v != nil {
			metricNamespace = *v
		}

		if rule.MetricTrigger.DividePerInstance != nil {
			dividePerInstance = *rule.MetricTrigger.DividePerInstance
		}

		metricTriggers = append(metricTriggers, map[string]interface{}{
			"metric_name":              rule.MetricTrigger.MetricName,
			"metric_namespace":         metricNamespace,
			"metric_resource_id":       rule.MetricTrigger.MetricResourceUri,
			"time_grain":               rule.MetricTrigger.TimeGrain,
			"statistic":                string(rule.MetricTrigger.Statistic),
			"time_window":              rule.MetricTrigger.TimeWindow,
			"time_aggregation":         string(rule.MetricTrigger.TimeAggregation),
			"operator":                 string(rule.MetricTrigger.Operator),
			"threshold":                rule.MetricTrigger.Threshold,
			"dimensions":               flattenAzureRmMonitorAutoScaleSettingRulesDimensions(rule.MetricTrigger.Dimensions),
			"divide_by_instance_count": dividePerInstance,
		})

		result["metric_trigger"] = metricTriggers

		scaleActions := make([]interface{}, 0)
		v := rule.ScaleAction
		action := make(map[string]interface{})

		action["direction"] = string(v.Direction)
		action["type"] = string(v.Type)
		action["cooldown"] = v.Cooldown

		if val := v.Value; val != nil && *val != "" {
			i, err := strconv.Atoi(*val)
			if err != nil {
				return nil, fmt.Errorf("`value` %q was not convertable to an int: %s", *val, err)
			}
			action["value"] = i
		}

		scaleActions = append(scaleActions, action)

		result["scale_action"] = scaleActions

		results = append(results, result)
	}

	return results, nil
}

func flattenAzureRmMonitorAutoScaleSettingFixedDate(input *autoscalesettings.TimeWindow) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	if timezone := input.TimeZone; timezone != nil {
		result["timezone"] = *timezone
	}
	result["start"] = input.Start
	result["end"] = input.End

	return []interface{}{result}
}

func flattenAzureRmMonitorAutoScaleSettingRecurrence(input *autoscalesettings.Recurrence) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	schedule := input.Schedule
	result["timezone"] = schedule.TimeZone

	days := make([]string, 0)
	if s := schedule.Days; s != nil {
		days = s
	}
	result["days"] = days

	hours := make([]int, 0)
	if schedule.Hours != nil {
		for _, v := range schedule.Hours {
			hours = append(hours, int(v))
		}
	}
	result["hours"] = hours

	minutes := make([]int, 0)
	if schedule.Minutes != nil {
		for _, v := range schedule.Minutes {
			minutes = append(minutes, int(v))
		}
	}
	result["minutes"] = minutes

	return []interface{}{result}
}

func flattenAzureRmMonitorAutoScaleSettingNotification(notifications *[]autoscalesettings.AutoscaleNotification) []interface{} {
	results := make([]interface{}, 0)

	if notifications == nil {
		return results
	}

	for _, notification := range *notifications {
		result := make(map[string]interface{})

		emails := make([]interface{}, 0)
		if email := notification.Email; email != nil {
			block := make(map[string]interface{})

			if send := email.SendToSubscriptionAdministrator; send != nil {
				block["send_to_subscription_administrator"] = *send
			}

			if send := email.SendToSubscriptionCoAdministrators; send != nil {
				block["send_to_subscription_co_administrator"] = *send
			}

			customEmails := make([]interface{}, 0)
			if custom := email.CustomEmails; custom != nil {
				for _, v := range *custom {
					customEmails = append(customEmails, v)
				}
			}
			block["custom_emails"] = customEmails

			emails = append(emails, block)
		}
		result["email"] = emails

		webhooks := make([]interface{}, 0)
		if hooks := notification.WebHooks; hooks != nil {
			for _, v := range *hooks {
				hook := make(map[string]interface{})

				if v.ServiceUri != nil {
					hook["service_uri"] = *v.ServiceUri
				}

				props := make(map[string]string)
				if webHookProps := v.Properties; webHookProps != nil {
					for key, value := range *v.Properties {
						props[key] = value

					}
					hook["properties"] = props
					webhooks = append(webhooks, hook)
				}
			}
		}

		result["webhook"] = webhooks

		results = append(results, result)
	}
	return results
}

func flattenAzureRmMonitorAutoScaleSettingRulesDimensions(dimensions *[]autoscalesettings.ScaleRuleMetricDimension) []interface{} {
	results := make([]interface{}, 0)

	if dimensions == nil {
		return results
	}

	for _, dimension := range *dimensions {

		results = append(results, map[string]interface{}{
			"name":     dimension.DimensionName,
			"operator": string(dimension.Operator),
			"values":   dimension.Values,
		})
	}
	return results
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
