package azurerm

import (
	"bytes"
	"sort"
	"strconv"
	"time"

	"fmt"

	"net/http"

	"github.com/Azure/azure-sdk-for-go/arm/insights"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArmAutoscaleSettings() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmAutoscaleSettingsCreateOrUpdate,
		Read:   resourceArmAutoscaleSettingsRead,
		Update: resourceArmAutoscaleSettingsCreateOrUpdate,
		Delete: resourceArmAutoscaleSettingsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"location": locationSchema(),
			"profile": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 20, // https://msdn.microsoft.com/en-us/library/azure/dn931928.aspx
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"capacity": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"minimum": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"maximum": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"default": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
							Set: resourceAzureRmAutoscaleHash,
						},
						"rule": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							MaxItems: 10, // https://msdn.microsoft.com/en-us/library/azure/dn931928.aspx
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_trigger": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"metric_name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"metric_resource_uri": {
													Type:     schema.TypeString,
													Required: true,
												},
												"time_grain": {
													Type:     schema.TypeString,
													Required: true,
												},
												"statistic": {
													Type:     schema.TypeString,
													Required: true,
												},
												"time_window": {
													Type:     schema.TypeString,
													Required: true,
												},
												"time_aggregation": {
													Type:     schema.TypeString,
													Required: true,
												},
												"operator": {
													Type:     schema.TypeString,
													Required: true,
												},
												"threshold": {
													Type:     schema.TypeInt,
													Required: true,
												},
											},
										},
										Set: resourceAzureRmAutoscaleHash,
									},
									"scale_action": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"direction": {
													Type:     schema.TypeString,
													Required: true,
												},
												"type": {
													Type:     schema.TypeString,
													Required: true,
												},
												"value": {
													Type:     schema.TypeString,
													Required: true,
												},
												"cooldown": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
										Set: resourceAzureRmAutoscaleHash,
									},
								},
							},
						},
						"fixed_date": {
							Type:          schema.TypeSet,
							Optional:      true,
							MinItems:      1,
							MaxItems:      1,
							ConflictsWith: []string{"profile.recurrence"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time_zone": {
										Type:     schema.TypeString,
										Required: true,
									},
									"start": {
										Type:     schema.TypeString,
										Required: true,
									},
									"end": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
							Set: resourceAzureRmAutoscaleHash,
						},
						"recurrence": {
							Type:          schema.TypeSet,
							Optional:      true,
							MinItems:      1,
							MaxItems:      1,
							ConflictsWith: []string{"profile.fixed_date"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"frequency": {
										Type:     schema.TypeString,
										Required: true,
									},
									"schedule": {
										Type:     schema.TypeSet,
										Required: true,
										MinItems: 1,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"time_zone": {
													Type:     schema.TypeString,
													Required: true,
												},
												"days": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"hours": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
												"minutes": {
													Type:     schema.TypeList,
													Required: true,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
											},
										},
										Set: resourceAzureRmAutoscaleHash,
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
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operation": {
							Type:     schema.TypeString,
							Required: true,
						},
						"email": {
							Type:     schema.TypeSet,
							Required: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"send_to_subscription_administrator": {
										Type:     schema.TypeBool,
										Required: true,
									},
									"send_to_subscription_co_administrator": {
										Type:     schema.TypeBool,
										Required: true,
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
							Set: resourceAzureRmAutoscaleHash,
						},
						"webhook": {
							Type:     schema.TypeSet,
							Optional: true,
							MinItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_uri": {
										Type:     schema.TypeString,
										Required: true,
									},
									"properties": {
										Type:     schema.TypeMap,
										Optional: true,
									},
								},
							},
							Set: resourceAzureRmAutoscaleHash,
						},
					},
				},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"target_resource_uri": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceArmAutoscaleSettingsCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	asClient := armClient.autoscaleSettingsClient

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)
	resourceType := "Microsoft.Insights/autoscaleSettings"
	location := d.Get("location").(string)
	enabled := d.Get("enabled").(bool)
	targetResourceURI := d.Get("target_resource_uri").(string)
	tags := d.Get("tags").(map[string]interface{})
	expandedTags := expandTags(tags)

	profiles, err := expandAzureRmAutoscaleProfile(d)
	if err != nil {
		return err
	}

	notifications, err := expandAzureRmAutoscaleNotification(d)
	if err != nil {
		return err
	}

	autoscaleSettings := insights.AutoscaleSetting{
		Name:              &name,
		Enabled:           &enabled,
		TargetResourceURI: &targetResourceURI,
		Profiles:          &profiles,
		Notifications:     &notifications,
	}

	parameters := insights.AutoscaleSettingResource{
		Name:             &name,
		Type:             &resourceType,
		Location:         &location,
		Tags:             expandedTags,
		AutoscaleSetting: &autoscaleSettings,
	}

	result, err := asClient.CreateOrUpdate(resourceGroupName, name, parameters)
	if err != nil {
		return err
	}

	d.SetId(*result.ID)

	return resourceArmAutoscaleSettingsRead(d, meta)
}

func resourceArmAutoscaleSettingsRead(d *schema.ResourceData, meta interface{}) error {
	asClient := meta.(*ArmClient).autoscaleSettingsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["autoscalesettings"]

	if name == "" {
		return fmt.Errorf("Cannot find resource name in Resource ID for Autoscaling Settings")
	}

	result, err := asClient.Get(resGroup, name)
	if err != nil {
		if result.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Autoscaling Settings %s: %s", name, err)
	}

	autoscaleSettings := *result.AutoscaleSetting
	d.Set("name", result.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", result.Location)
	d.Set("enabled", autoscaleSettings.Enabled)
	d.Set("target_resource_uri", autoscaleSettings.TargetResourceURI)
	flattenAndSetTags(d, result.Tags)

	d.Set("profile", flattenAzureRmAutoscaleProfile(autoscaleSettings.Profiles))
	d.Set("notification", flattenAzureRmAutoscaleNotification(autoscaleSettings.Notifications))

	return nil
}

func resourceArmAutoscaleSettingsDelete(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	asClient := armClient.autoscaleSettingsClient

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroupName := id.ResourceGroup
	autoscaleSettingName := id.Path["autoscalesettings"]

	_, err = asClient.Delete(resourceGroupName, autoscaleSettingName)
	return err
}

func expandAzureRmAutoscaleProfile(d *schema.ResourceData) ([]insights.AutoscaleProfile, error) {
	profileData := d.Get("profile").([]interface{})
	profiles := make([]insights.AutoscaleProfile, 0, len(profileData))

	for _, profileConfig := range profileData {
		profile := profileConfig.(map[string]interface{})

		profileName := profile["name"].(string)
		capacity := expandAzureRmAutoscaleCapacity(profile)
		rules := expandAzureRmAutoscaleRule(profile)
		autoscaleProfile := insights.AutoscaleProfile{
			Name:     &profileName,
			Capacity: &capacity,
			Rules:    &rules,
		}
		fixedDate, fixedDateErr := expandAzureRmAutoscaleRuleTimeWindow(profile)
		recurrence, recurrenceErr := expandAzureRmAutoscaleRecurrence(profile)

		if fixedDateErr == nil {
			autoscaleProfile.FixedDate = &fixedDate
		} else if recurrenceErr == nil {
			autoscaleProfile.Recurrence = &recurrence
		}

		profiles = append(profiles, autoscaleProfile)
	}

	return profiles, nil
}

func expandAzureRmAutoscaleCapacity(config map[string]interface{}) insights.ScaleCapacity {
	capacitySet := config["capacity"].(*schema.Set).List()
	capacityConfig := capacitySet[0].(map[string]interface{})
	min := strconv.Itoa(capacityConfig["minimum"].(int))
	max := strconv.Itoa(capacityConfig["maximum"].(int))
	defaultValue := strconv.Itoa(capacityConfig["default"].(int))

	scaleCapacity := insights.ScaleCapacity{
		Minimum: &min,
		Maximum: &max,
		Default: &defaultValue,
	}

	return scaleCapacity
}

func expandAzureRmAutoscaleRule(config map[string]interface{}) []insights.ScaleRule {
	ruleSet := config["rule"].([]interface{})
	scaleRules := make([]insights.ScaleRule, 0, len(ruleSet))

	for _, ruleConfig := range ruleSet {
		rule := ruleConfig.(map[string]interface{})
		metricTrigger := expandAzureRmMetricTrigger(rule)
		scaleAction := expandAzureRmScaleAction(rule)

		scaleRule := insights.ScaleRule{
			MetricTrigger: &metricTrigger,
			ScaleAction:   &scaleAction,
		}
		scaleRules = append(scaleRules, scaleRule)
	}

	return scaleRules
}

func expandAzureRmMetricTrigger(config map[string]interface{}) insights.MetricTrigger {
	metricTriggerSet := config["metric_trigger"].(*schema.Set).List()
	metricTriggerConfig := metricTriggerSet[0].(map[string]interface{})
	metricName := metricTriggerConfig["metric_name"].(string)
	metricResourceURI := metricTriggerConfig["metric_resource_uri"].(string)
	timeGrain := metricTriggerConfig["time_grain"].(string)
	statistic := metricTriggerConfig["statistic"].(string)
	timeWindow := metricTriggerConfig["time_window"].(string)
	timeAggregation := metricTriggerConfig["time_aggregation"].(string)
	operator := metricTriggerConfig["operator"].(string)
	threshold := float64(metricTriggerConfig["threshold"].(int))

	return insights.MetricTrigger{
		MetricName:        &metricName,
		MetricResourceURI: &metricResourceURI,
		TimeGrain:         &timeGrain,
		Statistic:         insights.MetricStatisticType(statistic),
		TimeWindow:        &timeWindow,
		TimeAggregation:   insights.TimeAggregationType(timeAggregation),
		Operator:          insights.ComparisonOperationType(operator),
		Threshold:         &threshold,
	}
}

func expandAzureRmScaleAction(config map[string]interface{}) insights.ScaleAction {
	scaleActionSet := config["scale_action"].(*schema.Set).List()
	scaleActionConfig := scaleActionSet[0].(map[string]interface{})
	direction := scaleActionConfig["direction"].(string)
	scaleType := scaleActionConfig["type"].(string)
	value := scaleActionConfig["value"].(string)
	cooldown := scaleActionConfig["cooldown"].(string)

	return insights.ScaleAction{
		Direction: insights.ScaleDirection(direction),
		Type:      insights.ScaleType(scaleType),
		Value:     &value,
		Cooldown:  &cooldown,
	}
}

func expandAzureRmAutoscaleRuleTimeWindow(config map[string]interface{}) (insights.TimeWindow, error) {
	timeWindow := insights.TimeWindow{}
	fixedDateSet := config["fixed_date"].(*schema.Set).List()

	if fixedDateSet == nil || len(fixedDateSet) == 0 {
		return timeWindow, fmt.Errorf("fixed_date not defined")
	}

	fixedDateConfig := fixedDateSet[0].(map[string]interface{})
	timeZone := fixedDateConfig["time_zone"].(string)
	startString := fixedDateConfig["start"].(string)
	endString := fixedDateConfig["end"].(string)

	startTime, startTimeErr := date.ParseTime(time.RFC3339, startString)
	if startTimeErr != nil {
		return timeWindow, startTimeErr
	}

	endTime, endTimeErr := date.ParseTime(time.RFC3339, endString)
	if endTimeErr != nil {
		return timeWindow, endTimeErr
	}

	timeWindow.TimeZone = &timeZone
	timeWindow.Start = &date.Time{Time: startTime}
	timeWindow.End = &date.Time{Time: endTime}

	return timeWindow, nil
}

func expandAzureRmAutoscaleRecurrence(config map[string]interface{}) (insights.Recurrence, error) {
	recurrence := insights.Recurrence{}
	recurrenceSet := config["recurrence"].(*schema.Set).List()

	if recurrenceSet == nil || len(recurrenceSet) == 0 {
		return recurrence, fmt.Errorf("recurrence not defined")
	}

	recurrenceConfig := recurrenceSet[0].(map[string]interface{})
	schedule := expandAzureRmAutoscaleRecurrentSchedule(recurrenceConfig)
	recurrence.Frequency = insights.RecurrenceFrequency(recurrenceConfig["frequency"].(string))
	recurrence.Schedule = &schedule
	return recurrence, nil
}

func expandAzureRmAutoscaleRecurrentSchedule(config map[string]interface{}) insights.RecurrentSchedule {
	scheduleSet := config["schedule"].(*schema.Set).List()
	scheduleConfig := scheduleSet[0].(map[string]interface{})
	timeZone := scheduleConfig["time_zone"].(string)

	daysConfig := scheduleConfig["days"].([]interface{})
	days := make([]string, len(daysConfig))
	for i, v := range daysConfig {
		days[i] = v.(string)
	}

	hoursConfig := scheduleConfig["hours"].([]interface{})
	hours := make([]int32, len(hoursConfig))
	for i, v := range hoursConfig {
		hours[i] = int32(v.(int))
	}

	minutesConfig := scheduleConfig["minutes"].([]interface{})
	minutes := make([]int32, len(minutesConfig))
	for i, v := range minutesConfig {
		minutes[i] = int32(v.(int))
	}

	return insights.RecurrentSchedule{
		TimeZone: &timeZone,
		Days:     &days,
		Hours:    &hours,
		Minutes:  &minutes,
	}
}

func expandAzureRmAutoscaleNotification(d *schema.ResourceData) ([]insights.AutoscaleNotification, error) {
	r, ok := d.GetOk("notification")
	if !ok {
		return nil, nil
	}

	notificationData := r.([]interface{})
	notifications := make([]insights.AutoscaleNotification, 0, len(notificationData))

	for _, item := range notificationData {
		notificationConfig := item.(map[string]interface{})
		operation := notificationConfig["operation"].(string)
		email := expandAzureRmAutoscaleEmailNotification(notificationConfig)
		webhooks := expandAzureRmAutoscaleWebhook(notificationConfig)

		notification := insights.AutoscaleNotification{
			Operation: &operation,
			Email:     &email,
			Webhooks:  &webhooks,
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func expandAzureRmAutoscaleEmailNotification(config map[string]interface{}) insights.EmailNotification {
	emailSet := config["email"].(*schema.Set).List()
	emailConfig := emailSet[0].(map[string]interface{})
	sendToAdmin := emailConfig["send_to_subscription_administrator"].(bool)
	sendToCoAdmin := emailConfig["send_to_subscription_co_administrator"].(bool)
	customEmailsConfig := emailConfig["custom_emails"].([]interface{})
	customEmails := make([]string, len(customEmailsConfig))
	for i, v := range customEmailsConfig {
		customEmails[i] = v.(string)
	}

	email := insights.EmailNotification{
		SendToSubscriptionAdministrator:    &sendToAdmin,
		SendToSubscriptionCoAdministrators: &sendToCoAdmin,
		CustomEmails:                       &customEmails,
	}

	return email
}

func expandAzureRmAutoscaleWebhook(config map[string]interface{}) []insights.WebhookNotification {
	r := config["webhook"]
	if r == nil {
		return nil
	}

	webhookData := r.(*schema.Set).List()
	webhooks := make([]insights.WebhookNotification, 0, len(webhookData))

	for _, item := range webhookData {
		webhookConfig := item.(map[string]interface{})
		serviceURI := webhookConfig["service_uri"].(string)
		webhook := insights.WebhookNotification{
			ServiceURI: &serviceURI,
		}

		p := webhookConfig["properties"]
		if p != nil {
			propertiesConfig := p.(map[string]interface{})

			properties := make(map[string]*string, len(propertiesConfig))
			for k, v := range propertiesConfig {
				value := fmt.Sprintf("%v", v)
				properties[k] = &value
			}

			webhook.Properties = &properties
		}

		webhooks = append(webhooks, webhook)
	}

	return webhooks
}

// func flattenAzureRmAutoscaleProfile(profiles *[]insights.AutoscaleProfile) *schema.Set {
// 	results := make([]interface{}, 0, len(*profiles))

// 	for _, profile := range *profiles {
// 		profileConfig := make(map[string]interface{})
// 		profileConfig["name"] = *profile.Name
// 		profileConfig["capacity"] = flattenAzureRmAutoscaleCapacity(profile)
// 		profileConfig["rule"] = flattenAzureRmAutoscaleRule(profile)

// 		if profile.FixedDate != nil {
// 			profileConfig["fixed_date"] = flattenAzureRmAutoscaleFixedDate(profile)
// 		}

// 		if profile.Recurrence != nil {
// 			profileConfig["recurrence"] = flattenAzureRmAutoscaleRecurrence(profile)
// 		}

// 		results = append(results, profileConfig)
// 	}

// 	return schema.NewSet(resourceAzureRmAutoscaleHash, results)
// }

func flattenAzureRmAutoscaleProfile(profiles *[]insights.AutoscaleProfile) []interface{} {
	results := make([]interface{}, 0, len(*profiles))

	for _, profile := range *profiles {
		profileConfig := make(map[string]interface{})
		profileConfig["name"] = *profile.Name
		profileConfig["capacity"] = flattenAzureRmAutoscaleCapacity(profile)
		profileConfig["rule"] = flattenAzureRmAutoscaleRule(profile)

		if profile.FixedDate != nil {
			profileConfig["fixed_date"] = flattenAzureRmAutoscaleFixedDate(profile)
		}

		if profile.Recurrence != nil {
			profileConfig["recurrence"] = flattenAzureRmAutoscaleRecurrence(profile)
		}

		results = append(results, profileConfig)
	}

	return results
}

func flattenAzureRmAutoscaleCapacity(profile insights.AutoscaleProfile) *schema.Set {
	capacity := make(map[string]interface{})
	capacity["minimum"], _ = strconv.Atoi(*profile.Capacity.Minimum)
	capacity["maximum"], _ = strconv.Atoi(*profile.Capacity.Maximum)
	capacity["default"], _ = strconv.Atoi(*profile.Capacity.Default)

	return schema.NewSet(resourceAzureRmAutoscaleHash, []interface{}{capacity})
}

// func flattenAzureRmAutoscaleCapacity(profile insights.AutoscaleProfile) []interface{} {
// 	capacity := make(map[string]interface{})
// 	capacity["minimum"], _ = strconv.Atoi(*profile.Capacity.Minimum)
// 	capacity["maximum"], _ = strconv.Atoi(*profile.Capacity.Maximum)
// 	capacity["default"], _ = strconv.Atoi(*profile.Capacity.Default)

// 	return []interface{}{capacity}
// }

// func flattenAzureRmAutoscaleRule(profile insights.AutoscaleProfile) *schema.Set {
// 	rules := make([]interface{}, 0, len(*profile.Rules))

// 	for _, rule := range *profile.Rules {
// 		ruleConfig := make(map[string]interface{})
// 		ruleConfig["metric_trigger"] = flattenAzureRmAutoscaleMetricTrigger(rule)
// 		ruleConfig["scale_action"] = flattenAzureRmAutoscaleAction(rule)

// 		rules = append(rules, ruleConfig)
// 	}

// 	return schema.NewSet(resourceAzureRmAutoscaleHash, rules)
// }

func flattenAzureRmAutoscaleRule(profile insights.AutoscaleProfile) []interface{} {
	rules := make([]interface{}, 0, len(*profile.Rules))

	for _, rule := range *profile.Rules {
		ruleConfig := make(map[string]interface{})
		ruleConfig["metric_trigger"] = flattenAzureRmAutoscaleMetricTrigger(rule)
		ruleConfig["scale_action"] = flattenAzureRmAutoscaleAction(rule)

		rules = append(rules, ruleConfig)
	}

	return rules
}

func flattenAzureRmAutoscaleMetricTrigger(rule insights.ScaleRule) *schema.Set {
	metricTrigger := make(map[string]interface{})
	metricTrigger["metric_name"] = *rule.MetricTrigger.MetricName
	metricTrigger["metric_resource_uri"] = *rule.MetricTrigger.MetricResourceURI
	metricTrigger["time_grain"] = *rule.MetricTrigger.TimeGrain
	metricTrigger["statistic"] = string(rule.MetricTrigger.Statistic)
	metricTrigger["time_window"] = *rule.MetricTrigger.TimeWindow
	metricTrigger["time_aggregation"] = string(rule.MetricTrigger.TimeAggregation)
	metricTrigger["operator"] = string(rule.MetricTrigger.Operator)
	metricTrigger["threshold"] = int(*rule.MetricTrigger.Threshold)

	return schema.NewSet(resourceAzureRmAutoscaleHash, []interface{}{metricTrigger})
}

// func flattenAzureRmAutoscaleMetricTrigger(rule insights.ScaleRule) []interface{} {
// 	metricTrigger := make(map[string]interface{})
// 	metricTrigger["metric_name"] = *rule.MetricTrigger.MetricName
// 	metricTrigger["metric_resource_uri"] = *rule.MetricTrigger.MetricResourceURI
// 	metricTrigger["time_grain"] = *rule.MetricTrigger.TimeGrain
// 	metricTrigger["statistic"] = string(rule.MetricTrigger.Statistic)
// 	metricTrigger["time_window"] = *rule.MetricTrigger.TimeWindow
// 	metricTrigger["time_aggregation"] = string(rule.MetricTrigger.TimeAggregation)
// 	metricTrigger["operator"] = string(rule.MetricTrigger.Operator)
// 	metricTrigger["threshold"] = int(*rule.MetricTrigger.Threshold)

// 	return []interface{}{metricTrigger}
// }

func flattenAzureRmAutoscaleAction(rule insights.ScaleRule) *schema.Set {
	scaleAction := make(map[string]interface{})
	scaleAction["direction"] = string(rule.ScaleAction.Direction)
	scaleAction["type"] = string(rule.ScaleAction.Type)
	scaleAction["value"] = *rule.ScaleAction.Value
	scaleAction["cooldown"] = *rule.ScaleAction.Cooldown

	return schema.NewSet(resourceAzureRmAutoscaleHash, []interface{}{scaleAction})
}

// func flattenAzureRmAutoscaleAction(rule insights.ScaleRule) []interface{} {
// 	scaleAction := make(map[string]interface{})
// 	scaleAction["direction"] = string(rule.ScaleAction.Direction)
// 	scaleAction["type"] = string(rule.ScaleAction.Type)
// 	scaleAction["value"] = *rule.ScaleAction.Value
// 	scaleAction["cooldown"] = *rule.ScaleAction.Cooldown

// 	return []interface{}{scaleAction}
// }

func flattenAzureRmAutoscaleFixedDate(profile insights.AutoscaleProfile) *schema.Set {
	fixedDate := make(map[string]interface{})
	fixedDate["time_zone"] = *profile.FixedDate.TimeZone
	fixedDate["start"] = profile.FixedDate.Start.String()
	fixedDate["end"] = profile.FixedDate.End.String()

	return schema.NewSet(resourceAzureRmAutoscaleHash, []interface{}{fixedDate})
}

// func flattenAzureRmAutoscaleFixedDate(profile insights.AutoscaleProfile) []interface{} {
// 	fixedDate := make(map[string]interface{})
// 	fixedDate["time_zone"] = *profile.FixedDate.TimeZone
// 	fixedDate["start"] = profile.FixedDate.Start.String()
// 	fixedDate["end"] = profile.FixedDate.End.String()

// 	return []interface{}{fixedDate}
// }

func flattenAzureRmAutoscaleRecurrence(profile insights.AutoscaleProfile) *schema.Set {
	recurrence := make(map[string]interface{})
	recurrence["frequency"] = string(profile.Recurrence.Frequency)
	recurrence["schedule"] = flattenAzureRmAutoscaleSchedule(*profile.Recurrence)

	return schema.NewSet(resourceAzureRmAutoscaleHash, []interface{}{recurrence})
}

// func flattenAzureRmAutoscaleRecurrence(profile insights.AutoscaleProfile) []interface{} {
// 	recurrence := make(map[string]interface{})
// 	recurrence["frequency"] = string(profile.Recurrence.Frequency)
// 	recurrence["schedule"] = flattenAzureRmAutoscaleSchedule(*profile.Recurrence)

// 	return []interface{}{recurrence}
// }

func flattenAzureRmAutoscaleSchedule(recurrence insights.Recurrence) *schema.Set {
	schedule := make(map[string]interface{})
	schedule["time_zone"] = *recurrence.Schedule.TimeZone

	days := make([]interface{}, len(*recurrence.Schedule.Days))
	for k, v := range *recurrence.Schedule.Days {
		days[k] = v
	}
	schedule["days"] = days

	hours := make([]interface{}, len(*recurrence.Schedule.Hours))
	for k, v := range *recurrence.Schedule.Hours {
		hours[k] = v
	}
	schedule["hours"] = hours

	minutes := make([]interface{}, len(*recurrence.Schedule.Minutes))
	for k, v := range *recurrence.Schedule.Minutes {
		minutes[k] = v
	}
	schedule["minutes"] = minutes

	return schema.NewSet(resourceAzureRmAutoscaleHash, []interface{}{schedule})
}

// func flattenAzureRmAutoscaleSchedule(recurrence insights.Recurrence) []interface{} {
// 	schedule := make(map[string]interface{})
// 	schedule["time_zone"] = *recurrence.Schedule.TimeZone

// 	days := make([]interface{}, len(*recurrence.Schedule.Days))
// 	for k, v := range *recurrence.Schedule.Days {
// 		days[k] = v
// 	}
// 	schedule["days"] = days

// 	hours := make([]interface{}, len(*recurrence.Schedule.Hours))
// 	for k, v := range *recurrence.Schedule.Hours {
// 		hours[k] = v
// 	}
// 	schedule["hours"] = hours

// 	minutes := make([]interface{}, len(*recurrence.Schedule.Minutes))
// 	for k, v := range *recurrence.Schedule.Minutes {
// 		minutes[k] = v
// 	}
// 	schedule["minutes"] = minutes

// 	return []interface{}{schedule}
// }

// func flattenAzureRmAutoscaleNotification(notifications *[]insights.AutoscaleNotification) *schema.Set {
// 	results := make([]interface{}, 0, len(*notifications))

// 	for _, notification := range *notifications {
// 		notificationConfig := make(map[string]interface{})
// 		notificationConfig["operation"] = *notification.Operation
// 		notificationConfig["email"] = flattenAzureRmAutoscaleEmailNotification(notification)
// 		if *notification.Webhooks != nil {
// 			notificationConfig["webhook"] = flattenAzureRmAutoscaleWebhook(notification)
// 		}
// 		results = append(results, notificationConfig)
// 	}

// 	return schema.NewSet(resourceAzureRmAutoscaleHash, results)
// }

func flattenAzureRmAutoscaleNotification(notifications *[]insights.AutoscaleNotification) []interface{} {
	results := make([]interface{}, 0, len(*notifications))

	for _, notification := range *notifications {
		notificationConfig := make(map[string]interface{})
		notificationConfig["operation"] = *notification.Operation
		notificationConfig["email"] = flattenAzureRmAutoscaleEmailNotification(notification)
		if *notification.Webhooks != nil {
			notificationConfig["webhook"] = flattenAzureRmAutoscaleWebhook(notification)
		}
		results = append(results, notificationConfig)
	}

	return results
}

func flattenAzureRmAutoscaleEmailNotification(notification insights.AutoscaleNotification) *schema.Set {
	email := make(map[string]interface{})
	email["send_to_subscription_administrator"] = *notification.Email.SendToSubscriptionAdministrator
	email["send_to_subscription_co_administrator"] = *notification.Email.SendToSubscriptionCoAdministrators

	if *notification.Email.CustomEmails != nil {
		email["custom_emails"] = *notification.Email.CustomEmails
	}

	return schema.NewSet(resourceAzureRmAutoscaleHash, []interface{}{email})
}

// func flattenAzureRmAutoscaleEmailNotification(notification insights.AutoscaleNotification) []interface{} {
// 	email := make(map[string]interface{})
// 	email["send_to_subscription_administrator"] = *notification.Email.SendToSubscriptionAdministrator
// 	email["send_to_subscription_co_administrator"] = *notification.Email.SendToSubscriptionCoAdministrators

// 	if *notification.Email.CustomEmails != nil {
// 		email["custom_emails"] = *notification.Email.CustomEmails
// 	}

// 	return []interface{}{email}
// }

func flattenAzureRmAutoscaleWebhook(notification insights.AutoscaleNotification) *schema.Set {
	set := &schema.Set{
		F: resourceAzureRmAutoscaleHash,
	}

	for _, v := range *notification.Webhooks {
		webhook := map[string]interface{}{}
		webhook["service_uri"] = *v.ServiceURI

		if *v.Properties != nil {
			properties := map[string]interface{}{}
			for key, value := range *v.Properties {
				properties[key] = *value
			}
			webhook["properties"] = properties
		}
		set.Add(webhook)
	}

	return set
}

// func flattenAzureRmAutoscaleWebhook(notification insights.AutoscaleNotification) []interface{} {
// 	webhooks := make([]interface{}, len(*notification.Webhooks))

// 	for _, v := range *notification.Webhooks {
// 		webhook := map[string]interface{}{}
// 		webhook["service_uri"] = *v.ServiceURI

// 		if *v.Properties != nil {
// 			properties := map[string]interface{}{}
// 			for key, value := range *v.Properties {
// 				properties[key] = *value
// 			}
// 			webhook["properties"] = properties
// 		}
// 		webhooks = append(webhooks, webhook)
// 	}

// 	return webhooks
// }

func resourceAzureRmAutoscaleHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	var keys []string

	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		buf.WriteString(fmt.Sprintf("%s:%v;", k, m[k]))
	}

	return hashcode.String(buf.String())
}
