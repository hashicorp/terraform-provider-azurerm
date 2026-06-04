// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"strings"
	"testing"

	cdnrules "github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-09-01/rules"
)

func TestCdnFrontDoorValidateActionDefinitions(t *testing.T) {
	tests := []struct {
		name        string
		urlRewrite  int
		urlRedirect int
		routeConfig int
		total       int
		errContains string
	}{
		{
			name:  "valid counts",
			total: 5,
		},
		{
			name:        "duplicate url rewrite",
			urlRewrite:  2,
			total:       2,
			errContains: "url_rewrite_action",
		},
		{
			name:        "duplicate url redirect",
			urlRedirect: 2,
			total:       2,
			errContains: "url_redirect_action",
		},
		{
			name:        "duplicate route configuration override",
			routeConfig: 2,
			total:       2,
			errContains: "route_configuration_override_action",
		},
		{
			name:        "rewrite and redirect together",
			urlRewrite:  1,
			urlRedirect: 1,
			total:       2,
			errContains: "both present",
		},
		{
			name:        "too many actions",
			total:       6,
			errContains: "up to 5 match actions",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CdnFrontDoorValidateActionDefinitions(test.urlRewrite, test.urlRedirect, test.routeConfig, test.total)
			assertValidationError(t, err, test.errContains)
		})
	}
}

func TestCdnFrontDoorValidateHeaderAction(t *testing.T) {
	tests := []struct {
		name         string
		blockName    string
		headerAction string
		value        string
		errContains  string
	}{
		{
			name:         "append requires value",
			blockName:    "request_header_action",
			headerAction: string(cdnrules.HeaderActionAppend),
			errContains:  "cannot be empty",
		},
		{
			name:         "overwrite requires value",
			blockName:    "response_header_action",
			headerAction: string(cdnrules.HeaderActionOverwrite),
			errContains:  "cannot be empty",
		},
		{
			name:         "delete must not set value",
			blockName:    "request_header_action",
			headerAction: string(cdnrules.HeaderActionDelete),
			value:        "x-test",
			errContains:  "must be empty",
		},
		{
			name:         "append with value valid",
			blockName:    "request_header_action",
			headerAction: string(cdnrules.HeaderActionAppend),
			value:        "x-test",
		},
		{
			name:         "delete without value valid",
			blockName:    "response_header_action",
			headerAction: string(cdnrules.HeaderActionDelete),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CdnFrontDoorValidateHeaderAction(test.blockName, test.headerAction, test.value)
			assertValidationError(t, err, test.errContains)
		})
	}
}

func TestCdnFrontDoorValidateConditionMatchValues(t *testing.T) {
	tests := []struct {
		name        string
		configName  string
		operator    string
		matchValues []string
		errContains string
	}{
		{
			name:        "empty operator invalid",
			configName:  "request_header_condition",
			errContains: "no 'operator' value has been set",
		},
		{
			name:        "any with values invalid",
			configName:  "request_header_condition",
			operator:    "Any",
			matchValues: []string{"x-test"},
			errContains: "must not be set",
		},
		{
			name:        "non any requires values",
			configName:  "request_header_condition",
			operator:    "Equal",
			errContains: "must be set",
		},
		{
			name:        "equal with values valid",
			configName:  "request_header_condition",
			operator:    "Equal",
			matchValues: []string{"x-test"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CdnFrontDoorValidateConditionMatchValues(test.configName, test.operator, test.matchValues)
			assertValidationError(t, err, test.errContains)
		})
	}
}

func TestCdnFrontDoorValidateGeoMatchCountryCodes(t *testing.T) {
	tests := []struct {
		name        string
		matchValues []string
		errContains string
	}{
		{
			name:        "valid country codes",
			matchValues: []string{"US", "DE"},
		},
		{
			name:        "lowercase country code invalid",
			matchValues: []string{"us"},
			errContains: "valid country code",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CdnFrontDoorValidateGeoMatchCountryCodes("remote_address_condition", test.matchValues)
			assertValidationError(t, err, test.errContains)
		})
	}
}

func TestCdnFrontDoorValidateCIDRMatchValues(t *testing.T) {
	tests := []struct {
		name        string
		matchValues []string
		errContains string
	}{
		{
			name:        "valid CIDRs",
			matchValues: []string{"192.168.0.1/24", "192.168.1.1/24"},
		},
		{
			name:        "invalid CIDR",
			matchValues: []string{"192.168.0.1:80/24"},
			errContains: "must be a valid IPv4 or IPv6 CIDR",
		},
		{
			name:        "overlapping CIDRs",
			matchValues: []string{"192.168.0.1/24", "192.168.0.1/26"},
			errContains: "invalid",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CdnFrontDoorValidateCIDRMatchValues("remote_address_condition", test.matchValues)
			assertValidationError(t, err, test.errContains)
		})
	}
}

func TestCdnFrontDoorValidateRouteConfigurationOverrideAction(t *testing.T) {
	tests := []struct {
		name        string
		input       CdnFrontDoorRouteConfigurationOverrideInput
		errContains string
	}{
		{
			name:  "empty block valid",
			input: CdnFrontDoorRouteConfigurationOverrideInput{},
		},
		{
			name: "origin group requires forwarding protocol",
			input: CdnFrontDoorRouteConfigurationOverrideInput{
				OriginGroupID: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Cdn/profiles/p/originGroups/g",
			},
			errContains: "forwarding_protocol' field must be set",
		},
		{
			name: "forwarding protocol requires origin group",
			input: CdnFrontDoorRouteConfigurationOverrideInput{
				ForwardingProtocol: string(cdnrules.ForwardingProtocolHTTPSOnly),
			},
			errContains: "cannot define the 'forwarding_protocol'",
		},
		{
			name: "disabled cache cannot define query string behavior",
			input: CdnFrontDoorRouteConfigurationOverrideInput{
				CacheBehavior:              string(cdnrules.RuleIsCompressionEnabledDisabled),
				QueryStringCachingBehavior: string(cdnrules.RuleQueryStringCachingBehaviorUseQueryString),
			},
			errContains: "cache_behavior' is set to 'Disabled'",
		},
		{
			name: "missing cache behavior invalid",
			input: CdnFrontDoorRouteConfigurationOverrideInput{
				QueryStringCachingBehavior: string(cdnrules.RuleQueryStringCachingBehaviorUseQueryString),
			},
			errContains: "cache_behavior' field must be set",
		},
		{
			name: "missing query string behavior invalid",
			input: CdnFrontDoorRouteConfigurationOverrideInput{
				CacheBehavior: string(cdnrules.RuleCacheBehaviorOverrideAlways),
			},
			errContains: "query_string_caching_behavior' field must be set",
		},
		{
			name: "honor origin cannot define cache duration",
			input: CdnFrontDoorRouteConfigurationOverrideInput{
				CacheBehavior:              string(cdnrules.RuleCacheBehaviorHonorOrigin),
				QueryStringCachingBehavior: string(cdnrules.RuleQueryStringCachingBehaviorUseQueryString),
				CacheDuration:              "00:10:00",
			},
			errContains: "must not be set if the 'cache_behavior' is 'HonorOrigin'",
		},
		{
			name: "include specified requires parameters",
			input: CdnFrontDoorRouteConfigurationOverrideInput{
				CacheBehavior:              string(cdnrules.RuleCacheBehaviorOverrideAlways),
				QueryStringCachingBehavior: string(cdnrules.RuleQueryStringCachingBehaviorIncludeSpecifiedQueryStrings),
				CacheDuration:              "00:10:00",
			},
			errContains: "query_string_parameters' cannot be empty",
		},
		{
			name: "use query strings forbids parameters",
			input: CdnFrontDoorRouteConfigurationOverrideInput{
				CacheBehavior:              string(cdnrules.RuleCacheBehaviorOverrideAlways),
				QueryStringCachingBehavior: string(cdnrules.RuleQueryStringCachingBehaviorUseQueryString),
				QueryStringParameters:      []string{"foo"},
				CacheDuration:              "00:10:00",
			},
			errContains: "must not be set if the'query_string_caching_behavior'",
		},
		{
			name: "valid cache override",
			input: CdnFrontDoorRouteConfigurationOverrideInput{
				OriginGroupID:              "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test/providers/Microsoft.Cdn/profiles/p/originGroups/g",
				ForwardingProtocol:         string(cdnrules.ForwardingProtocolHTTPSOnly),
				CacheBehavior:              string(cdnrules.RuleCacheBehaviorOverrideAlways),
				QueryStringCachingBehavior: string(cdnrules.RuleQueryStringCachingBehaviorUseQueryString),
				CacheDuration:              "00:10:00",
				CompressionEnabled:         true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CdnFrontDoorValidateRouteConfigurationOverrideAction(test.input)
			assertValidationError(t, err, test.errContains)
		})
	}
}

func assertValidationError(t *testing.T, err error, errContains string) {
	t.Helper()

	if errContains == "" {
		if err != nil {
			t.Fatalf("expected no error but got %q", err)
		}
		return
	}

	if err == nil {
		t.Fatalf("expected error containing %q but got nil", errContains)
	}

	if !strings.Contains(err.Error(), errContains) {
		t.Fatalf("expected error containing %q but got %q", errContains, err)
	}
}
