package monitor

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2021-07-01-preview/insights"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
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
			_, err := parse.AutoscaleSettingID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AutoscaleSettingUpgradeV0ToV1{},
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
														string(insights.ScaleTypeServiceAllowedNextValue),
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
		},
	}
}

func resourceMonitorAutoScaleSettingCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.AutoscaleSettingsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewAutoscaleSettingID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Monitor %s: %+v", id, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_monitor_autoscale_setting", *existing.ID)
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
	expandedTags := tags.Expand(t)

	parameters := insights.AutoscaleSettingResource{
		Location: utils.String(location),
		AutoscaleSetting: &insights.AutoscaleSetting{
			Enabled:           &enabled,
			Profiles:          profiles,
			Notifications:     notifications,
			TargetResourceURI: &targetResourceId,
		},
		Tags: expandedTags,
	}

	if _, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, parameters); err != nil {
		return fmt.Errorf("creating Monitor %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMonitorAutoScaleSettingRead(d, meta)
}

func resourceMonitorAutoScaleSettingRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.AutoscaleSettingsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutoscaleSettingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] AutoScale Setting %q (Resource Group %q) was not found - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Monitor %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("enabled", resp.Enabled)
	d.Set("target_resource_id", resp.TargetResourceURI)

	profile, err := flattenAzureRmMonitorAutoScaleSettingProfile(resp.Profiles)
	if err != nil {
		return fmt.Errorf("flattening `profile` of %s: %+v", *id, err)
	}
	if err = d.Set("profile", profile); err != nil {
		return fmt.Errorf("setting `profile` of %s: %+v", *id, err)
	}

	notifications := flattenAzureRmMonitorAutoScaleSettingNotification(resp.Notifications)
	if err = d.Set("notification", notifications); err != nil {
		return fmt.Errorf("setting `notification` of %s: %+v", *id, err)
	}

	// Return a new tag map filtered by the specified tag names.
	tagMap := tags.Filter(resp.Tags, "$type")
	return tags.FlattenAndSet(d, tagMap)
}

func resourceMonitorAutoScaleSettingDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.AutoscaleSettingsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutoscaleSettingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("deleting Monitor %s: %+v", *id, err)
		}
	}

	return nil
}

func expandAzureRmMonitorAutoScaleSettingProfile(input []interface{}) (*[]insights.AutoscaleProfile, error) {
	results := make([]insights.AutoscaleProfile, 0)

	for _, v := range input {
		raw := v.(map[string]interface{})

		name := raw["name"].(string)

		// this is Required, so we don't need to check for optionals here
		capacitiesRaw := raw["capacity"].([]interface{})
		capacityRaw := capacitiesRaw[0].(map[string]interface{})
		capacity := insights.ScaleCapacity{
			Minimum: utils.String(strconv.Itoa(capacityRaw["minimum"].(int))),
			Maximum: utils.String(strconv.Itoa(capacityRaw["maximum"].(int))),
			Default: utils.String(strconv.Itoa(capacityRaw["default"].(int))),
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

		result := insights.AutoscaleProfile{
			Name:       utils.String(name),
			Capacity:   &capacity,
			FixedDate:  fixedDate,
			Recurrence: recurrence,
			Rules:      rules,
		}
		results = append(results, result)
	}

	return &results, nil
}

func expandAzureRmMonitorAutoScaleSettingRule(input []interface{}) *[]insights.ScaleRule {
	rules := make([]insights.ScaleRule, 0)

	for _, v := range input {
		ruleRaw := v.(map[string]interface{})

		triggersRaw := ruleRaw["metric_trigger"].([]interface{})
		triggerRaw := triggersRaw[0].(map[string]interface{})
		metricTrigger := insights.MetricTrigger{
			MetricName:        utils.String(triggerRaw["metric_name"].(string)),
			MetricNamespace:   utils.String(triggerRaw["metric_namespace"].(string)),
			MetricResourceURI: utils.String(triggerRaw["metric_resource_id"].(string)),
			TimeGrain:         utils.String(triggerRaw["time_grain"].(string)),
			Statistic:         insights.MetricStatisticType(triggerRaw["statistic"].(string)),
			TimeWindow:        utils.String(triggerRaw["time_window"].(string)),
			TimeAggregation:   insights.TimeAggregationType(triggerRaw["time_aggregation"].(string)),
			Operator:          insights.ComparisonOperationType(triggerRaw["operator"].(string)),
			Threshold:         utils.Float(triggerRaw["threshold"].(float64)),
			Dimensions:        expandAzureRmMonitorAutoScaleSettingRuleDimensions(triggerRaw["dimensions"].([]interface{})),
			DividePerInstance: utils.Bool(triggerRaw["divide_by_instance_count"].(bool)),
		}

		actionsRaw := ruleRaw["scale_action"].([]interface{})
		actionRaw := actionsRaw[0].(map[string]interface{})
		scaleAction := insights.ScaleAction{
			Direction: insights.ScaleDirection(actionRaw["direction"].(string)),
			Type:      insights.ScaleType(actionRaw["type"].(string)),
			Value:     utils.String(strconv.Itoa(actionRaw["value"].(int))),
			Cooldown:  utils.String(actionRaw["cooldown"].(string)),
		}

		rule := insights.ScaleRule{
			MetricTrigger: &metricTrigger,
			ScaleAction:   &scaleAction,
		}

		rules = append(rules, rule)
	}

	return &rules
}

func expandAzureRmMonitorAutoScaleSettingFixedDate(input []interface{}) (*insights.TimeWindow, error) {
	if len(input) == 0 {
		return nil, nil
	}

	raw := input[0].(map[string]interface{})

	startString := raw["start"].(string)
	startTime, err := date.ParseTime(time.RFC3339, startString)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse `start` time %q as an RFC3339 date: %+v", startString, err)
	}
	endString := raw["end"].(string)
	endTime, err := date.ParseTime(time.RFC3339, endString)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse `end` time %q as an RFC3339 date: %+v", endString, err)
	}

	timeZone := raw["timezone"].(string)
	timeWindow := insights.TimeWindow{
		TimeZone: utils.String(timeZone),
		Start: &date.Time{
			Time: startTime,
		},
		End: &date.Time{
			Time: endTime,
		},
	}
	return &timeWindow, nil
}

func expandAzureRmMonitorAutoScaleSettingRecurrence(input []interface{}) *insights.Recurrence {
	if len(input) == 0 {
		return nil
	}

	recurrenceRaw := input[0].(map[string]interface{})

	timeZone := recurrenceRaw["timezone"].(string)
	days := make([]string, 0)
	for _, dayItem := range recurrenceRaw["days"].([]interface{}) {
		days = append(days, dayItem.(string))
	}

	hours := make([]int32, 0)
	for _, hourItem := range recurrenceRaw["hours"].([]interface{}) {
		hours = append(hours, int32(hourItem.(int)))
	}

	minutes := make([]int32, 0)
	for _, minuteItem := range recurrenceRaw["minutes"].([]interface{}) {
		minutes = append(minutes, int32(minuteItem.(int)))
	}

	return &insights.Recurrence{
		// API docs say this has to be `Week`.
		Frequency: insights.RecurrenceFrequencyWeek,
		Schedule: &insights.RecurrentSchedule{
			TimeZone: utils.String(timeZone),
			Days:     &days,
			Hours:    &hours,
			Minutes:  &minutes,
		},
	}
}

func expandAzureRmMonitorAutoScaleSettingNotifications(input []interface{}) *[]insights.AutoscaleNotification {
	notifications := make([]insights.AutoscaleNotification, 0)

	for _, v := range input {
		notificationRaw := v.(map[string]interface{})

		configsRaw := notificationRaw["webhook"].([]interface{})
		webhooks := expandAzureRmMonitorAutoScaleSettingNotificationWebhook(configsRaw)

		notification := insights.AutoscaleNotification{
			Operation: utils.String("scale"),
			Webhooks:  webhooks,
		}

		emailsRaw := notificationRaw["email"].([]interface{})
		if len(emailsRaw) > 0 && emailsRaw[0] != nil {
			notification.Email = expandAzureRmMonitorAutoScaleSettingNotificationEmail(emailsRaw[0].(map[string]interface{}))
		}

		notifications = append(notifications, notification)
	}

	return &notifications
}

func expandAzureRmMonitorAutoScaleSettingNotificationEmail(input map[string]interface{}) *insights.EmailNotification {
	customEmails := make([]string, 0)
	if v, ok := input["custom_emails"]; ok {
		for _, item := range v.([]interface{}) {
			customEmails = append(customEmails, item.(string))
		}
	}

	email := insights.EmailNotification{
		CustomEmails:                       &customEmails,
		SendToSubscriptionAdministrator:    utils.Bool(input["send_to_subscription_administrator"].(bool)),
		SendToSubscriptionCoAdministrators: utils.Bool(input["send_to_subscription_co_administrator"].(bool)),
	}

	return &email
}

func expandAzureRmMonitorAutoScaleSettingNotificationWebhook(input []interface{}) *[]insights.WebhookNotification {
	webhooks := make([]insights.WebhookNotification, 0)

	for _, v := range input {
		if v == nil {
			continue
		}
		webhookRaw := v.(map[string]interface{})

		webhook := insights.WebhookNotification{
			ServiceURI: utils.String(webhookRaw["service_uri"].(string)),
		}

		if props, ok := webhookRaw["properties"]; ok {
			properties := make(map[string]*string)
			for key, value := range props.(map[string]interface{}) {
				properties[key] = utils.String(value.(string))
			}

			webhook.Properties = properties
		}

		webhooks = append(webhooks, webhook)
	}

	return &webhooks
}

func expandAzureRmMonitorAutoScaleSettingRuleDimensions(input []interface{}) *[]insights.ScaleRuleMetricDimension {
	dimensions := make([]insights.ScaleRuleMetricDimension, 0)

	for _, v := range input {
		if v == nil {
			continue
		}
		dimensionRaw := v.(map[string]interface{})

		dimension := insights.ScaleRuleMetricDimension{
			DimensionName: utils.String(dimensionRaw["name"].(string)),
			Operator:      insights.ScaleRuleMetricDimensionOperationType(dimensionRaw["operator"].(string)),
			Values:        utils.ExpandStringSlice(dimensionRaw["values"].([]interface{})),
		}

		dimensions = append(dimensions, dimension)
	}

	return &dimensions
}

func flattenAzureRmMonitorAutoScaleSettingProfile(profiles *[]insights.AutoscaleProfile) ([]interface{}, error) {
	if profiles == nil {
		return []interface{}{}, nil
	}

	results := make([]interface{}, 0)
	for _, profile := range *profiles {
		result := make(map[string]interface{})

		if name := profile.Name; name != nil {
			result["name"] = *name
		}

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

func flattenAzureRmMonitorAutoScaleSettingCapacity(input *insights.ScaleCapacity) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	result := make(map[string]interface{})

	if minStr := input.Minimum; minStr != nil {
		min, err := strconv.Atoi(*minStr)
		if err != nil {
			return nil, fmt.Errorf("converting Minimum Scale Capacity %q to an int: %+v", *minStr, err)
		}
		result["minimum"] = min
	}

	if maxStr := input.Maximum; maxStr != nil {
		max, err := strconv.Atoi(*maxStr)
		if err != nil {
			return nil, fmt.Errorf("converting Maximum Scale Capacity %q to an int: %+v", *maxStr, err)
		}
		result["maximum"] = max
	}

	if defaultCapacityStr := input.Default; defaultCapacityStr != nil {
		defaultCapacity, err := strconv.Atoi(*defaultCapacityStr)
		if err != nil {
			return nil, fmt.Errorf("converting Default Scale Capacity %q to an int: %+v", *defaultCapacityStr, err)
		}
		result["default"] = defaultCapacity
	}

	return []interface{}{result}, nil
}

func flattenAzureRmMonitorAutoScaleSettingRules(input *[]insights.ScaleRule) ([]interface{}, error) {
	if input == nil {
		return []interface{}{}, nil
	}

	results := make([]interface{}, 0)
	for _, rule := range *input {
		result := make(map[string]interface{})

		metricTriggers := make([]interface{}, 0)
		if trigger := rule.MetricTrigger; trigger != nil {
			var metricName, metricNamespace, metricId, timeGrain, timeWindow string
			var dividePerInstance bool
			var threshold float64
			if trigger.MetricName != nil {
				metricName = *trigger.MetricName
			}

			if v := trigger.MetricNamespace; v != nil {
				metricNamespace = *v
			}

			if trigger.MetricResourceURI != nil {
				metricId = *trigger.MetricResourceURI
			}

			if trigger.TimeGrain != nil {
				timeGrain = *trigger.TimeGrain
			}

			if trigger.TimeWindow != nil {
				timeWindow = *trigger.TimeWindow
			}

			if trigger.Threshold != nil {
				threshold = *trigger.Threshold
			}

			if trigger.DividePerInstance != nil {
				dividePerInstance = *trigger.DividePerInstance
			}

			metricTriggers = append(metricTriggers, map[string]interface{}{
				"metric_name":              metricName,
				"metric_namespace":         metricNamespace,
				"metric_resource_id":       metricId,
				"time_grain":               timeGrain,
				"statistic":                string(trigger.Statistic),
				"time_window":              timeWindow,
				"time_aggregation":         string(trigger.TimeAggregation),
				"operator":                 string(trigger.Operator),
				"threshold":                threshold,
				"dimensions":               flattenAzureRmMonitorAutoScaleSettingRulesDimensions(trigger.Dimensions),
				"divide_by_instance_count": dividePerInstance,
			})
		}

		result["metric_trigger"] = metricTriggers

		scaleActions := make([]interface{}, 0)
		if v := rule.ScaleAction; v != nil {
			action := make(map[string]interface{})

			action["direction"] = string(v.Direction)
			action["type"] = string(v.Type)

			if v.Cooldown != nil {
				action["cooldown"] = *v.Cooldown
			}

			if val := v.Value; val != nil && *val != "" {
				i, err := strconv.Atoi(*val)
				if err != nil {
					return nil, fmt.Errorf("`value` %q was not convertable to an int: %s", *val, err)
				}
				action["value"] = i
			}

			scaleActions = append(scaleActions, action)
		}

		result["scale_action"] = scaleActions

		results = append(results, result)
	}

	return results, nil
}

func flattenAzureRmMonitorAutoScaleSettingFixedDate(input *insights.TimeWindow) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	if timezone := input.TimeZone; timezone != nil {
		result["timezone"] = *timezone
	}

	if start := input.Start; start != nil {
		result["start"] = start.String()
	}

	if end := input.End; end != nil {
		result["end"] = end.String()
	}

	return []interface{}{result}
}

func flattenAzureRmMonitorAutoScaleSettingRecurrence(input *insights.Recurrence) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	result := make(map[string]interface{})

	if schedule := input.Schedule; schedule != nil {
		if timezone := schedule.TimeZone; timezone != nil {
			result["timezone"] = *timezone
		}

		days := make([]string, 0)
		if s := schedule.Days; s != nil {
			days = *s
		}
		result["days"] = days

		hours := make([]int, 0)
		if schedule.Hours != nil {
			for _, v := range *schedule.Hours {
				hours = append(hours, int(v))
			}
		}
		result["hours"] = hours

		minutes := make([]int, 0)
		if schedule.Minutes != nil {
			for _, v := range *schedule.Minutes {
				minutes = append(minutes, int(v))
			}
		}
		result["minutes"] = minutes
	}

	return []interface{}{result}
}

func flattenAzureRmMonitorAutoScaleSettingNotification(notifications *[]insights.AutoscaleNotification) []interface{} {
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
		if hooks := notification.Webhooks; hooks != nil {
			for _, v := range *hooks {
				hook := make(map[string]interface{})

				if v.ServiceURI != nil {
					hook["service_uri"] = *v.ServiceURI
				}

				props := make(map[string]string)
				for key, value := range v.Properties {
					if value != nil {
						props[key] = *value
					}
				}
				hook["properties"] = props
				webhooks = append(webhooks, hook)
			}
		}

		result["webhook"] = webhooks

		results = append(results, result)
	}
	return results
}

func flattenAzureRmMonitorAutoScaleSettingRulesDimensions(dimensions *[]insights.ScaleRuleMetricDimension) []interface{} {
	results := make([]interface{}, 0)

	if dimensions == nil {
		return results
	}

	for _, dimension := range *dimensions {
		var name string

		if v := dimension.DimensionName; v != nil {
			name = *v
		}

		results = append(results, map[string]interface{}{
			"name":     name,
			"operator": string(dimension.Operator),
			"values":   utils.FlattenStringSlice(dimension.Values),
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
