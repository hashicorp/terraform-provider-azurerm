// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_firewall_policy

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2025-03-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandCdnFrontDoorFirewallCustomRules(input []any) *webapplicationfirewallpolicies.CustomRuleList {
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
		matchConditions := expandCdnFrontDoorFirewallMatchConditions(custom["match_condition"].([]any))
		action := custom["action"].(string)

		output = append(output, webapplicationfirewallpolicies.CustomRule{
			Name:                       pointer.To(name),
			Priority:                   priority,
			EnabledState:               pointer.To(enabled),
			RuleType:                   webapplicationfirewallpolicies.RuleType(ruleType),
			RateLimitDurationInMinutes: pointer.To(rateLimitDurationInMinutes),
			RateLimitThreshold:         pointer.To(rateLimitThreshold),
			MatchConditions:            matchConditions,
			Action:                     webapplicationfirewallpolicies.ActionType(action),
		})
	}

	return &webapplicationfirewallpolicies.CustomRuleList{
		Rules: &output,
	}
}

func expandCdnFrontDoorFirewallMatchConditions(input []any) []webapplicationfirewallpolicies.MatchCondition {
	result := make([]webapplicationfirewallpolicies.MatchCondition, 0)
	if len(input) == 0 {
		return nil
	}

	for _, v := range input {
		match := v.(map[string]any)

		matchVariable := match["match_variable"].(string)
		selector := match["selector"].(string)
		operator := match["operator"].(string)
		matchValues := pluginsdk.ExpandStringSlice(match["match_values"].([]any))
		transforms := match["transforms"].([]any)

		matchCondition := webapplicationfirewallpolicies.MatchCondition{
			Operator:        webapplicationfirewallpolicies.Operator(operator),
			NegateCondition: pointer.To(match["negation_condition"].(bool)),
			MatchValue:      *matchValues,
			Transforms:      expandCdnFrontDoorFirewallTransforms(transforms),
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

func expandCdnFrontDoorFirewallTransforms(input []any) *[]webapplicationfirewallpolicies.TransformType {
	result := make([]webapplicationfirewallpolicies.TransformType, 0)
	if len(input) == 0 {
		return nil
	}

	for _, v := range input {
		result = append(result, webapplicationfirewallpolicies.TransformType(v.(string)))
	}

	return &result
}

func expandCdnFrontDoorFirewallManagedRules(input []any) (*webapplicationfirewallpolicies.ManagedRuleSetList, error) {
	if len(input) == 0 {
		return nil, nil
	}

	result := make([]webapplicationfirewallpolicies.ManagedRuleSet, 0)
	for _, mr := range input {
		managedRule := mr.(map[string]any)

		ruleType := managedRule["type"].(string)
		version := managedRule["version"].(string)
		action := managedRule["action"].(string)
		overrides := managedRule["override"].([]any)
		exclusions := expandCdnFrontDoorFirewallManagedRuleGroupExclusion(managedRule["exclusion"].([]any))

		fVersion := 1.0
		if v, err := strconv.ParseFloat(version, 64); err == nil {
			fVersion = v
		}

		// NOTE: The API is deferring the version range from the rule type name
		// 'DefaultRuleSet' is < 1.1 and 'Microsoft_DefaultRuleSet' >= 1.1
		// 'AnomalyScoring' action only valid on 2.0 and above
		if ruleType == "DefaultRuleSet" && fVersion > 1.0 {
			return nil, fmt.Errorf("the managed rule set type %q and version %q is not supported. If you wish to use the 'DefaultRuleSet' type please update your 'version' field to be '1.0' or 'preview-0.1', got %q", ruleType, version, version)
		} else if ruleType == "Microsoft_DefaultRuleSet" && fVersion < 1.1 {
			return nil, fmt.Errorf("the managed rule set type %q and version %q is not supported. If you wish to use the 'Microsoft_DefaultRuleSet' type please update your 'version' field to be '1.1', '2.0' or '2.1', got %q", ruleType, version, version)
		}

		ruleGroupOverrides, err := expandCdnFrontDoorFirewallManagedRuleGroupOverride(overrides, version, fVersion, ruleType)
		if err != nil {
			return nil, err
		}

		managedRuleSet := webapplicationfirewallpolicies.ManagedRuleSet{
			Exclusions:         exclusions,
			RuleSetVersion:     version,
			RuleGroupOverrides: ruleGroupOverrides,
			RuleSetType:        ruleType,
		}

		if action != "" {
			managedRuleSet.RuleSetAction = pointer.ToEnum[webapplicationfirewallpolicies.ManagedRuleSetActionType](action)
		}

		result = append(result, managedRuleSet)
	}

	return &webapplicationfirewallpolicies.ManagedRuleSetList{
		ManagedRuleSets: &result,
	}, nil
}

func expandCdnFrontDoorFirewallManagedRuleGroupExclusion(input []any) *[]webapplicationfirewallpolicies.ManagedRuleExclusion {
	results := make([]webapplicationfirewallpolicies.ManagedRuleExclusion, 0)
	if len(input) == 0 {
		return nil
	}

	for _, v := range input {
		exclusion := v.(map[string]any)

		matchVariable := exclusion["match_variable"].(string)
		operator := exclusion["operator"].(string)
		selector := exclusion["selector"].(string)

		results = append(results, webapplicationfirewallpolicies.ManagedRuleExclusion{
			MatchVariable:         webapplicationfirewallpolicies.ManagedRuleExclusionMatchVariable(matchVariable),
			SelectorMatchOperator: webapplicationfirewallpolicies.ManagedRuleExclusionSelectorMatchOperator(operator),
			Selector:              selector,
		})
	}

	return &results
}

func expandCdnFrontDoorFirewallManagedRuleGroupOverride(input []any, versionRaw string, version float64, ruleType string) (*[]webapplicationfirewallpolicies.ManagedRuleGroupOverride, error) {
	result := make([]webapplicationfirewallpolicies.ManagedRuleGroupOverride, 0)
	if len(input) == 0 {
		return nil, nil
	}

	for _, v := range input {
		override := v.(map[string]any)

		exclusions := expandCdnFrontDoorFirewallManagedRuleGroupExclusion(override["exclusion"].([]any))
		ruleGroupName := override["rule_group_name"].(string)
		rules, err := expandCdnFrontDoorFirewallRuleOverride(override["rule"].([]any), versionRaw, version, ruleType)
		if err != nil {
			return nil, err
		}

		result = append(result, webapplicationfirewallpolicies.ManagedRuleGroupOverride{
			Exclusions:    exclusions,
			RuleGroupName: ruleGroupName,
			Rules:         rules,
		})
	}

	return &result, nil
}

func expandCdnFrontDoorFirewallRuleOverride(input []any, versionRaw string, version float64, ruleType string) (*[]webapplicationfirewallpolicies.ManagedRuleOverride, error) {
	result := make([]webapplicationfirewallpolicies.ManagedRuleOverride, 0)
	if len(input) == 0 {
		return nil, nil
	}

	for _, v := range input {
		rule := v.(map[string]any)

		enabled := webapplicationfirewallpolicies.ManagedRuleEnabledStateDisabled
		if rule["enabled"].(bool) {
			enabled = webapplicationfirewallpolicies.ManagedRuleEnabledStateEnabled
		}

		ruleId := rule["rule_id"].(string)
		action := webapplicationfirewallpolicies.ActionType(rule["action"].(string))

		// NOTE: Default Rule Sets(DRS) 2.0 and above rules only use action type of 'AnomalyScoring' or 'Log'. Issues 19088 and 19561
		// This will still work for bot rules as well since it will be the default value of 1.0
		switch {
		case version < 2.0 && action == webapplicationfirewallpolicies.ActionTypeAnomalyScoring:
			return nil, fmt.Errorf("%q is only valid in managed rules where 'type' is DRS and `version` is '2.0' or above, got %q", webapplicationfirewallpolicies.ActionTypeAnomalyScoring, versionRaw)

		case version >= 2.0 && action != webapplicationfirewallpolicies.ActionTypeAnomalyScoring && action != webapplicationfirewallpolicies.ActionTypeLog:
			return nil, fmt.Errorf("the managed rules 'action' field must be set to 'AnomalyScoring' or 'Log' if the managed rule is DRS 2.0 or above, got %q", action)

		case !strings.Contains(strings.ToLower(ruleType), "botmanagerruleset") && action == webapplicationfirewallpolicies.ActionTypeJSChallenge:
			return nil, fmt.Errorf("%q is only valid if the managed rules 'type' is 'Microsoft_BotManagerRuleSet', got %q", webapplicationfirewallpolicies.ActionTypeJSChallenge, ruleType)
		}

		exclusions := expandCdnFrontDoorFirewallManagedRuleGroupExclusion(rule["exclusion"].([]any))

		result = append(result, webapplicationfirewallpolicies.ManagedRuleOverride{
			RuleId:       ruleId,
			EnabledState: pointer.To(enabled),
			Action:       pointer.To(action),
			Exclusions:   exclusions,
		})
	}

	return &result, nil
}

func expandCdnFrontDoorFirewallLogScrubbingPolicy(input []any) (*webapplicationfirewallpolicies.PolicySettingsLogScrubbing, error) {
	if len(input) == 0 {
		return nil, nil
	}

	inputRaw := input[0].(map[string]any)

	policyEnabled := webapplicationfirewallpolicies.WebApplicationFirewallScrubbingStateDisabled
	if inputRaw["enabled"].(bool) {
		policyEnabled = webapplicationfirewallpolicies.WebApplicationFirewallScrubbingStateEnabled
	}

	scrubbingRules, err := expandCdnFrontDoorFirewallScrubbingRules(inputRaw["scrubbing_rule"].([]any))
	if err != nil {
		return nil, err
	}

	return &webapplicationfirewallpolicies.PolicySettingsLogScrubbing{
		State:          pointer.To(policyEnabled),
		ScrubbingRules: scrubbingRules,
	}, nil
}

func expandCdnFrontDoorFirewallScrubbingRules(input []any) (*[]webapplicationfirewallpolicies.WebApplicationFirewallScrubbingRules, error) {
	if len(input) == 0 {
		return nil, nil
	}

	scrubbingRules := make([]webapplicationfirewallpolicies.WebApplicationFirewallScrubbingRules, 0)

	for _, rule := range input {
		v := rule.(map[string]any)
		var item webapplicationfirewallpolicies.WebApplicationFirewallScrubbingRules

		enabled := webapplicationfirewallpolicies.ScrubbingRuleEntryStateDisabled
		if value := v["enabled"].(bool); value {
			enabled = webapplicationfirewallpolicies.ScrubbingRuleEntryStateEnabled
		}

		item.State = pointer.To(enabled)
		item.MatchVariable = webapplicationfirewallpolicies.ScrubbingRuleEntryMatchVariable(v["match_variable"].(string))
		item.SelectorMatchOperator = webapplicationfirewallpolicies.ScrubbingRuleEntryMatchOperator(v["operator"].(string))

		if selector, ok := v["selector"]; ok {
			item.Selector = pointer.To(selector.(string))
		}

		// NOTE: Validate the rules configuration...
		switch {
		case item.MatchVariable == webapplicationfirewallpolicies.ScrubbingRuleEntryMatchVariableRequestIPAddress || item.MatchVariable == webapplicationfirewallpolicies.ScrubbingRuleEntryMatchVariableRequestUri:
			// NOTE: 'RequestIPAddress' and 'RequestUri' 'match_variable's can only use the 'EqualsAny' 'operator'...
			if item.SelectorMatchOperator != webapplicationfirewallpolicies.ScrubbingRuleEntryMatchOperatorEqualsAny {
				return nil, fmt.Errorf("the %q 'match_variable' must use the %q 'operator', got %q", item.MatchVariable, webapplicationfirewallpolicies.ScrubbingRuleEntryMatchOperatorEqualsAny, item.SelectorMatchOperator)
			}

		case item.SelectorMatchOperator == webapplicationfirewallpolicies.ScrubbingRuleEntryMatchOperatorEquals:
			// NOTE: If the 'operator' is set to 'Equals' the 'selector' cannot be 'nil'...
			if pointer.From(item.Selector) == "" {
				return nil, fmt.Errorf("the 'selector' field must be set when the %q 'operator' is used, got %q", webapplicationfirewallpolicies.ScrubbingRuleEntryMatchOperatorEquals, "nil")
			}

		case item.SelectorMatchOperator == webapplicationfirewallpolicies.ScrubbingRuleEntryMatchOperatorEqualsAny:
			// NOTE: If the 'operator' is set to 'EqualsAny' the 'selector' must be 'nil'...
			if pointer.From(item.Selector) != "" {
				return nil, fmt.Errorf("the 'selector' field cannot be set when the %q 'operator' is used, got %q", webapplicationfirewallpolicies.ScrubbingRuleEntryMatchOperatorEqualsAny, pointer.From(item.Selector))
			}
		}

		scrubbingRules = append(scrubbingRules, item)
	}

	return pointer.To(scrubbingRules), nil
}
