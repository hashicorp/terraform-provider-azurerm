package cdnfrontdoorruleactions

import (
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontdoorActionParameters struct {
	Name       cdn.NameBasicDeliveryRuleAction
	TypeName   string
	ConfigName string
}

type CdnFrontdoorActionMappings struct {
	RouteConfigurationOverride CdnFrontdoorActionParameters
	RequestHeader              CdnFrontdoorActionParameters
	ResponseHeader             CdnFrontdoorActionParameters
	URLRedirect                CdnFrontdoorActionParameters
	URLRewrite                 CdnFrontdoorActionParameters
}

func InitializeCdnFrontdoorActionMappings() *CdnFrontdoorActionMappings {
	m := new(CdnFrontdoorActionMappings)

	m.RouteConfigurationOverride = CdnFrontdoorActionParameters{
		Name:       cdn.NameBasicDeliveryRuleActionNameRouteConfigurationOverride,
		TypeName:   "DeliveryRuleRouteConfigurationOverrideActionParameters",
		ConfigName: "route_configuration_override_action",
	}

	m.RequestHeader = CdnFrontdoorActionParameters{
		Name:       cdn.NameBasicDeliveryRuleActionNameModifyRequestHeader,
		TypeName:   "DeliveryRuleHeaderActionParameters",
		ConfigName: "request_header_action",
	}

	m.ResponseHeader = CdnFrontdoorActionParameters{
		Name:       cdn.NameBasicDeliveryRuleActionNameModifyResponseHeader,
		TypeName:   "DeliveryRuleHeaderActionParameters",
		ConfigName: "response_header_action",
	}

	m.URLRedirect = CdnFrontdoorActionParameters{
		Name:       cdn.NameBasicDeliveryRuleActionNameURLRedirect,
		TypeName:   "DeliveryRuleUrlRedirectActionParameters",
		ConfigName: "url_redirect_action",
	}

	m.URLRewrite = CdnFrontdoorActionParameters{
		Name:       cdn.NameBasicDeliveryRuleActionNameURLRedirect,
		TypeName:   "DeliveryRuleUrlRewriteActionParameters",
		ConfigName: "url_rewrite_action",
	}

	return m
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

func ExpandCdnFrontdoorRequestHeaderAction(input []interface{}) (*[]cdn.BasicDeliveryRuleAction, error) {
	output := make([]cdn.BasicDeliveryRuleAction, 0)

	m := InitializeCdnFrontdoorActionMappings()

	for _, v := range input {
		item := v.(map[string]interface{})

		requestHeaderAction := cdn.DeliveryRuleRequestHeaderAction{
			Name: m.RequestHeader.Name,
			Parameters: &cdn.HeaderActionParameters{
				TypeName:     &m.RequestHeader.TypeName,
				HeaderAction: cdn.HeaderAction(item["header_action"].(string)),
				HeaderName:   utils.String(item["header_name"].(string)),
				Value:        utils.String(item["value"].(string)),
			},
		}

		if headerValue := *requestHeaderAction.Parameters.Value; headerValue == "" {
			if requestHeaderAction.Parameters.HeaderAction == cdn.HeaderActionOverwrite || requestHeaderAction.Parameters.HeaderAction == cdn.HeaderActionAppend {
				return nil, fmt.Errorf("the %q block is not valid, %q can not be empty if the %q is set to %q or %q", m.RequestHeader.ConfigName, "value", "header_action", "Append", "Overwrite")
			}
		} else {
			if requestHeaderAction.Parameters.HeaderAction == cdn.HeaderActionDelete {
				return nil, fmt.Errorf("the %q block is not valid, %q must be empty if the %q is set to %q", m.RequestHeader.ConfigName, "value", "header_action", "Delete")
			}
		}

		output = append(output, requestHeaderAction)
	}

	return &output, nil
}

func ExpandCdnFrontdoorResponseHeaderAction(input []interface{}) (*[]cdn.BasicDeliveryRuleAction, error) {
	output := make([]cdn.BasicDeliveryRuleAction, 0)

	m := InitializeCdnFrontdoorActionMappings()

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
				return nil, fmt.Errorf("the %q block is not valid, %q can not be empty if the %q is set to %q or %q", m.ResponseHeader.ConfigName, "value", "header_action", "Append", "Overwrite")
			}
		} else {
			if responseHeaderAction.Parameters.HeaderAction == cdn.HeaderActionDelete {
				return nil, fmt.Errorf("the %q block is not valid, %q must be empty if the %q is set to %q", m.ResponseHeader.ConfigName, "value", "header_action", "Delete")
			}
		}

		output = append(output, responseHeaderAction)
	}

	return &output, nil
}

func ExpandCdnFrontdoorUrlRedirectAction(input []interface{}) (*[]cdn.BasicDeliveryRuleAction, error) {
	output := make([]cdn.BasicDeliveryRuleAction, 0)

	m := InitializeCdnFrontdoorActionMappings()

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

func ExpandCdnFrontdoorUrlRewriteAction(input []interface{}) (*[]cdn.BasicDeliveryRuleAction, error) {
	output := make([]cdn.BasicDeliveryRuleAction, 0)

	m := InitializeCdnFrontdoorActionMappings()

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

func ExpandCdnFrontdoorRouteConfigurationOverrideAction(input []interface{}) (*[]cdn.BasicDeliveryRuleAction, error) {
	output := make([]cdn.BasicDeliveryRuleAction, 0)

	m := InitializeCdnFrontdoorActionMappings()

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

		// RuleQueryStringCachingBehavior
		cacheConfiguration := &cdn.CacheConfiguration{
			QueryStringCachingBehavior: cdn.RuleQueryStringCachingBehavior(item["query_string_caching_behavior"].(string)),
			QueryParameters:            expandStringSliceToCsvFormat(item["query_string_parameters"].([]interface{})),
			IsCompressionEnabled:       compressionEnabled,
			CacheBehavior:              cdn.RuleCacheBehavior(item["cache_behavior"].(string)),
			CacheDuration:              utils.String(item["cache_duration"].(string)),
		}

		routeConfigurationOverrideAction := cdn.DeliveryRuleRouteConfigurationOverrideAction{
			Name: m.RouteConfigurationOverride.Name,
			Parameters: &cdn.RouteConfigurationOverrideActionParameters{
				TypeName:            utils.String(m.RouteConfigurationOverride.TypeName),
				OriginGroupOverride: originGroupOverride,
				CacheConfiguration:  cacheConfiguration,
			},
		}

		queryStringCachingBehavior := cacheConfiguration.QueryStringCachingBehavior
		if queryParameters := cacheConfiguration.QueryParameters; queryParameters == nil {
			if queryStringCachingBehavior == cdn.RuleQueryStringCachingBehaviorIncludeSpecifiedQueryStrings || queryStringCachingBehavior == cdn.RuleQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings {
				return nil, fmt.Errorf("the %q block is not valid, %q can not be empty if the %q is set to %q or %q", m.RouteConfigurationOverride.ConfigName, "query_string_parameters", "query_string_caching_behavior", "IncludeSpecifiedQueryStrings", "IgnoreSpecifiedQueryStrings")
			}
		} else {
			if queryStringCachingBehavior == cdn.RuleQueryStringCachingBehaviorUseQueryString || queryStringCachingBehavior == cdn.RuleQueryStringCachingBehaviorIgnoreQueryString {
				return nil, fmt.Errorf("the %q block is not valid, %q must not be set if the %q is set to %q or %q", m.RouteConfigurationOverride.ConfigName, "query_string_parameters", "query_string_caching_behavior", "UseQueryStrings", "IgnoreQueryStrings")
			}
		}

		output = append(output, routeConfigurationOverrideAction)
	}

	return &output, nil
}

func FlattenCdnFrontdoorRequestHeaderAction(input cdn.BasicDeliveryRuleAction) (map[string]interface{}, error) {
	action, ok := input.AsDeliveryRuleRequestHeaderAction()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule request header action")
	}

	return flattenCdnFrontdoorHeaderAction(action.Parameters), nil
}

func FlattenCdnFrontdoorResponseHeaderAction(input cdn.BasicDeliveryRuleAction) (map[string]interface{}, error) {
	action, ok := input.AsDeliveryRuleResponseHeaderAction()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule reesponse header action")
	}

	return flattenCdnFrontdoorHeaderAction(action.Parameters), nil
}

func flattenCdnFrontdoorHeaderAction(input *cdn.HeaderActionParameters) map[string]interface{} {
	action := ""
	name := ""
	value := ""

	if params := input; params != nil {
		action = string(params.HeaderAction)
		name = *params.HeaderName
		value = *params.Value
	}

	return map[string]interface{}{
		"header_action": action,
		"header_name":   name,
		"value":         value,
	}
}

func FlattenCdnFrontdoorUrlRedirectAction(input cdn.BasicDeliveryRuleAction) (map[string]interface{}, error) {
	action, ok := input.AsURLRedirectAction()
	if !ok {
		return nil, fmt.Errorf("expected a URL redirect action")
	}

	destinationHost := ""
	destinationPath := ""
	queryString := ""
	destinationProtocol := ""
	redirectType := ""
	fragment := ""

	if params := action.Parameters; params != nil {
		destinationHost = *params.CustomHostname
		destinationPath = *params.CustomPath
		queryString = *params.CustomQueryString
		destinationProtocol = string(params.DestinationProtocol)
		redirectType = string(params.RedirectType)
		fragment = *params.CustomFragment
	}

	return map[string]interface{}{
		"destination_hostname": destinationHost,
		"destination_path":     destinationPath,
		"query_string":         queryString,
		"redirect_protocol":    destinationProtocol,
		"redirect_type":        redirectType,
		"destination_fragment": fragment,
	}, nil
}

func FlattenCdnFrontdoorUrlRewriteAction(input cdn.BasicDeliveryRuleAction) (map[string]interface{}, error) {
	action, ok := input.AsURLRewriteAction()
	if !ok {
		return nil, fmt.Errorf("expected a URL redirect action")
	}

	destination := ""
	preservePath := false
	sourcePattern := ""

	if params := action.Parameters; params != nil {
		destination = *params.Destination
		preservePath = *params.PreserveUnmatchedPath
		sourcePattern = *params.SourcePattern
	}

	return map[string]interface{}{
		"destination":             destination,
		"preserve_unmatched_path": preservePath,
		"source_pattern":          sourcePattern,
	}, nil
}

func FlattenCdnFrontdoorRouteConfigurationOverrideAction(input cdn.BasicDeliveryRuleAction) (map[string]interface{}, error) {
	action, ok := input.AsDeliveryRuleRouteConfigurationOverrideAction()
	if !ok {
		return nil, fmt.Errorf("expected a route configuration override action")
	}

	queryStringCachingBehavior := ""
	cacheBehavior := ""
	compressionEnabled := false
	cacheDuration := ""
	queryParameters := make([]interface{}, 0)
	forwardingProtocol := ""
	originGroupId := ""

	if params := action.Parameters; params != nil {
		queryStringCachingBehavior = string(params.CacheConfiguration.QueryStringCachingBehavior)
		cacheBehavior = string(params.CacheConfiguration.CacheBehavior)
		compressionEnabled = (params.CacheConfiguration.IsCompressionEnabled == cdn.RuleIsCompressionEnabledEnabled)
		cacheDuration = *params.CacheConfiguration.CacheDuration
		queryParameters = flattenCsvToStringSlice(params.CacheConfiguration.QueryParameters)
		forwardingProtocol = string(params.OriginGroupOverride.ForwardingProtocol)
		originGroupId = *params.OriginGroupOverride.OriginGroup.ID
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
