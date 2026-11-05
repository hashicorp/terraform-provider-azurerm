// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_firewall_policy

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2025-03-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
)

func flattenCdnFrontDoorFirewallCustomRules(input *webapplicationfirewallpolicies.CustomRuleList) []any {
	if input == nil || input.Rules == nil {
		return []any{}
	}

	results := make([]any, 0)
	for _, v := range *input.Rules {
		action := string(v.Action)
		priority := int(v.Priority)
		ruleType := string(v.RuleType)

		enabled := false
		if v.EnabledState != nil {
			enabled = pointer.From(v.EnabledState) == webapplicationfirewallpolicies.CustomRuleEnabledStateEnabled
		}

		rateLimitDurationInMinutes := 0
		if v.RateLimitDurationInMinutes != nil {
			rateLimitDurationInMinutes = int(*v.RateLimitDurationInMinutes)
		}

		rateLimitThreshold := 0
		if v.RateLimitThreshold != nil {
			rateLimitThreshold = int(*v.RateLimitThreshold)
		}

		results = append(results, map[string]any{
			"action":                         action,
			"enabled":                        enabled,
			"match_condition":                flattenCdnFrontDoorFirewallMatchConditions(v.MatchConditions),
			"rate_limit_duration_in_minutes": rateLimitDurationInMinutes,
			"rate_limit_threshold":           rateLimitThreshold,
			"priority":                       priority,
			"name":                           pointer.From(v.Name),
			"type":                           ruleType,
		})
	}

	return results
}

func flattenCdnFrontDoorFirewallMatchConditions(input []webapplicationfirewallpolicies.MatchCondition) []any {
	if input == nil {
		return []any{}
	}

	results := make([]any, 0)
	for _, v := range input {
		results = append(results, map[string]any{
			"match_variable":     string(v.MatchVariable),
			"match_values":       v.MatchValue,
			"negation_condition": pointer.From(v.NegateCondition),
			"operator":           string(v.Operator),
			"selector":           pointer.From(v.Selector),
			"transforms":         helpers.FlattenTransformSlice(v.Transforms),
		})
	}

	return results
}

func flattenCdnFrontDoorFirewallManagedRules(input *webapplicationfirewallpolicies.ManagedRuleSetList) []any {
	if input == nil || input.ManagedRuleSets == nil {
		return []any{}
	}

	results := make([]any, 0)
	for _, r := range *input.ManagedRuleSets {
		ruleSetType := r.RuleSetType
		ruleSetVersion := r.RuleSetVersion
		ruleSetAction := pointer.FromEnum(r.RuleSetAction)

		results = append(results, map[string]any{
			"exclusion": flattenCdnFrontDoorFirewallExclusions(r.Exclusions),
			"override":  flattenCdnFrontDoorFirewallOverrides(r.RuleGroupOverrides),
			"type":      ruleSetType,
			"version":   ruleSetVersion,
			"action":    ruleSetAction,
		})
	}

	return results
}

func flattenCdnFrontDoorFirewallExclusions(input *[]webapplicationfirewallpolicies.ManagedRuleExclusion) []any {
	if input == nil {
		return []any{}
	}

	results := make([]any, 0)
	for _, v := range *input {
		matchVariable := string(v.MatchVariable)
		operator := string(v.SelectorMatchOperator)
		selector := v.Selector

		results = append(results, map[string]any{
			"match_variable": matchVariable,
			"operator":       operator,
			"selector":       selector,
		})
	}

	return results
}

func flattenCdnFrontDoorFirewallOverrides(input *[]webapplicationfirewallpolicies.ManagedRuleGroupOverride) []any {
	if input == nil {
		return []any{}
	}

	results := make([]any, 0)
	for _, v := range *input {
		ruleGroupName := v.RuleGroupName

		results = append(results, map[string]any{
			"rule_group_name": ruleGroupName,
			"exclusion":       flattenCdnFrontDoorFirewallExclusions(v.Exclusions),
			"rule":            flattenCdnFrontDoorFirewallRules(v.Rules),
		})
	}

	return results
}

func flattenCdnFrontDoorFirewallRules(input *[]webapplicationfirewallpolicies.ManagedRuleOverride) []any {
	if input == nil {
		return []any{}
	}

	results := make([]any, 0)
	for _, v := range *input {
		action := webapplicationfirewallpolicies.ActionTypeAnomalyScoring
		if v.Action != nil {
			action = pointer.From(v.Action)
		}
		enabled := pointer.From(v.EnabledState) == webapplicationfirewallpolicies.ManagedRuleEnabledStateEnabled
		ruleId := v.RuleId

		results = append(results, map[string]any{
			"action":    action,
			"enabled":   enabled,
			"exclusion": flattenCdnFrontDoorFirewallExclusions(v.Exclusions),
			"rule_id":   ruleId,
		})
	}

	return results
}

func flattenCdnFrontDoorFirewallLogScrubbingPolicy(input *webapplicationfirewallpolicies.PolicySettingsLogScrubbing) []any {
	if input == nil {
		return make([]any, 0)
	}

	result := make(map[string]any)
	result["enabled"] = pointer.From(input.State) == webapplicationfirewallpolicies.WebApplicationFirewallScrubbingStateEnabled
	result["scrubbing_rule"] = flattenCdnFrontDoorFirewallLogScrubbingRules(input.ScrubbingRules)

	return []any{result}
}

func flattenCdnFrontDoorFirewallLogScrubbingRules(scrubbingRules *[]webapplicationfirewallpolicies.WebApplicationFirewallScrubbingRules) any {
	result := make([]any, 0)

	if scrubbingRules == nil || len(*scrubbingRules) == 0 {
		return result
	}

	for _, scrubbingRule := range *scrubbingRules {
		item := map[string]any{}
		item["enabled"] = pointer.From(scrubbingRule.State) == webapplicationfirewallpolicies.ScrubbingRuleEntryStateEnabled
		item["match_variable"] = scrubbingRule.MatchVariable
		item["operator"] = scrubbingRule.SelectorMatchOperator
		item["selector"] = pointer.From(scrubbingRule.Selector)

		result = append(result, item)
	}

	return result
}
