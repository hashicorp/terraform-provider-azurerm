// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package frontdooractions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-09-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	cdnvalidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontDoorActionParameters struct {
	Name       rules.DeliveryRuleActionName
	TypeName   string
	ConfigName string
}

type CdnFrontDoorActionMappings struct {
	RouteConfigurationOverride CdnFrontDoorActionParameters
	RequestHeader              CdnFrontDoorActionParameters
	ResponseHeader             CdnFrontDoorActionParameters
	URLRedirect                CdnFrontDoorActionParameters
	URLRewrite                 CdnFrontDoorActionParameters
}

func InitializeCdnFrontDoorActionMappings() *CdnFrontDoorActionMappings {
	m := CdnFrontDoorActionMappings{}

	m.RouteConfigurationOverride = CdnFrontDoorActionParameters{
		Name:       rules.DeliveryRuleActionNameRouteConfigurationOverride,
		TypeName:   string(rules.DeliveryRuleActionParametersTypeDeliveryRuleRouteConfigurationOverrideActionParameters),
		ConfigName: "route_configuration_override_action",
	}

	m.RequestHeader = CdnFrontDoorActionParameters{
		Name:       rules.DeliveryRuleActionNameModifyRequestHeader,
		TypeName:   string(rules.DeliveryRuleActionParametersTypeDeliveryRuleHeaderActionParameters),
		ConfigName: "request_header_action",
	}

	m.ResponseHeader = CdnFrontDoorActionParameters{
		Name:       rules.DeliveryRuleActionNameModifyResponseHeader,
		TypeName:   string(rules.DeliveryRuleActionParametersTypeDeliveryRuleHeaderActionParameters),
		ConfigName: "response_header_action",
	}

	m.URLRedirect = CdnFrontDoorActionParameters{
		Name:       rules.DeliveryRuleActionNameURLRedirect,
		TypeName:   string(rules.DeliveryRuleActionParametersTypeDeliveryRuleURLRedirectActionParameters),
		ConfigName: "url_redirect_action",
	}

	m.URLRewrite = CdnFrontDoorActionParameters{
		Name:       rules.DeliveryRuleActionNameURLRewrite,
		TypeName:   string(rules.DeliveryRuleActionParametersTypeDeliveryRuleURLRewriteActionParameters),
		ConfigName: "url_rewrite_action",
	}

	return &m
}

func expandStringSliceToCsvFormat(input []interface{}) *string {
	if len(input) == 0 {
		return nil
	}

	v := utils.ExpandStringSlice(input)
	csv := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(*v)), ","), "[]")

	return &csv
}

func flattenCsvToStringSlice(input *string) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	v := strings.Split(*input, ",")

	for _, s := range v {
		results = append(results, s)
	}

	return results
}

func ExpandCdnFrontDoorRequestHeaderAction(input []interface{}) (*[]rules.DeliveryRuleAction, error) {
	output := make([]rules.DeliveryRuleAction, 0)

	m := InitializeCdnFrontDoorActionMappings()

	for _, v := range input {
		item := v.(map[string]interface{})

		value := item["value"].(string)

		requestHeaderAction := rules.DeliveryRuleRequestHeaderAction{
			Name: m.RequestHeader.Name,
			Parameters: rules.HeaderActionParameters{
				TypeName:     rules.DeliveryRuleActionParametersType(m.RequestHeader.TypeName),
				HeaderAction: rules.HeaderAction(item["header_action"].(string)),
				HeaderName:   item["header_name"].(string),
				Value:        pointer.To(value),
			},
		}

		if err := cdnvalidate.CdnFrontDoorValidateHeaderAction("request_header_action", item["header_action"].(string), value); err != nil {
			return nil, err
		}

		output = append(output, requestHeaderAction)
	}

	return &output, nil
}

func ExpandCdnFrontDoorResponseHeaderAction(input []interface{}) (*[]rules.DeliveryRuleAction, error) {
	output := make([]rules.DeliveryRuleAction, 0)

	m := InitializeCdnFrontDoorActionMappings()

	for _, v := range input {
		item := v.(map[string]interface{})

		responseHeaderAction := rules.DeliveryRuleResponseHeaderAction{
			Name: m.ResponseHeader.Name,
			Parameters: rules.HeaderActionParameters{
				TypeName:     rules.DeliveryRuleActionParametersType(m.ResponseHeader.TypeName),
				HeaderAction: rules.HeaderAction(item["header_action"].(string)),
				HeaderName:   item["header_name"].(string),
				Value:        pointer.To(item["value"].(string)),
			},
		}

		if err := cdnvalidate.CdnFrontDoorValidateHeaderAction("response_header_action", item["header_action"].(string), pointer.From(responseHeaderAction.Parameters.Value)); err != nil {
			return nil, err
		}

		output = append(output, responseHeaderAction)
	}

	return &output, nil
}

func ExpandCdnFrontDoorUrlRedirectAction(input []interface{}) (*[]rules.DeliveryRuleAction, error) {
	output := make([]rules.DeliveryRuleAction, 0)

	m := InitializeCdnFrontDoorActionMappings()

	for _, v := range input {
		item := v.(map[string]interface{})

		urlRedirectAction := rules.URLRedirectAction{
			Name: m.URLRedirect.Name,
			Parameters: rules.URLRedirectActionParameters{
				TypeName:            rules.DeliveryRuleActionParametersType(m.URLRedirect.TypeName),
				RedirectType:        rules.RedirectType(item["redirect_type"].(string)),
				DestinationProtocol: pointer.To(rules.DestinationProtocol(item["redirect_protocol"].(string))),
				CustomPath:          pointer.To(item["destination_path"].(string)),
				CustomHostname:      pointer.To(item["destination_hostname"].(string)),
				CustomQueryString:   pointer.To(item["query_string"].(string)),
				CustomFragment:      pointer.To(item["destination_fragment"].(string)),
			},
		}

		output = append(output, urlRedirectAction)
	}

	return &output, nil
}

func ExpandCdnFrontDoorUrlRewriteAction(input []interface{}) (*[]rules.DeliveryRuleAction, error) {
	output := make([]rules.DeliveryRuleAction, 0)

	m := InitializeCdnFrontDoorActionMappings()

	for _, v := range input {
		item := v.(map[string]interface{})

		urlRedirectAction := rules.URLRewriteAction{
			Name: m.URLRewrite.Name,
			Parameters: rules.URLRewriteActionParameters{
				TypeName:              rules.DeliveryRuleActionParametersType(m.URLRewrite.TypeName),
				Destination:           item["destination"].(string),
				PreserveUnmatchedPath: pointer.To(item["preserve_unmatched_path"].(bool)),
				SourcePattern:         item["source_pattern"].(string),
			},
		}

		output = append(output, urlRedirectAction)
	}

	return &output, nil
}

func ExpandCdnFrontDoorRouteConfigurationOverrideAction(input []interface{}) (*[]rules.DeliveryRuleAction, error) {
	output := make([]rules.DeliveryRuleAction, 0)
	m := InitializeCdnFrontDoorActionMappings()

	for _, v := range input {
		item := v.(map[string]interface{})
		queryStringParameters := utils.ExpandStringSlice(item["query_string_parameters"].([]interface{}))

		validationInput := cdnvalidate.CdnFrontDoorRouteConfigurationOverrideInput{
			OriginGroupID:              item["cdn_frontdoor_origin_group_id"].(string),
			ForwardingProtocol:         item["forwarding_protocol"].(string),
			QueryStringCachingBehavior: item["query_string_caching_behavior"].(string),
			QueryStringParameters:      pointer.From(queryStringParameters),
			CompressionEnabled:         item["compression_enabled"].(bool),
			CacheBehavior:              item["cache_behavior"].(string),
			CacheDuration:              item["cache_duration"].(string),
		}
		if err := cdnvalidate.CdnFrontDoorValidateRouteConfigurationOverrideAction(validationInput); err != nil {
			return nil, err
		}

		var originGroupOverride rules.OriginGroupOverride
		var cacheConfiguration rules.CacheConfiguration

		originGroupIdRaw := item["cdn_frontdoor_origin_group_id"].(string)
		protocol := item["forwarding_protocol"].(string)
		cacheBehavior := item["cache_behavior"].(string)
		compressionEnabled := rules.RuleIsCompressionEnabledEnabled
		queryStringCachingBehavior := item["query_string_caching_behavior"].(string)
		cacheDuration := item["cache_duration"].(string)

		if !item["compression_enabled"].(bool) {
			compressionEnabled = rules.RuleIsCompressionEnabledDisabled
		}

		if originGroupIdRaw != "" {
			originGroupOverride = rules.OriginGroupOverride{
				OriginGroup: &rules.ResourceReference{
					Id: pointer.To(originGroupIdRaw),
				},
				ForwardingProtocol: pointer.To(rules.ForwardingProtocol(protocol)),
			}
		}

		if cacheBehavior != string(rules.RuleIsCompressionEnabledDisabled) {
			cacheConfiguration = rules.CacheConfiguration{
				QueryStringCachingBehavior: pointer.To(rules.RuleQueryStringCachingBehavior(queryStringCachingBehavior)),
				QueryParameters:            expandStringSliceToCsvFormat(item["query_string_parameters"].([]interface{})),
				IsCompressionEnabled:       pointer.To(compressionEnabled),
				CacheBehavior:              pointer.To(rules.RuleCacheBehavior(cacheBehavior)),
			}

			if cacheDuration != "" {
				cacheConfiguration.CacheDuration = pointer.To(cacheDuration)
			}
		}

		routeConfigurationOverrideAction := rules.DeliveryRuleRouteConfigurationOverrideAction{
			Parameters: rules.RouteConfigurationOverrideActionParameters{
				TypeName: rules.DeliveryRuleActionParametersType(m.RouteConfigurationOverride.TypeName),
			},
		}

		if originGroupOverride.OriginGroup != nil {
			routeConfigurationOverrideAction.Parameters.OriginGroupOverride = &originGroupOverride
		}

		if cacheConfiguration.CacheDuration != nil || pointer.From(cacheConfiguration.CacheBehavior) == rules.RuleCacheBehaviorHonorOrigin {
			routeConfigurationOverrideAction.Parameters.CacheConfiguration = &cacheConfiguration
		}

		output = append(output, routeConfigurationOverrideAction)
	}

	return &output, nil
}

func FlattenRequestHeaderAction(input rules.DeliveryRuleRequestHeaderAction) map[string]interface{} {
	var value string

	params := input.Parameters
	action := string(params.HeaderAction)
	name := params.HeaderName

	if params.Value != nil {
		value = *params.Value
	}

	return map[string]interface{}{
		"header_action": action,
		"header_name":   name,
		"value":         value,
	}
}

func FlattenResponseHeaderAction(input rules.DeliveryRuleResponseHeaderAction) map[string]interface{} {
	var value string

	params := input.Parameters
	action := string(params.HeaderAction)
	name := params.HeaderName

	if params.Value != nil {
		value = *params.Value
	}

	return map[string]interface{}{
		"header_action": action,
		"header_name":   name,
		"value":         value,
	}
}

func FlattenCdnFrontDoorUrlRedirectAction(input rules.URLRedirectAction) map[string]interface{} {
	var destinationHost string
	var destinationPath string
	var queryString string
	var fragment string

	params := input.Parameters

	if params.CustomHostname != nil {
		destinationHost = *params.CustomHostname
	}
	if params.CustomPath != nil {
		destinationPath = *params.CustomPath
	}
	if params.CustomQueryString != nil {
		queryString = *params.CustomQueryString
	}

	destinationProtocol := string(pointer.From(params.DestinationProtocol))
	redirectType := string(params.RedirectType)

	if params.CustomFragment != nil {
		fragment = *params.CustomFragment
	}

	return map[string]interface{}{
		"destination_hostname": destinationHost,
		"destination_path":     destinationPath,
		"query_string":         queryString,
		"redirect_protocol":    destinationProtocol,
		"redirect_type":        redirectType,
		"destination_fragment": fragment,
	}
}

func FlattenCdnFrontDoorUrlRewriteAction(input rules.URLRewriteAction) map[string]interface{} {
	params := input.Parameters

	destination := params.Destination
	preservePath := *params.PreserveUnmatchedPath
	sourcePattern := params.SourcePattern

	return map[string]interface{}{
		"destination":             destination,
		"preserve_unmatched_path": preservePath,
		"source_pattern":          sourcePattern,
	}
}

func FlattenCdnFrontDoorRouteConfigurationOverrideAction(input rules.DeliveryRuleRouteConfigurationOverrideAction) (map[string]interface{}, error) {
	var queryStringCachingBehavior string
	var cacheBehavior string
	var compressionEnabled bool
	var cacheDuration string
	var forwardingProtocol string
	var originGroupId string
	queryParameters := make([]interface{}, 0)

	params := input.Parameters

	if config := params.CacheConfiguration; config != nil {
		queryStringCachingBehavior = string(pointer.From(config.QueryStringCachingBehavior))
		cacheBehavior = string(pointer.From(config.CacheBehavior))
		compressionEnabled = pointer.From(config.IsCompressionEnabled) == rules.RuleIsCompressionEnabledEnabled
		queryParameters = flattenCsvToStringSlice(config.QueryParameters)

		if config.CacheDuration != nil {
			cacheDuration = *config.CacheDuration
		}
	} else {
		cacheBehavior = string(rules.RuleIsCompressionEnabledDisabled)
	}

	if override := params.OriginGroupOverride; override != nil {
		forwardingProtocol = string(pointer.From(override.ForwardingProtocol))

		// NOTE: Need to parse this ID insensitively to normalize it because if you modified this
		// resource in the Azure Portal the 'resourceGroup' element comes back as 'resourcegroup'
		// not 'resourceGroup'.
		originGroup, err := parse.FrontDoorOriginGroupIDInsensitively(*override.OriginGroup.Id)
		if err != nil {
			return nil, err
		}

		originGroupId = originGroup.ID()
	}

	return map[string]interface{}{
		"query_string_caching_behavior": queryStringCachingBehavior,
		"cache_behavior":                cacheBehavior,
		"compression_enabled":           compressionEnabled,
		"cache_duration":                cacheDuration,
		"query_string_parameters":       queryParameters,
		"forwarding_protocol":           forwardingProtocol,
		"cdn_frontdoor_origin_group_id": originGroupId,
	}, nil
}
