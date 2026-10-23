// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_rule

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigins"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// Schema

func cdnFrontDoorRuleActionModifyHeaderSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"operator": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(rules.PossibleValuesForHeaderAction(), false),
				},
				"header_name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"header_value": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func cdnFrontDoorRuleTransformsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		MaxItems: 4,
		Elem: &pluginsdk.Schema{
			Type: pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice(rules.PossibleValuesForTransform(),
				false),
		},
	}
}

// Expand / Flatten functions

func expandCdnFrontDoorRuleActions(input []helpers.CdnFrontDoorRuleActionsModel) ([]rules.DeliveryRuleAction, error) {
	results := make([]rules.DeliveryRuleAction, 0)
	if len(input) == 0 {
		return results, nil
	}

	actions := input[0]
	if err := helpers.ValidateCdnFrontDoorRuleActionCounts(helpers.CountCdnFrontDoorBatchRuleActions(input)); err != nil {
		return nil, err
	}

	for _, item := range actions.URLRedirect {
		results = append(results, rules.URLRedirectAction{
			Name: rules.DeliveryRuleActionNameURLRedirect,
			Parameters: rules.URLRedirectActionParameters{
				TypeName:            rules.DeliveryRuleActionParametersTypeDeliveryRuleURLRedirectActionParameters,
				RedirectType:        rules.RedirectType(item.RedirectType),
				DestinationProtocol: pointer.ToEnum[rules.DestinationProtocol](item.RedirectProtocol),
				CustomPath:          pointer.To(item.DestinationPath),
				CustomHostname:      pointer.To(item.DestinationHostName),
				CustomQueryString:   pointer.To(item.QueryString),
				CustomFragment:      pointer.To(item.DestinationFragment),
			},
		})
	}

	for _, item := range actions.URLRewrite {
		results = append(results, rules.URLRewriteAction{
			Name: rules.DeliveryRuleActionNameURLRewrite,
			Parameters: rules.URLRewriteActionParameters{
				TypeName:              rules.DeliveryRuleActionParametersTypeDeliveryRuleURLRewriteActionParameters,
				SourcePattern:         item.SourcePattern,
				Destination:           item.DestinationPath,
				PreserveUnmatchedPath: pointer.To(item.PreserveUnmatchedPathEnabled),
			},
		})
	}

	for _, item := range actions.ModifyRequestHeader {
		if err := helpers.ValidateCdnFrontDoorRuleModifyHeaderAction("modify_request_header", item.Operator, item.HeaderValue); err != nil {
			return nil, err
		}
		results = append(results, rules.DeliveryRuleRequestHeaderAction{
			Name: rules.DeliveryRuleActionNameModifyRequestHeader,
			Parameters: rules.HeaderActionParameters{
				TypeName:     rules.DeliveryRuleActionParametersTypeDeliveryRuleHeaderActionParameters,
				HeaderAction: rules.HeaderAction(item.Operator),
				HeaderName:   item.HeaderName,
				Value:        pointer.To(item.HeaderValue),
			},
		})
	}

	for _, item := range actions.ModifyResponseHeader {
		if err := helpers.ValidateCdnFrontDoorRuleModifyHeaderAction("modify_response_header", item.Operator, item.HeaderValue); err != nil {
			return nil, err
		}
		results = append(results, rules.DeliveryRuleResponseHeaderAction{
			Name: rules.DeliveryRuleActionNameModifyResponseHeader,
			Parameters: rules.HeaderActionParameters{
				TypeName:     rules.DeliveryRuleActionParametersTypeDeliveryRuleHeaderActionParameters,
				HeaderAction: rules.HeaderAction(item.Operator),
				HeaderName:   item.HeaderName,
				Value:        pointer.To(item.HeaderValue),
			},
		})
	}

	for _, item := range actions.RouteConfigurationOverride {
		expandedCache, err := expandCdnFrontDoorRuleRouteConfigurationOverrideCaching(item.Caching)
		if err != nil {
			return nil, fmt.Errorf("expanding `route_configuration_override.0.caching`: %+v", err)
		}

		results = append(results, rules.DeliveryRuleRouteConfigurationOverrideAction{
			Name: rules.DeliveryRuleActionNameRouteConfigurationOverride,
			Parameters: rules.RouteConfigurationOverrideActionParameters{
				TypeName:            rules.DeliveryRuleActionParametersTypeDeliveryRuleRouteConfigurationOverrideActionParameters,
				CacheConfiguration:  expandedCache,
				OriginGroupOverride: expandCdnFrontDoorRuleRouteConfigurationOverrideOriginGroup(item.OriginGroup),
			},
		})
	}

	return results, nil
}

func flattenCdnFrontDoorRuleActions(input *[]rules.DeliveryRuleAction) ([]helpers.CdnFrontDoorRuleActionsModel, error) {
	if input == nil {
		return []helpers.CdnFrontDoorRuleActionsModel{}, nil
	}

	results := make([]helpers.CdnFrontDoorRuleActionsModel, 0, len(*input))
	actions := helpers.CdnFrontDoorRuleActionsModel{}

	for _, action := range *input {
		switch a := action.(type) {
		case rules.URLRedirectAction:
			actions.URLRedirect = append(actions.URLRedirect, helpers.CdnFrontDoorRuleURLRedirectActionModel{
				RedirectType:        string(a.Parameters.RedirectType),
				RedirectProtocol:    pointer.FromEnum(a.Parameters.DestinationProtocol),
				DestinationPath:     pointer.From(a.Parameters.CustomPath),
				DestinationHostName: pointer.From(a.Parameters.CustomHostname),
				QueryString:         pointer.From(a.Parameters.CustomQueryString),
				DestinationFragment: pointer.From(a.Parameters.CustomFragment),
			})
		case rules.URLRewriteAction:
			actions.URLRewrite = append(actions.URLRewrite, helpers.CdnFrontDoorRuleURLRewriteActionModel{
				SourcePattern:                a.Parameters.SourcePattern,
				DestinationPath:              a.Parameters.Destination,
				PreserveUnmatchedPathEnabled: pointer.From(a.Parameters.PreserveUnmatchedPath),
			})
		case rules.DeliveryRuleRequestHeaderAction:
			actions.ModifyRequestHeader = append(actions.ModifyRequestHeader, helpers.CdnFrontDoorRuleHeaderActionModel{
				Operator:    string(a.Parameters.HeaderAction),
				HeaderName:  a.Parameters.HeaderName,
				HeaderValue: pointer.From(a.Parameters.Value),
			})
		case rules.DeliveryRuleResponseHeaderAction:
			actions.ModifyResponseHeader = append(actions.ModifyResponseHeader, helpers.CdnFrontDoorRuleHeaderActionModel{
				Operator:    string(a.Parameters.HeaderAction),
				HeaderName:  a.Parameters.HeaderName,
				HeaderValue: pointer.From(a.Parameters.Value),
			})
		case rules.DeliveryRuleRouteConfigurationOverrideAction:
			flattened, err := flattenCdnFrontDoorRuleRouteConfigurationOverrideAction(a.Parameters)
			if err != nil {
				return results, err
			}
			actions.RouteConfigurationOverride = append(actions.RouteConfigurationOverride, flattened)
		default:
			return results, fmt.Errorf("unsupported action (`%s`) encountered", a.DeliveryRuleAction().Name)
		}
	}

	return append(results, actions), nil
}

func flattenCdnFrontDoorRuleRouteConfigurationOverrideAction(input rules.RouteConfigurationOverrideActionParameters) (helpers.CdnFrontDoorRuleRouteConfigurationOverrideActionModel, error) {
	result := helpers.CdnFrontDoorRuleRouteConfigurationOverrideActionModel{
		Caching: flattenCdnFrontDoorRuleRouteConfigurationOverrideCaching(input.CacheConfiguration),
	}

	originGroup, err := flattenCdnFrontDoorRuleRouteConfigurationOverrideOriginGroup(input.OriginGroupOverride)
	if err != nil {
		return result, err
	}
	result.OriginGroup = originGroup

	return result, nil
}

func expandCdnFrontDoorRuleRouteConfigurationOverrideCaching(input []helpers.CdnFrontDoorRuleRouteConfigurationOverrideCachingModel) (*rules.CacheConfiguration, error) {
	if len(input) == 0 {
		return nil, nil
	}

	v := input[0]
	if v.Behaviour == helpers.RuleCacheBehaviorDisabled {
		if v.Duration != "" {
			return nil, fmt.Errorf("when `route_configuration_override.caching.behaviour` is set to `Disabled`, you cannot define `route_configuration_override.caching.duration`, got `%s`", v.Duration)
		}

		if len(v.QueryStringParameters) != 0 {
			return nil, fmt.Errorf("when `route_configuration_override.caching.behaviour` is set to `Disabled`, you cannot define any `route_configuration_override.caching.query_string_parameters`, got `%d` parameters", len(v.QueryStringParameters))
		}

		if v.QueryStringBehaviour != "" {
			return nil, fmt.Errorf("when `route_configuration_override.caching.behaviour` is set to `Disabled`, you cannot define `route_configuration_override.caching.query_string_behaviour`, got `%s`", v.QueryStringBehaviour)
		}

		if v.CompressionEnabled {
			return nil, fmt.Errorf("when `route_configuration_override.caching.behaviour` is set to `Disabled`, you cannot define `route_configuration_override.caching.compression_enabled`, got `%t`", v.CompressionEnabled)
		}

		return nil, nil
	}

	if v.QueryStringBehaviour == "" {
		return nil, fmt.Errorf("when `route_configuration_override.caching.behaviour` is set to `%s`, you must also define `route_configuration_override.caching.query_string_behaviour`", v.Behaviour)
	}

	// `HonorOrigin` must not carry an explicit cache duration.
	if v.Behaviour == string(rules.RuleCacheBehaviorHonorOrigin) {
		if v.Duration != "" {
			return nil, fmt.Errorf("when `route_configuration_override.caching.behaviour` is set to `%s`, you cannot define `route_configuration_override.caching.duration`, got `%s`", rules.RuleCacheBehaviorHonorOrigin, v.Duration)
		}
	} else if v.Duration == "" {
		return nil, fmt.Errorf("when `route_configuration_override.caching.behaviour` is set to `%s`, you must also define `route_configuration_override.caching.duration`", v.Behaviour)
	}

	switch rules.RuleQueryStringCachingBehavior(v.QueryStringBehaviour) {
	case rules.RuleQueryStringCachingBehaviorIncludeSpecifiedQueryStrings, rules.RuleQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings:
		if len(v.QueryStringParameters) == 0 {
			return nil, fmt.Errorf("when `route_configuration_override.caching.query_string_behaviour` is set to `%s`, you must also define one or more `route_configuration_override.caching.query_string_parameters`", v.QueryStringBehaviour)
		}
	case rules.RuleQueryStringCachingBehaviorUseQueryString, rules.RuleQueryStringCachingBehaviorIgnoreQueryString:
		if len(v.QueryStringParameters) > 0 {
			return nil, fmt.Errorf("when `route_configuration_override.caching.query_string_behaviour` is set to `%s`, you cannot define `route_configuration_override.caching.query_string_parameters`", v.QueryStringBehaviour)
		}
	}

	compressionEnabled := rules.RuleIsCompressionEnabledDisabled
	if v.CompressionEnabled {
		compressionEnabled = rules.RuleIsCompressionEnabledEnabled
	}

	queryParameters := (*string)(nil)
	if len(v.QueryStringParameters) > 0 {
		queryParameters = pointer.To(strings.Join(v.QueryStringParameters, ","))
	}

	return &rules.CacheConfiguration{
		CacheBehavior:              pointer.ToEnum[rules.RuleCacheBehavior](v.Behaviour),
		CacheDuration:              pointer.ToOrNil(v.Duration),
		IsCompressionEnabled:       pointer.To(compressionEnabled),
		QueryParameters:            queryParameters,
		QueryStringCachingBehavior: pointer.ToOrNil(rules.RuleQueryStringCachingBehavior(v.QueryStringBehaviour)),
	}, nil
}

func flattenCdnFrontDoorRuleRouteConfigurationOverrideCaching(input *rules.CacheConfiguration) []helpers.CdnFrontDoorRuleRouteConfigurationOverrideCachingModel {
	result := make([]helpers.CdnFrontDoorRuleRouteConfigurationOverrideCachingModel, 0)
	if input == nil {
		// The API treats omission as disabled, so we'll need to set `Disabled` back into state if it's nil
		return append(result, helpers.CdnFrontDoorRuleRouteConfigurationOverrideCachingModel{
			Behaviour: string(rules.RuleIsCompressionEnabledDisabled),
		})
	}

	v := helpers.CdnFrontDoorRuleRouteConfigurationOverrideCachingModel{
		Behaviour:            pointer.FromEnum(input.CacheBehavior),
		Duration:             pointer.FromEnum(input.CacheDuration),
		CompressionEnabled:   pointer.From(input.IsCompressionEnabled) == rules.RuleIsCompressionEnabledEnabled,
		QueryStringBehaviour: pointer.FromEnum(input.QueryStringCachingBehavior),
	}

	if input.QueryParameters != nil {
		v.QueryStringParameters = strings.Split(*input.QueryParameters, ",")
	}

	return append(result, v)
}

func expandCdnFrontDoorRuleRouteConfigurationOverrideOriginGroup(input []helpers.CdnFrontDoorRuleRouteConfigurationOverrideOriginGroupModel) *rules.OriginGroupOverride {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	return &rules.OriginGroupOverride{
		ForwardingProtocol: pointer.ToEnum[rules.ForwardingProtocol](v.ForwardingProtocol),
		OriginGroup: &rules.ResourceReference{
			Id: pointer.To(v.CdnFrontdoorOriginGroupID),
		},
	}
}

func flattenCdnFrontDoorRuleRouteConfigurationOverrideOriginGroup(input *rules.OriginGroupOverride) ([]helpers.CdnFrontDoorRuleRouteConfigurationOverrideOriginGroupModel, error) {
	result := make([]helpers.CdnFrontDoorRuleRouteConfigurationOverrideOriginGroupModel, 0)
	if input == nil {
		return result, nil
	}

	v := helpers.CdnFrontDoorRuleRouteConfigurationOverrideOriginGroupModel{
		ForwardingProtocol: pointer.FromEnum(input.ForwardingProtocol),
	}

	if input.OriginGroup != nil && input.OriginGroup.Id != nil {
		originGroupID, err := afdorigins.ParseOriginGroupIDInsensitively(*input.OriginGroup.Id)
		if err != nil {
			return result, err
		}
		v.CdnFrontdoorOriginGroupID = originGroupID.ID()
	}

	return append(result, v), nil
}

func expandCdnFrontDoorRuleConditions(input []helpers.CdnFrontDoorRuleConditionsModel) ([]rules.DeliveryRuleCondition, error) {
	results := make([]rules.DeliveryRuleCondition, 0)
	if len(input) == 0 {
		return results, nil
	}

	conditions := input[0]

	expanded, err := expandCdnFrontDoorRuleConditionBaseModel(conditions.RemoteAddress, "remote_address", expandCdnFrontDoorRuleRemoteAddressCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorRuleConditionWithTransformsModel(conditions.QueryString, "query_string", expandCdnFrontDoorRuleQueryStringCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorRuleConditionWithTransformsModel(conditions.RequestURL, "request_url", expandCdnFrontDoorRuleRequestURLCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorRuleConditionWithTransformsModel(conditions.RequestBody, "request_body", expandCdnFrontDoorRuleRequestBodyCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorRuleConditionWithTransformsModel(conditions.RequestPath, "request_path", expandCdnFrontDoorRuleRequestPathCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorRuleConditionWithTransformsModel(conditions.RequestFileExtension, "request_file_extension", expandCdnFrontDoorRuleRequestFileExtensionCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorRuleConditionWithTransformsModel(conditions.RequestFilename, "request_filename", expandCdnFrontDoorRuleRequestFilenameCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorRuleConditionBaseModel(conditions.SocketAddress, "socket_address", expandCdnFrontDoorRuleSocketAddressCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorRuleConditionBaseModel(conditions.ClientPort, "client_port", expandCdnFrontDoorRuleClientPortCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorRuleConditionBaseModel(conditions.ServerPort, "server_port", expandCdnFrontDoorRuleServerPortCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorRuleConditionWithTransformsModel(conditions.HostName, "host_name", expandCdnFrontDoorRuleHostNameCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorRuleConditionBaseModel(conditions.RequestMethod, "request_method", expandCdnFrontDoorRuleRequestMethodCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorRuleConditionWithNameAndTransformsModel(conditions.PostArgs, "post_argument", expandCdnFrontDoorRulePostArgsCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorRuleConditionWithNameAndTransformsModel(conditions.RequestHeader, "request_header", expandCdnFrontDoorRuleRequestHeaderCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorRuleConditionBaseModel(conditions.RequestScheme, "request_scheme", expandCdnFrontDoorRuleRequestSchemeCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorRuleConditionBaseModel(conditions.HTTPVersion, "http_version", expandCdnFrontDoorRuleHTTPVersionCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorRuleConditionWithNameAndTransformsModel(conditions.RequestCookies, "request_cookies", expandCdnFrontDoorRuleCookiesCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorRuleConditionBaseModel(conditions.DeviceType, "device_type", expandCdnFrontDoorRuleDeviceTypeCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorRuleConditionBaseModel(conditions.SSLProtocol, "ssl_protocol", expandCdnFrontDoorRuleSSLProtocolCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	if len(results) > 10 {
		return nil, fmt.Errorf("the `conditions` block may only contain up to 10 match conditions, got %d", len(results))
	}

	return results, nil
}

func flattenCdnFrontDoorRuleConditions(input *[]rules.DeliveryRuleCondition) ([]helpers.CdnFrontDoorRuleConditionsModel, error) {
	if input == nil || len(*input) == 0 {
		return []helpers.CdnFrontDoorRuleConditionsModel{}, nil
	}

	conditions := helpers.CdnFrontDoorRuleConditionsModel{}
	for _, condition := range *input {
		switch c := condition.(type) {
		case rules.DeliveryRuleRemoteAddressCondition:
			conditions.RemoteAddress = append(conditions.RemoteAddress, helpers.FlattenCdnFrontDoorRuleConditionBaseModel(string(c.Parameters.Operator), pointer.From(c.Parameters.MatchValues), c.Parameters.NegateCondition))
		case rules.DeliveryRuleRequestMethodCondition:
			conditions.RequestMethod = append(conditions.RequestMethod, helpers.FlattenCdnFrontDoorRuleConditionBaseModel(string(c.Parameters.Operator), pointer.FromEnumSlice(c.Parameters.MatchValues), c.Parameters.NegateCondition))
		case rules.DeliveryRuleQueryStringCondition:
			conditions.QueryString = append(conditions.QueryString, helpers.FlattenCdnFrontDoorBatchRuleSetConditionWithTransformsModel(string(c.Parameters.Operator), pointer.From(c.Parameters.MatchValues), pointer.FromEnumSlice(c.Parameters.Transforms), c.Parameters.NegateCondition))
		case rules.DeliveryRulePostArgsCondition:
			conditions.PostArgs = append(conditions.PostArgs, helpers.FlattenCdnFrontDoorRuleConditionWithNameAndTransformsModel(string(c.Parameters.Operator), c.Parameters.Selector, pointer.From(c.Parameters.MatchValues), pointer.FromEnumSlice(c.Parameters.Transforms), c.Parameters.NegateCondition))
		case rules.DeliveryRuleRequestUriCondition:
			conditions.RequestURL = append(conditions.RequestURL, helpers.FlattenCdnFrontDoorBatchRuleSetConditionWithTransformsModel(string(c.Parameters.Operator), pointer.From(c.Parameters.MatchValues), pointer.FromEnumSlice(c.Parameters.Transforms), c.Parameters.NegateCondition))
		case rules.DeliveryRuleRequestHeaderCondition:
			conditions.RequestHeader = append(conditions.RequestHeader, helpers.FlattenCdnFrontDoorRuleConditionWithNameAndTransformsModel(string(c.Parameters.Operator), c.Parameters.Selector, pointer.From(c.Parameters.MatchValues), pointer.FromEnumSlice(c.Parameters.Transforms), c.Parameters.NegateCondition))
		case rules.DeliveryRuleRequestBodyCondition:
			conditions.RequestBody = append(conditions.RequestBody, helpers.FlattenCdnFrontDoorBatchRuleSetConditionWithTransformsModel(string(c.Parameters.Operator), pointer.From(c.Parameters.MatchValues), pointer.FromEnumSlice(c.Parameters.Transforms), c.Parameters.NegateCondition))
		case rules.DeliveryRuleRequestSchemeCondition:
			conditions.RequestScheme = append(conditions.RequestScheme, helpers.FlattenCdnFrontDoorRuleConditionBaseModel(string(c.Parameters.Operator), pointer.FromEnumSlice(c.Parameters.MatchValues), c.Parameters.NegateCondition))
		case rules.DeliveryRuleURLPathCondition:
			conditions.RequestPath = append(conditions.RequestPath, helpers.FlattenCdnFrontDoorBatchRuleSetConditionWithTransformsModel(string(c.Parameters.Operator), pointer.From(c.Parameters.MatchValues), pointer.FromEnumSlice(c.Parameters.Transforms), c.Parameters.NegateCondition))
		case rules.DeliveryRuleURLFileExtensionCondition:
			conditions.RequestFileExtension = append(conditions.RequestFileExtension, helpers.FlattenCdnFrontDoorBatchRuleSetConditionWithTransformsModel(string(c.Parameters.Operator), pointer.From(c.Parameters.MatchValues), pointer.FromEnumSlice(c.Parameters.Transforms), c.Parameters.NegateCondition))
		case rules.DeliveryRuleURLFileNameCondition:
			conditions.RequestFilename = append(conditions.RequestFilename, helpers.FlattenCdnFrontDoorBatchRuleSetConditionWithTransformsModel(string(c.Parameters.Operator), pointer.From(c.Parameters.MatchValues), pointer.FromEnumSlice(c.Parameters.Transforms), c.Parameters.NegateCondition))
		case rules.DeliveryRuleHTTPVersionCondition:
			conditions.HTTPVersion = append(conditions.HTTPVersion, helpers.FlattenCdnFrontDoorRuleConditionBaseModel(string(c.Parameters.Operator), pointer.FromEnumSlice(c.Parameters.MatchValues), c.Parameters.NegateCondition))
		case rules.DeliveryRuleCookiesCondition:
			conditions.RequestCookies = append(conditions.RequestCookies, helpers.FlattenCdnFrontDoorRuleConditionWithNameAndTransformsModel(string(c.Parameters.Operator), c.Parameters.Selector, pointer.From(c.Parameters.MatchValues), pointer.FromEnumSlice(c.Parameters.Transforms), c.Parameters.NegateCondition))
		case rules.DeliveryRuleIsDeviceCondition:
			conditions.DeviceType = append(conditions.DeviceType, helpers.FlattenCdnFrontDoorRuleConditionBaseModel(string(c.Parameters.Operator), pointer.FromEnumSlice(c.Parameters.MatchValues), c.Parameters.NegateCondition))
		case rules.DeliveryRuleSocketAddrCondition:
			conditions.SocketAddress = append(conditions.SocketAddress, helpers.FlattenCdnFrontDoorRuleConditionBaseModel(string(c.Parameters.Operator), pointer.From(c.Parameters.MatchValues), c.Parameters.NegateCondition))
		case rules.DeliveryRuleClientPortCondition:
			conditions.ClientPort = append(conditions.ClientPort, helpers.FlattenCdnFrontDoorRuleConditionBaseModel(string(c.Parameters.Operator), pointer.From(c.Parameters.MatchValues), c.Parameters.NegateCondition))
		case rules.DeliveryRuleServerPortCondition:
			conditions.ServerPort = append(conditions.ServerPort, helpers.FlattenCdnFrontDoorRuleConditionBaseModel(string(c.Parameters.Operator), pointer.From(c.Parameters.MatchValues), c.Parameters.NegateCondition))
		case rules.DeliveryRuleHostNameCondition:
			conditions.HostName = append(conditions.HostName, helpers.FlattenCdnFrontDoorBatchRuleSetConditionWithTransformsModel(string(c.Parameters.Operator), pointer.From(c.Parameters.MatchValues), pointer.FromEnumSlice(c.Parameters.Transforms), c.Parameters.NegateCondition))
		case rules.DeliveryRuleSslProtocolCondition:
			conditions.SSLProtocol = append(conditions.SSLProtocol, helpers.FlattenCdnFrontDoorRuleConditionBaseModel(string(c.Parameters.Operator), pointer.FromEnumSlice(c.Parameters.MatchValues), c.Parameters.NegateCondition))
		default:
			return []helpers.CdnFrontDoorRuleConditionsModel{}, fmt.Errorf("unsupported condition (`%s`) encountered", condition.DeliveryRuleCondition().Name)
		}
	}

	return []helpers.CdnFrontDoorRuleConditionsModel{conditions}, nil
}

func expandCdnFrontDoorRuleClientPortCondition(input helpers.CdnFrontDoorRuleConditionBaseModel) (rules.DeliveryRuleCondition, error) {
	operator, negated := helpers.ExpandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := helpers.ValidateCdnFrontDoorRuleConditionValues("client_port", operator, input.Values); err != nil {
		return nil, err
	}

	return rules.DeliveryRuleClientPortCondition{
		Name: rules.MatchVariableClientPort,
		Parameters: rules.ClientPortMatchConditionParameters{
			TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleClientPortConditionParameters,
			Operator:        rules.ClientPortOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
		},
	}, nil
}

func expandCdnFrontDoorRuleCookiesCondition(input helpers.CdnFrontDoorRuleConditionWithNameAndTransformsModel) (rules.DeliveryRuleCondition, error) {
	operator, negated := helpers.ExpandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := helpers.ValidateCdnFrontDoorRuleConditionValues("request_cookies", operator, input.Values); err != nil {
		return nil, err
	}

	return rules.DeliveryRuleCookiesCondition{
		Name: rules.MatchVariableCookies,
		Parameters: rules.CookiesMatchConditionParameters{
			TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleCookiesConditionParameters,
			Selector:        pointer.To(input.Name),
			Operator:        rules.CookiesOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
			Transforms:      pointer.ToEnumSlice[rules.Transform](input.Transforms),
		},
	}, nil
}

func expandCdnFrontDoorRuleHostNameCondition(input helpers.CdnFrontDoorRuleConditionWithTransformsModel) (rules.DeliveryRuleCondition, error) {
	operator, negated := helpers.ExpandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := helpers.ValidateCdnFrontDoorRuleConditionValues("host_name", operator, input.Values); err != nil {
		return nil, err
	}

	return rules.DeliveryRuleHostNameCondition{
		Name: rules.MatchVariableHostName,
		Parameters: rules.HostNameMatchConditionParameters{
			TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleHostNameConditionParameters,
			Operator:        rules.HostNameOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
			Transforms:      pointer.ToEnumSlice[rules.Transform](input.Transforms),
		},
	}, nil
}

func expandCdnFrontDoorRuleHTTPVersionCondition(input helpers.CdnFrontDoorRuleConditionBaseModel) (rules.DeliveryRuleCondition, error) {
	operator, negated := helpers.ExpandCdnFrontDoorRuleConditionOperator(input.Operator)

	return rules.DeliveryRuleHTTPVersionCondition{
		Name: rules.MatchVariableHTTPVersion,
		Parameters: rules.HTTPVersionMatchConditionParameters{
			TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleHTTPVersionConditionParameters,
			Operator:        rules.HTTPVersionOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
		},
	}, nil
}

func expandCdnFrontDoorRuleDeviceTypeCondition(input helpers.CdnFrontDoorRuleConditionBaseModel) (rules.DeliveryRuleCondition, error) {
	operator, negated := helpers.ExpandCdnFrontDoorRuleConditionOperator(input.Operator)

	return rules.DeliveryRuleIsDeviceCondition{
		Name: rules.MatchVariableIsDevice,
		Parameters: rules.IsDeviceMatchConditionParameters{
			TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleIsDeviceConditionParameters,
			Operator:        rules.IsDeviceOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.ToEnumSlice[rules.IsDeviceMatchValue](input.Values),
		},
	}, nil
}

func expandCdnFrontDoorRulePostArgsCondition(input helpers.CdnFrontDoorRuleConditionWithNameAndTransformsModel) (rules.DeliveryRuleCondition, error) {
	operator, negated := helpers.ExpandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := helpers.ValidateCdnFrontDoorRuleConditionValues("post_argument", operator, input.Values); err != nil {
		return nil, err
	}

	return rules.DeliveryRulePostArgsCondition{
		Name: rules.MatchVariablePostArgs,
		Parameters: rules.PostArgsMatchConditionParameters{
			TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRulePostArgsConditionParameters,
			Selector:        pointer.To(input.Name),
			Operator:        rules.PostArgsOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
			Transforms:      pointer.ToEnumSlice[rules.Transform](input.Transforms),
		},
	}, nil
}

func expandCdnFrontDoorRuleQueryStringCondition(input helpers.CdnFrontDoorRuleConditionWithTransformsModel) (rules.DeliveryRuleCondition, error) {
	operator, negated := helpers.ExpandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := helpers.ValidateCdnFrontDoorRuleConditionValues("query_string", operator, input.Values); err != nil {
		return nil, err
	}

	return rules.DeliveryRuleQueryStringCondition{
		Name: rules.MatchVariableQueryString,
		Parameters: rules.QueryStringMatchConditionParameters{
			TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleQueryStringConditionParameters,
			Operator:        rules.QueryStringOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
			Transforms:      pointer.ToEnumSlice[rules.Transform](input.Transforms),
		},
	}, nil
}

func expandCdnFrontDoorRuleRemoteAddressCondition(input helpers.CdnFrontDoorRuleConditionBaseModel) (rules.DeliveryRuleCondition, error) {
	operator, negated := helpers.ExpandCdnFrontDoorRuleConditionOperator(input.Operator)

	switch rules.RemoteAddressOperator(operator) {
	case rules.RemoteAddressOperatorGeoMatch:
		for _, v := range input.Values {
			if ok, _ := validation.RegExHelper(v, "values", `^[A-Z]{2}$`); !ok {
				return nil, fmt.Errorf("when `conditions.remote_address.operator` is `%s` the values in `conditions.remote_address.values` must be valid country codes consisting of 2 uppercase characters, got `%s`", input.Operator, v)
			}
		}
	case rules.RemoteAddressOperatorIPMatch:
		values := make([]any, 0, len(input.Values))
		for _, matchValue := range input.Values {
			values = append(values, matchValue)
			if _, errs := validate.FrontDoorRuleCidrIsValid(matchValue, "values"); len(errs) > 0 {
				return nil, fmt.Errorf("when `conditions.remote_address.operator` is `%s` the values in `conditions.remote_address.values` must be valid IPv4 or IPv6 CIDRs, got `%s`", input.Operator, matchValue)
			}
		}

		if _, errs := validate.FrontDoorRuleCidrOverlap(values, "values"); len(errs) > 0 {
			return nil, fmt.Errorf("`remote_address` is invalid: %+v", errs[0])
		}
	}

	return rules.DeliveryRuleRemoteAddressCondition{
		Name: rules.MatchVariableRemoteAddress,
		Parameters: rules.RemoteAddressMatchConditionParameters{
			TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleRemoteAddressConditionParameters,
			Operator:        rules.RemoteAddressOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
		},
	}, nil
}

func expandCdnFrontDoorRuleRequestBodyCondition(input helpers.CdnFrontDoorRuleConditionWithTransformsModel) (rules.DeliveryRuleCondition, error) {
	operator, negated := helpers.ExpandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := helpers.ValidateCdnFrontDoorRuleConditionValues("request_body", operator, input.Values); err != nil {
		return nil, err
	}

	return rules.DeliveryRuleRequestBodyCondition{
		Name: rules.MatchVariableRequestBody,
		Parameters: rules.RequestBodyMatchConditionParameters{
			TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleRequestBodyConditionParameters,
			Operator:        rules.RequestBodyOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
			Transforms:      pointer.ToEnumSlice[rules.Transform](input.Transforms),
		},
	}, nil
}

func expandCdnFrontDoorRuleRequestHeaderCondition(input helpers.CdnFrontDoorRuleConditionWithNameAndTransformsModel) (rules.DeliveryRuleCondition, error) {
	operator, negated := helpers.ExpandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := helpers.ValidateCdnFrontDoorRuleConditionValues("request_header", operator, input.Values); err != nil {
		return nil, err
	}

	return rules.DeliveryRuleRequestHeaderCondition{
		Name: rules.MatchVariableRequestHeader,
		Parameters: rules.RequestHeaderMatchConditionParameters{
			TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleRequestHeaderConditionParameters,
			Selector:        pointer.To(input.Name),
			Operator:        rules.RequestHeaderOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
			Transforms:      pointer.ToEnumSlice[rules.Transform](input.Transforms),
		},
	}, nil
}

func expandCdnFrontDoorRuleRequestMethodCondition(input helpers.CdnFrontDoorRuleConditionBaseModel) (rules.DeliveryRuleCondition, error) {
	operator, negated := helpers.ExpandCdnFrontDoorRuleConditionOperator(input.Operator)

	return rules.DeliveryRuleRequestMethodCondition{
		Name: rules.MatchVariableRequestMethod,
		Parameters: rules.RequestMethodMatchConditionParameters{
			TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleRequestMethodConditionParameters,
			Operator:        rules.RequestMethodOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.ToEnumSlice[rules.RequestMethodMatchValue](input.Values),
		},
	}, nil
}

func expandCdnFrontDoorRuleRequestSchemeCondition(input helpers.CdnFrontDoorRuleConditionBaseModel) (rules.DeliveryRuleCondition, error) {
	operator, negated := helpers.ExpandCdnFrontDoorRuleConditionOperator(input.Operator)

	return rules.DeliveryRuleRequestSchemeCondition{
		Name: rules.MatchVariableRequestScheme,
		Parameters: rules.RequestSchemeMatchConditionParameters{
			TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleRequestSchemeConditionParameters,
			Operator:        rules.RequestSchemeMatchConditionParametersOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.ToEnumSlice[rules.RequestSchemeMatchValue](input.Values),
		},
	}, nil
}

func expandCdnFrontDoorRuleRequestURLCondition(input helpers.CdnFrontDoorRuleConditionWithTransformsModel) (rules.DeliveryRuleCondition, error) {
	operator, negated := helpers.ExpandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := helpers.ValidateCdnFrontDoorRuleConditionValues("request_url", operator, input.Values); err != nil {
		return nil, err
	}

	return rules.DeliveryRuleRequestUriCondition{
		Name: rules.MatchVariableRequestUri,
		Parameters: rules.RequestUriMatchConditionParameters{
			TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleRequestUriConditionParameters,
			Operator:        rules.RequestUriOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
			Transforms:      pointer.ToEnumSlice[rules.Transform](input.Transforms),
		},
	}, nil
}

func expandCdnFrontDoorRuleServerPortCondition(input helpers.CdnFrontDoorRuleConditionBaseModel) (rules.DeliveryRuleCondition, error) {
	operator, negated := helpers.ExpandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := helpers.ValidateCdnFrontDoorRuleConditionValues("server_port", operator, input.Values); err != nil {
		return nil, err
	}

	return rules.DeliveryRuleServerPortCondition{
		Name: rules.MatchVariableServerPort,
		Parameters: rules.ServerPortMatchConditionParameters{
			TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleServerPortConditionParameters,
			Operator:        rules.ServerPortOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
		},
	}, nil
}

func expandCdnFrontDoorRuleSocketAddressCondition(input helpers.CdnFrontDoorRuleConditionBaseModel) (rules.DeliveryRuleCondition, error) {
	operator, negated := helpers.ExpandCdnFrontDoorRuleConditionOperator(input.Operator)

	return rules.DeliveryRuleSocketAddrCondition{
		Name: rules.MatchVariableSocketAddr,
		Parameters: rules.SocketAddrMatchConditionParameters{
			TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleSocketAddrConditionParameters,
			Operator:        rules.SocketAddrOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
		},
	}, nil
}

func expandCdnFrontDoorRuleSSLProtocolCondition(input helpers.CdnFrontDoorRuleConditionBaseModel) (rules.DeliveryRuleCondition, error) {
	operator, negated := helpers.ExpandCdnFrontDoorRuleConditionOperator(input.Operator)

	return rules.DeliveryRuleSslProtocolCondition{
		Name: rules.MatchVariableSslProtocol,
		Parameters: rules.SslProtocolMatchConditionParameters{
			TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleSslProtocolConditionParameters,
			Operator:        rules.SslProtocolOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.ToEnumSlice[rules.SslProtocol](input.Values),
		},
	}, nil
}

func expandCdnFrontDoorRuleRequestFileExtensionCondition(input helpers.CdnFrontDoorRuleConditionWithTransformsModel) (rules.DeliveryRuleCondition, error) {
	operator, negated := helpers.ExpandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := helpers.ValidateCdnFrontDoorRuleConditionValues("request_file_extension", operator, input.Values); err != nil {
		return nil, err
	}

	return rules.DeliveryRuleURLFileExtensionCondition{
		Name: rules.MatchVariableURLFileExtension,
		Parameters: rules.URLFileExtensionMatchConditionParameters{
			TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleURLFileExtensionMatchConditionParameters,
			Operator:        rules.URLFileExtensionOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
			Transforms:      pointer.ToEnumSlice[rules.Transform](input.Transforms),
		},
	}, nil
}

func expandCdnFrontDoorRuleRequestFilenameCondition(input helpers.CdnFrontDoorRuleConditionWithTransformsModel) (rules.DeliveryRuleCondition, error) {
	operator, negated := helpers.ExpandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := helpers.ValidateCdnFrontDoorRuleConditionValues("request_filename", operator, input.Values); err != nil {
		return nil, err
	}

	return rules.DeliveryRuleURLFileNameCondition{
		Name: rules.MatchVariableURLFileName,
		Parameters: rules.URLFileNameMatchConditionParameters{
			TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleURLFilenameConditionParameters,
			Operator:        rules.URLFileNameOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
			Transforms:      pointer.ToEnumSlice[rules.Transform](input.Transforms),
		},
	}, nil
}

func expandCdnFrontDoorRuleRequestPathCondition(input helpers.CdnFrontDoorRuleConditionWithTransformsModel) (rules.DeliveryRuleCondition, error) {
	operator, negated := helpers.ExpandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := helpers.ValidateCdnFrontDoorRuleConditionValues("request_path", operator, input.Values); err != nil {
		return nil, err
	}

	return rules.DeliveryRuleURLPathCondition{
		Name: rules.MatchVariableURLPath,
		Parameters: rules.URLPathMatchConditionParameters{
			TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleURLPathMatchConditionParameters,
			Operator:        rules.URLPathOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
			Transforms:      pointer.ToEnumSlice[rules.Transform](input.Transforms),
		},
	}, nil
}

func expandCdnFrontDoorRuleConditionBaseModel(input []helpers.CdnFrontDoorRuleConditionBaseModel, key string, expand func(helpers.CdnFrontDoorRuleConditionBaseModel) (rules.DeliveryRuleCondition, error)) ([]rules.DeliveryRuleCondition, error) {
	result := make([]rules.DeliveryRuleCondition, 0, len(input))
	for _, c := range input {
		expandedCondition, err := expand(c)
		if err != nil {
			return nil, fmt.Errorf("expanding `%s`: %+v", key, err)
		}
		result = append(result, expandedCondition)
	}
	return result, nil
}

func expandCdnFrontDoorRuleConditionWithNameAndTransformsModel(input []helpers.CdnFrontDoorRuleConditionWithNameAndTransformsModel, key string, expand func(helpers.CdnFrontDoorRuleConditionWithNameAndTransformsModel) (rules.DeliveryRuleCondition, error)) ([]rules.DeliveryRuleCondition, error) {
	result := make([]rules.DeliveryRuleCondition, 0, len(input))
	for _, c := range input {
		expandedCondition, err := expand(c)
		if err != nil {
			return nil, fmt.Errorf("expanding `%s`: %+v", key, err)
		}
		result = append(result, expandedCondition)
	}
	return result, nil
}

func expandCdnFrontDoorRuleConditionWithTransformsModel(input []helpers.CdnFrontDoorRuleConditionWithTransformsModel, key string, expand func(helpers.CdnFrontDoorRuleConditionWithTransformsModel) (rules.DeliveryRuleCondition, error)) ([]rules.DeliveryRuleCondition, error) {
	result := make([]rules.DeliveryRuleCondition, 0, len(input))
	for _, c := range input {
		expandedCondition, err := expand(c)
		if err != nil {
			return nil, fmt.Errorf("expanding `%s`: %+v", key, err)
		}
		result = append(result, expandedCondition)
	}
	return result, nil
}
