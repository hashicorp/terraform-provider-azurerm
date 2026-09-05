// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rulesets"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	cdnValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// This file contains common models, schema, and functions used by both `azurerm_cdn_frontdoor_rule` and `azurerm_cdn_frontdoor_batch_rule_set`.

const (
	cdnFrontDoorRuleConditionOperatorNegatedPrefix = "Not"

	RuleCacheBehaviorDisabled = "Disabled"
)

func PossibleValuesForRuleCacheBehavior() []string {
	return []string{
		string(rules.RuleCacheBehaviorHonorOrigin),
		string(rules.RuleCacheBehaviorOverrideAlways),
		string(rules.RuleCacheBehaviorOverrideIfOriginMissing),
		// Exposed `Disabled` as a valid value for provider issue #19008.
		RuleCacheBehaviorDisabled,
	}
}

// Schema

func cdnFrontDoorRuleValuesRequiredSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MinItems: 1,
		MaxItems: 25,
		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func cdnFrontDoorRuleValuesOptionalSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 25,
		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func cdnFrontDoorRuleConditionNameSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:         pluginsdk.TypeString,
		Required:     true,
		ValidateFunc: validation.StringIsNotEmpty,
	}
}

func cdnFrontDoorRuleHttpVersionValuesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Required: true,
		MinItems: 1,
		MaxItems: 4,
		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice([]string{"2.0", "1.1", "1.0", "0.9"}, false),
		},
	}
}

func cdnFrontDoorServerPortValuesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		MaxItems: 2,
		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice([]string{"80", "443"}, false),
		},
	}
}

func cdnFrontDoorSocketAddressValuesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MinItems: 1,
		MaxItems: 25,
		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: cdnValidate.FrontDoorRuleCidrIsValid,
		},
	}
}

// Models

type CdnFrontDoorRuleActionsModel struct {
	URLRedirect                []CdnFrontDoorRuleURLRedirectActionModel                `tfschema:"url_redirect"`
	URLRewrite                 []CdnFrontDoorRuleURLRewriteActionModel                 `tfschema:"url_rewrite"`
	ModifyRequestHeader        []CdnFrontDoorRuleHeaderActionModel                     `tfschema:"modify_request_header"`
	ModifyResponseHeader       []CdnFrontDoorRuleHeaderActionModel                     `tfschema:"modify_response_header"`
	RouteConfigurationOverride []CdnFrontDoorRuleRouteConfigurationOverrideActionModel `tfschema:"route_configuration_override"`
}

type CdnFrontDoorRuleURLRedirectActionModel struct {
	RedirectType        string `tfschema:"redirect_type"`
	RedirectProtocol    string `tfschema:"redirect_protocol"`
	DestinationPath     string `tfschema:"destination_path"`
	DestinationHostName string `tfschema:"destination_host_name"`
	QueryString         string `tfschema:"query_string"`
	DestinationFragment string `tfschema:"destination_fragment"`
}

type CdnFrontDoorRuleURLRewriteActionModel struct {
	SourcePattern                string `tfschema:"source_pattern"`
	DestinationPath              string `tfschema:"destination_path"`
	PreserveUnmatchedPathEnabled bool   `tfschema:"preserve_unmatched_path_enabled"`
}

type CdnFrontDoorRuleHeaderActionModel struct {
	Operator    string `tfschema:"operator"`
	HeaderName  string `tfschema:"header_name"`
	HeaderValue string `tfschema:"header_value"`
}

type CdnFrontDoorRuleRouteConfigurationOverrideActionModel struct {
	OriginGroup []CdnFrontDoorRuleRouteConfigurationOverrideOriginGroupModel `tfschema:"origin_group"`
	Caching     []CdnFrontDoorRuleRouteConfigurationOverrideCachingModel     `tfschema:"caching"`
}

type CdnFrontDoorRuleRouteConfigurationOverrideOriginGroupModel struct {
	CdnFrontdoorOriginGroupID string `tfschema:"cdn_frontdoor_origin_group_id"`
	ForwardingProtocol        string `tfschema:"forwarding_protocol"`
}

type CdnFrontDoorRuleRouteConfigurationOverrideCachingModel struct {
	Behaviour             string   `tfschema:"behaviour"`
	Duration              string   `tfschema:"duration"`
	CompressionEnabled    bool     `tfschema:"compression_enabled"`
	QueryStringBehaviour  string   `tfschema:"query_string_behaviour"`
	QueryStringParameters []string `tfschema:"query_string_parameters"`
}

type CdnFrontDoorRuleConditionsModel struct {
	RemoteAddress        []CdnFrontDoorRuleConditionBaseModel                  `tfschema:"remote_address"`
	RequestMethod        []CdnFrontDoorRuleConditionBaseModel                  `tfschema:"request_method"`
	QueryString          []CdnFrontDoorRuleConditionWithTransformsModel        `tfschema:"query_string"`
	PostArgs             []CdnFrontDoorRuleConditionWithNameAndTransformsModel `tfschema:"post_argument"`
	RequestURL           []CdnFrontDoorRuleConditionWithTransformsModel        `tfschema:"request_url"`
	RequestHeader        []CdnFrontDoorRuleConditionWithNameAndTransformsModel `tfschema:"request_header"`
	RequestBody          []CdnFrontDoorRuleConditionWithTransformsModel        `tfschema:"request_body"`
	RequestScheme        []CdnFrontDoorRuleConditionBaseModel                  `tfschema:"request_scheme"`
	RequestPath          []CdnFrontDoorRuleConditionWithTransformsModel        `tfschema:"request_path"`
	RequestFileExtension []CdnFrontDoorRuleConditionWithTransformsModel        `tfschema:"request_file_extension"`
	RequestFilename      []CdnFrontDoorRuleConditionWithTransformsModel        `tfschema:"request_filename"`
	HTTPVersion          []CdnFrontDoorRuleConditionBaseModel                  `tfschema:"http_version"`
	RequestCookies       []CdnFrontDoorRuleConditionWithNameAndTransformsModel `tfschema:"request_cookies"`
	DeviceType           []CdnFrontDoorRuleConditionBaseModel                  `tfschema:"device_type"`
	SocketAddress        []CdnFrontDoorRuleConditionBaseModel                  `tfschema:"socket_address"`
	ClientPort           []CdnFrontDoorRuleConditionBaseModel                  `tfschema:"client_port"`
	ServerPort           []CdnFrontDoorRuleConditionBaseModel                  `tfschema:"server_port"`
	HostName             []CdnFrontDoorRuleConditionWithTransformsModel        `tfschema:"host_name"`
	SSLProtocol          []CdnFrontDoorRuleConditionBaseModel                  `tfschema:"ssl_protocol"`
}

type CdnFrontDoorRuleConditionBaseModel struct {
	Operator string   `tfschema:"operator"`
	Values   []string `tfschema:"values"`
}

func flattenCdnFrontDoorRuleConditionBaseModel(operator string, values []string, negated *bool) CdnFrontDoorRuleConditionBaseModel {
	return CdnFrontDoorRuleConditionBaseModel{
		Operator: flattenCdnFrontDoorRuleConditionOperator(operator, pointer.From(negated)),
		Values:   values,
	}
}

type CdnFrontDoorRuleConditionWithTransformsModel struct {
	Operator   string   `tfschema:"operator"`
	Values     []string `tfschema:"values"`
	Transforms []string `tfschema:"transforms"`
}

func flattenCdnFrontDoorBatchRuleSetConditionWithTransformsModel(operator string, values, transforms []string, negated *bool) CdnFrontDoorRuleConditionWithTransformsModel {
	return CdnFrontDoorRuleConditionWithTransformsModel{
		Operator:   flattenCdnFrontDoorRuleConditionOperator(operator, pointer.From(negated)),
		Values:     values,
		Transforms: transforms,
	}
}

type CdnFrontDoorRuleConditionWithNameAndTransformsModel struct {
	Name       string   `tfschema:"name"`
	Operator   string   `tfschema:"operator"`
	Values     []string `tfschema:"values"`
	Transforms []string `tfschema:"transforms"`
}

func flattenCdnFrontDoorRuleConditionWithNameAndTransformsModel(operator string, name *string, values, transforms []string, negated *bool) CdnFrontDoorRuleConditionWithNameAndTransformsModel {
	return CdnFrontDoorRuleConditionWithNameAndTransformsModel{
		Name:       pointer.From(name),
		Operator:   flattenCdnFrontDoorRuleConditionOperator(operator, pointer.From(negated)),
		Values:     values,
		Transforms: transforms,
	}
}

// Non-model specific expand/flatten functions

func expandCdnFrontDoorRuleConditionOperator(input string) (string, bool) {
	negated := false
	if strings.HasPrefix(input, cdnFrontDoorRuleConditionOperatorNegatedPrefix) {
		negated = true
		input = strings.TrimPrefix(input, cdnFrontDoorRuleConditionOperatorNegatedPrefix)
	}
	return input, negated
}

func flattenCdnFrontDoorRuleConditionOperator(input string, negated bool) string {
	result := input
	if negated {
		result = cdnFrontDoorRuleConditionOperatorNegatedPrefix + result
	}
	return result
}

// Validation

func validateCdnFrontDoorRuleName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if m, _ := validate.RegExHelper(i, k, `^[a-zA-Z][\da-zA-Z]{0,259}$`); !m {
		return nil, []error{fmt.Errorf(`%q must be between 1 and 260 characters in length, begin with a letter and may contain only letters and numbers, got %q`, k, v)}
	}

	return nil, nil
}

func validateCdnFrontDoorRuleActionCounts(urlRewriteCount, urlRedirectCount, routeConfigurationOverrideCount, totalCount int) error {
	if totalCount == 0 {
		return errors.New("the `actions` block must define at least one action")
	}

	if urlRedirectCount > 0 && urlRewriteCount > 0 {
		return errors.New("cannot specify both `url_redirect` and the `url_rewrite` in the `actions` block")
	}

	if urlRedirectCount > 0 && routeConfigurationOverrideCount > 0 {
		return errors.New("cannot specify both `url_redirect` and the `route_configuration_override` in the `actions` block")
	}

	if totalCount > 5 {
		return fmt.Errorf("the `actions` block may only contain up to 5 actions, got %d", totalCount)
	}

	return nil
}

func validateCdnFrontDoorRuleConditionValues(configName, operator string, matchValues []string) error {
	if operator == "" {
		return fmt.Errorf("`%s` is invalid: no `operator` value has been set, got `%s`", configName, operator)
	}

	// There are multiple condition-specific `Any` operators in the API surface, but they all
	// resolve to the same literal value.
	if operator == "Any" && len(matchValues) > 0 {
		return fmt.Errorf("when `conditions.%[1]s.operator` is set to `Any`, `conditions.%[1]s.values` cannot be defined", configName)
	}

	if operator != "Any" && len(matchValues) == 0 {
		return fmt.Errorf("when `conditions.%[1]s.operator` is set to `%[2]s`, `conditions.%[1]s.values` must set one or more values", configName, operator)
	}

	return nil
}

func validateCdnFrontDoorRuleModifyHeaderAction(blockName, headerAction, value string) error {
	if value == "" {
		if headerAction == string(rulesets.HeaderActionOverwrite) || headerAction == string(rulesets.HeaderActionAppend) {
			return fmt.Errorf("the `%s` block is not valid, `header_value` cannot be empty if the `operator` is set to `Append` or `Overwrite`", blockName)
		}
	} else if headerAction == string(rulesets.HeaderActionDelete) {
		return fmt.Errorf("the `%s` block is not valid, `header_value` must be empty if the `operator` is set to `Delete`", blockName)
	}

	return nil
}

func validateCdnFrontDoorUrlRedirectActionQueryString(i interface{}, k string) (_ []string, errors []error) {
	// Query string must be in <key>=<value> format. ? and & will be added automatically so do not include them.
	// NOTE: the 2048 character limit matches the service code validation logic for this field
	return validation.All(
		validation.StringIsNotEmpty,
		validation.StringDoesNotMatch(regexp.MustCompile(`^\?`), "must not start with the '?' character, it will be automatically added by Frontdoor"),
		validation.StringLenBetween(1, 2048),
	)(i, k)
}

func validateCdnFrontDoorCacheDuration(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if strings.HasPrefix(v, "0.") {
		return nil, []error{fmt.Errorf(`%q must not begin with %q if the duration is less than 1 day. If the %q is less than 1 day it should be in the HH:MM:SS format, got %q`, k, "0.", k, v)}
	}

	// fix for issue #22668
	durationParts := strings.Split(v, ".")

	if len(durationParts) > 1 {
		days, err := strconv.Atoi(durationParts[0])
		if err != nil {
			return nil, []error{fmt.Errorf(`%q 'days' segment is invalid, the 'days' segment must be a valid number and have a value that is between 1 and 365, got %q`, k, v)}
		}

		if days > 365 {
			return nil, []error{fmt.Errorf(`%q must be in the d.HH:MM:SS or HH:MM:SS format and must be equal to or lower than %q, got %q`, k, "365.23:59:59", v)}
		}
	}

	// the old regular expression was broken because it wouldn't allow the value in the tens
	// position to be greater than 6 and the ones position greater than 5
	if m, _ := validate.RegExHelper(i, k, `^([1-9]|([1-9][0-9])|([1-3][0-9][0-9])).((?:[01]\d|2[0123]):(?:[012345]\d):(?:[012345]\d))$|^((?:[01]\d|2[0123]):(?:[012345]\d):(?:[012345]\d))$`); !m {
		return nil, []error{fmt.Errorf(`%q must be in the d.HH:MM:SS or HH:MM:SS format and must be equal to or lower than %q, got %q`, k, "365.23:59:59", v)}
	}

	return nil, nil
}

// Helpers

func cdnFrontDoorRuleConditionOperatorPossibleValues(values []string) []string {
	result := make([]string, 0, len(values)*2)
	for _, value := range values {
		// For each SDK constant, we'll add a negated version in favour of exposing another property
		// this matches portal more closely, and avoids exposing a number of `operator` fields with only a single allowed value.
		result = append(result, value, cdnFrontDoorRuleConditionOperatorNegatedPrefix+value)
	}
	return result
}

func countCdnFrontDoorBatchRuleActions(input []CdnFrontDoorRuleActionsModel) (urlRewrite, urlRedirect, routeConfigurationOverride, total int) {
	if len(input) == 0 {
		return
	}

	actions := input[0]
	urlRewrite = len(actions.URLRewrite)
	urlRedirect = len(actions.URLRedirect)
	routeConfigurationOverride = len(actions.RouteConfigurationOverride)
	total = urlRewrite + urlRedirect + routeConfigurationOverride + len(actions.ModifyRequestHeader) + len(actions.ModifyResponseHeader)

	return
}
