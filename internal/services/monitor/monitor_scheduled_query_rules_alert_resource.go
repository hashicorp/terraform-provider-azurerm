// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2018-04-16/scheduledqueryrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMonitorScheduledQueryRulesAlert() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorScheduledQueryRulesAlertCreateUpdate,
		Read:   resourceMonitorScheduledQueryRulesAlertRead,
		Update: resourceMonitorScheduledQueryRulesAlertCreateUpdate,
		Delete: resourceMonitorScheduledQueryRulesAlertDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := scheduledqueryrules.ParseScheduledQueryRuleID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.ScheduledQueryRulesAlertUpgradeV0ToV1{},
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
				ValidateFunc: validation.StringDoesNotContainAny("<>*%&:\\?+/"),
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"authorized_resource_ids": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				MaxItems: 100,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
			},
			"action": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"action_group": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: azure.ValidateResourceID,
							},
						},
						"custom_webhook_payload": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsJSON,
						},
						"email_subject": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},
			"data_source_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},
			"auto_mitigation_enabled": {
				Type:          pluginsdk.TypeBool,
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"throttling"},
			},
			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 4096),
			},
			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},
			"frequency": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(5, 1440),
			},
			"query": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"query_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "ResultCount",
				ValidateFunc: validation.StringInSlice([]string{
					"ResultCount",
					"Number",
				}, false),
			},
			"severity": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 4),
			},
			"throttling": {
				Type:          pluginsdk.TypeInt,
				Optional:      true,
				ValidateFunc:  validation.IntBetween(0, 10000),
				ConflictsWith: []string{"auto_mitigation_enabled"},
			},
			"time_window": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(5, 2880),
			},
			"trigger": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"metric_trigger": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"metric_trigger_type": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"Consecutive",
											"Total",
										}, false),
									},
									"operator": {
										Type:     pluginsdk.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"GreaterThan",
											"GreaterThanOrEqual",
											"LessThan",
											"LessThanOrEqual",
											"Equal",
										}, false),
									},
									"threshold": {
										Type:         pluginsdk.TypeFloat,
										Required:     true,
										ValidateFunc: validate.ScheduledQueryRulesAlertThreshold,
									},
									"metric_column": {
										Type:         pluginsdk.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
								},
							},
						},
						"operator": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"GreaterThan",
								"GreaterThanOrEqual",
								"LessThan",
								"LessThanOrEqual",
								"Equal",
							}, false),
						},
						"threshold": {
							Type:         pluginsdk.TypeFloat,
							Required:     true,
							ValidateFunc: validate.ScheduledQueryRulesAlertThreshold,
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceMonitorScheduledQueryRulesAlertCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	action := expandMonitorScheduledQueryRulesAlertingAction(d)
	schedule := expandMonitorScheduledQueryRulesAlertSchedule(d)
	client := meta.(*clients.Client).Monitor.ScheduledQueryRulesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := scheduledqueryrules.NewScheduledQueryRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	frequency := d.Get("frequency").(int)
	timeWindow := d.Get("time_window").(int)
	if timeWindow < frequency {
		return fmt.Errorf("in parameter values for %s: time_window must be greater than or equal to frequency", id)
	}

	query := d.Get("query").(string)
	_, ok := d.GetOk("metric_trigger")
	if ok {
		if !(strings.Contains(query, "summarize") &&
			strings.Contains(query, "AggregatedValue") &&
			strings.Contains(query, "bin")) {
			return fmt.Errorf("in parameter values for %s: query must contain summarize, AggregatedValue, and bin when metric_trigger is specified", id)
		}
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Monitor %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_monitor_scheduled_query_rules_alert", id.ID())
		}
	}

	autoMitigate := d.Get("auto_mitigation_enabled").(bool)
	description := d.Get("description").(string)
	enabledRaw := d.Get("enabled").(bool)

	enabled := scheduledqueryrules.EnabledTrue
	if !enabledRaw {
		enabled = scheduledqueryrules.EnabledFalse
	}

	location := azure.NormalizeLocation(d.Get("location"))

	source := expandMonitorScheduledQueryRulesCommonSource(d)

	t := d.Get("tags").(map[string]interface{})

	parameters := scheduledqueryrules.LogSearchRuleResource{
		Location: location,
		Properties: scheduledqueryrules.LogSearchRule{
			Description:  utils.String(description),
			Enabled:      pointer.To(enabled),
			Source:       source,
			Schedule:     schedule,
			Action:       action,
			AutoMitigate: utils.Bool(autoMitigate),
		},
		Tags: utils.ExpandPtrMapStringString(t),
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating or updating Monitor %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMonitorScheduledQueryRulesAlertRead(d, meta)
}

func resourceMonitorScheduledQueryRulesAlertRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ScheduledQueryRulesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := scheduledqueryrules.ParseScheduledQueryRuleID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("getting Monitor %s: %+v", *id, err)
	}

	d.Set("name", id.ScheduledQueryRuleName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", azure.NormalizeLocation(model.Location))

		props := model.Properties
		d.Set("auto_mitigation_enabled", props.AutoMitigate)
		d.Set("description", props.Description)
		if props.Enabled != nil && *props.Enabled == scheduledqueryrules.EnabledTrue {
			d.Set("enabled", true)
		} else {
			d.Set("enabled", false)
		}

		action, ok := props.Action.(scheduledqueryrules.AlertingAction)
		if !ok {
			return fmt.Errorf("wrong action type in %s: %T", *id, props.Action)
		}
		if err = d.Set("action", flattenAzureRmScheduledQueryRulesAlertAction(action.AznsAction)); err != nil {
			return fmt.Errorf("setting `action`: %+v", err)
		}
		severity, err := strconv.Atoi(string(action.Severity))
		if err != nil {
			return fmt.Errorf("converting action.Severity %q to int in %s: %+v", action.Severity, *id, err)
		}
		d.Set("severity", severity)
		d.Set("throttling", action.ThrottlingInMin)
		if err = d.Set("trigger", flattenAzureRmScheduledQueryRulesAlertTrigger(action.Trigger)); err != nil {
			return fmt.Errorf("setting `trigger`: %+v", err)
		}

		if schedule := props.Schedule; schedule != nil {
			d.Set("frequency", schedule.FrequencyInMinutes)
			d.Set("time_window", schedule.TimeWindowInMinutes)
		}

		d.Set("authorized_resource_ids", utils.FlattenStringSlice(props.Source.AuthorizedResources))
		d.Set("data_source_id", props.Source.DataSourceId)
		d.Set("query", props.Source.Query)
		d.Set("query_type", string(pointer.From(props.Source.QueryType)))

		if err = d.Set("tags", utils.FlattenPtrMapStringString(model.Tags)); err != nil {
			return err
		}
	}

	return nil
}

func resourceMonitorScheduledQueryRulesAlertDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.ScheduledQueryRulesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := scheduledqueryrules.ParseScheduledQueryRuleID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting Monitor %s: %+v", *id, err)
		}
	}

	return nil
}

func expandMonitorScheduledQueryRulesAlertingAction(d *pluginsdk.ResourceData) *scheduledqueryrules.AlertingAction {
	alertActionRaw := d.Get("action").([]interface{})
	alertAction := expandMonitorScheduledQueryRulesAlertAction(alertActionRaw)
	severityRaw := d.Get("severity").(int)
	severity := strconv.Itoa(severityRaw)

	triggerRaw := d.Get("trigger").([]interface{})
	trigger := expandMonitorScheduledQueryRulesAlertTrigger(triggerRaw)

	action := scheduledqueryrules.AlertingAction{
		AznsAction: alertAction,
		Severity:   scheduledqueryrules.AlertSeverity(severity),
		Trigger:    trigger,
	}

	if throttling, ok := d.Get("throttling").(int); ok && throttling != 0 {
		action.ThrottlingInMin = utils.Int64(int64(throttling))
	}

	return &action
}

func expandMonitorScheduledQueryRulesAlertAction(input []interface{}) *scheduledqueryrules.AzNsActionGroup {
	result := scheduledqueryrules.AzNsActionGroup{}

	if len(input) == 0 {
		return &result
	}
	for _, item := range input {
		if item == nil {
			continue
		}

		v, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		actionGroups := v["action_group"].(*pluginsdk.Set).List()
		result.ActionGroup = utils.ExpandStringSlice(actionGroups)
		result.EmailSubject = utils.String(v["email_subject"].(string))
		if v := v["custom_webhook_payload"].(string); v != "" {
			result.CustomWebhookPayload = utils.String(v)
		}
	}

	return &result
}

func expandMonitorScheduledQueryRulesAlertMetricTrigger(input []interface{}) *scheduledqueryrules.LogMetricTrigger {
	if len(input) == 0 {
		return nil
	}

	result := scheduledqueryrules.LogMetricTrigger{}
	for _, item := range input {
		if item == nil {
			continue
		}
		v, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		result.ThresholdOperator = pointer.To(scheduledqueryrules.ConditionalOperator(v["operator"].(string)))
		result.Threshold = utils.Float(v["threshold"].(float64))
		result.MetricTriggerType = pointer.To(scheduledqueryrules.MetricTriggerType(v["metric_trigger_type"].(string)))
		result.MetricColumn = utils.String(v["metric_column"].(string))
	}

	return &result
}

func expandMonitorScheduledQueryRulesAlertSchedule(d *pluginsdk.ResourceData) *scheduledqueryrules.Schedule {
	frequency := d.Get("frequency").(int)
	timeWindow := d.Get("time_window").(int)

	schedule := scheduledqueryrules.Schedule{
		FrequencyInMinutes:  int64(frequency),
		TimeWindowInMinutes: int64(timeWindow),
	}

	return &schedule
}

func expandMonitorScheduledQueryRulesAlertTrigger(input []interface{}) scheduledqueryrules.TriggerCondition {
	result := scheduledqueryrules.TriggerCondition{}
	if len(input) == 0 {
		return result
	}

	for _, item := range input {
		if item == nil {
			continue
		}
		v, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		metricTriggerRaw := v["metric_trigger"].([]interface{})

		result.ThresholdOperator = scheduledqueryrules.ConditionalOperator(v["operator"].(string))
		result.Threshold = v["threshold"].(float64)
		result.MetricTrigger = expandMonitorScheduledQueryRulesAlertMetricTrigger(metricTriggerRaw)
	}

	return result
}
