// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_application_firewall_policy

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandWebApplicationFirewallPolicyWebApplicationFirewallCustomRule(input []any) *[]webapplicationfirewallpolicies.WebApplicationFirewallCustomRule {
	results := make([]webapplicationfirewallpolicies.WebApplicationFirewallCustomRule, 0)
	for _, item := range input {
		v := item.(map[string]any)
		name := v["name"].(string)
		priority := v["priority"].(int)
		ruleType := v["rule_type"].(string)
		matchConditions := v["match_conditions"].([]any)
		action := v["action"].(string)

		enabled := webapplicationfirewallpolicies.WebApplicationFirewallStateEnabled
		if value, ok := v["enabled"].(bool); ok && !value {
			enabled = webapplicationfirewallpolicies.WebApplicationFirewallStateDisabled
		}

		result := webapplicationfirewallpolicies.WebApplicationFirewallCustomRule{
			State:           pointer.To(enabled),
			Action:          webapplicationfirewallpolicies.WebApplicationFirewallAction(action),
			MatchConditions: expandWebApplicationFirewallPolicyMatchCondition(matchConditions),
			Name:            pointer.To(name),
			Priority:        int64(priority),
			RuleType:        webapplicationfirewallpolicies.WebApplicationFirewallRuleType(ruleType),
		}

		if rateLimitDuration, ok := v["rate_limit_duration"]; ok && rateLimitDuration.(string) != "" {
			result.RateLimitDuration = pointer.ToEnum[webapplicationfirewallpolicies.ApplicationGatewayFirewallRateLimitDuration](rateLimitDuration.(string))
		}

		if rateLimitThreshHold, ok := v["rate_limit_threshold"]; ok && rateLimitThreshHold.(int) > 0 {
			result.RateLimitThreshold = pointer.To(int64(rateLimitThreshHold.(int)))
		}

		if groupBy, ok := v["group_rate_limit_by"]; ok && groupBy.(string) != "" {
			groups := []webapplicationfirewallpolicies.GroupByUserSession{
				{
					GroupByVariables: []webapplicationfirewallpolicies.GroupByVariable{
						{
							VariableName: webapplicationfirewallpolicies.ApplicationGatewayFirewallUserSessionVariable(groupBy.(string)),
						},
					},
				},
			}
			result.GroupByUserSession = &groups
		}

		results = append(results, result)
	}
	return &results
}

func expandWebApplicationFirewallPolicyPolicySettings(input []any) *webapplicationfirewallpolicies.PolicySettings {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]any)

	enabled := webapplicationfirewallpolicies.WebApplicationFirewallEnabledStateDisabled
	if value, ok := v["enabled"].(bool); ok && value {
		enabled = webapplicationfirewallpolicies.WebApplicationFirewallEnabledStateEnabled
	}
	mode := v["mode"].(string)
	requestBodyCheck := v["request_body_check"].(bool)
	requestBodyEnforcement := v["request_body_enforcement"].(bool)
	fileUploadEnforcement := v["file_upload_enforcement"].(bool)
	maxRequestBodySizeInKb := v["max_request_body_size_in_kb"].(int)
	fileUploadLimitInMb := v["file_upload_limit_in_mb"].(int)

	result := webapplicationfirewallpolicies.PolicySettings{
		State:                             pointer.To(enabled),
		Mode:                              pointer.ToEnum[webapplicationfirewallpolicies.WebApplicationFirewallMode](mode),
		FileUploadEnforcement:             pointer.To(fileUploadEnforcement),
		RequestBodyCheck:                  pointer.To(requestBodyCheck),
		RequestBodyEnforcement:            pointer.To(requestBodyEnforcement),
		MaxRequestBodySizeInKb:            pointer.To(int64(maxRequestBodySizeInKb)),
		FileUploadLimitInMb:               pointer.To(int64(fileUploadLimitInMb)),
		LogScrubbing:                      expandWebApplicationFirewallPolicyLogScrubbing(v["log_scrubbing"].([]any)),
		RequestBodyInspectLimitInKB:       pointer.To(int64(v["request_body_inspect_limit_in_kb"].(int))),
		JsChallengeCookieExpirationInMins: pointer.To(int64(v["js_challenge_cookie_expiration_in_minutes"].(int))),
	}

	return &result
}

func expandWebApplicationFirewallPolicyLogScrubbing(input []any) *webapplicationfirewallpolicies.PolicySettingsLogScrubbing {
	if len(input) == 0 {
		return nil
	}

	var res webapplicationfirewallpolicies.PolicySettingsLogScrubbing
	v := input[0].(map[string]any)
	state := webapplicationfirewallpolicies.WebApplicationFirewallScrubbingStateDisabled
	if value, ok := v["enabled"].(bool); ok && value {
		state = webapplicationfirewallpolicies.WebApplicationFirewallScrubbingStateEnabled
	}
	res.State = &state

	res.ScrubbingRules = expandWebApplicationPolicyScrubbingRules(v["rule"].([]any))

	return &res
}

func expandWebApplicationPolicyScrubbingRules(input []any) *[]webapplicationfirewallpolicies.WebApplicationFirewallScrubbingRules {
	if len(input) == 0 {
		return nil
	}

	res := make([]webapplicationfirewallpolicies.WebApplicationFirewallScrubbingRules, 0)
	for _, rule := range input {
		v := rule.(map[string]any)
		var item webapplicationfirewallpolicies.WebApplicationFirewallScrubbingRules
		state := webapplicationfirewallpolicies.ScrubbingRuleEntryStateDisabled
		if value, ok := v["enabled"].(bool); ok && value {
			state = webapplicationfirewallpolicies.ScrubbingRuleEntryStateEnabled
		}
		item.State = &state
		item.MatchVariable = webapplicationfirewallpolicies.ScrubbingRuleEntryMatchVariable(v["match_variable"].(string))
		item.SelectorMatchOperator = webapplicationfirewallpolicies.ScrubbingRuleEntryMatchOperator(v["selector_match_operator"].(string))
		if val, ok := v["selector"]; ok {
			item.Selector = pointer.To(val.(string))
		}

		res = append(res, item)
	}
	return &res
}

func expandWebApplicationFirewallPolicyManagedRulesDefinition(input []any) *webapplicationfirewallpolicies.ManagedRulesDefinition {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]any)

	return &webapplicationfirewallpolicies.ManagedRulesDefinition{
		Exceptions:      expandWebApplicationFirewallPolicyExceptions(v["exception"].([]any)),
		Exclusions:      expandWebApplicationFirewallPolicyExclusions(v["exclusion"].([]any)),
		ManagedRuleSets: *expandWebApplicationFirewallPolicyManagedRuleSet(v["managed_rule_set"].([]any)),
	}
}

func expandWebApplicationFirewallPolicyExceptions(input []any) *[]webapplicationfirewallpolicies.ExceptionEntry {
	results := make([]webapplicationfirewallpolicies.ExceptionEntry, 0)
	for _, item := range input {
		v := item.(map[string]any)
		result := webapplicationfirewallpolicies.ExceptionEntry{
			MatchVariable:      webapplicationfirewallpolicies.ExceptionEntryMatchVariable(v["match_variable"].(string)),
			ValueMatchOperator: webapplicationfirewallpolicies.ExceptionEntryValueMatchOperator(v["value_match_operator"].(string)),
		}
		if managedRuleSets := v["exception_rule_set"].([]any); len(managedRuleSets) > 0 {
			result.ExceptionManagedRuleSets = expandWebApplicationFirewallPolicyExceptionManagedRuleSets(managedRuleSets)
		}
		if selector := v["selector"].(string); selector != "" {
			result.Selector = pointer.To(selector)
		}
		if selectorMatchOperator := v["selector_match_operator"].(string); selectorMatchOperator != "" {
			result.SelectorMatchOperator = pointer.ToEnum[webapplicationfirewallpolicies.ExceptionEntrySelectorMatchOperator](selectorMatchOperator)
		}
		if values := v["values"].([]any); len(values) > 0 {
			result.Values = pluginsdk.ExpandStringSlice(values)
		}
		results = append(results, result)
	}
	return &results
}

func expandWebApplicationFirewallPolicyExceptionManagedRuleSets(input []any) *[]webapplicationfirewallpolicies.ExclusionManagedRuleSet {
	results := make([]webapplicationfirewallpolicies.ExclusionManagedRuleSet, 0)
	for _, item := range input {
		v := item.(map[string]any)
		result := webapplicationfirewallpolicies.ExclusionManagedRuleSet{
			RuleSetType:    v["type"].(string),
			RuleSetVersion: v["version"].(string),
		}
		if ruleGroups := v["rule_group"].([]any); len(ruleGroups) > 0 {
			result.RuleGroups = expandWebApplicationFirewallPolicyExceptionManagedRuleGroups(ruleGroups)
		}
		results = append(results, result)
	}
	return &results
}

func expandWebApplicationFirewallPolicyExceptionManagedRuleGroups(input []any) *[]webapplicationfirewallpolicies.ExclusionManagedRuleGroup {
	results := make([]webapplicationfirewallpolicies.ExclusionManagedRuleGroup, 0)
	for _, item := range input {
		v := item.(map[string]any)
		result := webapplicationfirewallpolicies.ExclusionManagedRuleGroup{
			RuleGroupName: v["rule_group_name"].(string),
		}
		if rules := v["rules"].([]any); len(rules) > 0 {
			result.Rules = expandWebApplicationFirewallPolicyManagedRuleIDs(rules)
		}
		results = append(results, result)
	}
	return &results
}

func expandWebApplicationFirewallPolicyManagedRuleIDs(input []any) *[]webapplicationfirewallpolicies.ExclusionManagedRule {
	results := make([]webapplicationfirewallpolicies.ExclusionManagedRule, 0)
	for _, item := range input {
		ruleID := item.(string)

		result := webapplicationfirewallpolicies.ExclusionManagedRule{
			RuleId: ruleID,
		}

		results = append(results, result)
	}
	return &results
}

func expandWebApplicationFirewallPolicyExclusionManagedRuleGroup(input []any) *[]webapplicationfirewallpolicies.ExclusionManagedRuleGroup {
	results := make([]webapplicationfirewallpolicies.ExclusionManagedRuleGroup, 0)
	for _, item := range input {
		v := item.(map[string]any)

		ruleGroupName := v["rule_group_name"].(string)

		result := webapplicationfirewallpolicies.ExclusionManagedRuleGroup{
			RuleGroupName: ruleGroupName,
		}

		if excludedRules := v["excluded_rules"].([]any); len(excludedRules) > 0 {
			result.Rules = expandWebApplicationFirewallPolicyManagedRuleIDs(excludedRules)
		}

		results = append(results, result)
	}
	return &results
}

func expandWebApplicationFirewallPolicyExclusionManagedRuleSet(input []any) *[]webapplicationfirewallpolicies.ExclusionManagedRuleSet {
	results := make([]webapplicationfirewallpolicies.ExclusionManagedRuleSet, 0)
	for _, item := range input {
		v := item.(map[string]any)

		ruleSetType := v["type"].(string)
		ruleSetVersion := v["version"].(string)
		ruleGroups := make([]any, 0)
		if value, exists := v["rule_group"]; exists {
			ruleGroups = value.([]any)
		}
		result := webapplicationfirewallpolicies.ExclusionManagedRuleSet{
			RuleSetType:    ruleSetType,
			RuleSetVersion: ruleSetVersion,
			RuleGroups:     expandWebApplicationFirewallPolicyExclusionManagedRuleGroup(ruleGroups),
		}

		results = append(results, result)
	}
	return &results
}

func expandWebApplicationFirewallPolicyExclusions(input []any) *[]webapplicationfirewallpolicies.OwaspCrsExclusionEntry {
	results := make([]webapplicationfirewallpolicies.OwaspCrsExclusionEntry, 0)
	for _, item := range input {
		v := item.(map[string]any)

		matchVariable := v["match_variable"].(string)
		selectorMatchOperator := v["selector_match_operator"].(string)
		selector := v["selector"].(string)
		exclusionManagedRuleSets := v["excluded_rule_set"].([]any)

		result := webapplicationfirewallpolicies.OwaspCrsExclusionEntry{
			MatchVariable:            webapplicationfirewallpolicies.OwaspCrsExclusionEntryMatchVariable(matchVariable),
			SelectorMatchOperator:    webapplicationfirewallpolicies.OwaspCrsExclusionEntrySelectorMatchOperator(selectorMatchOperator),
			Selector:                 selector,
			ExclusionManagedRuleSets: expandWebApplicationFirewallPolicyExclusionManagedRuleSet(exclusionManagedRuleSets),
		}

		results = append(results, result)
	}
	return &results
}

func expandWebApplicationFirewallPolicyManagedRuleSet(input []any) *[]webapplicationfirewallpolicies.ManagedRuleSet {
	results := make([]webapplicationfirewallpolicies.ManagedRuleSet, 0)

	for _, item := range input {
		v := item.(map[string]any)
		ruleSetType := v["type"].(string)
		ruleSetVersion := v["version"].(string)
		ruleGroupOverrides := []any{}
		if value, exists := v["rule_group_override"]; exists {
			ruleGroupOverrides = value.([]any)
		}

		expandedRuleGroupOverrides := expandWebApplicationFirewallPolicyRuleGroupOverrides(ruleGroupOverrides)

		result := webapplicationfirewallpolicies.ManagedRuleSet{
			RuleSetType:        ruleSetType,
			RuleSetVersion:     ruleSetVersion,
			RuleGroupOverrides: expandedRuleGroupOverrides,
		}

		results = append(results, result)
	}
	return &results
}

func expandWebApplicationFirewallPolicyRuleGroupOverrides(input []any) *[]webapplicationfirewallpolicies.ManagedRuleGroupOverride {
	results := make([]webapplicationfirewallpolicies.ManagedRuleGroupOverride, 0)
	for _, item := range input {
		v := item.(map[string]any)

		ruleGroupName := v["rule_group_name"].(string)

		result := webapplicationfirewallpolicies.ManagedRuleGroupOverride{
			RuleGroupName: ruleGroupName,
		}

		if rules := v["rule"].([]any); len(rules) > 0 {
			result.Rules = expandWebApplicationFirewallPolicyOverrideRules(rules)
		}

		results = append(results, result)
	}

	return &results
}

func expandWebApplicationFirewallPolicyOverrideRules(input []any) *[]webapplicationfirewallpolicies.ManagedRuleOverride {
	results := make([]webapplicationfirewallpolicies.ManagedRuleOverride, 0)
	for _, item := range input {
		v := item.(map[string]any)
		state := webapplicationfirewallpolicies.ManagedRuleEnabledStateDisabled
		if v["enabled"].(bool) {
			state = webapplicationfirewallpolicies.ManagedRuleEnabledStateEnabled
		}

		result := webapplicationfirewallpolicies.ManagedRuleOverride{
			RuleId: v["id"].(string),
			State:  pointer.To(state),
		}

		action := v["action"].(string)
		if action != "" {
			result.Action = pointer.ToEnum[webapplicationfirewallpolicies.ActionType](action)
		}

		results = append(results, result)
	}

	return &results
}

func expandWebApplicationFirewallPolicyMatchCondition(input []any) []webapplicationfirewallpolicies.MatchCondition {
	results := make([]webapplicationfirewallpolicies.MatchCondition, 0)
	for _, item := range input {
		v := item.(map[string]any)
		matchVariables := v["match_variables"].([]any)
		operator := v["operator"].(string)
		negationCondition := v["negation_condition"].(bool)
		matchValues := v["match_values"].([]any)
		transformsRaw := v["transforms"].(*pluginsdk.Set).List()

		var transforms []webapplicationfirewallpolicies.WebApplicationFirewallTransform
		for _, trans := range transformsRaw {
			transforms = append(transforms, webapplicationfirewallpolicies.WebApplicationFirewallTransform(trans.(string)))
		}
		result := webapplicationfirewallpolicies.MatchCondition{
			MatchValues:      pointer.From(pluginsdk.ExpandStringSlice(matchValues)),
			MatchVariables:   expandWebApplicationFirewallPolicyMatchVariable(matchVariables),
			NegationConditon: pointer.To(negationCondition),
			Operator:         webapplicationfirewallpolicies.WebApplicationFirewallOperator(operator),
			Transforms:       &transforms,
		}

		results = append(results, result)
	}
	return results
}

func expandWebApplicationFirewallPolicyMatchVariable(input []any) []webapplicationfirewallpolicies.MatchVariable {
	results := make([]webapplicationfirewallpolicies.MatchVariable, 0)
	for _, item := range input {
		v := item.(map[string]any)
		variableName := v["variable_name"].(string)
		selector := v["selector"].(string)

		result := webapplicationfirewallpolicies.MatchVariable{
			Selector:     pointer.To(selector),
			VariableName: webapplicationfirewallpolicies.WebApplicationFirewallMatchVariable(variableName),
		}

		results = append(results, result)
	}
	return results
}
