// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2019-05-05-preview/actionrules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2019-05-05-preview/alertsmanagements"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	scheduleDateLayout     = "01/02/2006"
	scheduleTimeLayout     = "15:04:05"
	scheduleDateTimeLayout = scheduleDateLayout + " " + scheduleTimeLayout
)

var weekDays = []string{
	"Sunday",
	"Monday",
	"Tuesday",
	"Wednesday",
	"Thursday",
	"Friday",
	"Saturday",
}

var weekDayMap = map[string]int{
	"Sunday":    0,
	"Monday":    1,
	"Tuesday":   2,
	"Wednesday": 3,
	"Thursday":  4,
	"Friday":    5,
	"Saturday":  6,
}

func schemaActionRuleConditions() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"alert_context": schemaActionRuleCondition(
					[]string{
						string(actionrules.OperatorEquals),
						string(actionrules.OperatorNotEquals),
						string(actionrules.OperatorContains),
						string(actionrules.OperatorDoesNotContain),
					},
					nil,
				),

				"alert_rule_id": schemaActionRuleCondition(
					[]string{
						string(actionrules.OperatorEquals),
						string(actionrules.OperatorNotEquals),
						string(actionrules.OperatorContains),
						string(actionrules.OperatorDoesNotContain),
					}, nil,
				),

				"description": schemaActionRuleCondition(
					[]string{
						string(actionrules.OperatorEquals),
						string(actionrules.OperatorNotEquals),
						string(actionrules.OperatorContains),
						string(actionrules.OperatorDoesNotContain),
					},
					nil,
				),

				"monitor": schemaActionRuleCondition(
					[]string{
						string(actionrules.OperatorEquals),
						string(actionrules.OperatorNotEquals),
					},
					[]string{
						string(alertsmanagements.MonitorConditionFired),
						string(alertsmanagements.MonitorConditionResolved),
					},
				),

				"monitor_service": schemaActionRuleCondition(
					[]string{
						string(actionrules.OperatorEquals),
						string(actionrules.OperatorNotEquals),
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
						"Custom",
						"Azure Stack Edge",
						"Data Box Gateway",
						"Log Analytics",
						"Platform",
						"Resource Health",
						"Smart Detector",
						"Log Alerts V2",
						"VM Insights - Health",
						"Azure Stack Hub",
						"Health Platform",
					},
				),

				"severity": schemaActionRuleCondition(
					[]string{
						string(actionrules.OperatorEquals),
						string(actionrules.OperatorNotEquals),
					},
					[]string{
						string(actionrules.SeveritySevZero),
						string(actionrules.SeveritySevOne),
						string(actionrules.SeveritySevTwo),
						string(actionrules.SeveritySevThree),
						string(actionrules.SeveritySevFour),
					},
				),

				"target_resource_type": schemaActionRuleCondition(
					[]string{
						string(actionrules.OperatorEquals),
						string(actionrules.OperatorNotEquals),
					},
					nil,
				),
			},
		},
	}
}

func schemaActionRuleCondition(operatorValidateItems, valuesValidateItems []string) *pluginsdk.Schema {
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
					Type:     pluginsdk.TypeSet,
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

func expandActionRuleCondition(input []interface{}) *actionrules.Condition {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &actionrules.Condition{
		Operator: pointer.To(actionrules.Operator(v["operator"].(string))),
		Values:   utils.ExpandStringSlice(v["values"].(*pluginsdk.Set).List()),
	}
}

func expandActionRuleScope(input []interface{}) *actionrules.Scope {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &actionrules.Scope{
		ScopeType: pointer.To(actionrules.ScopeType(v["type"].(string))),
		Values:    utils.ExpandStringSlice(v["resource_ids"].(*pluginsdk.Set).List()),
	}
}

func expandActionRuleConditions(input []interface{}) *actionrules.Conditions {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	return &actionrules.Conditions{
		AlertContext:       expandActionRuleCondition(v["alert_context"].([]interface{})),
		AlertRuleId:        expandActionRuleCondition(v["alert_rule_id"].([]interface{})),
		Description:        expandActionRuleCondition(v["description"].([]interface{})),
		MonitorCondition:   expandActionRuleCondition(v["monitor"].([]interface{})),
		MonitorService:     expandActionRuleCondition(v["monitor_service"].([]interface{})),
		Severity:           expandActionRuleCondition(v["severity"].([]interface{})),
		TargetResourceType: expandActionRuleCondition(v["target_resource_type"].([]interface{})),
	}
}

func flattenActionRuleCondition(input *actionrules.Condition) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var operator string
	if input.Operator != nil {
		operator = string(*input.Operator)
	}
	return []interface{}{
		map[string]interface{}{
			"operator": operator,
			"values":   utils.FlattenStringSlice(input.Values),
		},
	}
}

func flattenActionRuleScope(input *actionrules.Scope) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var scopeType actionrules.ScopeType
	if input.ScopeType != nil {
		scopeType = *input.ScopeType
	}
	return []interface{}{
		map[string]interface{}{
			"type":         scopeType,
			"resource_ids": utils.FlattenStringSlice(input.Values),
		},
	}
}

func flattenActionRuleConditions(input *actionrules.Conditions) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}
	return []interface{}{
		map[string]interface{}{
			"alert_context":        flattenActionRuleCondition(input.AlertContext),
			"alert_rule_id":        flattenActionRuleCondition(input.AlertRuleId),
			"description":          flattenActionRuleCondition(input.Description),
			"monitor":              flattenActionRuleCondition(input.MonitorCondition),
			"monitor_service":      flattenActionRuleCondition(input.MonitorService),
			"severity":             flattenActionRuleCondition(input.Severity),
			"target_resource_type": flattenActionRuleCondition(input.TargetResourceType),
		},
	}
}

func importMonitorActionRule(actionRuleType actionrules.ActionRuleType) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) (data []*pluginsdk.ResourceData, err error) {
		id, err := actionrules.ParseActionRuleID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).Monitor.ActionRulesClient

		actionRule, err := client.GetByName(ctx, *id)
		if err != nil {
			return nil, fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if actionRule.Model == nil {
			return nil, fmt.Errorf("retrieving %s: `properties` was nil", id)
		}

		var t actionrules.ActionRuleType
		switch actionRule.Model.Properties.(type) {
		case actionrules.Suppression:
			t = actionrules.ActionRuleTypeSuppression
		case actionrules.ActionGroup:
			t = actionrules.ActionRuleTypeActionGroup
		case actionrules.Diagnostics:
			t = actionrules.ActionRuleTypeDiagnostics
		}

		if t != actionRuleType {
			return nil, fmt.Errorf("%s has mismatched kind, expected: %q, got %q", id, actionRuleType, t)
		}

		return []*pluginsdk.ResourceData{d}, nil
	}
}
