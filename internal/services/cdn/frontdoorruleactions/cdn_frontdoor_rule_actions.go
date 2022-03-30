package cdnfrontdoorruleactions

import (
	"fmt"
	"strings"

	track1 "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontdoorActionParameters struct {
	Name       track1.NameBasicDeliveryRuleAction
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
		Name:       track1.NameBasicDeliveryRuleActionNameRouteConfigurationOverride,
		TypeName:   "DeliveryRuleRouteConfigurationOverrideActionParameters",
		ConfigName: "route_configuration_override_action",
	}

	m.RequestHeader = CdnFrontdoorActionParameters{
		Name:       track1.NameBasicDeliveryRuleActionNameModifyRequestHeader,
		TypeName:   "DeliveryRuleHeaderActionParameters",
		ConfigName: "request_header_action",
	}

	m.ResponseHeader = CdnFrontdoorActionParameters{
		Name:       track1.NameBasicDeliveryRuleActionNameModifyResponseHeader,
		TypeName:   "DeliveryRuleHeaderActionParameters",
		ConfigName: "response_header_action",
	}

	m.URLRedirect = CdnFrontdoorActionParameters{
		Name:       track1.NameBasicDeliveryRuleActionNameURLRedirect,
		TypeName:   "DeliveryRuleUrlRedirectActionParameters",
		ConfigName: "url_redirect_action",
	}

	m.URLRewrite = CdnFrontdoorActionParameters{
		Name:       track1.NameBasicDeliveryRuleActionNameURLRedirect,
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

func ExpandCdnFrontdoorRequestHeaderAction(input []interface{}) (*[]track1.BasicDeliveryRuleAction, error) {
	output := make([]track1.BasicDeliveryRuleAction, 0)

	m := InitializeCdnFrontdoorActionMappings()

	for _, v := range input {
		item := v.(map[string]interface{})

		requestHeaderAction := track1.DeliveryRuleRequestHeaderAction{
			Name: m.RequestHeader.Name,
			Parameters: &track1.HeaderActionParameters{
				TypeName:     &m.RequestHeader.TypeName,
				HeaderAction: track1.HeaderAction(item["header_action"].(string)),
				HeaderName:   utils.String(item["header_name"].(string)),
				Value:        utils.String(item["value"].(string)),
			},
		}

		if headerValue := *requestHeaderAction.Parameters.Value; headerValue == "" {
			if requestHeaderAction.Parameters.HeaderAction == track1.HeaderActionOverwrite || requestHeaderAction.Parameters.HeaderAction == track1.HeaderActionAppend {
				return nil, fmt.Errorf("the %q block is not valid, %q can not be empty if the %q is set to %q or %q", m.RequestHeader.ConfigName, "value", "header_action", "Append", "Overwrite")
			}
		} else {
			if requestHeaderAction.Parameters.HeaderAction == track1.HeaderActionDelete {
				return nil, fmt.Errorf("the %q block is not valid, %q must be empty if the %q is set to %q", m.RequestHeader.ConfigName, "value", "header_action", "Delete")
			}
		}

		output = append(output, requestHeaderAction)
	}

	return &output, nil
}

func ExpandCdnFrontdoorResponseHeaderAction(input []interface{}) (*[]track1.BasicDeliveryRuleAction, error) {
	output := make([]track1.BasicDeliveryRuleAction, 0)

	m := InitializeCdnFrontdoorActionMappings()

	for _, v := range input {
		item := v.(map[string]interface{})

		responseHeaderAction := track1.DeliveryRuleResponseHeaderAction{
			Name: m.ResponseHeader.Name,
			Parameters: &track1.HeaderActionParameters{
				TypeName:     utils.String(m.ResponseHeader.TypeName),
				HeaderAction: track1.HeaderAction(item["header_action"].(string)),
				HeaderName:   utils.String(item["header_name"].(string)),
				Value:        utils.String(item["value"].(string)),
			},
		}

		if headerValue := *responseHeaderAction.Parameters.Value; headerValue == "" {
			if responseHeaderAction.Parameters.HeaderAction == track1.HeaderActionOverwrite || responseHeaderAction.Parameters.HeaderAction == track1.HeaderActionAppend {
				return nil, fmt.Errorf("the %q block is not valid, %q can not be empty if the %q is set to %q or %q", m.ResponseHeader.ConfigName, "value", "header_action", "Append", "Overwrite")
			}
		} else {
			if responseHeaderAction.Parameters.HeaderAction == track1.HeaderActionDelete {
				return nil, fmt.Errorf("the %q block is not valid, %q must be empty if the %q is set to %q", m.ResponseHeader.ConfigName, "value", "header_action", "Delete")
			}
		}

		output = append(output, responseHeaderAction)
	}

	return &output, nil
}

func ExpandCdnFrontdoorUrlRedirectAction(input []interface{}) (*[]track1.BasicDeliveryRuleAction, error) {
	output := make([]track1.BasicDeliveryRuleAction, 0)

	m := InitializeCdnFrontdoorActionMappings()

	for _, v := range input {
		item := v.(map[string]interface{})

		urlRedirectAction := track1.URLRedirectAction{
			Name: m.URLRedirect.Name,
			Parameters: &track1.URLRedirectActionParameters{
				TypeName:            utils.String(m.URLRedirect.TypeName),
				RedirectType:        track1.RedirectType(item["redirect_type"].(string)),
				DestinationProtocol: track1.DestinationProtocol(item["redirect_protocol"].(string)),
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

func ExpandCdnFrontdoorUrlRewriteAction(input []interface{}) (*[]track1.BasicDeliveryRuleAction, error) {
	output := make([]track1.BasicDeliveryRuleAction, 0)

	m := InitializeCdnFrontdoorActionMappings()

	for _, v := range input {
		item := v.(map[string]interface{})

		urlRedirectAction := track1.URLRewriteAction{
			Name: m.URLRewrite.Name,
			Parameters: &track1.URLRewriteActionParameters{
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

func ExpandCdnFrontdoorRouteConfigurationOverrideAction(input []interface{}) (*[]track1.BasicDeliveryRuleAction, error) {
	output := make([]track1.BasicDeliveryRuleAction, 0)

	m := InitializeCdnFrontdoorActionMappings()

	for _, v := range input {
		item := v.(map[string]interface{})

		originGroupOverride := &track1.OriginGroupOverride{
			OriginGroup: &track1.ResourceReference{
				ID: utils.String(item["cdn_frontdoor_origin_group_id"].(string)),
			},
			ForwardingProtocol: track1.ForwardingProtocol(item["forwarding_protocol"].(string)),
		}

		compressionEnabled := track1.RuleIsCompressionEnabledEnabled
		if !item["compression_enabled"].(bool) {
			compressionEnabled = track1.RuleIsCompressionEnabledDisabled
		}

		// RuleQueryStringCachingBehavior
		cacheConfiguration := &track1.CacheConfiguration{
			QueryStringCachingBehavior: track1.RuleQueryStringCachingBehavior(item["query_string_caching_behavior"].(string)),
			QueryParameters:            expandStringSliceToCsvFormat(item["query_string_parameters"].([]interface{})),
			IsCompressionEnabled:       compressionEnabled,
			CacheBehavior:              track1.RuleCacheBehavior(item["cache_behavior"].(string)),
			CacheDuration:              utils.String(item["cache_duration"].(string)),
		}

		routeConfigurationOverrideAction := track1.DeliveryRuleRouteConfigurationOverrideAction{
			Name: m.RouteConfigurationOverride.Name,
			Parameters: &track1.RouteConfigurationOverrideActionParameters{
				TypeName:            utils.String(m.RouteConfigurationOverride.TypeName),
				OriginGroupOverride: originGroupOverride,
				CacheConfiguration:  cacheConfiguration,
			},
		}

		queryStringCachingBehavior := cacheConfiguration.QueryStringCachingBehavior
		if queryParameters := cacheConfiguration.QueryParameters; queryParameters == nil {
			if queryStringCachingBehavior == track1.RuleQueryStringCachingBehaviorIncludeSpecifiedQueryStrings || queryStringCachingBehavior == track1.RuleQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings {
				return nil, fmt.Errorf("the %q block is not valid, %q can not be empty if the %q is set to %q or %q", m.RouteConfigurationOverride.ConfigName, "query_string_parameters", "query_string_caching_behavior", "IncludeSpecifiedQueryStrings", "IgnoreSpecifiedQueryStrings")
			}
		} else {
			if queryStringCachingBehavior == track1.RuleQueryStringCachingBehaviorUseQueryString || queryStringCachingBehavior == track1.RuleQueryStringCachingBehaviorIgnoreQueryString {
				return nil, fmt.Errorf("the %q block is not valid, %q must not be set if the %q is set to %q or %q", m.RouteConfigurationOverride.ConfigName, "query_string_parameters", "query_string_caching_behavior", "UseQueryStrings", "IgnoreQueryStrings")
			}
		}

		output = append(output, routeConfigurationOverrideAction)
	}

	return &output, nil
}

func FlattenCdnFrontdoorRequestHeaderAction(input track1.BasicDeliveryRuleAction) (map[string]interface{}, error) {
	action, ok := input.AsDeliveryRuleRequestHeaderAction()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule request header action")
	}

	return flattenCdnFrontdoorHeaderAction(action.Parameters), nil
}

func FlattenCdnFrontdoorResponseHeaderAction(input track1.BasicDeliveryRuleAction) (map[string]interface{}, error) {
	action, ok := input.AsDeliveryRuleResponseHeaderAction()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule reesponse header action")
	}

	return flattenCdnFrontdoorHeaderAction(action.Parameters), nil
}

func flattenCdnFrontdoorHeaderAction(input *track1.HeaderActionParameters) map[string]interface{} {
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

func FlattenCdnFrontdoorUrlRedirectAction(input track1.BasicDeliveryRuleAction) (map[string]interface{}, error) {
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

func FlattenCdnFrontdoorUrlRewriteAction(input track1.BasicDeliveryRuleAction) (map[string]interface{}, error) {
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

func FlattenCdnFrontdoorRouteConfigurationOverrideAction(input track1.BasicDeliveryRuleAction) (map[string]interface{}, error) {
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
		compressionEnabled = (params.CacheConfiguration.IsCompressionEnabled == track1.RuleIsCompressionEnabledEnabled)
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
