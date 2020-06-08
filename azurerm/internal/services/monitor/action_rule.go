package monitor

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/alertsmanagement/mgmt/2019-05-05/alertsmanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/monitor/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
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

func schemaActionRuleConditions() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"alert_context": schemaActionRuleCondition(
					[]string{
						string(alertsmanagement.Equals),
						string(alertsmanagement.NotEquals),
						string(alertsmanagement.Contains),
						string(alertsmanagement.DoesNotContain),
					},
					nil,
				),

				"alert_rule_id": schemaActionRuleCondition(
					[]string{
						string(alertsmanagement.Equals),
						string(alertsmanagement.NotEquals),
						string(alertsmanagement.Contains),
						string(alertsmanagement.DoesNotContain),
					}, nil,
				),

				"description": schemaActionRuleCondition(
					[]string{
						string(alertsmanagement.Equals),
						string(alertsmanagement.NotEquals),
						string(alertsmanagement.Contains),
						string(alertsmanagement.DoesNotContain),
					},
					nil,
				),

				"monitor": schemaActionRuleCondition(
					[]string{
						string(alertsmanagement.Equals),
						string(alertsmanagement.NotEquals),
					},
					[]string{
						string(alertsmanagement.Fired),
						string(alertsmanagement.Resolved),
					},
				),

				"monitor_service": schemaActionRuleCondition(
					[]string{
						string(alertsmanagement.Equals),
						string(alertsmanagement.NotEquals),
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
						"Data Box Edge",
						"Data Box Gateway",
						"Health Platform",
						"Log Analytics",
						"Platform",
						"Resource Health",
					},
				),

				"severity": schemaActionRuleCondition(
					[]string{
						string(alertsmanagement.Equals),
						string(alertsmanagement.NotEquals),
					},
					[]string{
						string(alertsmanagement.Sev0),
						string(alertsmanagement.Sev1),
						string(alertsmanagement.Sev2),
						string(alertsmanagement.Sev3),
						string(alertsmanagement.Sev4),
					},
				),

				"target_resource_type": schemaActionRuleCondition(
					[]string{
						string(alertsmanagement.Equals),
						string(alertsmanagement.NotEquals),
					},
					nil,
				),
			},
		},
	}
}

func schemaActionRuleCondition(operatorValidateItems, valuesValidateItems []string) *schema.Schema {
	operatorValidateFunc := validation.StringIsNotEmpty
	valuesValidateFunc := validation.StringIsNotEmpty
	if len(operatorValidateItems) > 0 {
		operatorValidateFunc = validation.StringInSlice(operatorValidateItems, false)
	}
	if len(valuesValidateItems) > 0 {
		valuesValidateFunc = validation.StringInSlice(valuesValidateItems, false)
	}

	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"operator": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: operatorValidateFunc,
				},

				"values": {
					Type:     schema.TypeSet,
					Required: true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: valuesValidateFunc,
					},
				},
			},
		},
	}
}

func expandArmActionRuleCondition(input []interface{}) *alertsmanagement.Condition {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &alertsmanagement.Condition{
		Operator: alertsmanagement.Operator(v["operator"].(string)),
		Values:   utils.ExpandStringSlice(v["values"].(*schema.Set).List()),
	}
}

func expandArmActionRuleScope(input []interface{}) *alertsmanagement.Scope {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &alertsmanagement.Scope{
		ScopeType: alertsmanagement.ScopeType(v["type"].(string)),
		Values:    utils.ExpandStringSlice(v["resource_ids"].(*schema.Set).List()),
	}
}

func expandArmActionRuleConditions(input []interface{}) *alertsmanagement.Conditions {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	return &alertsmanagement.Conditions{
		AlertContext:       expandArmActionRuleCondition(v["alert_context"].([]interface{})),
		AlertRuleID:        expandArmActionRuleCondition(v["alert_rule_id"].([]interface{})),
		Description:        expandArmActionRuleCondition(v["description"].([]interface{})),
		MonitorCondition:   expandArmActionRuleCondition(v["monitor"].([]interface{})),
		MonitorService:     expandArmActionRuleCondition(v["monitor_service"].([]interface{})),
		Severity:           expandArmActionRuleCondition(v["severity"].([]interface{})),
		TargetResourceType: expandArmActionRuleCondition(v["target_resource_type"].([]interface{})),
	}
}

func flattenArmActionRuleCondition(input *alertsmanagement.Condition) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var operator string
	if input.Operator != "" {
		operator = string(input.Operator)
	}
	return []interface{}{
		map[string]interface{}{
			"operator": operator,
			"values":   utils.FlattenStringSlice(input.Values),
		},
	}
}

func flattenArmActionRuleScope(input *alertsmanagement.Scope) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var scopeType alertsmanagement.ScopeType
	if input.ScopeType != "" {
		scopeType = input.ScopeType
	}
	return []interface{}{
		map[string]interface{}{
			"type":         scopeType,
			"resource_ids": utils.FlattenStringSlice(input.Values),
		},
	}
}

func flattenArmActionRuleConditions(input *alertsmanagement.Conditions) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}
	return []interface{}{
		map[string]interface{}{
			"alert_context":        flattenArmActionRuleCondition(input.AlertContext),
			"alert_rule_id":        flattenArmActionRuleCondition(input.AlertRuleID),
			"description":          flattenArmActionRuleCondition(input.Description),
			"monitor":              flattenArmActionRuleCondition(input.MonitorCondition),
			"monitor_service":      flattenArmActionRuleCondition(input.MonitorService),
			"severity":             flattenArmActionRuleCondition(input.Severity),
			"target_resource_type": flattenArmActionRuleCondition(input.TargetResourceType),
		},
	}
}

func importMonitorActionRule(actionRuleType alertsmanagement.Type) func(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
	return func(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
		id, err := parse.ActionRuleID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).Monitor.ActionRulesClient
		ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
		defer cancel()

		actionRule, err := client.GetByName(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return nil, fmt.Errorf("retrieving Monitor Action Rule %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		if actionRule.Properties == nil {
			return nil, fmt.Errorf("retrieving Monitor Action Rule %q (Resource Group %q): `properties` was nil", id.Name, id.ResourceGroup)
		}

		var t alertsmanagement.Type
		switch actionRule.Properties.(type) {
		case alertsmanagement.Suppression:
			t = alertsmanagement.TypeSuppression
		case alertsmanagement.ActionGroup:
			t = alertsmanagement.TypeActionGroup
		case alertsmanagement.Diagnostics:
			t = alertsmanagement.TypeDiagnostics
		}

		if t != actionRuleType {
			return nil, fmt.Errorf("Monitor Action Rule has mismatched kind, expected: %q, got %q", actionRuleType, t)
		}

		return []*schema.ResourceData{d}, nil
	}
}
