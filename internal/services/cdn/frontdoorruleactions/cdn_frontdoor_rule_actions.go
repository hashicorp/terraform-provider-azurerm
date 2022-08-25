package CdnFrontDoorruleactions

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
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
				return nil, fmt.Errorf("the 'request_header_action' block is not valid, 'value' can not be empty if the 'header_action' is set to 'Append' or 'Overwrite'")
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
				return nil, fmt.Errorf("the 'response_header_action' block is not valid, 'value' can not be empty if the 'header_action' is set to 'Append' or 'Overwrite'")
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

		originGroupOverride := &cdn.OriginGroupOverride{
			OriginGroup: &cdn.ResourceReference{
				ID: utils.String(item["cdn_frontdoor_origin_group_id"].(string)),
			},
			ForwardingProtocol: cdn.ForwardingProtocol(item["forwarding_protocol"].(string)),
		}

		compressionEnabled := cdn.RuleIsCompressionEnabledEnabled
		if !item["compression_enabled"].(bool) {
			compressionEnabled = cdn.RuleIsCompressionEnabledDisabled
		}

		cacheConfiguration := &cdn.CacheConfiguration{
			QueryStringCachingBehavior: cdn.RuleQueryStringCachingBehavior(item["query_string_caching_behavior"].(string)),
			QueryParameters:            expandStringSliceToCsvFormat(item["query_string_parameters"].([]interface{})),
			IsCompressionEnabled:       compressionEnabled,
			CacheBehavior:              cdn.RuleCacheBehavior(item["cache_behavior"].(string)),
			CacheDuration:              utils.String(item["cache_duration"].(string)),
		}

		routeConfigurationOverrideAction := cdn.DeliveryRuleRouteConfigurationOverrideAction{
			Parameters: &cdn.RouteConfigurationOverrideActionParameters{
				TypeName:            utils.String(m.RouteConfigurationOverride.TypeName),
				OriginGroupOverride: originGroupOverride,
				CacheConfiguration:  cacheConfiguration,
			},
		}

		queryStringCachingBehavior := cacheConfiguration.QueryStringCachingBehavior
		if queryParameters := cacheConfiguration.QueryParameters; queryParameters == nil {
			if queryStringCachingBehavior == cdn.RuleQueryStringCachingBehaviorIncludeSpecifiedQueryStrings || queryStringCachingBehavior == cdn.RuleQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings {
				return nil, fmt.Errorf("the 'route_configuration_override_action' block is not valid, 'query_string_parameters' can not be empty if the 'query_string_caching_behavior' is set to 'IncludeSpecifiedQueryStrings' or 'IgnoreSpecifiedQueryStrings'")
			}
		} else {
			if queryStringCachingBehavior == cdn.RuleQueryStringCachingBehaviorUseQueryString || queryStringCachingBehavior == cdn.RuleQueryStringCachingBehaviorIgnoreQueryString {
				return nil, fmt.Errorf("the 'route_configuration_override_action' block is not valid, 'query_string_parameters' must not be set if the'query_string_caching_behavior' is set to 'UseQueryStrings' or 'IgnoreQueryStrings'")
			}
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

func FlattenCdnFrontDoorRouteConfigurationOverrideAction(input cdn.DeliveryRuleRouteConfigurationOverrideAction) map[string]interface{} {
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
			cacheDuration = *config.CacheDuration
			queryParameters = flattenCsvToStringSlice(config.QueryParameters)
		}

		if override := params.OriginGroupOverride; override != nil {
			forwardingProtocol = string(override.ForwardingProtocol)

			if group := override.OriginGroup; group != nil && group.ID != nil {
				originGroupId = *group.ID
			}
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
	}
}
