package monitor

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/monitor/mgmt/2019-06-01/insights"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceMonitorAutoScaleSetting() *schema.Resource {
	return &schema.Resource{
		Create: resourceMonitorAutoScaleSettingCreateUpdate,
		Read:   resourceMonitorAutoScaleSettingRead,
		Update: resourceMonitorAutoScaleSettingCreateUpdate,
		Delete: resourceMonitorAutoScaleSettingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"target_resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"profile": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 20,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"capacity": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"minimum": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(0, 1000),
									},
									"maximum": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(0, 1000),
									},
									"default": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(0, 1000),
									},
								},
							},
						},
						"rule": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 10,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_trigger": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"metric_name": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},
												"metric_resource_id": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: azure.ValidateResourceID,
												},
												"time_grain": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validate.ISO8601Duration,
												},
												"statistic": {
													Type:     schema.TypeString,
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
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validate.ISO8601Duration,
												},
												"time_aggregation": {
													Type:     schema.TypeString,
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
													Type:     schema.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(insights.Equals),
														string(insights.GreaterThan),
														string(insights.GreaterThanOrEqual),
														string(insights.LessThan),
														string(insights.LessThanOrEqual),
														string(insights.NotEquals),
													}, true),
													DiffSuppressFunc: suppress.CaseDifference,
												},
												"threshold": {
													Type:     schema.TypeFloat,
													Required: true,
												},

												"metric_namespace": {
													Type:         schema.TypeString,
													Optional:     true,
													ValidateFunc: validation.StringIsNotEmpty,
												},

												"dimensions": {
													Type:     schema.TypeList,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:         schema.TypeString,
																Required:     true,
																ValidateFunc: validation.StringIsNotEmpty,
															},

															"operator": {
																Type:     schema.TypeString,
																Required: true,
																ValidateFunc: validation.StringInSlice([]string{
																	string(insights.ScaleRuleMetricDimensionOperationTypeEquals),
																	string(insights.ScaleRuleMetricDimensionOperationTypeNotEquals),
																}, false),
															},

															"values": {
																Type:     schema.TypeList,
																Required: true,
																Elem: &schema.Schema{
																	Type:         schema.TypeString,
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
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"direction": {
													Type:     schema.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(insights.ScaleDirectionDecrease),
														string(insights.ScaleDirectionIncrease),
													}, true),
													DiffSuppressFunc: suppress.CaseDifference,
												},
												"type": {
													Type:     schema.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(insights.ChangeCount),
														string(insights.ExactCount),
														string(insights.PercentChangeCount),
													}, true),
													DiffSuppressFunc: suppress.CaseDifference,
												},
												"value": {
													Type:         schema.TypeInt,
													Required:     true,
													ValidateFunc: validation.IntAtLeast(0),
												},
												"cooldown": {
													Type:         schema.TypeString,
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
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"timezone": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "UTC",
										ValidateFunc: validateMonitorAutoScaleSettingsTimeZone(),
									},
									"start": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.IsRFC3339Time,
									},
									"end": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.IsRFC3339Time,
									},
								},
							},
						},
						"recurrence": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"timezone": {
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "UTC",
										ValidateFunc: validateMonitorAutoScaleSettingsTimeZone(),
									},
									"days": {
										Type:     schema.TypeList,
										Required: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
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
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Schema{
											Type:         schema.TypeInt,
											ValidateFunc: validation.IntBetween(0, 23),
										},
									},
									"minutes": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Schema{
											Type:         schema.TypeInt,
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
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"email": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"send_to_subscription_administrator": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"send_to_subscription_co_administrator": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"custom_emails": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"webhook": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_uri": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"properties": {
										Type:     schema.TypeMap,
										Optional: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceMonitorAutoScaleSettingCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.AutoscaleSettingsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Monitor AutoScale Setting %q (Resource Group %q): %s", name, resourceGroup, err)
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
		return fmt.Errorf("Error expanding `profile`: %+v", err)
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

	if _, err = client.CreateOrUpdate(ctx, resourceGroup, name, parameters); err != nil {
		return fmt.Errorf("Error creating AutoScale Setting %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving AutoScale Setting %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("AutoScale Setting %q (Resource Group %q) has no ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceMonitorAutoScaleSettingRead(d, meta)
}

func resourceMonitorAutoScaleSettingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.AutoscaleSettingsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["autoscalesettings"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] AutoScale Setting %q (Resource Group %q) was not found - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading AutoScale Setting %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("enabled", resp.Enabled)
	d.Set("target_resource_id", resp.TargetResourceURI)

	profile, err := flattenAzureRmMonitorAutoScaleSettingProfile(resp.Profiles)
	if err != nil {
		return fmt.Errorf("Error flattening `profile` of Autoscale Setting %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = d.Set("profile", profile); err != nil {
		return fmt.Errorf("Error setting `profile` of Autoscale Setting %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	notifications := flattenAzureRmMonitorAutoScaleSettingNotification(resp.Notifications)
	if err = d.Set("notification", notifications); err != nil {
		return fmt.Errorf("Error setting `notification` of Autoscale Setting %q (resource group %q): %+v", name, resourceGroup, err)
	}

	// Return a new tag map filtered by the specified tag names.
	tagMap := tags.Filter(resp.Tags, "$type")
	return tags.FlattenAndSet(d, tagMap)
}

func resourceMonitorAutoScaleSettingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.AutoscaleSettingsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["autoscalesettings"]

	resp, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if !response.WasNotFound(resp.Response) {
			return fmt.Errorf("Error deleting AutoScale Setting %q (Resource Group %q): %+v", name, resourceGroup, err)
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
			return nil, fmt.Errorf("Error expanding `fixed_date`: %+v", err)
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
			return nil, fmt.Errorf("Error flattening `capacity`: %+v", err)
		}
		result["capacity"] = capacity

		result["fixed_date"] = flattenAzureRmMonitorAutoScaleSettingFixedDate(profile.FixedDate)
		result["recurrence"] = flattenAzureRmMonitorAutoScaleSettingRecurrence(profile.Recurrence)

		rule, err := flattenAzureRmMonitorAutoScaleSettingRules(profile.Rules)
		if err != nil {
			return nil, fmt.Errorf("Error flattening Rule: %s", err)
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
			return nil, fmt.Errorf("Error converting Minimum Scale Capacity %q to an int: %+v", *minStr, err)
		}
		result["minimum"] = min
	}

	if maxStr := input.Maximum; maxStr != nil {
		max, err := strconv.Atoi(*maxStr)
		if err != nil {
			return nil, fmt.Errorf("Error converting Maximum Scale Capacity %q to an int: %+v", *maxStr, err)
		}
		result["maximum"] = max
	}

	if defaultCapacityStr := input.Default; defaultCapacityStr != nil {
		defaultCapacity, err := strconv.Atoi(*defaultCapacityStr)
		if err != nil {
			return nil, fmt.Errorf("Error converting Default Scale Capacity %q to an int: %+v", *defaultCapacityStr, err)
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

			metricTriggers = append(metricTriggers, map[string]interface{}{
				"metric_name":        metricName,
				"metric_namespace":   metricNamespace,
				"metric_resource_id": metricId,
				"time_grain":         timeGrain,
				"statistic":          string(trigger.Statistic),
				"time_window":        timeWindow,
				"time_aggregation":   string(trigger.TimeAggregation),
				"operator":           string(trigger.Operator),
				"threshold":          threshold,
				"dimensions":         flattenAzureRmMonitorAutoScaleSettingRulesDimensions(trigger.Dimensions),
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

func validateMonitorAutoScaleSettingsTimeZone() schema.SchemaValidateFunc {
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
