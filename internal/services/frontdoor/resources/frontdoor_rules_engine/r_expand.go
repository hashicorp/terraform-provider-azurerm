// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package frontdoor_rules_engine

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-05-01/frontdoors"
)

func expandFrontDoorRulesEngineAction(input []any) frontdoors.RulesEngineAction {
	if len(input) == 0 || input[0] == nil {
		return frontdoors.RulesEngineAction{}
	}

	ruleAction := input[0].(map[string]any)

	requestHeaderActions := ruleAction["request_header"].([]any)
	responseHeaderActions := ruleAction["response_header"].([]any)

	return frontdoors.RulesEngineAction{
		RequestHeaderActions:  expandHeaderAction(requestHeaderActions),
		ResponseHeaderActions: expandHeaderAction(responseHeaderActions),
	}
}

func expandHeaderAction(input []any) *[]frontdoors.HeaderAction {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	output := make([]frontdoors.HeaderAction, 0)

	for _, a := range input {
		action := a.(map[string]any)

		headerName := action["header_name"].(string)
		value := action["value"].(string)
		headerActionType := action["header_action_type"].(string)

		frontdoorRulesEngineRuleHeaderAction := frontdoors.HeaderAction{
			HeaderName:       headerName,
			Value:            pointer.To(value),
			HeaderActionType: frontdoors.HeaderActionType(headerActionType),
		}

		output = append(output, frontdoorRulesEngineRuleHeaderAction)
	}

	return &output
}

func expandFrontDoorRulesEngineRules(input []any) *[]frontdoors.RulesEngineRule {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	output := make([]frontdoors.RulesEngineRule, 0)

	for _, r := range input {
		rule := r.(map[string]any)

		ruleName := rule["name"].(string)
		priority := int64(rule["priority"].(int))
		actions := rule["action"].([]any)
		matchConditions := rule["match_condition"].([]any)

		frontdoorRulesEngineRule := frontdoors.RulesEngineRule{
			Name:            ruleName,
			Priority:        priority,
			Action:          expandFrontDoorRulesEngineAction(actions),
			MatchConditions: expandFrontDoorRulesEngineMatchCondition(matchConditions),
		}

		output = append(output, frontdoorRulesEngineRule)
	}
	return &output
}

func expandFrontDoorRulesEngineMatchCondition(input []any) *[]frontdoors.RulesEngineMatchCondition {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	output := make([]frontdoors.RulesEngineMatchCondition, 0)

	for _, c := range input {
		condition := c.(map[string]any)

		selector := condition["selector"].(string)
		matchVariable := condition["variable"].(string)
		operator := condition["operator"].(string)
		transform := condition["transform"].([]any)
		matchValue := condition["value"].([]any)

		matchValueArray := make([]string, 0)
		for _, v := range matchValue {
			matchValueArray = append(matchValueArray, v.(string))
		}

		matchCondition := frontdoors.RulesEngineMatchCondition{
			RulesEngineMatchVariable: frontdoors.RulesEngineMatchVariable(matchVariable),
			Selector:                 pointer.To(selector),
			RulesEngineOperator:      frontdoors.RulesEngineOperator(operator),
			NegateCondition:          pointer.To(condition["negate_condition"].(bool)),
			RulesEngineMatchValue:    matchValueArray,
			Transforms:               expandFrontDoorRulesEngineMatchConditionTransform(transform),
		}
		output = append(output, matchCondition)
	}
	return &output
}

func expandFrontDoorRulesEngineMatchConditionTransform(input []any) *[]frontdoors.Transform {
	if len(input) == 0 || input[0] == nil {
		return &[]frontdoors.Transform{}
	}

	output := make([]frontdoors.Transform, 0)

	for _, t := range input {
		result := frontdoors.Transform(t.(string))

		output = append(output, result)
	}
	return &output
}
