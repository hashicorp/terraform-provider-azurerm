package azurerm

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/monitor/mgmt/2018-03-01/insights"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmAutoScaleSetting() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutoscaleSettingCreateOrUpdate,
		Read:   resourceArmAutoscaleSettingRead,
		Update: resourceArmAutoscaleSettingCreateOrUpdate,
		Delete: resourceArmAutoscaleSettingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"location": locationSchema(),

			"target_resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"profile": {
				Type:     schema.TypeList,
				Required: true,
				// https://docs.microsoft.com/en-us/rest/api/monitor/autoscalesettings/createorupdate#default
				MaxItems: 20,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
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
										ValidateFunc: validation.IntBetween(1, 40),
									},
									"maximum": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 40),
									},
									"default": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 40),
									},
								},
							},
						},
						"rule": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 10, // https://docs.microsoft.com/en-us/rest/api/monitor/autoscalesettings/createorupdate#autoscaleprofile
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
													ValidateFunc: validation.NoZeroValues,
												},
												"metric_resource_id": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: azure.ValidateResourceID,
												},
												"time_grain": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validateIso8601Duration(),
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
													DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
												},
												"time_window": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validateIso8601Duration(),
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
													}, true),
													DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
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
													DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
												},
												"threshold": {
													Type:     schema.TypeFloat,
													Required: true,
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
													DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
												},
												"type": {
													Type:     schema.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(insights.ChangeCount),
														string(insights.ExactCount),
														string(insights.PercentChangeCount),
													}, true),
													DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
												},
												"value": {
													Type:         schema.TypeInt,
													Required:     true,
													ValidateFunc: validation.IntAtLeast(0),
												},
												"cooldown": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validateIso8601Duration(),
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
									"time_zone": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateAutoScaleSettingsTimeZone(),
									},
									"start": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateRFC3339Date,
									},
									"end": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateRFC3339Date,
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
									"time_zone": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateAutoScaleSettingsTimeZone(),
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
											DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
										},
									},
									"hours": {
										Type:     schema.TypeList,
										Required: true,
										MinItems: 1,
										Elem: &schema.Schema{
											Type:         schema.TypeInt,
											ValidateFunc: validation.IntBetween(0, 23),
										},
									},
									"minutes": {
										Type:     schema.TypeList,
										Required: true,
										MinItems: 1,
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

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
										// TODO: does this want to be a Set?
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
										ValidateFunc: validation.NoZeroValues,
									},
									"properties": {
										Type:     schema.TypeMap,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmAutoscaleSettingCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).autoscaleSettingsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))

	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	enabled := d.Get("enabled").(bool)
	targetResourceId := d.Get("target_resource_id").(string)

	parameters := insights.AutoscaleSettingResource{
		Location: utils.String(location),
		AutoscaleSetting: &insights.AutoscaleSetting{
			Enabled:           &enabled,
			TargetResourceURI: &targetResourceId,
		},
		Tags: expandedTags,
	}

	if v, ok := d.GetOk("profile"); ok {
		profiles, err := expandAzureRmAutoscaleSettingProfile(v.([]interface{}))
		if err != nil {
			return fmt.Errorf("Error expanding `profile`: %+v", err)
		}
		parameters.AutoscaleSetting.Profiles = profiles
	}
	if v, ok := d.GetOk("notification"); ok {
		parameters.AutoscaleSetting.Notifications = expandAzureRmAutoscaleSettingNotifications(v.([]interface{}))
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroupName, name, parameters); err != nil {
		return fmt.Errorf("Error creating Autoscale Setting %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}

	read, err := client.Get(ctx, resourceGroupName, name)
	if err != nil {
		return fmt.Errorf("Error reading Autoscale Setting %q (Resource Group %q): %+v", name, resourceGroupName, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Autoscale Setting %q (Resource Group %q) ID is empty", name, resourceGroupName)
	}

	d.SetId(*read.ID)

	return resourceArmAutoscaleSettingRead(d, meta)
}

func resourceArmAutoscaleSettingRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).autoscaleSettingsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing Autoscale Setting resource ID: %+v", err)
	}
	resGroup := id.ResourceGroup
	name := id.Path["autoscalesettings"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Autoscale Setting %q (resource group %q): %+v", name, resGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	d.Set("enabled", resp.Enabled)
	d.Set("target_resource_id", resp.TargetResourceURI)

	profile, err := flattenAzureRmAutoscaleSettingProfile(resp.Profiles)
	if err != nil {
		return fmt.Errorf("Error flattening `profile` of Autoscale Setting %q (Resource Group %q): %+v", name, resGroup, err)
	}
	if err = d.Set("profile", profile); err != nil {
		return fmt.Errorf("Error setting `profile` of Autoscale Setting %q (Resource Group %q): %+v", name, resGroup, err)
	}

	notifications := flattenAzureRmAutoscaleSettingNotification(resp.Notifications)
	if err = d.Set("notification", notifications); err != nil {
		return fmt.Errorf("Error setting `notification` of Autoscale Setting %q (resource group %q): %+v", name, resGroup, err)
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func resourceArmAutoscaleSettingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).autoscaleSettingsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing Autoscale Setting resource ID: %+v", err)
	}
	resGroup := id.ResourceGroup
	name := id.Path["autoscalesettings"]

	_, err = client.Delete(ctx, resGroup, name)

	// TODO: handle the error here

	return err
}

func expandAzureRmAutoscaleSettingProfile(v []interface{}) (*[]insights.AutoscaleProfile, error) {
	profiles := make([]insights.AutoscaleProfile, 0)
	for _, profileItem := range v {
		val := profileItem.(map[string]interface{})

		capacity := val["capacity"].([]interface{})[0].(map[string]interface{})
		profile := insights.AutoscaleProfile{
			Name: utils.String(val["name"].(string)),
			Capacity: &insights.ScaleCapacity{
				Minimum: utils.String(strconv.Itoa(capacity["minimum"].(int))),
				Maximum: utils.String(strconv.Itoa(capacity["maximum"].(int))),
				Default: utils.String(strconv.Itoa(capacity["default"].(int))),
			},
		}
		if v, ok := val["rule"]; ok {
			profile.Rules = expandAzureRmAutoscaleSettingRule(v.([]interface{}))
		}
		if v, ok := val["fixed_date"]; ok {
			vObj := v.([]interface{})[0]
			fixedDate, err := expandAzureRmAutoscaleSettingFixedDate(vObj.(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			profile.FixedDate = fixedDate
		}
		if v, ok := val["recurrence"]; ok {
			vObj := v.([]interface{})[0]
			profile.Recurrence = expandAzureRmAutoscaleSettingRecurrence(vObj.(map[string]interface{}))
		}

		profiles = append(profiles, profile)
	}
	return &profiles, nil
}

func expandAzureRmAutoscaleSettingRule(v []interface{}) *[]insights.ScaleRule {
	rules := make([]insights.ScaleRule, 0)
	for _, ruleItem := range v {
		ruleValue := ruleItem.(map[string]interface{})

		metric := ruleValue["metric_trigger"].([]interface{})[0].(map[string]interface{})
		scale := ruleValue["scale_action"].([]interface{})[0].(map[string]interface{})
		rule := insights.ScaleRule{
			MetricTrigger: &insights.MetricTrigger{
				MetricName:        utils.String(metric["metric_name"].(string)),
				MetricResourceURI: utils.String(metric["metric_resource_id"].(string)),
				TimeGrain:         utils.String(metric["time_grain"].(string)),
				Statistic:         insights.MetricStatisticType(metric["statistic"].(string)),
				TimeWindow:        utils.String(metric["time_window"].(string)),
				TimeAggregation:   insights.TimeAggregationType(metric["time_aggregation"].(string)),
				Operator:          insights.ComparisonOperationType(metric["operator"].(string)),
				Threshold:         utils.Float(metric["threshold"].(float64)),
			},
			ScaleAction: &insights.ScaleAction{
				Direction: insights.ScaleDirection(scale["direction"].(string)),
				Type:      insights.ScaleType(scale["type"].(string)),
				Value:     utils.String(scale["value"].(string)),
				Cooldown:  utils.String(scale["cooldown"].(string)),
			},
		}

		rules = append(rules, rule)
	}
	return &rules
}

func expandAzureRmAutoscaleSettingFixedDate(o map[string]interface{}) (*insights.TimeWindow, error) {
	startTime, err := date.ParseTime(time.RFC3339, o["start"].(string))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse `start` field: %+v", err)
	}
	endTime, err := date.ParseTime(time.RFC3339, o["end"].(string))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse `end` field: %+v", err)
	}
	return &insights.TimeWindow{
		TimeZone: utils.String(o["time_zone"].(string)),
		Start:    &date.Time{Time: startTime},
		End:      &date.Time{Time: endTime},
	}, nil
}

func expandAzureRmAutoscaleSettingRecurrence(o map[string]interface{}) *insights.Recurrence {
	days := make([]string, 0)
	for _, dayItem := range o["days"].([]interface{}) {
		days = append(days, dayItem.(string))
	}
	hours := make([]int32, 0)
	for _, hourItem := range o["hours"].([]interface{}) {
		hours = append(hours, hourItem.(int32))
	}
	minutes := make([]int32, 0)
	for _, minuteItem := range o["minutes"].([]interface{}) {
		minutes = append(minutes, minuteItem.(int32))
	}
	return &insights.Recurrence{
		Frequency: insights.RecurrenceFrequencyWeek,
		Schedule: &insights.RecurrentSchedule{
			TimeZone: utils.String(o["time_zone"].(string)),
			Days:     &days,
			Hours:    &hours,
			Minutes:  &minutes,
		},
	}
}

func expandAzureRmAutoscaleSettingNotifications(v []interface{}) *[]insights.AutoscaleNotification {
	notifications := make([]insights.AutoscaleNotification, 0)
	for _, notificationItem := range v {
		val := notificationItem.(map[string]interface{})
		notification := insights.AutoscaleNotification{
			Operation: utils.String("scale"),
		}
		if v, ok := val["email"]; ok {
			vObj := v.([]interface{})[0]
			notification.Email = expandAzureRmAutoscaleSettingNotificationEmail(vObj.(map[string]interface{}))
		}
		if v, ok := val["config"]; ok {
			notification.Webhooks = expandAzureRmAutoscaleSettingNotificationWebhook(v.([]interface{}))
		}
		notifications = append(notifications, notification)
	}
	return &notifications
}

func expandAzureRmAutoscaleSettingNotificationEmail(o map[string]interface{}) *insights.EmailNotification {
	email := insights.EmailNotification{
		SendToSubscriptionAdministrator:    utils.Bool(o["send_to_subscription_administrator"].(bool)),
		SendToSubscriptionCoAdministrators: utils.Bool(o["send_to_subscription_co_administrator"].(bool)),
	}
	if v, ok := o["custom_emails"]; ok {
		customEmails := make([]string, 0)
		for _, item := range v.([]interface{}) {
			customEmails = append(customEmails, item.(string))
		}
		email.CustomEmails = &customEmails
	}
	return &email
}

func expandAzureRmAutoscaleSettingNotificationWebhook(v []interface{}) *[]insights.WebhookNotification {
	webhooks := make([]insights.WebhookNotification, 0)
	for _, webhookItem := range v {
		val := webhookItem.(map[string]interface{})
		webhook := insights.WebhookNotification{
			ServiceURI: utils.String(val["service_uri"].(string)),
		}
		if v, ok := val["properties"]; ok {
			properties := make(map[string]*string)
			for key, value := range v.(map[string]interface{}) {
				properties[key] = utils.String(value.(string))
			}
			webhook.Properties = properties
		}
		webhooks = append(webhooks, webhook)
	}
	return &webhooks
}

func flattenAzureRmAutoscaleSettingProfile(profiles *[]insights.AutoscaleProfile) ([]interface{}, error) {
	result := make([]interface{}, 0)
	for _, profile := range *profiles {
		v := make(map[string]interface{})
		v["name"] = *profile.Name
		vCap, err := flattenAzureRmAutoscaleSettingCapacity(profile.Capacity)
		if err != nil {
			return nil, err
		}
		v["capacity"] = vCap
		if profile.Rules != nil {
			v["rule"] = flattenAzureRmAutoscaleSettingRules(profile.Rules)
		}
		if profile.FixedDate != nil {
			v["fixed_date"] = flattenAzureRmAutoscaleSettingFixedDate(profile.FixedDate)
		}
		if profile.Recurrence != nil {
			v["recurrence"] = flattenAzureRmAutoscaleSettingRecurrence(profile.Recurrence)
		}
		result = append(result, v)
	}
	return result, nil
}

func flattenAzureRmAutoscaleSettingCapacity(capacity *insights.ScaleCapacity) ([]interface{}, error) {
	result := make(map[string]interface{})
	vMin, err := strconv.Atoi(*capacity.Minimum)
	if err != nil {
		return nil, err
	}
	result["minimum"] = vMin
	vMax, err := strconv.Atoi(*capacity.Maximum)
	if err != nil {
		return nil, err
	}
	result["maximum"] = vMax
	vDef, err := strconv.Atoi(*capacity.Default)
	if err != nil {
		return nil, err
	}
	result["default"] = vDef
	return []interface{}{result}, nil
}

func flattenAzureRmAutoscaleSettingRules(rules *[]insights.ScaleRule) []interface{} {
	result := make([]interface{}, 0)
	for _, rule := range *rules {
		vRule := make(map[string]interface{})

		vMetric := make(map[string]interface{})
		vMetric["metric_name"] = *rule.MetricTrigger.MetricName
		vMetric["metric_resource_id"] = *rule.MetricTrigger.MetricResourceURI
		vMetric["time_grain"] = *rule.MetricTrigger.TimeGrain
		vMetric["statistic"] = string(rule.MetricTrigger.Statistic)
		vMetric["time_window"] = *rule.MetricTrigger.TimeWindow
		vMetric["time_aggregation"] = string(rule.MetricTrigger.TimeAggregation)
		vMetric["operator"] = string(rule.MetricTrigger.Operator)
		vMetric["threshold"] = *rule.MetricTrigger.Threshold
		vRule["metric_trigger"] = []interface{}{vMetric}

		vScale := make(map[string]interface{})
		vScale["direction"] = string(rule.ScaleAction.Direction)
		vScale["type"] = string(rule.ScaleAction.Type)
		vScale["value"] = *rule.ScaleAction.Value
		vScale["cooldown"] = *rule.ScaleAction.Cooldown
		vRule["scale_action"] = []interface{}{vScale}

		result = append(result, vRule)
	}
	return result
}

func flattenAzureRmAutoscaleSettingFixedDate(fixedDate *insights.TimeWindow) []interface{} {
	result := make(map[string]interface{})
	result["time_zone"] = *fixedDate.TimeZone
	result["start"] = fixedDate.Start.String()
	result["end"] = fixedDate.End.String()
	return []interface{}{result}
}

func flattenAzureRmAutoscaleSettingRecurrence(recurrence *insights.Recurrence) []interface{} {
	result := make(map[string]interface{})
	result["time_zone"] = *recurrence.Schedule.TimeZone

	days := make([]string, 0)
	for _, v := range *recurrence.Schedule.Days {
		days = append(days, v)
	}
	result["days"] = days

	hours := make([]int, 0)
	for _, v := range *recurrence.Schedule.Hours {
		hours = append(hours, int(v))
	}
	result["hours"] = hours

	minutes := make([]int, 0)
	for _, v := range *recurrence.Schedule.Minutes {
		minutes = append(minutes, int(v))
	}
	result["minutes"] = minutes

	return []interface{}{result}
}

func flattenAzureRmAutoscaleSettingNotification(notifications *[]insights.AutoscaleNotification) []interface{} {
	results := make([]interface{}, 0)

	if notifications == nil {
		return results
	}

	for _, notification := range *notifications {
		result := make(map[string]interface{})

		emails := make([]interface{}, 0)
		if email := notification.Email; email != nil {
			result := make(map[string]interface{}, 0)

			if send := email.SendToSubscriptionAdministrator; send != nil {
				result["send_to_subscription_administrator"] = *send
			}

			if send := email.SendToSubscriptionCoAdministrators; send != nil {
				result["send_to_subscription_co_administrator"] = *send
			}

			if custom := email.CustomEmails; custom != nil {
				result["custom_emails"] = *custom
			}

			emails = append(emails, result)
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

func validateAutoScaleSettingsTimeZone() schema.SchemaValidateFunc {
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
