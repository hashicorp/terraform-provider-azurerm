// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"errors"
	"fmt"

	cdnrules "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-09-01/rules"
	azureValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
)

type CdnFrontDoorRouteConfigurationOverrideInput struct {
	OriginGroupID              string
	ForwardingProtocol         string
	QueryStringCachingBehavior string
	QueryStringParameters      []string
	CompressionEnabled         bool
	CacheBehavior              string
	CacheDuration              string
}

func CdnFrontDoorValidateActionDefinitions(urlRewriteCount, urlRedirectCount, routeConfigurationOverrideCount, totalCount int) error {
	if totalCount == 0 {
		return errors.New("the `actions` block must define at least one action")
	}

	if urlRewriteCount > 1 {
		return fmt.Errorf("the `url_rewrite_action` is only allowed once in the `actions` match block, got %d", urlRewriteCount)
	}

	if urlRedirectCount > 1 {
		return fmt.Errorf("the `url_redirect_action` is only allowed once in the `actions` match block, got %d", urlRedirectCount)
	}

	if routeConfigurationOverrideCount > 1 {
		return fmt.Errorf("the `route_configuration_override_action` is only allowed once in the `actions` match block, got %d", routeConfigurationOverrideCount)
	}

	if urlRedirectCount > 0 && urlRewriteCount > 0 {
		return errors.New("the `url_redirect_action` and the `url_rewrite_action` are both present in the `actions` block which is invalid")
	}

	if totalCount > 5 {
		return fmt.Errorf("the `actions` match block may only contain up to 5 match actions, got %d", totalCount)
	}

	return nil
}

func CdnFrontDoorValidateHeaderAction(blockName, headerAction, value string) error {
	if value == "" {
		if headerAction == string(cdnrules.HeaderActionOverwrite) || headerAction == string(cdnrules.HeaderActionAppend) {
			return fmt.Errorf("the `%s` block is not valid, `value` cannot be empty if the `header_action` is set to `Append` or `Overwrite`", blockName)
		}
	} else if headerAction == string(cdnrules.HeaderActionDelete) {
		return fmt.Errorf("the `%s` block is not valid, `value` must be empty if the `header_action` is set to `Delete`", blockName)
	}

	return nil
}

func CdnFrontDoorValidateConditionMatchValues(configName, operator string, matchValues []string) error {
	if operator == "" {
		return fmt.Errorf("`%s` is invalid: no `operator` value has been set, got `%s`", configName, operator)
	}

	// There are multiple condition-specific `Any` operators in the API surface, but they all
	// resolve to the same literal value. Keep that Azure API quirk documented here so the
	// shared validator preserves the original reasoning from the per-resource implementation.
	if operator == "Any" && len(matchValues) > 0 {
		return fmt.Errorf("`%s` is invalid: the `match_values` field must not be set if the condition `operator` is set to `Any`", configName)
	}

	if operator != "Any" && len(matchValues) == 0 {
		return fmt.Errorf("`%s` is invalid: the `match_values` field must be set if the condition `operator` is not set to `Any`", configName)
	}

	return nil
}

func CdnFrontDoorValidateGeoMatchCountryCodes(configName string, matchValues []string) error {
	for _, matchValue := range matchValues {
		if ok, _ := azureValidate.RegExHelper(matchValue, "match_values", `^[A-Z]{2}$`); !ok {
			return fmt.Errorf("`%s` is invalid: when the `operator` is set to `GeoMatch` the value must be a valid country code consisting of 2 uppercase characters, got `%s`", configName, matchValue)
		}
	}

	return nil
}

func CdnFrontDoorValidateCIDRMatchValues(configName string, matchValues []string) error {
	rawMatchValues := make([]interface{}, 0, len(matchValues))
	for _, matchValue := range matchValues {
		rawMatchValues = append(rawMatchValues, matchValue)
		if _, err := FrontDoorRuleCidrIsValid(matchValue, "match_values"); err != nil {
			return fmt.Errorf("`%s` is invalid: when the `operator` is set to `IPMatch` the `match_values` must be a valid IPv4 or IPv6 CIDR, got `%s`", configName, matchValue)
		}
	}

	if _, err := FrontDoorRuleCidrOverlap(rawMatchValues, "match_values"); err != nil {
		return fmt.Errorf("`%s` is invalid: %+v", configName, err)
	}

	return nil
}

func CdnFrontDoorValidateRouteConfigurationOverrideAction(input CdnFrontDoorRouteConfigurationOverrideInput) error {
	// It is valid to omit the origin-group override entirely, but if no origin group is
	// provided Azure also forbids setting a forwarding protocol. Keep that dependency
	// documented at the shared validation point so it is not lost when callers are refactored.
	if input.OriginGroupID != "" {
		if input.ForwardingProtocol == "" {
			return errors.New("the `route_configuration_override_action` block is not valid, the `forwarding_protocol` field must be set")
		}
	} else if input.ForwardingProtocol != "" {
		return fmt.Errorf("the `route_configuration_override_action` block is not valid, if the `cdn_frontdoor_origin_group_id` is not set you cannot define the `forwarding_protocol`, got `%s`", input.ForwardingProtocol)
	}

	hasCacheConfiguration := input.CacheBehavior != "" || input.QueryStringCachingBehavior != "" || len(input.QueryStringParameters) > 0 || input.CacheDuration != "" || input.CompressionEnabled
	if !hasCacheConfiguration {
		return nil
	}

	if input.CacheBehavior == string(cdnrules.RuleIsCompressionEnabledDisabled) {
		if input.QueryStringCachingBehavior != "" {
			return fmt.Errorf("the `route_configuration_override_action` block is not valid, if the `cache_behavior` is set to `Disabled` you cannot define the `query_string_caching_behavior`, got `%s`", input.QueryStringCachingBehavior)
		}

		if len(input.QueryStringParameters) != 0 {
			return fmt.Errorf("the `route_configuration_override_action` block is not valid, if the `cache_behavior` is set to `Disabled` you cannot define the `query_string_parameters`, got %d", len(input.QueryStringParameters))
		}

		if input.CacheDuration != "" {
			return fmt.Errorf("the `route_configuration_override_action` block is not valid, if the `cache_behavior` is set to `Disabled` you cannot define the `cache_duration`, got `%s`", input.CacheDuration)
		}

		return nil
	}

	if input.CacheBehavior == "" {
		return errors.New("the `route_configuration_override_action` block is not valid, the `cache_behavior` field must be set")
	}

	if input.QueryStringCachingBehavior == "" {
		return errors.New("the `route_configuration_override_action` block is not valid, the `query_string_caching_behavior` field must be set")
	}

	// `HonorOrigin` must not carry an explicit cache duration. Keep this service quirk
	// documented at the shared validation point so the history from issue #19311 stays
	// attached to the behavior after refactors.
	if input.CacheBehavior != string(cdnrules.RuleCacheBehaviorHonorOrigin) {
		if input.CacheDuration == "" {
			return errors.New("the `route_configuration_override_action` block is not valid, the `cache_duration` field must be set")
		}
	} else if input.CacheDuration != "" {
		return errors.New("the `route_configuration_override_action` block is not valid, the `cache_duration` field must not be set if the `cache_behavior` is `HonorOrigin`")
	}

	if (input.QueryStringCachingBehavior == string(cdnrules.RuleQueryStringCachingBehaviorIncludeSpecifiedQueryStrings) || input.QueryStringCachingBehavior == string(cdnrules.RuleQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings)) && len(input.QueryStringParameters) == 0 {
		return errors.New("the `route_configuration_override_action` block is not valid, `query_string_parameters` cannot be empty if the `query_string_caching_behavior` is set to `IncludeSpecifiedQueryStrings` or `IgnoreSpecifiedQueryStrings`")
	}

	if len(input.QueryStringParameters) > 0 && (input.QueryStringCachingBehavior == string(cdnrules.RuleQueryStringCachingBehaviorUseQueryString) || input.QueryStringCachingBehavior == string(cdnrules.RuleQueryStringCachingBehaviorIgnoreQueryString)) {
		return errors.New("the `route_configuration_override_action` block is not valid, `query_string_parameters` must not be set if the `query_string_caching_behavior` is set to `UseQueryString` or `IgnoreQueryString`")
	}

	return nil
}
