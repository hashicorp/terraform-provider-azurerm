package monitor

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/sdk/2021-08-08/alertsmanagement"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMonitorAlertProcessingRule() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorAlertProcessingRuleCreate,
		Read:   resourceMonitorAlertProcessingRuleRead,
		Update: resourceMonitorAlertProcessingRuleUpdate,
		Delete: resourceMonitorAlertProcessingRuleDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := alertsmanagement.ParseActionRuleID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ActionRuleName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"action": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice(
								alertsmanagement.PossibleValuesForActionType(), false),
						},
						"add_action_group_ids": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validate.ActionGroupID,
							},
						},
					},
				},
			},

			"scopes": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"condition": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"alert_context": schemaAlertProcessingRuleCondition(
							[]string{
								string(alertsmanagement.OperatorEquals),
								string(alertsmanagement.OperatorNotEquals),
								string(alertsmanagement.OperatorContains),
								string(alertsmanagement.OperatorDoesNotContain),
							}, nil,
						),
						"alert_rule_id": schemaAlertProcessingRuleCondition(
							[]string{
								string(alertsmanagement.OperatorEquals),
								string(alertsmanagement.OperatorNotEquals),
								string(alertsmanagement.OperatorContains),
								string(alertsmanagement.OperatorDoesNotContain),
							}, nil,
						),
						"alert_rule_name": schemaAlertProcessingRuleCondition(
							[]string{
								string(alertsmanagement.OperatorEquals),
								string(alertsmanagement.OperatorNotEquals),
								string(alertsmanagement.OperatorContains),
								string(alertsmanagement.OperatorDoesNotContain),
							}, nil,
						),
						"description": schemaAlertProcessingRuleCondition(
							[]string{
								string(alertsmanagement.OperatorEquals),
								string(alertsmanagement.OperatorNotEquals),
								string(alertsmanagement.OperatorContains),
								string(alertsmanagement.OperatorDoesNotContain),
							},
							nil,
						),
						"monitor_condition": schemaAlertProcessingRuleCondition(
							[]string{
								string(alertsmanagement.OperatorEquals),
								string(alertsmanagement.OperatorNotEquals),
							},
							[]string{
								"Fired",
								"Resolved",
							},
						),
						"monitor_service": schemaAlertProcessingRuleCondition(
							[]string{
								string(alertsmanagement.OperatorEquals),
								string(alertsmanagement.OperatorNotEquals),
							},
							// the supported type list is not consistent with the swagger and sdk
							// https://github.com/Azure/azure-rest-api-specs/issues/9076
							// directly use string constant
							[]string{
								"ActivityLog Administrative",
								"ActivityLog Autoscale",
								"ActivityLog Policy",
								"ActivityLog Recommendation",
								"ActivityLog Security",
								"Application Insights",
								"Azure Backup",
								"Azure Stack Edge",
								"Azure Stack Hub",
								"Custom",
								"Data Box Gateway",
								"Health Platform",
								"Log Alerts V2",
								"Log Analytics",
								"Platform",
								"Prometheus",
								"Resource Health",
								"Smart Detector",
								"VM Insights - Health",
							},
						),
						"severity": schemaAlertProcessingRuleCondition(
							[]string{
								string(alertsmanagement.OperatorEquals),
								string(alertsmanagement.OperatorNotEquals),
							},
							[]string{
								"Sev0",
								"Sev1",
								"Sev2",
								"Sev3",
								"Sev4",
							},
						),
						"signal_type": schemaAlertProcessingRuleCondition(
							[]string{
								string(alertsmanagement.OperatorEquals),
								string(alertsmanagement.OperatorNotEquals),
							},
							[]string{
								"Metric",
								"Log",
								"Unknown",
								"Health",
							},
						),
						"target_resource": schemaAlertProcessingRuleCondition(
							[]string{
								string(alertsmanagement.OperatorEquals),
								string(alertsmanagement.OperatorNotEquals),
								string(alertsmanagement.OperatorContains),
								string(alertsmanagement.OperatorDoesNotContain),
							}, nil,
						),
						"target_resource_group": schemaAlertProcessingRuleCondition(
							[]string{
								string(alertsmanagement.OperatorEquals),
								string(alertsmanagement.OperatorNotEquals),
								string(alertsmanagement.OperatorContains),
								string(alertsmanagement.OperatorDoesNotContain),
							}, nil,
						),
						"target_resource_type": schemaAlertProcessingRuleCondition(
							[]string{
								string(alertsmanagement.OperatorEquals),
								string(alertsmanagement.OperatorNotEquals),
							},
							nil,
						),
					},
				},
			},

			"schedule": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"effective_from": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validate.AlertProcessingRuleScheduleTime(),
						},
						"effective_until": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validate.AlertProcessingRuleScheduleTime(),
						},
						"time_zone": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      "UTC",
							ValidateFunc: validate.AlertProcessingRuleScheduleTimeZone(),
						},
						"recurrence": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"daily": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"start_time": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ValidateFunc: validate.AlertProcessingRuleScheduleDayTime(),
												},
												"end_time": {
													Type:         pluginsdk.TypeString,
													Required:     true,
													ValidateFunc: validate.AlertProcessingRuleScheduleDayTime(),
												},
											},
										},
									},
									"weekly": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"start_time": {
													Type:         pluginsdk.TypeString,
													Optional:     true,
													ValidateFunc: validate.AlertProcessingRuleScheduleDayTime(),
												},
												"end_time": {
													Type:         pluginsdk.TypeString,
													Optional:     true,
													ValidateFunc: validate.AlertProcessingRuleScheduleDayTime(),
												},
												"days_of_week": {
													Type:     pluginsdk.TypeList,
													Required: true,
													MinItems: 1,
													Elem: &pluginsdk.Schema{
														Type:         pluginsdk.TypeString,
														ValidateFunc: validation.IsDayOfTheWeek(false),
													},
												},
											},
										},
									},
									"monthly": {
										Type:     pluginsdk.TypeList,
										Optional: true,
										Elem: &pluginsdk.Resource{
											Schema: map[string]*pluginsdk.Schema{
												"start_time": {
													Type:         pluginsdk.TypeString,
													Optional:     true,
													ValidateFunc: validate.AlertProcessingRuleScheduleDayTime(),
												},
												"end_time": {
													Type:         pluginsdk.TypeString,
													Optional:     true,
													ValidateFunc: validate.AlertProcessingRuleScheduleDayTime(),
												},
												"days_of_month": {
													Type:     pluginsdk.TypeList,
													Required: true,
													MinItems: 1,
													Elem: &pluginsdk.Schema{
														Type:         pluginsdk.TypeInt,
														ValidateFunc: validation.IntBetween(1, 31),
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceMonitorAlertProcessingRuleCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.AlertsManagementClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := alertsmanagement.NewActionRuleID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.AlertProcessingRulesGetByName(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Monitor %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_monitor_alert_processing_rule", id.ID())
		}
	}

	actions, err := expandAlertProcessingRuleActions(d.Get("action").([]interface{}))
	if err != nil {
		return err
	}
	alertProcessingRule := alertsmanagement.AlertProcessingRule{
		// Location support "global" only
		Location: "global",
		Properties: &alertsmanagement.AlertProcessingRuleProperties{
			Actions:     actions,
			Conditions:  expandAlertProcessingRuleConditions(d.Get("condition").([]interface{})),
			Description: utils.String(d.Get("description").(string)),
			Enabled:     utils.Bool(d.Get("enabled").(bool)),
			Schedule:    expandAlertProcessingRuleSchedule(d.Get("schedule").([]interface{})),
			Scopes:      *utils.ExpandStringSlice(d.Get("scopes").([]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.AlertProcessingRulesCreateOrUpdate(ctx, id, alertProcessingRule); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMonitorAlertProcessingRuleRead(d, meta)
}

func resourceMonitorAlertProcessingRuleRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.AlertsManagementClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := alertsmanagement.ParseActionRuleID(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", id.AlertProcessingRuleName)
	d.Set("resource_group_name", id.ResourceGroupName)

	resp, err := client.AlertProcessingRulesGetByName(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Alert Processing Rule %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Monitor %s: %+v", *id, err)
	}

	if model := resp.Model; model != nil {
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}

		if props := model.Properties; props != nil {
			d.Set("description", props.Description)
			d.Set("enabled", props.Enabled)

			if err := d.Set("action", flattenAlertProcessingRuleActions(props.Actions)); err != nil {
				return fmt.Errorf("setting actions: %+v", err)
			}
			if err := d.Set("scopes", utils.FlattenStringSlice(&props.Scopes)); err != nil {
				return fmt.Errorf("setting scope: %+v", err)
			}
			if err := d.Set("condition", flattenAlertProcessingRuleConditions(props.Conditions)); err != nil {
				return fmt.Errorf("setting condition: %+v", err)
			}
			if err := d.Set("schedule", flattenAlertProcessingRuleSchedule(props.Schedule)); err != nil {
				return fmt.Errorf("setting schedule: %+v", err)
			}
		}
	}

	return nil
}

func resourceMonitorAlertProcessingRuleUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.AlertsManagementClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := alertsmanagement.ParseActionRuleID(d.Id())
	if err != nil {
		return err
	}

	actions, err := expandAlertProcessingRuleActions(d.Get("action").([]interface{}))
	if err != nil {
		return err
	}
	alertProcessingRule := alertsmanagement.AlertProcessingRule{
		// Location support "global" only
		Location: "global",
		Properties: &alertsmanagement.AlertProcessingRuleProperties{
			Actions:     actions,
			Conditions:  expandAlertProcessingRuleConditions(d.Get("condition").([]interface{})),
			Description: utils.String(d.Get("description").(string)),
			Enabled:     utils.Bool(d.Get("enabled").(bool)),
			Schedule:    expandAlertProcessingRuleSchedule(d.Get("schedule").([]interface{})),
			Scopes:      *utils.ExpandStringSlice(d.Get("scopes").([]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.AlertProcessingRulesCreateOrUpdate(ctx, *id, alertProcessingRule); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceMonitorAlertProcessingRuleRead(d, meta)
}

func resourceMonitorAlertProcessingRuleDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.AlertsManagementClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := alertsmanagement.ParseActionRuleID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.AlertProcessingRulesDelete(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	return nil
}

func schemaAlertProcessingRuleCondition(operatorValidateItems, valuesValidateItems []string) *pluginsdk.Schema {
	operatorValidateFunc := validation.StringIsNotEmpty
	valuesValidateFunc := validation.StringIsNotEmpty
	if len(operatorValidateItems) > 0 {
		operatorValidateFunc = validation.StringInSlice(operatorValidateItems, false)
	}
	if len(valuesValidateItems) > 0 {
		valuesValidateFunc = validation.StringInSlice(valuesValidateItems, false)
	}

	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"operator": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: operatorValidateFunc,
				},

				"values": {
					Type:     pluginsdk.TypeList,
					Required: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: valuesValidateFunc,
					},
				},
			},
		},
	}
}

func expandAlertProcessingRuleActions(input []interface{}) ([]alertsmanagement.Action, error) {
	if len(input) == 0 {
		return nil, nil
	}

	v := input[0].(map[string]interface{})
	actions := make([]alertsmanagement.Action, 0)
	var action alertsmanagement.Action
	addActionGroupIds := v["add_action_group_ids"].([]interface{})

	switch v["type"].(string) {
	case string(alertsmanagement.ActionTypeAddActionGroups):
		if len(addActionGroupIds) == 0 {
			return nil, fmt.Errorf("add_action_group_ids must be provided with action type \"AddActionGroups\"")
		}
		action = alertsmanagement.AddActionGroups{
			ActionGroupIds: *utils.ExpandStringSlice(v["add_action_group_ids"].([]interface{})),
		}
	case string(alertsmanagement.ActionTypeRemoveAllActionGroups):
		if len(addActionGroupIds) != 0 {
			return nil, fmt.Errorf("add_action_group_ids should not be set with action type \"RemoveAllActionGroups\"")
		}
		action = alertsmanagement.RemoveAllActionGroups{}
	}

	actions = append(actions, action)
	return actions, nil
}

func expandAlertProcessingRuleConditions(input []interface{}) *[]alertsmanagement.Condition {
	if len(input) == 0 {
		return nil
	}

	conditions := make([]alertsmanagement.Condition, 0)
	v := input[0].(map[string]interface{})
	for key, item := range v {
		field := parseField(key)
		if field != nil && item != nil {
			prop := item.([]interface{})
			if len(prop) != 0 {
				props := prop[0].(map[string]interface{})
				operator := alertsmanagement.Operator(props["operator"].(string))
				condition := alertsmanagement.Condition{
					Field:    field,
					Operator: &operator,
					Values:   utils.ExpandStringSlice(props["values"].([]interface{})),
				}
				conditions = append(conditions, condition)
			}
		}
	}

	return &conditions
}

func parseField(input string) *alertsmanagement.Field {
	vals := map[string]alertsmanagement.Field{
		"alert_context":         alertsmanagement.FieldAlertContext,
		"alert_rule_id":         alertsmanagement.FieldAlertRuleId,
		"alert_rule_name":       alertsmanagement.FieldAlertRuleName,
		"description":           alertsmanagement.FieldDescription,
		"monitor_condition":     alertsmanagement.FieldMonitorCondition,
		"monitor_service":       alertsmanagement.FieldMonitorService,
		"severity":              alertsmanagement.FieldSeverity,
		"signal_type":           alertsmanagement.FieldSignalType,
		"target_resource":       alertsmanagement.FieldTargetResource,
		"target_resource_group": alertsmanagement.FieldTargetResourceGroup,
		"target_resource_type":  alertsmanagement.FieldTargetResourceType,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v
	}
	return nil
}

func expandAlertProcessingRuleSchedule(input []interface{}) *alertsmanagement.Schedule {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})
	var effectiveFrom, effectiveUntil *string

	if ef, ok := v["effective_from"]; ok && ef.(string) != "" {
		effectiveFrom = utils.String(ef.(string))
	}

	if eu, ok := v["effective_until"]; ok && eu.(string) != "" {
		effectiveUntil = utils.String(eu.(string))
	}

	schedule := alertsmanagement.Schedule{
		EffectiveFrom:  effectiveFrom,
		EffectiveUntil: effectiveUntil,
		Recurrences:    expandAlertProcessingRuleScheduleRecurrences(v["recurrence"].([]interface{})),
		TimeZone:       utils.String(v["time_zone"].(string)),
	}

	return &schedule
}

func expandAlertProcessingRuleScheduleRecurrences(input []interface{}) *[]alertsmanagement.Recurrence {
	if len(input) == 0 {
		return nil
	}

	recurrences := make([]alertsmanagement.Recurrence, 0, len(input))
	v := input[0].(map[string]interface{})

	for _, item := range v["daily"].([]interface{}) {
		recurrences = append(recurrences, expandAlertProcessingRuleScheduleRecurrence(item, alertsmanagement.RecurrenceTypeDaily))
	}

	for _, item := range v["weekly"].([]interface{}) {
		recurrences = append(recurrences, expandAlertProcessingRuleScheduleRecurrence(item, alertsmanagement.RecurrenceTypeWeekly))
	}

	for _, item := range v["monthly"].([]interface{}) {
		recurrences = append(recurrences, expandAlertProcessingRuleScheduleRecurrence(item, alertsmanagement.RecurrenceTypeMonthly))
	}

	return &recurrences
}

func expandAlertProcessingRuleScheduleRecurrence(input interface{}, recurrenceType alertsmanagement.RecurrenceType) *alertsmanagement.Recurrence {
	if input == nil {
		return nil
	}

	v := input.(map[string]interface{})
	var recurrence alertsmanagement.Recurrence
	var startTime, endTime *string

	if st, ok := v["start_time"]; ok && st.(string) != "" {
		startTime = utils.String(st.(string))
	}

	if et, ok := v["end_time"]; ok && et.(string) != "" {
		endTime = utils.String(et.(string))
	}

	switch recurrenceType {
	case alertsmanagement.RecurrenceTypeDaily:
		recurrence = alertsmanagement.DailyRecurrence{
			StartTime: startTime,
			EndTime:   endTime,
		}
	case alertsmanagement.RecurrenceTypeWeekly:
		recurrence = alertsmanagement.WeeklyRecurrence{
			StartTime:  startTime,
			EndTime:    endTime,
			DaysOfWeek: *expandAlertProcessingRuleScheduleRecurrenceDaysOfWeek(v["days_of_week"].([]interface{})),
		}
	case alertsmanagement.RecurrenceTypeMonthly:
		recurrence = alertsmanagement.MonthlyRecurrence{
			StartTime:   startTime,
			EndTime:     endTime,
			DaysOfMonth: *expandAlertProcessingRuleScheduleRecurrenceDaysOfMonth(v["days_of_month"].([]interface{})),
		}
	}

	return &recurrence
}

func expandAlertProcessingRuleScheduleRecurrenceDaysOfWeek(input []interface{}) *[]alertsmanagement.DaysOfWeek {
	result := make([]alertsmanagement.DaysOfWeek, 0, len(input))
	for _, v := range input {
		result = append(result, alertsmanagement.DaysOfWeek(v.(string)))
	}

	return &result
}

func expandAlertProcessingRuleScheduleRecurrenceDaysOfMonth(input []interface{}) *[]int64 {
	result := make([]int64, len(input))
	for i, item := range input {
		result[i] = int64(item.(int))
	}

	return &result
}

func flattenAlertProcessingRuleActions(input []alertsmanagement.Action) []interface{} {
	if input == nil {
		return nil
	}
	result := make([]interface{}, 0)
	var actionType alertsmanagement.ActionType
	var actionGroupIds []string

	for _, item := range input {
		switch t := item.(type) {
		case alertsmanagement.AddActionGroups:
			actionType = alertsmanagement.ActionTypeAddActionGroups
			actionGroupIds = item.(alertsmanagement.AddActionGroups).ActionGroupIds
		case alertsmanagement.RemoveAllActionGroups:
			actionType = alertsmanagement.ActionTypeRemoveAllActionGroups
		default:
			log.Printf("[WARN] Alert Processing Rule got unsupported action type %v", t)
			continue
		}

		action := map[string]interface{}{
			"type":                 string(actionType),
			"add_action_group_ids": actionGroupIds,
		}
		result = append(result, action)
	}

	return result
}

func flattenAlertProcessingRuleConditions(input *[]alertsmanagement.Condition) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	condition := make(map[string]interface{})
	for _, item := range *input {
		if item.Operator != nil {
			field := flattenAlertProcessingRuleConditionsField(item.Field)
			condition[field] = []interface{}{
				map[string]interface{}{
					"operator": string(*item.Operator),
					"values":   utils.FlattenStringSlice(item.Values),
				},
			}
		}
	}

	return []interface{}{
		condition,
	}
}

func flattenAlertProcessingRuleConditionsField(input *alertsmanagement.Field) string {
	vals := map[alertsmanagement.Field]string{
		alertsmanagement.FieldAlertContext:        "alert_context",
		alertsmanagement.FieldAlertRuleId:         "alert_rule_id",
		alertsmanagement.FieldAlertRuleName:       "alert_rule_name",
		alertsmanagement.FieldDescription:         "description",
		alertsmanagement.FieldMonitorCondition:    "monitor_condition",
		alertsmanagement.FieldMonitorService:      "monitor_service",
		alertsmanagement.FieldSeverity:            "severity",
		alertsmanagement.FieldSignalType:          "signal_type",
		alertsmanagement.FieldTargetResource:      "target_resource",
		alertsmanagement.FieldTargetResourceGroup: "target_resource_group",
		alertsmanagement.FieldTargetResourceType:  "target_resource_type",
	}
	if v, ok := vals[*input]; ok {
		return v
	}

	return string(*input)
}

func flattenAlertProcessingRuleSchedule(input *alertsmanagement.Schedule) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	return []interface{}{
		map[string]interface{}{
			"effective_from":  flattenPtrString(input.EffectiveFrom),
			"effective_until": flattenPtrString(input.EffectiveUntil),
			"recurrence":      flattenAlertProcessingRuleRecurrences(input.Recurrences),
			"time_zone":       flattenPtrString(input.TimeZone),
		},
	}
}

func flattenAlertProcessingRuleRecurrences(input *[]alertsmanagement.Recurrence) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}
	var recurrenceDaily, recurrenceWeekly, recurrenceMonthly []interface{}
	for _, item := range *input {
		switch t := item.(type) {
		case alertsmanagement.DailyRecurrence:
			dailyRecurrence := item.(alertsmanagement.DailyRecurrence)
			recurrence := map[string]interface{}{
				"start_time": flattenPtrString(dailyRecurrence.StartTime),
				"end_time":   flattenPtrString(dailyRecurrence.EndTime),
			}
			recurrenceDaily = append(recurrenceDaily, recurrence)

		case alertsmanagement.WeeklyRecurrence:
			weeklyRecurrence := item.(alertsmanagement.WeeklyRecurrence)
			recurrence := map[string]interface{}{
				"days_of_week": flattenAlertProcessingRuleRecurrenceDaysOfWeek(&weeklyRecurrence.DaysOfWeek),
				"start_time":   flattenPtrString(weeklyRecurrence.StartTime),
				"end_time":     flattenPtrString(weeklyRecurrence.EndTime),
			}
			recurrenceWeekly = append(recurrenceWeekly, recurrence)

		case alertsmanagement.MonthlyRecurrence:
			monthlyRecurrence := item.(alertsmanagement.MonthlyRecurrence)
			recurrence := map[string]interface{}{
				"days_of_month": flattenAlertProcessingRuleRecurrenceDaysOfMonth(&monthlyRecurrence.DaysOfMonth),
				"start_time":    flattenPtrString(monthlyRecurrence.StartTime),
				"end_time":      flattenPtrString(monthlyRecurrence.EndTime),
			}
			recurrenceMonthly = append(recurrenceMonthly, recurrence)

		default:
			log.Printf("[WARN] Alert Processing Rule got unsupported recurrence type %v", t)
		}
	}
	recurrences := []interface{}{
		map[string]interface{}{
			"daily":   recurrenceDaily,
			"weekly":  recurrenceWeekly,
			"monthly": recurrenceMonthly,
		},
	}

	return recurrences
}

func flattenPtrString(input *string) string {
	if input == nil {
		return ""
	}

	return *input
}

func flattenAlertProcessingRuleRecurrenceDaysOfWeek(input *[]alertsmanagement.DaysOfWeek) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make([]interface{}, 0)
	for _, item := range *input {
		result = append(result, string(item))
	}

	return result
}

func flattenAlertProcessingRuleRecurrenceDaysOfMonth(input *[]int64) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make([]interface{}, 0)
	for _, item := range *input {
		result = append(result, item)
	}

	return result
}
