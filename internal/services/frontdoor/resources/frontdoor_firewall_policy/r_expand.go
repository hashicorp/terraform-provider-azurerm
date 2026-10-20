// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package frontdoor_firewall_policy

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-04-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandFrontDoorFirewallCustomRules(input []any) *webapplicationfirewallpolicies.CustomRuleList {
	if len(input) == 0 {
		return nil
	}

	output := make([]webapplicationfirewallpolicies.CustomRule, 0)

	for _, cr := range input {
		custom := cr.(map[string]any)

		enabled := webapplicationfirewallpolicies.CustomRuleEnabledStateDisabled
		if custom["enabled"].(bool) {
			enabled = webapplicationfirewallpolicies.CustomRuleEnabledStateEnabled
		}

		name := custom["name"].(string)
		priority := int64(custom["priority"].(int))
		ruleType := custom["type"].(string)
		rateLimitDurationInMinutes := int64(custom["rate_limit_duration_in_minutes"].(int))
		rateLimitThreshold := int64(custom["rate_limit_threshold"].(int))
		matchConditions := expandFrontDoorFirewallMatchConditions(custom["match_condition"].([]any))
		action := custom["action"].(string)

		customRule := webapplicationfirewallpolicies.CustomRule{
			Name:                       pointer.To(name),
			Priority:                   priority,
			EnabledState:               &enabled,
			RuleType:                   webapplicationfirewallpolicies.RuleType(ruleType),
			RateLimitDurationInMinutes: pointer.To(rateLimitDurationInMinutes),
			RateLimitThreshold:         pointer.To(rateLimitThreshold),
			MatchConditions:            matchConditions,
			Action:                     webapplicationfirewallpolicies.ActionType(action),
		}
		output = append(output, customRule)
	}

	return &webapplicationfirewallpolicies.CustomRuleList{
		Rules: &output,
	}
}

func expandFrontDoorFirewallMatchConditions(input []any) []webapplicationfirewallpolicies.MatchCondition {
	if len(input) == 0 {
		return nil
	}

	result := make([]webapplicationfirewallpolicies.MatchCondition, 0)

	for _, v := range input {
		match := v.(map[string]any)

		matchVariable := match["match_variable"].(string)
		selector := match["selector"].(string)
		operator := match["operator"].(string)
		matchValues := match["match_values"].([]any)
		transforms := match["transforms"].([]any)

		matchCondition := webapplicationfirewallpolicies.MatchCondition{
			Operator:        webapplicationfirewallpolicies.Operator(operator),
			NegateCondition: pointer.To(match["negation_condition"].(bool)),
			MatchValue:      *pluginsdk.ExpandStringSlice(matchValues),
			Transforms:      expandFrontDoorFirewallTransforms(transforms),
		}

		if matchVariable != "" {
			matchCondition.MatchVariable = webapplicationfirewallpolicies.MatchVariable(matchVariable)
		}
		if selector != "" {
			matchCondition.Selector = pointer.To(selector)
		}

		result = append(result, matchCondition)
	}

	return result
}

func expandFrontDoorFirewallTransforms(input []any) *[]webapplicationfirewallpolicies.TransformType {
	if len(input) == 0 {
		return nil
	}

	result := make([]webapplicationfirewallpolicies.TransformType, 0)
	for _, v := range input {
		result = append(result, webapplicationfirewallpolicies.TransformType(v.(string)))
	}

	return &result
}

func expandFrontDoorFirewallManagedRules(input []any) *webapplicationfirewallpolicies.ManagedRuleSetList {
	if len(input) == 0 {
		return nil
	}

	managedRules := make([]webapplicationfirewallpolicies.ManagedRuleSet, 0)

	for _, mr := range input {
		managedRule := mr.(map[string]any)

		ruleType := managedRule["type"].(string)
		version := managedRule["version"].(string)

		managedRuleSet := webapplicationfirewallpolicies.ManagedRuleSet{
			RuleSetType:    ruleType,
			RuleSetVersion: version,
		}

		if exclusions := expandFrontDoorFirewallManagedRuleGroupExclusion(managedRule["exclusion"].([]any)); exclusions != nil {
			managedRuleSet.Exclusions = exclusions
		}

		if ruleGroupOverrides := expandFrontDoorFirewallManagedRuleGroupOverride(managedRule["override"].([]any)); ruleGroupOverrides != nil {
			managedRuleSet.RuleGroupOverrides = ruleGroupOverrides
		}

		managedRules = append(managedRules, managedRuleSet)
	}

	return &webapplicationfirewallpolicies.ManagedRuleSetList{
		ManagedRuleSets: &managedRules,
	}
}

func expandFrontDoorFirewallManagedRuleGroupExclusion(input []any) *[]webapplicationfirewallpolicies.ManagedRuleExclusion {
	if len(input) == 0 {
		return nil
	}

	managedRuleExclusions := make([]webapplicationfirewallpolicies.ManagedRuleExclusion, 0)
	for _, v := range input {
		exclusion := v.(map[string]any)

		matchVariable := exclusion["match_variable"].(string)
		operator := exclusion["operator"].(string)
		selector := exclusion["selector"].(string)

		managedRuleExclusion := webapplicationfirewallpolicies.ManagedRuleExclusion{
			MatchVariable:         webapplicationfirewallpolicies.ManagedRuleExclusionMatchVariable(matchVariable),
			SelectorMatchOperator: webapplicationfirewallpolicies.ManagedRuleExclusionSelectorMatchOperator(operator),
			Selector:              selector,
		}

		managedRuleExclusions = append(managedRuleExclusions, managedRuleExclusion)
	}

	return &managedRuleExclusions
}

func expandFrontDoorFirewallManagedRuleGroupOverride(input []any) *[]webapplicationfirewallpolicies.ManagedRuleGroupOverride {
	if len(input) == 0 {
		return nil
	}

	managedRuleGroupOverrides := make([]webapplicationfirewallpolicies.ManagedRuleGroupOverride, 0)
	for _, v := range input {
		override := v.(map[string]any)

		ruleGroupName := override["rule_group_name"].(string)

		managedRuleGroupOverride := webapplicationfirewallpolicies.ManagedRuleGroupOverride{
			RuleGroupName: ruleGroupName,
		}

		if exclusions := expandFrontDoorFirewallManagedRuleGroupExclusion(override["exclusion"].([]any)); exclusions != nil {
			managedRuleGroupOverride.Exclusions = exclusions
		}

		if managedRuleOverride := expandFrontDoorFirewallRuleOverride(override["rule"].([]any)); managedRuleOverride != nil {
			managedRuleGroupOverride.Rules = managedRuleOverride
		}

		managedRuleGroupOverrides = append(managedRuleGroupOverrides, managedRuleGroupOverride)
	}

	return &managedRuleGroupOverrides
}

func expandFrontDoorFirewallRuleOverride(input []any) *[]webapplicationfirewallpolicies.ManagedRuleOverride {
	if len(input) == 0 {
		return nil
	}

	managedRuleOverrides := make([]webapplicationfirewallpolicies.ManagedRuleOverride, 0)
	for _, v := range input {
		rule := v.(map[string]any)

		enabled := webapplicationfirewallpolicies.ManagedRuleEnabledStateDisabled
		if rule["enabled"].(bool) {
			enabled = webapplicationfirewallpolicies.ManagedRuleEnabledStateEnabled
		}
		ruleId := rule["rule_id"].(string)

		managedRuleOverride := webapplicationfirewallpolicies.ManagedRuleOverride{
			RuleId:       ruleId,
			EnabledState: &enabled,
			Action:       pointer.ToEnum[webapplicationfirewallpolicies.ActionType](rule["action"].(string)),
		}

		if exclusions := expandFrontDoorFirewallManagedRuleGroupExclusion(rule["exclusion"].([]any)); exclusions != nil {
			managedRuleOverride.Exclusions = exclusions
		}

		managedRuleOverrides = append(managedRuleOverrides, managedRuleOverride)
	}

	return &managedRuleOverrides
}
