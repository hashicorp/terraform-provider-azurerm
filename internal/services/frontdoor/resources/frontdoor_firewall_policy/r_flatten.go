// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package frontdoor_firewall_policy

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/frontdoor/2020-04-01/webapplicationfirewallpolicies"
)

func flattenFrontDoorFirewallCustomRules(input *webapplicationfirewallpolicies.CustomRuleList) []any {
	if input == nil || input.Rules == nil {
		return make([]any, 0)
	}

	results := make([]any, 0)
	for _, r := range *input.Rules {
		output := make(map[string]any)

		output["name"] = r.Name
		output["type"] = string(r.RuleType)
		output["action"] = string(r.Action)
		if r.EnabledState != nil {
			output["enabled"] = *r.EnabledState == webapplicationfirewallpolicies.CustomRuleEnabledStateEnabled
		}
		output["match_condition"] = flattenFrontDoorFirewallMatchConditions(r.MatchConditions)
		output["priority"] = int(r.Priority)

		if v := r.RateLimitDurationInMinutes; v != nil {
			output["rate_limit_duration_in_minutes"] = int(*v)
		}

		if v := r.RateLimitThreshold; v != nil {
			output["rate_limit_threshold"] = int(*v)
		}

		results = append(results, output)
	}

	return results
}

func flattenFrontDoorFirewallMatchConditions(condition []webapplicationfirewallpolicies.MatchCondition) []any {
	if condition == nil {
		return make([]any, 0)
	}

	results := make([]any, 0)
	for _, c := range condition {
		output := make(map[string]any)

		output["match_variable"] = string(c.MatchVariable)
		output["operator"] = string(c.Operator)
		output["match_values"] = c.MatchValue
		output["transforms"] = FlattenTransformSlice(c.Transforms)

		if v := c.Selector; v != nil {
			output["selector"] = *v
		}

		if v := c.NegateCondition; v != nil {
			output["negation_condition"] = *v
		}

		results = append(results, output)
	}

	return results
}

func flattenFrontDoorFirewallManagedRules(input *webapplicationfirewallpolicies.ManagedRuleSetList) []any {
	if input == nil || input.ManagedRuleSets == nil {
		return make([]any, 0)
	}

	results := make([]any, 0)
	for _, r := range *input.ManagedRuleSets {
		output := make(map[string]any)

		output["type"] = r.RuleSetType

		output["version"] = r.RuleSetVersion

		if v := r.RuleGroupOverrides; v != nil {
			output["override"] = flattenFrontDoorFirewallOverrides(v)
		}

		if v := r.Exclusions; v != nil {
			output["exclusion"] = flattenFrontDoorFirewallExclusions(v)
		}

		results = append(results, output)
	}

	return results
}

func flattenFrontDoorFirewallExclusions(managedRuleExclusion *[]webapplicationfirewallpolicies.ManagedRuleExclusion) []any {
	if managedRuleExclusion == nil {
		return make([]any, 0)
	}

	results := make([]any, 0)
	for _, o := range *managedRuleExclusion {
		output := make(map[string]any)

		output["match_variable"] = o.MatchVariable
		output["operator"] = o.SelectorMatchOperator
		output["selector"] = o.Selector

		results = append(results, output)
	}

	return results
}

func flattenFrontDoorFirewallOverrides(groupOverride *[]webapplicationfirewallpolicies.ManagedRuleGroupOverride) []any {
	if groupOverride == nil {
		return make([]any, 0)
	}

	results := make([]any, 0)
	for _, o := range *groupOverride {
		output := make(map[string]any)

		output["rule_group_name"] = o.RuleGroupName

		if v := o.Exclusions; v != nil {
			output["exclusion"] = flattenFrontDoorFirewallExclusions(v)
		}

		if rules := o.Rules; rules != nil {
			output["rule"] = flattenArmFrontdoorFirewallRules(rules)
		}

		results = append(results, output)
	}

	return results
}

func flattenArmFrontdoorFirewallRules(override *[]webapplicationfirewallpolicies.ManagedRuleOverride) []any {
	if override == nil {
		return make([]any, 0)
	}

	results := make([]any, 0)
	for _, o := range *override {
		output := make(map[string]any)

		if o.EnabledState != nil {
			output["enabled"] = *o.EnabledState == webapplicationfirewallpolicies.ManagedRuleEnabledStateEnabled
		}
		if o.Action != nil {
			output["action"] = string(*o.Action)
		}

		output["rule_id"] = o.RuleId

		if v := o.Exclusions; v != nil {
			output["exclusion"] = flattenFrontDoorFirewallExclusions(v)
		}

		results = append(results, output)
	}

	return results
}
