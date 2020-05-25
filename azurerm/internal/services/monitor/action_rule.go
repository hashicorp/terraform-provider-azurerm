package monitor

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/Azure/azure-sdk-for-go/services/preview/alertsmanagement/mgmt/2019-05-05/alertsmanagement"
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

func schemaActionRuleCondition(operatorValidateFunc, valuesValidateFunc schema.SchemaValidateFunc) *schema.Schema {
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

func FlattenInt32Slice(input *[]int32) []interface{} {
	result := make([]interface{}, 0)
	if input != nil {
		for _, item := range *input {
			result = append(result, item)
		}
	}
	return result
}
