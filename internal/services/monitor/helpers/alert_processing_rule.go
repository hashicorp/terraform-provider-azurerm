// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2021-08-08/alertprocessingrules"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
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
	StartTime   string  `tfschema:"start_time"`
	EndTime     string  `tfschema:"end_time"`
	DaysOfMonth []int64 `tfschema:"days_of_month"`
}

func SchemaAlertProcessingRule() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.AlertProcessingRuleName,
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
					"alert_context": SchemaAlertProcessingRuleCondition(
						alertprocessingrules.PossibleValuesForOperator(), nil,
						[]string{
							"condition.0.alert_context", "condition.0.alert_rule_id", "condition.0.alert_rule_name",
							"condition.0.description", "condition.0.monitor_condition", "condition.0.monitor_service",
							"condition.0.severity", "condition.0.signal_type", "condition.0.target_resource",
							"condition.0.target_resource_group", "condition.0.target_resource_type",
						},
					),
					"alert_rule_id": SchemaAlertProcessingRuleCondition(
						alertprocessingrules.PossibleValuesForOperator(), nil,
						[]string{
							"condition.0.alert_context", "condition.0.alert_rule_id", "condition.0.alert_rule_name",
							"condition.0.description", "condition.0.monitor_condition", "condition.0.monitor_service",
							"condition.0.severity", "condition.0.signal_type", "condition.0.target_resource",
							"condition.0.target_resource_group", "condition.0.target_resource_type",
						},
					),
					"alert_rule_name": SchemaAlertProcessingRuleCondition(
						alertprocessingrules.PossibleValuesForOperator(), nil,
						[]string{
							"condition.0.alert_context", "condition.0.alert_rule_id", "condition.0.alert_rule_name",
							"condition.0.description", "condition.0.monitor_condition", "condition.0.monitor_service",
							"condition.0.severity", "condition.0.signal_type", "condition.0.target_resource",
							"condition.0.target_resource_group", "condition.0.target_resource_type",
						},
					),
					"description": SchemaAlertProcessingRuleCondition(
						alertprocessingrules.PossibleValuesForOperator(), nil,
						[]string{
							"condition.0.alert_context", "condition.0.alert_rule_id", "condition.0.alert_rule_name",
							"condition.0.description", "condition.0.monitor_condition", "condition.0.monitor_service",
							"condition.0.severity", "condition.0.signal_type", "condition.0.target_resource",
							"condition.0.target_resource_group", "condition.0.target_resource_type",
						},
					),
					"monitor_condition": SchemaAlertProcessingRuleCondition(
						[]string{
							string(alertprocessingrules.OperatorEquals),
							string(alertprocessingrules.OperatorNotEquals),
						},
						[]string{
							"Fired",
							"Resolved",
						},
						[]string{
							"condition.0.alert_context", "condition.0.alert_rule_id", "condition.0.alert_rule_name",
							"condition.0.description", "condition.0.monitor_condition", "condition.0.monitor_service",
							"condition.0.severity", "condition.0.signal_type", "condition.0.target_resource",
							"condition.0.target_resource_group", "condition.0.target_resource_type",
						},
					),
					"monitor_service": SchemaAlertProcessingRuleCondition(
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
						[]string{
							"condition.0.alert_context", "condition.0.alert_rule_id", "condition.0.alert_rule_name",
							"condition.0.description", "condition.0.monitor_condition", "condition.0.monitor_service",
							"condition.0.severity", "condition.0.signal_type", "condition.0.target_resource",
							"condition.0.target_resource_group", "condition.0.target_resource_type",
						},
					),
					"severity": SchemaAlertProcessingRuleCondition(
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
						[]string{
							"condition.0.alert_context", "condition.0.alert_rule_id", "condition.0.alert_rule_name",
							"condition.0.description", "condition.0.monitor_condition", "condition.0.monitor_service",
							"condition.0.severity", "condition.0.signal_type", "condition.0.target_resource",
							"condition.0.target_resource_group", "condition.0.target_resource_type",
						},
					),
					"signal_type": SchemaAlertProcessingRuleCondition(
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
						[]string{
							"condition.0.alert_context", "condition.0.alert_rule_id", "condition.0.alert_rule_name",
							"condition.0.description", "condition.0.monitor_condition", "condition.0.monitor_service",
							"condition.0.severity", "condition.0.signal_type", "condition.0.target_resource",
							"condition.0.target_resource_group", "condition.0.target_resource_type",
						},
					),
					"target_resource": SchemaAlertProcessingRuleCondition(
						alertprocessingrules.PossibleValuesForOperator(), nil,
						[]string{
							"condition.0.alert_context", "condition.0.alert_rule_id", "condition.0.alert_rule_name",
							"condition.0.description", "condition.0.monitor_condition", "condition.0.monitor_service",
							"condition.0.severity", "condition.0.signal_type", "condition.0.target_resource",
							"condition.0.target_resource_group", "condition.0.target_resource_type",
						},
					),
					"target_resource_group": SchemaAlertProcessingRuleCondition(
						alertprocessingrules.PossibleValuesForOperator(), nil,
						[]string{
							"condition.0.alert_context", "condition.0.alert_rule_id", "condition.0.alert_rule_name",
							"condition.0.description", "condition.0.monitor_condition", "condition.0.monitor_service",
							"condition.0.severity", "condition.0.signal_type", "condition.0.target_resource",
							"condition.0.target_resource_group", "condition.0.target_resource_type",
						},
					),
					"target_resource_type": SchemaAlertProcessingRuleCondition(
						[]string{
							string(alertprocessingrules.OperatorEquals),
							string(alertprocessingrules.OperatorNotEquals),
						},
						nil,
						[]string{
							"condition.0.alert_context", "condition.0.alert_rule_id", "condition.0.alert_rule_name",
							"condition.0.description", "condition.0.monitor_condition", "condition.0.monitor_service",
							"condition.0.severity", "condition.0.signal_type", "condition.0.target_resource",
							"condition.0.target_resource_group", "condition.0.target_resource_type",
						},
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

func SchemaAlertProcessingRuleCondition(operatorValidateItems, valuesValidateItems []string, atLeastOneOf []string) *pluginsdk.Schema {
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

func ExpandAlertProcessingRuleConditions(input []AlertProcessingRuleConditionModel) *[]alertprocessingrules.Condition {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	conditions := make([]alertprocessingrules.Condition, 0)
	ExpandAlertProcessingRuleSingleConditions(v.AlertContext, alertprocessingrules.FieldAlertContext, &conditions)
	ExpandAlertProcessingRuleSingleConditions(v.AlertRuleId, alertprocessingrules.FieldAlertRuleId, &conditions)
	ExpandAlertProcessingRuleSingleConditions(v.AlertRuleName, alertprocessingrules.FieldAlertRuleName, &conditions)
	ExpandAlertProcessingRuleSingleConditions(v.Description, alertprocessingrules.FieldDescription, &conditions)
	ExpandAlertProcessingRuleSingleConditions(v.MonitorCondition, alertprocessingrules.FieldMonitorCondition, &conditions)
	ExpandAlertProcessingRuleSingleConditions(v.MonitorService, alertprocessingrules.FieldMonitorService, &conditions)
	ExpandAlertProcessingRuleSingleConditions(v.Severity, alertprocessingrules.FieldSeverity, &conditions)
	ExpandAlertProcessingRuleSingleConditions(v.SignalType, alertprocessingrules.FieldSignalType, &conditions)
	ExpandAlertProcessingRuleSingleConditions(v.TargetResource, alertprocessingrules.FieldTargetResource, &conditions)
	ExpandAlertProcessingRuleSingleConditions(v.TargetResourceGroup, alertprocessingrules.FieldTargetResourceGroup, &conditions)
	ExpandAlertProcessingRuleSingleConditions(v.TargetResourceType, alertprocessingrules.FieldTargetResourceType, &conditions)

	return &conditions
}

func ExpandAlertProcessingRuleSingleConditions(input []AlertProcessingRuleSingleConditionModel, field alertprocessingrules.Field, conditions *[]alertprocessingrules.Condition) {
	if len(input) == 0 {
		return
	}

	for _, v := range input {
		*conditions = append(*conditions, alertprocessingrules.Condition{
			Field:    &field,
			Operator: pointer.ToEnum[alertprocessingrules.Operator](v.Operator),
			Values:   pointer.To(v.Values),
		})
	}
}

func ExpandAlertProcessingRuleSchedule(input []AlertProcessingRuleScheduleModel) *alertprocessingrules.Schedule {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	var effectiveFrom, effectiveUntil *string

	if v.EffectiveFrom != "" {
		effectiveFrom = pointer.To(v.EffectiveFrom)
	}

	if v.EffectiveUntil != "" {
		effectiveUntil = pointer.To(v.EffectiveUntil)
	}

	schedule := alertprocessingrules.Schedule{
		EffectiveFrom:  effectiveFrom,
		EffectiveUntil: effectiveUntil,
		Recurrences:    ExpandAlertProcessingRuleScheduleRecurrences(v.Recurrence),
		TimeZone:       pointer.To(v.TimeZone),
	}

	return &schedule
}

func ExpandAlertProcessingRuleScheduleRecurrences(input []AlertProcessingRuleRecurrenceModel) *[]alertprocessingrules.Recurrence {
	if len(input) == 0 {
		return nil
	}

	recurrences := make([]alertprocessingrules.Recurrence, 0, len(input))
	v := input[0]

	for _, item := range v.Daily {
		var startTime, endTime *string
		if item.StartTime != "" {
			startTime = pointer.To(item.StartTime)
		}
		if item.EndTime != "" {
			endTime = pointer.To(item.EndTime)
		}

		recurrences = append(recurrences, alertprocessingrules.DailyRecurrence{
			StartTime: startTime,
			EndTime:   endTime,
		})
	}

	for _, item := range v.Weekly {
		var startTime, endTime *string
		if item.StartTime != "" {
			startTime = pointer.To(item.StartTime)
		}
		if item.EndTime != "" {
			endTime = pointer.To(item.EndTime)
		}

		recurrences = append(recurrences, alertprocessingrules.WeeklyRecurrence{
			StartTime:  startTime,
			EndTime:    endTime,
			DaysOfWeek: *ExpandAlertProcessingRuleScheduleRecurrenceDaysOfWeek(item.DaysOfWeek),
		})
	}

	for _, item := range v.Monthly {
		var startTime, endTime *string
		if item.StartTime != "" {
			startTime = pointer.To(item.StartTime)
		}
		if item.EndTime != "" {
			endTime = pointer.To(item.EndTime)
		}

		recurrences = append(recurrences, alertprocessingrules.MonthlyRecurrence{
			StartTime:   startTime,
			EndTime:     endTime,
			DaysOfMonth: item.DaysOfMonth,
		})
	}

	return &recurrences
}

func ExpandAlertProcessingRuleScheduleRecurrenceDaysOfWeek(input []string) *[]alertprocessingrules.DaysOfWeek {
	result := make([]alertprocessingrules.DaysOfWeek, 0, len(input))
	for _, v := range input {
		result = append(result, alertprocessingrules.DaysOfWeek(v))
	}

	return &result
}

func FlattenAlertProcessingRuleConditions(input *[]alertprocessingrules.Condition) []AlertProcessingRuleConditionModel {
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

func FlattenAlertProcessingRuleSchedule(input *alertprocessingrules.Schedule) []AlertProcessingRuleScheduleModel {
	if input == nil {
		return make([]AlertProcessingRuleScheduleModel, 0)
	}

	return []AlertProcessingRuleScheduleModel{{
		EffectiveFrom:  FlattenPtrString(input.EffectiveFrom),
		EffectiveUntil: FlattenPtrString(input.EffectiveUntil),
		Recurrence:     FlattenAlertProcessingRuleRecurrences(input.Recurrences),
		TimeZone:       FlattenPtrString(input.TimeZone),
	}}
}

func FlattenAlertProcessingRuleRecurrences(input *[]alertprocessingrules.Recurrence) []AlertProcessingRuleRecurrenceModel {
	if input == nil {
		return make([]AlertProcessingRuleRecurrenceModel, 0)
	}

	recurrence := AlertProcessingRuleRecurrenceModel{}
	for _, item := range *input {
		switch t := item.(type) {
		case alertprocessingrules.DailyRecurrence:
			dailyRecurrence := item.(alertprocessingrules.DailyRecurrence)
			daily := AlertProcessingRuleDailyModel{
				StartTime: FlattenPtrString(dailyRecurrence.StartTime),
				EndTime:   FlattenPtrString(dailyRecurrence.EndTime),
			}
			recurrence.Daily = append(recurrence.Daily, daily)

		case alertprocessingrules.WeeklyRecurrence:
			weeklyRecurrence := item.(alertprocessingrules.WeeklyRecurrence)
			weekly := AlertProcessingRuleWeeklyModel{
				DaysOfWeek: FlattenAlertProcessingRuleRecurrenceDaysOfWeek(&weeklyRecurrence.DaysOfWeek),
				StartTime:  FlattenPtrString(weeklyRecurrence.StartTime),
				EndTime:    FlattenPtrString(weeklyRecurrence.EndTime),
			}
			recurrence.Weekly = append(recurrence.Weekly, weekly)

		case alertprocessingrules.MonthlyRecurrence:
			monthlyRecurrence := item.(alertprocessingrules.MonthlyRecurrence)
			monthly := AlertProcessingRuleMonthlyModel{
				DaysOfMonth: monthlyRecurrence.DaysOfMonth,
				StartTime:   FlattenPtrString(monthlyRecurrence.StartTime),
				EndTime:     FlattenPtrString(monthlyRecurrence.EndTime),
			}
			recurrence.Monthly = append(recurrence.Monthly, monthly)

		default:
			log.Printf("[WARN] Alert Processing Rule got unsupported recurrence type %v", t)
		}
	}

	return []AlertProcessingRuleRecurrenceModel{recurrence}
}

func FlattenPtrString(input *string) string {
	if input == nil {
		return ""
	}

	return *input
}

func FlattenAlertProcessingRuleRecurrenceDaysOfWeek(input *[]alertprocessingrules.DaysOfWeek) []string {
	if input == nil {
		return make([]string, 0)
	}

	result := make([]string, 0)
	for _, item := range *input {
		result = append(result, string(item))
	}

	return result
}
