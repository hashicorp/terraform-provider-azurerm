// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package route_map

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualwans"
)

func flattenRules(input *[]virtualwans.RouteMapRule) []Rule {
	if input == nil {
		return []Rule{}
	}

	rules := make([]Rule, 0, len(*input))
	for _, v := range *input {
		rule := Rule{
			Actions:       flattenActions(v.Actions),
			MatchCriteria: flattenCriteria(v.MatchCriteria),
		}

		if v.Name != nil {
			rule.Name = *v.Name
		}

		if v.NextStepIfMatched != nil {
			rule.NextStepIfMatched = pointer.From(v.NextStepIfMatched)
		}

		rules = append(rules, rule)
	}

	return rules
}

func flattenActions(input *[]virtualwans.Action) []Action {
	if input == nil {
		return []Action{}
	}

	actions := make([]Action, 0, len(*input))
	for _, v := range *input {
		action := Action{
			Parameters: flattenParameters(v.Parameters),
		}

		if v.Type != nil {
			action.Type = pointer.From(v.Type)
		}

		actions = append(actions, action)
	}

	return actions
}

func flattenParameters(input *[]virtualwans.Parameter) []Parameter {
	if input == nil {
		return []Parameter{}
	}

	parameters := make([]Parameter, 0, len(*input))
	for _, v := range *input {
		parameter := Parameter{}

		if v.AsPath != nil {
			parameter.AsPath = *v.AsPath
		}

		if v.Community != nil {
			parameter.Community = *v.Community
		}

		if v.RoutePrefix != nil {
			parameter.RoutePrefix = *v.RoutePrefix
		}

		parameters = append(parameters, parameter)
	}

	return parameters
}

func flattenCriteria(input *[]virtualwans.Criterion) []Criterion {
	if input == nil {
		return []Criterion{}
	}

	criteria := make([]Criterion, 0, len(*input))
	for _, v := range *input {
		criterion := Criterion{}

		if v.AsPath != nil {
			criterion.AsPath = *v.AsPath
		}

		if v.Community != nil {
			criterion.Community = *v.Community
		}

		if v.MatchCondition != nil {
			criterion.MatchCondition = pointer.From(v.MatchCondition)
		}

		if v.RoutePrefix != nil {
			criterion.RoutePrefix = *v.RoutePrefix
		}

		criteria = append(criteria, criterion)
	}

	return criteria
}
