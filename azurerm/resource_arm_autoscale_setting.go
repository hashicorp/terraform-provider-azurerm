package azurerm

import (
	"bytes"
	"sort"
	"strconv"
	"time"

	"fmt"

	"net/http"

	"regexp"

	"strings"

	"github.com/Azure/azure-sdk-for-go/services/monitor/mgmt/2017-05-01-preview/insights"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceArmAutoscaleSetting() *schema.Resource {
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
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"minimum": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 100),
									},
									"maximum": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 100),
									},
									"default": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 100),
									},
								},
							},
							Set: resourceAzureRmAutoscaleDefaultHash,
						},
						"rule": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 10, // https://msdn.microsoft.com/en-us/library/azure/dn931928.aspx
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_trigger": {
										Type:     schema.TypeSet,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"metric_name": {
													Type:     schema.TypeString,
													Required: true,
												},
												"metric_resource_id": {
													Type:     schema.TypeString,
													Required: true,
												},
												"time_grain": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: iso8601DurationString(),
												},
												"statistic": {
													Type:     schema.TypeString,
													Required: true,
													ValidateFunc: validation.StringInSlice([]string{
														string(insights.Average),
														string(insights.Max),
														string(insights.Min),
														string(insights.Sum),
													}, true),
													DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
												},
												"time_window": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: iso8601DurationString(),
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
													Type:     schema.TypeInt,
													Required: true,
												},
											},
										},
										Set: resourceAzureRmAutoscaleDefaultHash,
									},
									"scale_action": {
										Type:     schema.TypeSet,
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
														string(insights.PercentChangeCount),
													}, true),
													DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
												},
												"value": {
													Type:     schema.TypeString,
													Required: true,
												},
												"cooldown": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: iso8601DurationString(),
												},
											},
										},
										Set: resourceAzureRmAutoscaleDefaultHash,
									},
								},
							},
						},
						"fixed_date": {
							Type:          schema.TypeSet,
							Optional:      true,
							MaxItems:      1,
							ConflictsWith: []string{"profile.recurrence"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time_zone": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(listTimeZoneNames(), false),
									},
									"start": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: rfc3339String(),
									},
									"end": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: rfc3339String(),
									},
								},
							},
							Set: resourceAzureRmAutoscaleDefaultHash,
						},
						"recurrence": {
							Type:          schema.TypeSet,
							Optional:      true,
							MaxItems:      1,
							ConflictsWith: []string{"profile.fixed_date"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"frequency": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											string(insights.Week),
										}, true),
										DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
									},
									"schedule": {
										Type:     schema.TypeSet,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"time_zone": {
													Type:         schema.TypeString,
													Required:     true,
													ValidateFunc: validation.StringInSlice(listTimeZoneNames(), false),
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
										Set: resourceAzureRmAutoscaleDefaultHash,
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
						"operation": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"Scale"}, true),
						},
						"email": {
							Type:     schema.TypeSet,
							Required: true,
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
							Set: resourceAzureRmAutoscaleDefaultHash,
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
							Set: resourceAzureRmAutoscaleDefaultHash,
						},
					},
				},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"target_resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceArmAutoscaleSettingCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	asClient := armClient.autoscaleSettingsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)
	location := d.Get("location").(string)
	enabled := d.Get("enabled").(bool)
	targetResourceURI := d.Get("target_resource_id").(string)
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

	autoscaleSetting := insights.AutoscaleSetting{
		Name:              &name,
		Enabled:           &enabled,
		TargetResourceURI: &targetResourceURI,
		Profiles:          &profiles,
		Notifications:     &notifications,
	}

	parameters := insights.AutoscaleSettingResource{
		Name:             &name,
		Location:         &location,
		Tags:             expandedTags,
		AutoscaleSetting: &autoscaleSetting,
	}

	result, err := asClient.CreateOrUpdate(ctx, resourceGroupName, name, parameters)
	if err != nil {
		return err
	}

	d.SetId(*result.ID)

	return resourceArmAutoscaleSettingRead(d, meta)
}

func resourceArmAutoscaleSettingRead(d *schema.ResourceData, meta interface{}) error {
	asClient := meta.(*ArmClient).autoscaleSettingsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["autoscalesettings"]

	if name == "" {
		return fmt.Errorf("Cannot find resource name in Resource ID for Autoscale Setting")
	}

	result, err := asClient.Get(ctx, resGroup, name)
	if err != nil {
		if result.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Autoscale Setting %s: %+v", name, err)
	}

	if result.AutoscaleSetting == nil {
		return fmt.Errorf("Unexpected result when reading Autoscale Setting %s: %#v", name, result)
	}
	autoscaleSetting := *result.AutoscaleSetting
	d.Set("name", result.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("location", result.Location)
	d.Set("enabled", autoscaleSetting.Enabled)
	d.Set("target_resource_id", autoscaleSetting.TargetResourceURI)
	flattenAndSetTags(d, result.Tags)

	d.Set("profile", flattenAzureRmAutoscaleProfile(autoscaleSetting.Profiles))
	d.Set("notification", flattenAzureRmAutoscaleNotification(autoscaleSetting.Notifications))

	return nil
}

func resourceArmAutoscaleSettingDelete(d *schema.ResourceData, meta interface{}) error {
	armClient := meta.(*ArmClient)
	asClient := armClient.autoscaleSettingsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroupName := id.ResourceGroup
	autoscaleSettingName := id.Path["autoscalesettings"]

	_, err = asClient.Delete(ctx, resourceGroupName, autoscaleSettingName)
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

		fixedDate, fixedDateErr := expandAzureRmAutoscaleFixedDate(profile)
		recurrence, recurrenceErr := expandAzureRmAutoscaleRecurrence(profile)

		if fixedDate != nil && recurrence != nil {
			return nil, fmt.Errorf("Conflict between fixed_date and reucrrence in profile %s", profileName)
		}

		if fixedDateErr != nil {
			return nil, fixedDateErr
		}
		autoscaleProfile.FixedDate = fixedDate

		if recurrenceErr != nil {
			return nil, recurrenceErr
		}
		autoscaleProfile.Recurrence = recurrence

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
	metricResourceURI := metricTriggerConfig["metric_resource_id"].(string)
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

func expandAzureRmAutoscaleFixedDate(config map[string]interface{}) (*insights.TimeWindow, error) {
	timeWindow := &insights.TimeWindow{}
	fixedDateSet := config["fixed_date"].(*schema.Set).List()

	if fixedDateSet == nil || len(fixedDateSet) == 0 {
		return nil, nil
	}

	fixedDateConfig := fixedDateSet[0].(map[string]interface{})
	timeZone := fixedDateConfig["time_zone"].(string)
	startString := fixedDateConfig["start"].(string)
	endString := fixedDateConfig["end"].(string)

	startTime, startTimeErr := date.ParseTime(time.RFC3339, startString)
	if startTimeErr != nil {
		return nil, startTimeErr
	}

	endTime, endTimeErr := date.ParseTime(time.RFC3339, endString)
	if endTimeErr != nil {
		return nil, endTimeErr
	}

	timeWindow.TimeZone = &timeZone
	timeWindow.Start = &date.Time{Time: startTime}
	timeWindow.End = &date.Time{Time: endTime}

	return timeWindow, nil
}

func expandAzureRmAutoscaleRecurrence(config map[string]interface{}) (*insights.Recurrence, error) {
	recurrence := &insights.Recurrence{}
	recurrenceSet := config["recurrence"].(*schema.Set).List()

	if recurrenceSet == nil || len(recurrenceSet) == 0 {
		return nil, nil
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
		webhooks, err := expandAzureRmAutoscaleWebhook(notificationConfig)

		if err != nil {
			return nil, err
		}

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

func expandAzureRmAutoscaleWebhook(config map[string]interface{}) ([]insights.WebhookNotification, error) {
	r := config["webhook"]
	if r == nil {
		return nil, nil
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
				value, success := v.(string)

				if !success {
					return nil, fmt.Errorf("Expect string in webhook.properties values, got '%#v'", v)
				}

				properties[k] = &value
			}

			webhook.Properties = properties
		}

		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
}

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

	return schema.NewSet(resourceAzureRmAutoscaleDefaultHash, []interface{}{capacity})
}

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
	metricTrigger["metric_resource_id"] = *rule.MetricTrigger.MetricResourceURI
	metricTrigger["time_grain"] = *rule.MetricTrigger.TimeGrain
	metricTrigger["statistic"] = string(rule.MetricTrigger.Statistic)
	metricTrigger["time_window"] = *rule.MetricTrigger.TimeWindow
	metricTrigger["time_aggregation"] = string(rule.MetricTrigger.TimeAggregation)
	metricTrigger["operator"] = string(rule.MetricTrigger.Operator)
	metricTrigger["threshold"] = int(*rule.MetricTrigger.Threshold)

	return schema.NewSet(resourceAzureRmAutoscaleDefaultHash, []interface{}{metricTrigger})
}

func flattenAzureRmAutoscaleAction(rule insights.ScaleRule) *schema.Set {
	scaleAction := make(map[string]interface{})
	scaleAction["direction"] = string(rule.ScaleAction.Direction)
	scaleAction["type"] = string(rule.ScaleAction.Type)
	scaleAction["value"] = *rule.ScaleAction.Value
	scaleAction["cooldown"] = *rule.ScaleAction.Cooldown

	return schema.NewSet(resourceAzureRmAutoscaleDefaultHash, []interface{}{scaleAction})
}

func flattenAzureRmAutoscaleFixedDate(profile insights.AutoscaleProfile) *schema.Set {
	fixedDate := make(map[string]interface{})
	fixedDate["time_zone"] = *profile.FixedDate.TimeZone
	fixedDate["start"] = profile.FixedDate.Start.String()
	fixedDate["end"] = profile.FixedDate.End.String()

	return schema.NewSet(resourceAzureRmAutoscaleDefaultHash, []interface{}{fixedDate})
}

func flattenAzureRmAutoscaleRecurrence(profile insights.AutoscaleProfile) *schema.Set {
	recurrence := make(map[string]interface{})
	recurrence["frequency"] = string(profile.Recurrence.Frequency)
	recurrence["schedule"] = flattenAzureRmAutoscaleSchedule(*profile.Recurrence)

	return schema.NewSet(resourceAzureRmAutoscaleDefaultHash, []interface{}{recurrence})
}

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

	return schema.NewSet(resourceAzureRmAutoscaleDefaultHash, []interface{}{schedule})
}

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

	return schema.NewSet(resourceAzureRmAutoscaleDefaultHash, []interface{}{email})
}

func flattenAzureRmAutoscaleWebhook(notification insights.AutoscaleNotification) *schema.Set {
	set := &schema.Set{
		F: resourceAzureRmAutoscaleDefaultHash,
	}

	for _, v := range *notification.Webhooks {
		webhook := map[string]interface{}{}
		webhook["service_uri"] = *v.ServiceURI

		if v.Properties != nil {
			properties := map[string]interface{}{}
			for key, value := range v.Properties {
				properties[key] = *value
			}
			webhook["properties"] = properties
		}
		set.Add(webhook)
	}

	return set
}

func resourceAzureRmAutoscaleDefaultHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	var keys []string

	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		buf.WriteString(strings.ToLower(fmt.Sprintf("%s:%v;", k, m[k])))
	}

	return hashcode.String(buf.String())
}

func listTimeZoneNames() []string {
	return []string{
		"Dateline Standard Time", "UTC-11", "Hawaiian Standard Time", "Alaskan Standard Time",
		"Pacific Standard Time (Mexico)", "Pacific Standard Time", "US Mountain Standard Time", "Mountain Standard Time (Mexico)",
		"Mountain Standard Time", "Central America Standard Time", "Central Standard Time", "Central Standard Time (Mexico)",
		"Canada Central Standard Time", "SA Pacific Standard Time", "Eastern Standard Time", "US Eastern Standard Time",
		"Venezuela Standard Time", "Paraguay Standard Time", "Atlantic Standard Time", "Central Brazilian Standard Time",
		"SA Western Standard Time", "Pacific SA Standard Time", "Newfoundland Standard Time", "E. South America Standard Time",
		"Argentina Standard Time", "SA Eastern Standard Time", "Greenland Standard Time", "Montevideo Standard Time",
		"Bahia Standard Time", "UTC-02", "Mid-Atlantic Standard Time", "Azores Standard Time",
		"Cape Verde Standard Time", "Morocco Standard Time", "UTC", "GMT Standard Time",
		"Greenwich Standard Time", "W. Europe Standard Time", "Central Europe Standard Time", "Romance Standard Time",
		"Central European Standard Time", "W. Central Africa Standard Time", "Namibia Standard Time", "Jordan Standard Time",
		"GTB Standard Time", "Middle East Standard Time", "Egypt Standard Time", "Syria Standard Time",
		"E. Europe Standard Time", "South Africa Standard Time", "FLE Standard Time", "Turkey Standard Time",
		"Israel Standard Time", "Kaliningrad Standard Time", "Libya Standard Time", "Arabic Standard Time",
		"Arab Standard Time", "Belarus Standard Time", "Russian Standard Time", "E. Africa Standard Time",
		"Iran Standard Time", "Arabian Standard Time", "Azerbaijan Standard Time", "Russia Time Zone 3",
		"Mauritius Standard Time", "Georgian Standard Time", "Caucasus Standard Time", "Afghanistan Standard Time",
		"West Asia Standard Time", "Ekaterinburg Standard Time", "Pakistan Standard Time", "India Standard Time",
		"Sri Lanka Standard Time", "Nepal Standard Time", "Central Asia Standard Time", "Bangladesh Standard Time",
		"N. Central Asia Standard Time", "Myanmar Standard Time", "SE Asia Standard Time", "North Asia Standard Time",
		"China Standard Time", "North Asia East Standard Time", "Singapore Standard Time", "W. Australia Standard Time",
		"Taipei Standard Time", "Ulaanbaatar Standard Time", "Tokyo Standard Time", "Korea Standard Time",
		"Yakutsk Standard Time", "Cen. Australia Standard Time", "AUS Central Standard Time", "E. Australia Standard Time",
		"AUS Eastern Standard Time", "West Pacific Standard Time", "Tasmania Standard Time", "Magadan Standard Time",
		"Vladivostok Standard Time", "Russia Time Zone 10", "Central Pacific Standard Time", "Russia Time Zone 11",
		"New Zealand Standard Time", "UTC+12", "Fiji Standard Time", "Kamchatka Standard Time",
		"Tonga Standard Time", "Samoa Standard Time", "Line Islands Standard Time"}
}

func rfc3339String() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		_, parseErr := date.ParseTime(time.RFC3339, v)

		if parseErr != nil {
			es = append(es, fmt.Errorf("expected %s to be in RFC3339 format, got %s", k, v))
		}
		return
	}
}

func iso8601DurationString() schema.SchemaValidateFunc {
	return func(i interface{}, k string) (s []string, es []error) {
		v, ok := i.(string)
		if !ok {
			es = append(es, fmt.Errorf("expected type of %s to be string", k))
			return
		}

		matched, _ := regexp.MatchString(`^P(([0-9]+Y)?([0-9]+M)?([0-9]+W)?([0-9]+D)?(T([0-9]+H)?([0-9]+M)?([0-9]+(\.?[0-9]+)?S)?))?$`, v)

		if !matched {
			es = append(es, fmt.Errorf("expected %s to be in ISO 8601 duration format, got %s", k, v))
		}
		return
	}
}
