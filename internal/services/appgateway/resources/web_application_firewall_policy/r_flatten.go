// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package web_application_firewall_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/webapplicationfirewallpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceWebApplicationFirewallPolicySetFlatten(d *pluginsdk.ResourceData, id *webapplicationfirewallpolicies.ApplicationGatewayWebApplicationFirewallPolicyId, model *webapplicationfirewallpolicies.WebApplicationFirewallPolicy) error {
	d.Set("name", id.ApplicationGatewayWebApplicationFirewallPolicyName)
	d.Set("resource_group_name", id.ResourceGroupName)
	if model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		if prop := model.Properties; prop != nil {
			if err := d.Set("custom_rules", flattenWebApplicationFirewallPolicyWebApplicationFirewallCustomRule(prop.CustomRules)); err != nil {
				return fmt.Errorf("setting `custom_rules`: %+v", err)
			}
			if err := d.Set("policy_settings", flattenWebApplicationFirewallPolicyPolicySettings(prop.PolicySettings)); err != nil {
				return fmt.Errorf("setting `policy_settings`: %+v", err)
			}
			if err := d.Set("managed_rules", flattenWebApplicationFirewallPolicyManagedRulesDefinition(prop.ManagedRules)); err != nil {
				return fmt.Errorf("setting `managed_rules`: %+v", err)
			}
			if err := d.Set("http_listener_ids", flattenWebApplicationFirewallPoliciesSubResourcesToIDs(prop.HTTPListeners)); err != nil {
				return fmt.Errorf("setting `http_listeners`: %+v", err)
			}
			if err := d.Set("path_based_rule_ids", flattenWebApplicationFirewallPoliciesSubResourcesToIDs(prop.PathBasedRules)); err != nil {
				return fmt.Errorf("setting `path_based_rules`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}
	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenWebApplicationFirewallPolicyWebApplicationFirewallCustomRule(input *[]webapplicationfirewallpolicies.WebApplicationFirewallCustomRule) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]any)

		if name := item.Name; name != nil {
			v["name"] = *name
		}
		v["action"] = string(item.Action)
		v["enabled"] = pointer.From(item.State) == webapplicationfirewallpolicies.WebApplicationFirewallStateEnabled
		v["match_conditions"] = flattenWebApplicationFirewallPolicyMatchCondition(item.MatchConditions)
		v["priority"] = int(item.Priority)
		v["rule_type"] = string(item.RuleType)
		v["rate_limit_duration"] = pointer.From(item.RateLimitDuration)
		v["rate_limit_threshold"] = pointer.From(item.RateLimitThreshold)

		if item.GroupByUserSession != nil && len(*item.GroupByUserSession) > 0 {
			if groupVariable := (*item.GroupByUserSession)[0].GroupByVariables; len(groupVariable) > 0 {
				v["group_rate_limit_by"] = groupVariable[0].VariableName
			}
		}

		results = append(results, v)
	}

	return results
}

func flattenWebApplicationFirewallPolicyPolicySettings(input *webapplicationfirewallpolicies.PolicySettings) []any {
	if input == nil {
		return make([]any, 0)
	}

	result := make(map[string]any)

	result["enabled"] = pointer.From(input.State) == webapplicationfirewallpolicies.WebApplicationFirewallEnabledStateEnabled
	result["mode"] = pointer.FromEnum(input.Mode)
	result["request_body_check"] = input.RequestBodyCheck
	result["request_body_enforcement"] = input.RequestBodyEnforcement
	result["file_upload_enforcement"] = input.FileUploadEnforcement
	result["max_request_body_size_in_kb"] = int(pointer.From(input.MaxRequestBodySizeInKb))
	result["file_upload_limit_in_mb"] = int(pointer.From(input.FileUploadLimitInMb))
	result["log_scrubbing"] = flattenWebApplicationFirewallPolicyLogScrubbing(input.LogScrubbing)
	result["request_body_inspect_limit_in_kb"] = pointer.From(input.RequestBodyInspectLimitInKB)

	jsChallengeCookieExpirationInMins := 30 // default value is not returned, https://github.com/Azure/azure-rest-api-specs/issues/36330
	if v := pointer.From(input.JsChallengeCookieExpirationInMins); v != 0 {
		jsChallengeCookieExpirationInMins = int(v)
	}
	result["js_challenge_cookie_expiration_in_minutes"] = jsChallengeCookieExpirationInMins

	return []any{result}
}

func flattenWebApplicationFirewallPolicyLogScrubbing(input *webapplicationfirewallpolicies.PolicySettingsLogScrubbing) any {
	if input == nil {
		return make([]any, 0)
	}
	result := make(map[string]any)
	result["enabled"] = pointer.From(input.State) == webapplicationfirewallpolicies.WebApplicationFirewallScrubbingStateEnabled
	result["rule"] = flattenWebApplicationFirewallPolicyLogScrubbingRules(input.ScrubbingRules)
	return []any{result}
}

func flattenWebApplicationFirewallPolicyLogScrubbingRules(rules *[]webapplicationfirewallpolicies.WebApplicationFirewallScrubbingRules) any {
	result := make([]any, 0)
	if rules == nil || len(*rules) == 0 {
		return result
	}
	for _, rule := range *rules {
		item := map[string]any{}
		item["enabled"] = pointer.From(rule.State) == webapplicationfirewallpolicies.ScrubbingRuleEntryStateEnabled
		item["match_variable"] = rule.MatchVariable
		item["selector_match_operator"] = rule.SelectorMatchOperator
		item["selector"] = pointer.From(rule.Selector)
		result = append(result, item)
	}
	return &result
}

func flattenWebApplicationFirewallPolicyManagedRulesDefinition(input webapplicationfirewallpolicies.ManagedRulesDefinition) []any {
	results := make([]any, 0)

	v := make(map[string]any)

	v["exception"] = flattenWebApplicationFirewallPolicyExceptions(input.Exceptions)
	v["exclusion"] = flattenWebApplicationFirewallPolicyExclusions(input.Exclusions)
	v["managed_rule_set"] = flattenWebApplicationFirewallPolicyManagedRuleSets(input.ManagedRuleSets)

	results = append(results, v)

	return results
}

func flattenWebApplicationFirewallPolicyExceptionManagedRuleGroups(input *[]webapplicationfirewallpolicies.ExclusionManagedRuleGroup) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}
	for _, item := range *input {
		v := make(map[string]any)

		v["rule_group_name"] = item.RuleGroupName
		v["rules"] = flattenWebApplicationFirewallPolicyManagedRuleIDs(item.Rules)

		results = append(results, v)
	}
	return results
}

func flattenWebApplicationFirewallPolicyExceptionManagedRuleSets(input *[]webapplicationfirewallpolicies.ExclusionManagedRuleSet) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}
	for _, item := range *input {
		v := make(map[string]any)

		v["type"] = item.RuleSetType
		v["version"] = item.RuleSetVersion
		v["rule_group"] = flattenWebApplicationFirewallPolicyExceptionManagedRuleGroups(item.RuleGroups)

		results = append(results, v)
	}
	return results
}

func flattenWebApplicationFirewallPolicyExceptions(input *[]webapplicationfirewallpolicies.ExceptionEntry) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}
	for _, item := range *input {
		v := make(map[string]any)

		v["match_variable"] = string(item.MatchVariable)
		v["selector"] = pointer.From(item.Selector)
		v["selector_match_operator"] = pointer.FromEnum(item.SelectorMatchOperator)
		v["value_match_operator"] = string(item.ValueMatchOperator)
		v["values"] = pointer.From(item.Values)
		v["exception_rule_set"] = flattenWebApplicationFirewallPolicyExceptionManagedRuleSets(item.ExceptionManagedRuleSets)

		results = append(results, v)
	}
	return results
}

func flattenWebApplicationFirewallPolicyManagedRuleIDs(input *[]webapplicationfirewallpolicies.ExclusionManagedRule) []string {
	results := make([]string, 0)
	if input == nil || len(*input) == 0 {
		return results
	}

	for _, item := range *input {
		results = append(results, item.RuleId)
	}

	return results
}

func flattenWebApplicationFirewallPolicyExclusionManagedRuleGroups(input *[]webapplicationfirewallpolicies.ExclusionManagedRuleGroup) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]any)

		v["rule_group_name"] = item.RuleGroupName
		v["excluded_rules"] = flattenWebApplicationFirewallPolicyManagedRuleIDs(item.Rules)

		results = append(results, v)
	}
	return results
}

func flattenWebApplicationFirewallPolicyExclusionManagedRuleSets(input *[]webapplicationfirewallpolicies.ExclusionManagedRuleSet) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]any)

		v["type"] = item.RuleSetType
		v["version"] = item.RuleSetVersion
		v["rule_group"] = flattenWebApplicationFirewallPolicyExclusionManagedRuleGroups(item.RuleGroups)

		results = append(results, v)
	}
	return results
}

func flattenWebApplicationFirewallPolicyExclusions(input *[]webapplicationfirewallpolicies.OwaspCrsExclusionEntry) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]any)

		v["match_variable"] = string(item.MatchVariable)
		v["selector"] = item.Selector

		v["selector_match_operator"] = string(item.SelectorMatchOperator)
		v["excluded_rule_set"] = flattenWebApplicationFirewallPolicyExclusionManagedRuleSets(item.ExclusionManagedRuleSets)

		results = append(results, v)
	}
	return results
}

func flattenWebApplicationFirewallPolicyManagedRuleSets(input []webapplicationfirewallpolicies.ManagedRuleSet) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range input {
		v := make(map[string]any)

		v["type"] = item.RuleSetType
		v["version"] = item.RuleSetVersion
		v["rule_group_override"] = flattenWebApplicationFirewallPolicyRuleGroupOverrides(item.RuleGroupOverrides)

		results = append(results, v)
	}
	return results
}

func flattenWebApplicationFirewallPolicyRuleGroupOverrides(input *[]webapplicationfirewallpolicies.ManagedRuleGroupOverride) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]any)

		v["rule_group_name"] = item.RuleGroupName

		v["rule"] = flattenWebApplicationFirewallPolicyOverrideRules(item.Rules)

		results = append(results, v)
	}
	return results
}

func flattenWebApplicationFirewallPolicyOverrideRules(input *[]webapplicationfirewallpolicies.ManagedRuleOverride) []any {
	results := make([]any, 0)
	if input == nil || len(*input) == 0 {
		return results
	}

	for _, item := range *input {
		v := make(map[string]any)
		v["id"] = item.RuleId

		v["enabled"] = pointer.From(item.State) == webapplicationfirewallpolicies.ManagedRuleEnabledStateEnabled

		v["action"] = pointer.FromEnum(item.Action)

		results = append(results, v)
	}

	return results
}

func flattenWebApplicationFirewallPolicyMatchCondition(input []webapplicationfirewallpolicies.MatchCondition) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range input {
		v := make(map[string]any)

		var transforms []any
		if item.Transforms != nil {
			for _, trans := range *item.Transforms {
				transforms = append(transforms, string(trans))
			}
		}
		v["match_values"] = pluginsdk.FlattenSlice(pointer.To(item.MatchValues))
		v["match_variables"] = flattenWebApplicationFirewallPolicyMatchVariable(item.MatchVariables)
		if negationCondition := item.NegationConditon; negationCondition != nil {
			v["negation_condition"] = *negationCondition
		}
		v["operator"] = string(item.Operator)
		v["transforms"] = transforms

		results = append(results, v)
	}

	return results
}

func flattenWebApplicationFirewallPolicyMatchVariable(input []webapplicationfirewallpolicies.MatchVariable) []any {
	results := make([]any, 0)
	if input == nil {
		return results
	}

	for _, item := range input {
		v := make(map[string]any)

		if selector := item.Selector; selector != nil {
			v["selector"] = *selector
		}
		v["variable_name"] = string(item.VariableName)

		results = append(results, v)
	}

	return results
}

func flattenWebApplicationFirewallPoliciesSubResourcesToIDs(input *[]webapplicationfirewallpolicies.SubResource) []any {
	ids := make([]any, 0)
	if input == nil {
		return ids
	}

	for _, v := range *input {
		if v.Id == nil {
			continue
		}

		ids = append(ids, *v.Id)
	}

	return ids
}
