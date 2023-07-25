// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package CdnFrontDoorruleactions

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontDoorActionParameters struct {
	Name       cdn.NameBasicDeliveryRuleAction
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
		Name:       cdn.NameBasicDeliveryRuleActionNameRouteConfigurationOverride,
		TypeName:   "DeliveryRuleRouteConfigurationOverrideActionParameters",
		ConfigName: "route_configuration_override_action",
	}

	m.RequestHeader = CdnFrontDoorActionParameters{
		Name:       cdn.NameBasicDeliveryRuleActionNameModifyRequestHeader,
		TypeName:   "DeliveryRuleHeaderActionParameters",
		ConfigName: "request_header_action",
	}

	m.ResponseHeader = CdnFrontDoorActionParameters{
		Name:       cdn.NameBasicDeliveryRuleActionNameModifyResponseHeader,
		TypeName:   "DeliveryRuleHeaderActionParameters",
		ConfigName: "response_header_action",
	}

	m.URLRedirect = CdnFrontDoorActionParameters{
		Name:       cdn.NameBasicDeliveryRuleActionNameURLRedirect,
		TypeName:   "DeliveryRuleUrlRedirectActionParameters",
		ConfigName: "url_redirect_action",
	}

	m.URLRewrite = CdnFrontDoorActionParameters{
		Name:       cdn.NameBasicDeliveryRuleActionNameURLRedirect,
		TypeName:   "DeliveryRuleUrlRewriteActionParameters",
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

func ExpandCdnFrontDoorRequestHeaderAction(input []interface{}) (*[]cdn.BasicDeliveryRuleAction, error) {
	output := make([]cdn.BasicDeliveryRuleAction, 0)

	m := InitializeCdnFrontDoorActionMappings()

	for _, v := range input {
		item := v.(map[string]interface{})

		value := item["value"].(string)

		requestHeaderAction := cdn.DeliveryRuleRequestHeaderAction{
			Name: m.RequestHeader.Name,
			Parameters: &cdn.HeaderActionParameters{
				TypeName:     &m.RequestHeader.TypeName,
				HeaderAction: cdn.HeaderAction(item["header_action"].(string)),
				HeaderName:   utils.String(item["header_name"].(string)),
				Value:        utils.String(value),
			},
		}

		if value == "" {
			if requestHeaderAction.Parameters.HeaderAction == cdn.HeaderActionOverwrite || requestHeaderAction.Parameters.HeaderAction == cdn.HeaderActionAppend {
				return nil, fmt.Errorf("the 'request_header_action' block is not valid, 'value' cannot be empty if the 'header_action' is set to 'Append' or 'Overwrite'")
			}
		} else {
			if requestHeaderAction.Parameters.HeaderAction == cdn.HeaderActionDelete {
				return nil, fmt.Errorf("the 'request_header_action' block is not valid, 'value' must be empty if the 'header_action' is set to 'Delete'")
			}
		}

		output = append(output, requestHeaderAction)
	}

	return &output, nil
}

func ExpandCdnFrontDoorResponseHeaderAction(input []interface{}) (*[]cdn.BasicDeliveryRuleAction, error) {
	output := make([]cdn.BasicDeliveryRuleAction, 0)

	m := InitializeCdnFrontDoorActionMappings()

	for _, v := range input {
		item := v.(map[string]interface{})

		responseHeaderAction := cdn.DeliveryRuleResponseHeaderAction{
			Name: m.ResponseHeader.Name,
			Parameters: &cdn.HeaderActionParameters{
				TypeName:     utils.String(m.ResponseHeader.TypeName),
				HeaderAction: cdn.HeaderAction(item["header_action"].(string)),
				HeaderName:   utils.String(item["header_name"].(string)),
				Value:        utils.String(item["value"].(string)),
			},
		}

		if headerValue := *responseHeaderAction.Parameters.Value; headerValue == "" {
			if responseHeaderAction.Parameters.HeaderAction == cdn.HeaderActionOverwrite || responseHeaderAction.Parameters.HeaderAction == cdn.HeaderActionAppend {
				return nil, fmt.Errorf("the 'response_header_action' block is not valid, 'value' cannot be empty if the 'header_action' is set to 'Append' or 'Overwrite'")
			}
		} else {
			if responseHeaderAction.Parameters.HeaderAction == cdn.HeaderActionDelete {
				return nil, fmt.Errorf("the 'response_header_action' block is not valid, 'value' must be empty if the 'header_action' is set to 'Delete'")
			}
		}

		output = append(output, responseHeaderAction)
	}

	return &output, nil
}

func ExpandCdnFrontDoorUrlRedirectAction(input []interface{}) (*[]cdn.BasicDeliveryRuleAction, error) {
	output := make([]cdn.BasicDeliveryRuleAction, 0)

	m := InitializeCdnFrontDoorActionMappings()

	for _, v := range input {
		item := v.(map[string]interface{})

		urlRedirectAction := cdn.URLRedirectAction{
			Name: m.URLRedirect.Name,
			Parameters: &cdn.URLRedirectActionParameters{
				TypeName:            utils.String(m.URLRedirect.TypeName),
				RedirectType:        cdn.RedirectType(item["redirect_type"].(string)),
				DestinationProtocol: cdn.DestinationProtocol(item["redirect_protocol"].(string)),
				CustomPath:          utils.String(item["destination_path"].(string)),
				CustomHostname:      utils.String(item["destination_hostname"].(string)),
				CustomQueryString:   utils.String(item["query_string"].(string)),
				CustomFragment:      utils.String(item["destination_fragment"].(string)),
			},
		}

		output = append(output, urlRedirectAction)
	}

	return &output, nil
}

func ExpandCdnFrontDoorUrlRewriteAction(input []interface{}) (*[]cdn.BasicDeliveryRuleAction, error) {
	output := make([]cdn.BasicDeliveryRuleAction, 0)

	m := InitializeCdnFrontDoorActionMappings()

	for _, v := range input {
		item := v.(map[string]interface{})

		urlRedirectAction := cdn.URLRewriteAction{
			Name: m.URLRewrite.Name,
			Parameters: &cdn.URLRewriteActionParameters{
				TypeName:              utils.String(m.URLRewrite.TypeName),
				Destination:           utils.String(item["destination"].(string)),
				PreserveUnmatchedPath: utils.Bool(item["preserve_unmatched_path"].(bool)),
				SourcePattern:         utils.String(item["source_pattern"].(string)),
			},
		}

		output = append(output, urlRedirectAction)
	}

	return &output, nil
}

func ExpandCdnFrontDoorRouteConfigurationOverrideAction(input []interface{}) (*[]cdn.BasicDeliveryRuleAction, error) {
	output := make([]cdn.BasicDeliveryRuleAction, 0)
	m := InitializeCdnFrontDoorActionMappings()

	for _, v := range input {
		item := v.(map[string]interface{})

		var originGroupOverride cdn.OriginGroupOverride
		var cacheConfiguration cdn.CacheConfiguration

		originGroupIdRaw := item["cdn_frontdoor_origin_group_id"].(string)
		protocol := item["forwarding_protocol"].(string)
		cacheBehavior := item["cache_behavior"].(string)
		compressionEnabled := cdn.RuleIsCompressionEnabledEnabled
		queryStringCachingBehavior := item["query_string_caching_behavior"].(string)
		cacheDuration := item["cache_duration"].(string)

		if !item["compression_enabled"].(bool) {
			compressionEnabled = cdn.RuleIsCompressionEnabledDisabled
		}

		if !features.FourPointOhBeta() {
			// set the default value for forwarding protocol to avoid a breaking change...
			if protocol == "" {
				protocol = string(cdn.ForwardingProtocolMatchRequest)
			}
		}

		// NOTE: It is valid to not define the originGroupOverride in the Route Configuration Override Action
		// however, if you do not define the Origin Group ID you also cannot define the Forwarding Protocol either
		if originGroupIdRaw != "" {
			if protocol == "" {
				return nil, fmt.Errorf("the 'route_configuration_override_action' block is not valid, the 'forwarding_protocol' field must be set")
			}

			originGroupOverride = cdn.OriginGroupOverride{
				OriginGroup: &cdn.ResourceReference{
					ID: utils.String(originGroupIdRaw),
				},
				ForwardingProtocol: cdn.ForwardingProtocol(protocol),
			}
		} else if originGroupIdRaw == "" && item["forwarding_protocol"].(string) != "" {
			return nil, fmt.Errorf("the 'route_configuration_override_action' block is not valid, if the 'cdn_frontdoor_origin_group_id' is not set you cannot define the 'forwarding_protocol', got %q", protocol)
		}

		if cacheBehavior == string(cdn.RuleIsCompressionEnabledDisabled) {
			if queryStringCachingBehavior != "" {
				return nil, fmt.Errorf("the 'route_configuration_override_action' block is not valid, if the 'cache_behavior' is set to 'Disabled' you cannot define the 'query_string_caching_behavior', got %q", queryStringCachingBehavior)
			}

			if queryParameters := item["query_string_parameters"].([]interface{}); len(queryParameters) != 0 {
				return nil, fmt.Errorf("the 'route_configuration_override_action' block is not valid, if the 'cache_behavior' is set to 'Disabled' you cannot define the 'query_string_parameters', got %d", len(queryParameters))
			}

			if cacheDuration != "" {
				return nil, fmt.Errorf("the 'route_configuration_override_action' block is not valid, if the 'cache_behavior' is set to 'Disabled' you cannot define the 'cache_duration', got %q", cacheDuration)
			}
		} else {
			if !features.FourPointOhBeta() {
				// since 'cache_duration', 'query_string_caching_behavior' and 'cache_behavior' are optional create a default values
				// for those values if not set.
				if cacheBehavior == "" {
					cacheBehavior = string(cdn.RuleCacheBehaviorHonorOrigin)
				}

				if queryStringCachingBehavior == "" {
					queryStringCachingBehavior = string(cdn.RuleQueryStringCachingBehaviorIgnoreQueryString)
				}

				// NOTE: if the cacheBehavior is 'HonorOrigin' the cacheDuration must be null, issue #19311
				if cacheBehavior != string(cdn.RuleCacheBehaviorHonorOrigin) {
					if cacheDuration == "" {
						cacheDuration = "1.12:00:00"
					}
				} else if cacheDuration != "" {
					return nil, fmt.Errorf("the 'route_configuration_override_action' block is not valid, if the 'cache_behavior' field is set to 'HonorOrigin' the 'cache_duration' must not be set")
				}
			}

			if features.FourPointOhBeta() {
				if cacheBehavior == "" {
					return nil, fmt.Errorf("the 'route_configuration_override_action' block is not valid, the 'cache_behavior' field must be set")
				}

				if queryStringCachingBehavior == "" {
					return nil, fmt.Errorf("the 'route_configuration_override_action' block is not valid, the 'query_string_caching_behavior' field must be set")
				}

				// NOTE: if the cacheBehavior is 'HonorOrigin' cacheDuration must be null, issue #19311
				if cacheBehavior != string(cdn.RuleCacheBehaviorHonorOrigin) {
					if cacheDuration == "" {
						return nil, fmt.Errorf("the 'route_configuration_override_action' block is not valid, the 'cache_duration' field must be set")
					}
				} else if cacheDuration != "" {
					return nil, fmt.Errorf("the 'route_configuration_override_action' block is not valid, the 'cache_duration' field must not be set if the 'cache_behavior' is 'HonorOrigin'")
				}
			}

			cacheConfiguration = cdn.CacheConfiguration{
				QueryStringCachingBehavior: cdn.RuleQueryStringCachingBehavior(queryStringCachingBehavior),
				QueryParameters:            expandStringSliceToCsvFormat(item["query_string_parameters"].([]interface{})),
				IsCompressionEnabled:       compressionEnabled,
				CacheBehavior:              cdn.RuleCacheBehavior(cacheBehavior),
			}

			if cacheDuration != "" {
				cacheConfiguration.CacheDuration = utils.String(cacheDuration)
			}

			if queryParameters := cacheConfiguration.QueryParameters; queryParameters == nil {
				if cacheConfiguration.QueryStringCachingBehavior == cdn.RuleQueryStringCachingBehaviorIncludeSpecifiedQueryStrings || cacheConfiguration.QueryStringCachingBehavior == cdn.RuleQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings {
					return nil, fmt.Errorf("the 'route_configuration_override_action' block is not valid, 'query_string_parameters' cannot be empty if the 'query_string_caching_behavior' is set to 'IncludeSpecifiedQueryStrings' or 'IgnoreSpecifiedQueryStrings'")
				}
			} else {
				if cacheConfiguration.QueryStringCachingBehavior == cdn.RuleQueryStringCachingBehaviorUseQueryString || cacheConfiguration.QueryStringCachingBehavior == cdn.RuleQueryStringCachingBehaviorIgnoreQueryString {
					return nil, fmt.Errorf("the 'route_configuration_override_action' block is not valid, 'query_string_parameters' must not be set if the'query_string_caching_behavior' is set to 'UseQueryStrings' or 'IgnoreQueryStrings'")
				}
			}
		}

		routeConfigurationOverrideAction := cdn.DeliveryRuleRouteConfigurationOverrideAction{
			Parameters: &cdn.RouteConfigurationOverrideActionParameters{
				TypeName: utils.String(m.RouteConfigurationOverride.TypeName),
			},
		}

		if originGroupOverride.OriginGroup != nil {
			routeConfigurationOverrideAction.Parameters.OriginGroupOverride = &originGroupOverride
		}

		if cacheConfiguration.CacheDuration != nil || cacheConfiguration.CacheBehavior == cdn.RuleCacheBehaviorHonorOrigin {
			routeConfigurationOverrideAction.Parameters.CacheConfiguration = &cacheConfiguration
		}

		output = append(output, routeConfigurationOverrideAction)
	}

	return &output, nil
}

func FlattenHeaderActionParameters(input *cdn.HeaderActionParameters) map[string]interface{} {
	action := ""
	name := ""
	value := ""

	if params := input; params != nil {
		action = string(params.HeaderAction)
		if params.HeaderName != nil {
			name = *params.HeaderName
		}
		if params.Value != nil {
			value = *params.Value
		}
	}

	return map[string]interface{}{
		"header_action": action,
		"header_name":   name,
		"value":         value,
	}
}

func FlattenCdnFrontDoorUrlRedirectAction(input cdn.URLRedirectAction) map[string]interface{} {
	destinationHost := ""
	destinationPath := ""
	queryString := ""
	destinationProtocol := ""
	redirectType := ""
	fragment := ""

	if params := input.Parameters; params != nil {
		if params.CustomHostname != nil {
			destinationHost = *params.CustomHostname
		}
		if params.CustomPath != nil {
			destinationPath = *params.CustomPath
		}
		if params.CustomQueryString != nil {
			queryString = *params.CustomQueryString
		}
		destinationProtocol = string(params.DestinationProtocol)
		redirectType = string(params.RedirectType)
		if params.CustomFragment != nil {
			fragment = *params.CustomFragment
		}
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

func FlattenCdnFrontDoorUrlRewriteAction(input cdn.URLRewriteAction) map[string]interface{} {
	destination := ""
	preservePath := false
	sourcePattern := ""
	if params := input.Parameters; params != nil {
		if params.Destination != nil {
			destination = *params.Destination
		}
		if params.PreserveUnmatchedPath != nil {
			preservePath = *params.PreserveUnmatchedPath
		}
		if params.SourcePattern != nil {
			sourcePattern = *params.SourcePattern
		}
	}

	return map[string]interface{}{
		"destination":             destination,
		"preserve_unmatched_path": preservePath,
		"source_pattern":          sourcePattern,
	}
}

func FlattenCdnFrontDoorRouteConfigurationOverrideAction(input cdn.DeliveryRuleRouteConfigurationOverrideAction) (map[string]interface{}, error) {
	queryStringCachingBehavior := ""
	cacheBehavior := ""
	compressionEnabled := false
	cacheDuration := ""
	queryParameters := make([]interface{}, 0)
	forwardingProtocol := ""
	originGroupId := ""

	if params := input.Parameters; params != nil {
		if config := params.CacheConfiguration; config != nil {
			queryStringCachingBehavior = string(config.QueryStringCachingBehavior)
			cacheBehavior = string(config.CacheBehavior)
			compressionEnabled = config.IsCompressionEnabled == cdn.RuleIsCompressionEnabledEnabled
			queryParameters = flattenCsvToStringSlice(config.QueryParameters)

			if config.CacheDuration != nil {
				cacheDuration = *config.CacheDuration
			}
		} else {
			cacheBehavior = string(cdn.RuleIsCompressionEnabledDisabled)
		}

		if override := params.OriginGroupOverride; override != nil {
			forwardingProtocol = string(override.ForwardingProtocol)

			// NOTE: Need to normalize this ID here because if you modified this in portal the resourceGroup comes back as resourcegroup.
			// ignore the error here since it was set on the resource in Azure and we know it is valid.
			originGroup, err := parse.FrontDoorOriginGroupIDInsensitively(*override.OriginGroup.ID)
			if err != nil {
				return nil, err
			}

			originGroupId = originGroup.ID()
		}
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
