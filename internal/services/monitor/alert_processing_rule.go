// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2021-08-08/alertprocessingrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type AlertProcessingRuleConditionModel struct {
	AlertContext        []AlertProcessingRuleSingleConditionModel `tfschema:"alert_context"`
	AlertRuleId         []AlertProcessingRuleSingleConditionModel `tfschema:"alert_rule_id"`
	AlertRuleName       []AlertProcessingRuleSingleConditionModel `tfschema:"alert_rule_name"`
	Description         []AlertProcessingRuleSingleConditionModel `tfschema:"description"`
	MonitorCondition    []AlertProcessingRuleSingleConditionModel `tfschema:"monitor_condition"`
	MonitorService      []AlertProcessingRuleSingleConditionModel `tfschema:"monitor_service"`
	Severity            []AlertProcessingRuleSingleConditionModel `tfschema:"severity"`
	SignalType          []AlertProcessingRuleSingleConditionModel `tfschema:"signal_type"`
	TargetResource      []AlertProcessingRuleSingleConditionModel `tfschema:"target_resource"`
	TargetResourceGroup []AlertProcessingRuleSingleConditionModel `tfschema:"target_resource_group"`
	TargetResourceType  []AlertProcessingRuleSingleConditionModel `tfschema:"target_resource_type"`
}

type AlertProcessingRuleSingleConditionModel struct {
	Operator string   `tfschema:"operator"`
	Values   []string `tfschema:"values"`
}

type AlertProcessingRuleScheduleModel struct {
	EffectiveFrom  string                               `tfschema:"effective_from"`
	EffectiveUntil string                               `tfschema:"effective_until"`
	TimeZone       string                               `tfschema:"time_zone"`
	Recurrence     []AlertProcessingRuleRecurrenceModel `tfschema:"recurrence"`
}

type AlertProcessingRuleRecurrenceModel struct {
	Daily   []AlertProcessingRuleDailyModel   `tfschema:"daily"`
	Weekly  []AlertProcessingRuleWeeklyModel  `tfschema:"weekly"`
	Monthly []AlertProcessingRuleMonthlyModel `tfschema:"monthly"`
}

type AlertProcessingRuleDailyModel struct {
	StartTime string `tfschema:"start_time"`
	EndTime   string `tfschema:"end_time"`
}

type AlertProcessingRuleWeeklyModel struct {
	StartTime  string   `tfschema:"start_time"`
	EndTime    string   `tfschema:"end_time"`
	DaysOfWeek []string `tfschema:"days_of_week"`
}

type AlertProcessingRuleMonthlyModel struct {
	StartTime   string `tfschema:"start_time"`
	EndTime     string `tfschema:"end_time"`
	DaysOfMonth []int  `tfschema:"days_of_month"`
}

func schemaAlertProcessingRule() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ActionRuleName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

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
						alertprocessingrules.PossibleValuesForOperator(), nil,
						[]string{"condition.0.alert_context", "condition.0.alert_rule_id", "condition.0.alert_rule_name",
							"condition.0.description", "condition.0.monitor_condition", "condition.0.monitor_service",
							"condition.0.severity", "condition.0.signal_type", "condition.0.target_resource",
							"condition.0.target_resource_group", "condition.0.target_resource_type"},
					),
					"alert_rule_id": schemaAlertProcessingRuleCondition(
						alertprocessingrules.PossibleValuesForOperator(), nil,
						[]string{"condition.0.alert_context", "condition.0.alert_rule_id", "condition.0.alert_rule_name",
							"condition.0.description", "condition.0.monitor_condition", "condition.0.monitor_service",
							"condition.0.severity", "condition.0.signal_type", "condition.0.target_resource",
							"condition.0.target_resource_group", "condition.0.target_resource_type"},
					),
					"alert_rule_name": schemaAlertProcessingRuleCondition(
						alertprocessingrules.PossibleValuesForOperator(), nil,
						[]string{"condition.0.alert_context", "condition.0.alert_rule_id", "condition.0.alert_rule_name",
							"condition.0.description", "condition.0.monitor_condition", "condition.0.monitor_service",
							"condition.0.severity", "condition.0.signal_type", "condition.0.target_resource",
							"condition.0.target_resource_group", "condition.0.target_resource_type"},
					),
					"description": schemaAlertProcessingRuleCondition(
						alertprocessingrules.PossibleValuesForOperator(), nil,
						[]string{"condition.0.alert_context", "condition.0.alert_rule_id", "condition.0.alert_rule_name",
							"condition.0.description", "condition.0.monitor_condition", "condition.0.monitor_service",
							"condition.0.severity", "condition.0.signal_type", "condition.0.target_resource",
							"condition.0.target_resource_group", "condition.0.target_resource_type"},
					),
					"monitor_condition": schemaAlertProcessingRuleCondition(
						[]string{
							string(alertprocessingrules.OperatorEquals),
							string(alertprocessingrules.OperatorNotEquals),
						},
						[]string{
							"Fired",
							"Resolved",
						},
						[]string{"condition.0.alert_context", "condition.0.alert_rule_id", "condition.0.alert_rule_name",
							"condition.0.description", "condition.0.monitor_condition", "condition.0.monitor_service",
							"condition.0.severity", "condition.0.signal_type", "condition.0.target_resource",
							"condition.0.target_resource_group", "condition.0.target_resource_type"},
					),
					"monitor_service": schemaAlertProcessingRuleCondition(
						[]string{
							string(alertprocessingrules.OperatorEquals),
							string(alertprocessingrules.OperatorNotEquals),
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
						[]string{"condition.0.alert_context", "condition.0.alert_rule_id", "condition.0.alert_rule_name",
							"condition.0.description", "condition.0.monitor_condition", "condition.0.monitor_service",
							"condition.0.severity", "condition.0.signal_type", "condition.0.target_resource",
							"condition.0.target_resource_group", "condition.0.target_resource_type"},
					),
					"severity": schemaAlertProcessingRuleCondition(
						[]string{
							string(alertprocessingrules.OperatorEquals),
							string(alertprocessingrules.OperatorNotEquals),
						},
						[]string{
							"Sev0",
							"Sev1",
							"Sev2",
							"Sev3",
							"Sev4",
						},
						[]string{"condition.0.alert_context", "condition.0.alert_rule_id", "condition.0.alert_rule_name",
							"condition.0.description", "condition.0.monitor_condition", "condition.0.monitor_service",
							"condition.0.severity", "condition.0.signal_type", "condition.0.target_resource",
							"condition.0.target_resource_group", "condition.0.target_resource_type"},
					),
					"signal_type": schemaAlertProcessingRuleCondition(
						[]string{
							string(alertprocessingrules.OperatorEquals),
							string(alertprocessingrules.OperatorNotEquals),
						},
						[]string{
							"Metric",
							"Log",
							"Unknown",
							"Health",
						},
						[]string{"condition.0.alert_context", "condition.0.alert_rule_id", "condition.0.alert_rule_name",
							"condition.0.description", "condition.0.monitor_condition", "condition.0.monitor_service",
							"condition.0.severity", "condition.0.signal_type", "condition.0.target_resource",
							"condition.0.target_resource_group", "condition.0.target_resource_type"},
					),
					"target_resource": schemaAlertProcessingRuleCondition(
						alertprocessingrules.PossibleValuesForOperator(), nil,
						[]string{"condition.0.alert_context", "condition.0.alert_rule_id", "condition.0.alert_rule_name",
							"condition.0.description", "condition.0.monitor_condition", "condition.0.monitor_service",
							"condition.0.severity", "condition.0.signal_type", "condition.0.target_resource",
							"condition.0.target_resource_group", "condition.0.target_resource_type"},
					),
					"target_resource_group": schemaAlertProcessingRuleCondition(
						alertprocessingrules.PossibleValuesForOperator(), nil,
						[]string{"condition.0.alert_context", "condition.0.alert_rule_id", "condition.0.alert_rule_name",
							"condition.0.description", "condition.0.monitor_condition", "condition.0.monitor_service",
							"condition.0.severity", "condition.0.signal_type", "condition.0.target_resource",
							"condition.0.target_resource_group", "condition.0.target_resource_type"},
					),
					"target_resource_type": schemaAlertProcessingRuleCondition(
						[]string{
							string(alertprocessingrules.OperatorEquals),
							string(alertprocessingrules.OperatorNotEquals),
						},
						nil,
						[]string{"condition.0.alert_context", "condition.0.alert_rule_id", "condition.0.alert_rule_name",
							"condition.0.description", "condition.0.monitor_condition", "condition.0.monitor_service",
							"condition.0.severity", "condition.0.signal_type", "condition.0.target_resource",
							"condition.0.target_resource_group", "condition.0.target_resource_type"},
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
	}
}

func schemaAlertProcessingRuleCondition(operatorValidateItems, valuesValidateItems []string, atLeastOneOf []string) *pluginsdk.Schema {
	operatorValidateFunc := validation.StringIsNotEmpty
	valuesValidateFunc := validation.StringIsNotEmpty
	if len(operatorValidateItems) > 0 {
		operatorValidateFunc = validation.StringInSlice(operatorValidateItems, false)
	}
	if len(valuesValidateItems) > 0 {
		valuesValidateFunc = validation.StringInSlice(valuesValidateItems, false)
	}

	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeList,
		Optional:     true,
		MaxItems:     1,
		AtLeastOneOf: atLeastOneOf,
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

func expandAlertProcessingRuleConditions(input []AlertProcessingRuleConditionModel) *[]alertprocessingrules.Condition {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	conditions := make([]alertprocessingrules.Condition, 0)
	expandAlertProcessingRuleSingleConditions(v.AlertContext, alertprocessingrules.FieldAlertContext, &conditions)
	expandAlertProcessingRuleSingleConditions(v.AlertRuleId, alertprocessingrules.FieldAlertRuleId, &conditions)
	expandAlertProcessingRuleSingleConditions(v.AlertRuleName, alertprocessingrules.FieldAlertRuleName, &conditions)
	expandAlertProcessingRuleSingleConditions(v.Description, alertprocessingrules.FieldDescription, &conditions)
	expandAlertProcessingRuleSingleConditions(v.MonitorCondition, alertprocessingrules.FieldMonitorCondition, &conditions)
	expandAlertProcessingRuleSingleConditions(v.MonitorService, alertprocessingrules.FieldMonitorService, &conditions)
	expandAlertProcessingRuleSingleConditions(v.Severity, alertprocessingrules.FieldSeverity, &conditions)
	expandAlertProcessingRuleSingleConditions(v.SignalType, alertprocessingrules.FieldSignalType, &conditions)
	expandAlertProcessingRuleSingleConditions(v.TargetResource, alertprocessingrules.FieldTargetResource, &conditions)
	expandAlertProcessingRuleSingleConditions(v.TargetResourceGroup, alertprocessingrules.FieldTargetResourceGroup, &conditions)
	expandAlertProcessingRuleSingleConditions(v.TargetResourceType, alertprocessingrules.FieldTargetResourceType, &conditions)

	return &conditions
}
func expandAlertProcessingRuleSingleConditions(input []AlertProcessingRuleSingleConditionModel, field alertprocessingrules.Field, conditions *[]alertprocessingrules.Condition) {
	if len(input) == 0 {
		return
	}

	for _, v := range input {
		operator := alertprocessingrules.Operator(v.Operator)
		values := v.Values
		*conditions = append(*conditions, alertprocessingrules.Condition{
			Field:    &field,
			Operator: &operator,
			Values:   &values,
		})
	}
}

func expandAlertProcessingRuleSchedule(input []AlertProcessingRuleScheduleModel) *alertprocessingrules.Schedule {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	var effectiveFrom, effectiveUntil *string

	if v.EffectiveFrom != "" {
		effectiveFrom = utils.String(v.EffectiveFrom)
	}

	if v.EffectiveUntil != "" {
		effectiveUntil = utils.String(v.EffectiveUntil)
	}

	schedule := alertprocessingrules.Schedule{
		EffectiveFrom:  effectiveFrom,
		EffectiveUntil: effectiveUntil,
		Recurrences:    expandAlertProcessingRuleScheduleRecurrences(v.Recurrence),
		TimeZone:       utils.String(v.TimeZone),
	}

	return &schedule
}

func expandAlertProcessingRuleScheduleRecurrences(input []AlertProcessingRuleRecurrenceModel) *[]alertprocessingrules.Recurrence {
	if len(input) == 0 {
		return nil
	}

	recurrences := make([]alertprocessingrules.Recurrence, 0, len(input))
	v := input[0]

	for _, item := range v.Daily {
		var startTime, endTime *string
		if item.StartTime != "" {
			startTime = utils.String(item.StartTime)
		}
		if item.EndTime != "" {
			endTime = utils.String(item.EndTime)
		}

		recurrences = append(recurrences, alertprocessingrules.DailyRecurrence{
			StartTime: startTime,
			EndTime:   endTime,
		})
	}

	for _, item := range v.Weekly {
		var startTime, endTime *string
		if item.StartTime != "" {
			startTime = utils.String(item.StartTime)
		}
		if item.EndTime != "" {
			endTime = utils.String(item.EndTime)
		}

		recurrences = append(recurrences, alertprocessingrules.WeeklyRecurrence{
			StartTime:  startTime,
			EndTime:    endTime,
			DaysOfWeek: *expandAlertProcessingRuleScheduleRecurrenceDaysOfWeek(item.DaysOfWeek),
		})
	}

	for _, item := range v.Monthly {
		var startTime, endTime *string
		if item.StartTime != "" {
			startTime = utils.String(item.StartTime)
		}
		if item.EndTime != "" {
			endTime = utils.String(item.EndTime)
		}

		recurrences = append(recurrences, alertprocessingrules.MonthlyRecurrence{
			StartTime:   startTime,
			EndTime:     endTime,
			DaysOfMonth: *expandAlertProcessingRuleScheduleRecurrenceDaysOfMonth(item.DaysOfMonth),
		})
	}

	return &recurrences
}

func expandAlertProcessingRuleScheduleRecurrenceDaysOfWeek(input []string) *[]alertprocessingrules.DaysOfWeek {
	result := make([]alertprocessingrules.DaysOfWeek, 0, len(input))
	for _, v := range input {
		result = append(result, alertprocessingrules.DaysOfWeek(v))
	}

	return &result
}

func expandAlertProcessingRuleScheduleRecurrenceDaysOfMonth(input []int) *[]int64 {
	result := make([]int64, 0)
	for _, v := range input {
		result = append(result, int64(v))
	}

	return &result
}

func flattenAlertProcessingRuleAddActionGroupId(input []alertprocessingrules.Action) ([]string, error) {
	if len(input) != 1 {
		return make([]string, 0), fmt.Errorf("read add_action_group_ids, the action should contains 1 element, but get %d element", len(input))
	}

	if addActionGroups, ok := input[0].(alertprocessingrules.AddActionGroups); ok {
		return addActionGroups.ActionGroupIds, nil
	}

	return make([]string, 0), fmt.Errorf("read add_action_group_ids, get unsupported action type %v", input[0])
}

func flattenAlertProcessingRuleConditions(input *[]alertprocessingrules.Condition) []AlertProcessingRuleConditionModel {
	if input == nil {
		return make([]AlertProcessingRuleConditionModel, 0)
	}

	condition := AlertProcessingRuleConditionModel{}
	for _, item := range *input {
		if item.Field != nil && item.Operator != nil && item.Values != nil {
			switch *item.Field {
			case alertprocessingrules.FieldAlertContext:
				condition.AlertContext = []AlertProcessingRuleSingleConditionModel{{string(*item.Operator), *item.Values}}
			case alertprocessingrules.FieldAlertRuleId:
				condition.AlertRuleId = []AlertProcessingRuleSingleConditionModel{{string(*item.Operator), *item.Values}}
			case alertprocessingrules.FieldAlertRuleName:
				condition.AlertRuleName = []AlertProcessingRuleSingleConditionModel{{string(*item.Operator), *item.Values}}
			case alertprocessingrules.FieldDescription:
				condition.Description = []AlertProcessingRuleSingleConditionModel{{string(*item.Operator), *item.Values}}
			case alertprocessingrules.FieldMonitorCondition:
				condition.MonitorCondition = []AlertProcessingRuleSingleConditionModel{{string(*item.Operator), *item.Values}}
			case alertprocessingrules.FieldMonitorService:
				condition.MonitorService = []AlertProcessingRuleSingleConditionModel{{string(*item.Operator), *item.Values}}
			case alertprocessingrules.FieldSeverity:
				condition.Severity = []AlertProcessingRuleSingleConditionModel{{string(*item.Operator), *item.Values}}
			case alertprocessingrules.FieldSignalType:
				condition.SignalType = []AlertProcessingRuleSingleConditionModel{{string(*item.Operator), *item.Values}}
			case alertprocessingrules.FieldTargetResource:
				condition.TargetResource = []AlertProcessingRuleSingleConditionModel{{string(*item.Operator), *item.Values}}
			case alertprocessingrules.FieldTargetResourceGroup:
				condition.TargetResourceGroup = []AlertProcessingRuleSingleConditionModel{{string(*item.Operator), *item.Values}}
			case alertprocessingrules.FieldTargetResourceType:
				condition.TargetResourceType = []AlertProcessingRuleSingleConditionModel{{string(*item.Operator), *item.Values}}
			}
		}
	}

	return []AlertProcessingRuleConditionModel{condition}
}

func flattenAlertProcessingRuleSchedule(input *alertprocessingrules.Schedule) []AlertProcessingRuleScheduleModel {
	if input == nil {
		return make([]AlertProcessingRuleScheduleModel, 0)
	}

	return []AlertProcessingRuleScheduleModel{{
		EffectiveFrom:  flattenPtrString(input.EffectiveFrom),
		EffectiveUntil: flattenPtrString(input.EffectiveUntil),
		Recurrence:     flattenAlertProcessingRuleRecurrences(input.Recurrences),
		TimeZone:       flattenPtrString(input.TimeZone),
	}}
}

func flattenAlertProcessingRuleRecurrences(input *[]alertprocessingrules.Recurrence) []AlertProcessingRuleRecurrenceModel {
	if input == nil {
		return make([]AlertProcessingRuleRecurrenceModel, 0)
	}

	recurrence := AlertProcessingRuleRecurrenceModel{}
	for _, item := range *input {
		switch t := item.(type) {
		case alertprocessingrules.DailyRecurrence:
			dailyRecurrence := item.(alertprocessingrules.DailyRecurrence)
			daily := AlertProcessingRuleDailyModel{
				StartTime: flattenPtrString(dailyRecurrence.StartTime),
				EndTime:   flattenPtrString(dailyRecurrence.EndTime),
			}
			recurrence.Daily = append(recurrence.Daily, daily)

		case alertprocessingrules.WeeklyRecurrence:
			weeklyRecurrence := item.(alertprocessingrules.WeeklyRecurrence)
			weekly := AlertProcessingRuleWeeklyModel{
				DaysOfWeek: flattenAlertProcessingRuleRecurrenceDaysOfWeek(&weeklyRecurrence.DaysOfWeek),
				StartTime:  flattenPtrString(weeklyRecurrence.StartTime),
				EndTime:    flattenPtrString(weeklyRecurrence.EndTime),
			}
			recurrence.Weekly = append(recurrence.Weekly, weekly)

		case alertprocessingrules.MonthlyRecurrence:
			monthlyRecurrence := item.(alertprocessingrules.MonthlyRecurrence)
			monthly := AlertProcessingRuleMonthlyModel{
				DaysOfMonth: flattenAlertProcessingRuleRecurrenceDaysOfMonth(&monthlyRecurrence.DaysOfMonth),
				StartTime:   flattenPtrString(monthlyRecurrence.StartTime),
				EndTime:     flattenPtrString(monthlyRecurrence.EndTime),
			}
			recurrence.Monthly = append(recurrence.Monthly, monthly)

		default:
			log.Printf("[WARN] Alert Processing Rule got unsupported recurrence type %v", t)
		}
	}

	return []AlertProcessingRuleRecurrenceModel{recurrence}
}

func flattenPtrString(input *string) string {
	if input == nil {
		return ""
	}

	return *input
}

func flattenAlertProcessingRuleRecurrenceDaysOfWeek(input *[]alertprocessingrules.DaysOfWeek) []string {
	if input == nil {
		return make([]string, 0)
	}

	result := make([]string, 0)
	for _, item := range *input {
		result = append(result, string(item))
	}

	return result
}

func flattenAlertProcessingRuleRecurrenceDaysOfMonth(input *[]int64) []int {
	if input == nil {
		return make([]int, 0)
	}
	result := make([]int, 0)
	for _, v := range *input {
		result = append(result, int(v))
	}

	return result
}
