// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/afdorigins"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rulesets"
	helperValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type CdnFrontDoorBatchRuleSetModel struct {
	Name                  string                              `tfschema:"name"`
	CdnFrontdoorProfileID string                              `tfschema:"cdn_frontdoor_profile_id"`
	Rule                  []CdnFrontDoorBatchRuleSetRuleModel `tfschema:"rule"`
	ID                    string                              `tfschema:"id"`
}

type CdnFrontDoorBatchRuleSetRuleModel struct {
	Name             string                            `tfschema:"name"`
	BehaviourOnMatch string                            `tfschema:"behaviour_on_match"`
	Order            int64                             `tfschema:"order"`
	Actions          []CdnFrontDoorRuleActionsModel    `tfschema:"actions"`
	Conditions       []CdnFrontDoorRuleConditionsModel `tfschema:"conditions"`
}

func cdnFrontDoorBatchRuleSetActionModifyHeaderSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"operator": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(rulesets.PossibleValuesForHeaderAction(), false),
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

func cdnFrontDoorBatchRuleSetConditionTransformsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeSet,
		Optional: true,
		MaxItems: 4,
		Elem: &pluginsdk.Schema{
			Type:         pluginsdk.TypeString,
			ValidateFunc: validation.StringInSlice(rulesets.PossibleValuesForTransform(), false),
		},
	}
}

func expandCdnFrontDoorBatchRuleSetPayload(model CdnFrontDoorBatchRuleSetModel) (rulesets.RuleSet, error) {
	expandedRules, err := expandCdnFrontDoorBatchRuleSetRules(model.Rule)
	if err != nil {
		return rulesets.RuleSet{}, err
	}

	return rulesets.RuleSet{
		Properties: &rulesets.RuleSetProperties{
			BatchMode: pointer.To(true),
			Rules:     pointer.To(expandedRules),
		},
	}, nil
}

func expandCdnFrontDoorBatchRuleSetRules(input []CdnFrontDoorBatchRuleSetRuleModel) ([]rulesets.BatchRuleProperties, error) {
	results := make([]rulesets.BatchRuleProperties, 0, len(input))
	for _, item := range input {
		actions, err := expandCdnFrontDoorBatchRuleSetActions(item.Actions)
		if err != nil {
			return nil, err
		}

		conditions, err := expandCdnFrontDoorBatchRuleSetConditions(item.Conditions)
		if err != nil {
			return nil, err
		}

		results = append(results, rulesets.BatchRuleProperties{
			Actions:                 pointer.To(actions),
			Conditions:              pointer.To(conditions),
			MatchProcessingBehavior: pointer.ToEnum[rulesets.MatchProcessingBehavior](item.BehaviourOnMatch),
			Order:                   pointer.To(item.Order),
			RuleName:                item.Name,
		})
	}

	return results, nil
}

func flattenCdnFrontDoorBatchRuleSetRules(input *[]rulesets.BatchRuleProperties) ([]CdnFrontDoorBatchRuleSetRuleModel, error) {
	results := make([]CdnFrontDoorBatchRuleSetRuleModel, 0)
	if input == nil || len(*input) == 0 {
		return results, nil
	}

	for _, item := range *input {
		ruleState, err := flattenCdnFrontDoorBatchRuleSetRule(item)
		if err != nil {
			return results, err
		}
		results = append(results, ruleState)
	}

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Order < results[j].Order
	})

	return results, nil
}

func flattenCdnFrontDoorBatchRuleSetRule(input rulesets.BatchRuleProperties) (CdnFrontDoorBatchRuleSetRuleModel, error) {
	state := CdnFrontDoorBatchRuleSetRuleModel{
		BehaviourOnMatch: pointer.FromEnum(input.MatchProcessingBehavior),
		Name:             input.RuleName,
		Order:            pointer.From(input.Order),
	}

	actions, err := flattenCdnFrontDoorBatchRuleSetActions(input.Actions)
	if err != nil {
		return CdnFrontDoorBatchRuleSetRuleModel{}, fmt.Errorf("flattening `actions`: %+v", err)
	}
	state.Actions = actions

	conditions, err := flattenCdnFrontDoorBatchRuleSetConditions(input.Conditions)
	if err != nil {
		return CdnFrontDoorBatchRuleSetRuleModel{}, fmt.Errorf("flattening `conditions`: %+v", err)
	}
	state.Conditions = conditions

	return state, nil
}

func expandCdnFrontDoorBatchRuleSetActions(input []CdnFrontDoorRuleActionsModel) ([]rulesets.DeliveryRuleAction, error) {
	results := make([]rulesets.DeliveryRuleAction, 0)
	if len(input) == 0 {
		return results, nil
	}

	actions := input[0]
	if err := validateCdnFrontDoorRuleActionCounts(countCdnFrontDoorBatchRuleActions(input)); err != nil {
		return nil, err
	}

	for _, item := range actions.URLRedirect {
		results = append(results, rulesets.URLRedirectAction{
			Name: rulesets.DeliveryRuleActionNameURLRedirect,
			Parameters: rulesets.URLRedirectActionParameters{
				TypeName:            rulesets.DeliveryRuleActionParametersTypeDeliveryRuleURLRedirectActionParameters,
				RedirectType:        rulesets.RedirectType(item.RedirectType),
				DestinationProtocol: pointer.ToEnum[rulesets.DestinationProtocol](item.RedirectProtocol),
				CustomPath:          pointer.To(item.DestinationPath),
				CustomHostname:      pointer.To(item.DestinationHostName),
				CustomQueryString:   pointer.To(item.QueryString),
				CustomFragment:      pointer.To(item.DestinationFragment),
			},
		})
	}

	for _, item := range actions.URLRewrite {
		results = append(results, rulesets.URLRewriteAction{
			Name: rulesets.DeliveryRuleActionNameURLRewrite,
			Parameters: rulesets.URLRewriteActionParameters{
				TypeName:              rulesets.DeliveryRuleActionParametersTypeDeliveryRuleURLRewriteActionParameters,
				SourcePattern:         item.SourcePattern,
				Destination:           item.DestinationPath,
				PreserveUnmatchedPath: pointer.To(item.PreserveUnmatchedPathEnabled),
			},
		})
	}

	for _, item := range actions.ModifyRequestHeader {
		if err := validateCdnFrontDoorRuleModifyHeaderAction("modify_request_header", item.Operator, item.HeaderValue); err != nil {
			return nil, err
		}
		results = append(results, rulesets.DeliveryRuleRequestHeaderAction{
			Name: rulesets.DeliveryRuleActionNameModifyRequestHeader,
			Parameters: rulesets.HeaderActionParameters{
				TypeName:     rulesets.DeliveryRuleActionParametersTypeDeliveryRuleHeaderActionParameters,
				HeaderAction: rulesets.HeaderAction(item.Operator),
				HeaderName:   item.HeaderName,
				Value:        pointer.To(item.HeaderValue),
			},
		})
	}

	for _, item := range actions.ModifyResponseHeader {
		if err := validateCdnFrontDoorRuleModifyHeaderAction("modify_response_header", item.Operator, item.HeaderValue); err != nil {
			return nil, err
		}
		results = append(results, rulesets.DeliveryRuleResponseHeaderAction{
			Name: rulesets.DeliveryRuleActionNameModifyResponseHeader,
			Parameters: rulesets.HeaderActionParameters{
				TypeName:     rulesets.DeliveryRuleActionParametersTypeDeliveryRuleHeaderActionParameters,
				HeaderAction: rulesets.HeaderAction(item.Operator),
				HeaderName:   item.HeaderName,
				Value:        pointer.To(item.HeaderValue),
			},
		})
	}

	for _, item := range actions.RouteConfigurationOverride {
		expandedCache, err := expandCdnFrontDoorBatchRuleSetRouteConfigurationOverrideCaching(item.Caching)
		if err != nil {
			return nil, fmt.Errorf("expanding `route_configuration_override.0.caching`: %+v", err)
		}

		results = append(results, rulesets.DeliveryRuleRouteConfigurationOverrideAction{
			Name: rulesets.DeliveryRuleActionNameRouteConfigurationOverride,
			Parameters: rulesets.RouteConfigurationOverrideActionParameters{
				TypeName:            rulesets.DeliveryRuleActionParametersTypeDeliveryRuleRouteConfigurationOverrideActionParameters,
				CacheConfiguration:  expandedCache,
				OriginGroupOverride: expandCdnFrontDoorBatchRuleSetRouteConfigurationOverrideOriginGroup(item.OriginGroup),
			},
		})
	}

	return results, nil
}

func flattenCdnFrontDoorBatchRuleSetActions(input *[]rulesets.DeliveryRuleAction) ([]CdnFrontDoorRuleActionsModel, error) {
	if input == nil {
		return []CdnFrontDoorRuleActionsModel{}, nil
	}

	results := make([]CdnFrontDoorRuleActionsModel, 0, len(*input))
	actions := CdnFrontDoorRuleActionsModel{}

	for _, action := range *input {
		switch a := action.(type) {
		case rulesets.URLRedirectAction:
			actions.URLRedirect = append(actions.URLRedirect, CdnFrontDoorRuleURLRedirectActionModel{
				RedirectType:        string(a.Parameters.RedirectType),
				RedirectProtocol:    pointer.FromEnum(a.Parameters.DestinationProtocol),
				DestinationPath:     pointer.From(a.Parameters.CustomPath),
				DestinationHostName: pointer.From(a.Parameters.CustomHostname),
				QueryString:         pointer.From(a.Parameters.CustomQueryString),
				DestinationFragment: pointer.From(a.Parameters.CustomFragment),
			})
		case rulesets.URLRewriteAction:
			actions.URLRewrite = append(actions.URLRewrite, CdnFrontDoorRuleURLRewriteActionModel{
				SourcePattern:                a.Parameters.SourcePattern,
				DestinationPath:              a.Parameters.Destination,
				PreserveUnmatchedPathEnabled: pointer.From(a.Parameters.PreserveUnmatchedPath),
			})
		case rulesets.DeliveryRuleRequestHeaderAction:
			actions.ModifyRequestHeader = append(actions.ModifyRequestHeader, CdnFrontDoorRuleHeaderActionModel{
				Operator:    string(a.Parameters.HeaderAction),
				HeaderName:  a.Parameters.HeaderName,
				HeaderValue: pointer.From(a.Parameters.Value),
			})
		case rulesets.DeliveryRuleResponseHeaderAction:
			actions.ModifyResponseHeader = append(actions.ModifyResponseHeader, CdnFrontDoorRuleHeaderActionModel{
				Operator:    string(a.Parameters.HeaderAction),
				HeaderName:  a.Parameters.HeaderName,
				HeaderValue: pointer.From(a.Parameters.Value),
			})
		case rulesets.DeliveryRuleRouteConfigurationOverrideAction:
			flattened, err := flattenCdnFrontDoorBatchRuleSetRouteConfigurationOverrideAction(a.Parameters)
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

func flattenCdnFrontDoorBatchRuleSetRouteConfigurationOverrideAction(input rulesets.RouteConfigurationOverrideActionParameters) (CdnFrontDoorRuleRouteConfigurationOverrideActionModel, error) {
	result := CdnFrontDoorRuleRouteConfigurationOverrideActionModel{
		Caching: flattenCdnFrontDoorBatchRuleSetRouteConfigurationOverrideCaching(input.CacheConfiguration),
	}

	originGroup, err := flattenCdnFrontDoorBatchRuleSetRouteConfigurationOverrideOriginGroup(input.OriginGroupOverride)
	if err != nil {
		return result, err
	}
	result.OriginGroup = originGroup

	return result, nil
}

func expandCdnFrontDoorBatchRuleSetRouteConfigurationOverrideCaching(input []CdnFrontDoorRuleRouteConfigurationOverrideCachingModel) (*rulesets.CacheConfiguration, error) {
	if len(input) == 0 {
		return nil, nil
	}

	v := input[0]
	if v.Behaviour == RuleCacheBehaviorDisabled {
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
	if v.Behaviour == string(rulesets.RuleCacheBehaviorHonorOrigin) {
		if v.Duration != "" {
			return nil, fmt.Errorf("when `route_configuration_override.caching.behaviour` is set to `%s`, you cannot define `route_configuration_override.caching.duration`, got `%s`", rulesets.RuleCacheBehaviorHonorOrigin, v.Duration)
		}
	} else if v.Duration == "" {
		return nil, fmt.Errorf("when `route_configuration_override.caching.behaviour` is set to `%s`, you must also define `route_configuration_override.caching.duration`", v.Behaviour)
	}

	switch rulesets.RuleQueryStringCachingBehavior(v.QueryStringBehaviour) {
	case rulesets.RuleQueryStringCachingBehaviorIncludeSpecifiedQueryStrings, rulesets.RuleQueryStringCachingBehaviorIgnoreSpecifiedQueryStrings:
		if len(v.QueryStringParameters) == 0 {
			return nil, fmt.Errorf("when `route_configuration_override.caching.query_string_behaviour` is set to `%s`, you must also define one or more `route_configuration_override.caching.query_string_parameters`", v.QueryStringBehaviour)
		}
	case rulesets.RuleQueryStringCachingBehaviorUseQueryString, rulesets.RuleQueryStringCachingBehaviorIgnoreQueryString:
		if len(v.QueryStringParameters) > 0 {
			return nil, fmt.Errorf("when `route_configuration_override.caching.query_string_behaviour` is set to `%s`, you cannot define `route_configuration_override.caching.query_string_parameters`", v.QueryStringBehaviour)
		}
	}

	compressionEnabled := rulesets.RuleIsCompressionEnabledDisabled
	if v.CompressionEnabled {
		compressionEnabled = rulesets.RuleIsCompressionEnabledEnabled
	}

	queryParameters := (*string)(nil)
	if len(v.QueryStringParameters) > 0 {
		queryParameters = pointer.To(strings.Join(v.QueryStringParameters, ","))
	}

	return &rulesets.CacheConfiguration{
		CacheBehavior:              pointer.ToEnum[rulesets.RuleCacheBehavior](v.Behaviour),
		CacheDuration:              pointer.ToOrNil(v.Duration),
		IsCompressionEnabled:       pointer.To(compressionEnabled),
		QueryParameters:            queryParameters,
		QueryStringCachingBehavior: pointer.ToOrNil(rulesets.RuleQueryStringCachingBehavior(v.QueryStringBehaviour)),
	}, nil
}

func flattenCdnFrontDoorBatchRuleSetRouteConfigurationOverrideCaching(input *rulesets.CacheConfiguration) []CdnFrontDoorRuleRouteConfigurationOverrideCachingModel {
	result := make([]CdnFrontDoorRuleRouteConfigurationOverrideCachingModel, 0)
	if input == nil {
		// The API treats omission as disabled, so we'll need to set `Disabled` back into state if it's nil
		return append(result, CdnFrontDoorRuleRouteConfigurationOverrideCachingModel{
			Behaviour: string(rulesets.RuleIsCompressionEnabledDisabled),
		})
	}

	v := CdnFrontDoorRuleRouteConfigurationOverrideCachingModel{
		Behaviour:            pointer.FromEnum(input.CacheBehavior),
		Duration:             pointer.FromEnum(input.CacheDuration),
		CompressionEnabled:   pointer.From(input.IsCompressionEnabled) == rulesets.RuleIsCompressionEnabledEnabled,
		QueryStringBehaviour: pointer.FromEnum(input.QueryStringCachingBehavior),
	}

	if input.QueryParameters != nil {
		v.QueryStringParameters = strings.Split(*input.QueryParameters, ",")
	}

	return append(result, v)
}

func expandCdnFrontDoorBatchRuleSetRouteConfigurationOverrideOriginGroup(input []CdnFrontDoorRuleRouteConfigurationOverrideOriginGroupModel) *rulesets.OriginGroupOverride {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	return &rulesets.OriginGroupOverride{
		ForwardingProtocol: pointer.ToEnum[rulesets.ForwardingProtocol](v.ForwardingProtocol),
		OriginGroup: &rulesets.ResourceReference{
			Id: pointer.To(v.CdnFrontdoorOriginGroupID),
		},
	}
}

func flattenCdnFrontDoorBatchRuleSetRouteConfigurationOverrideOriginGroup(input *rulesets.OriginGroupOverride) ([]CdnFrontDoorRuleRouteConfigurationOverrideOriginGroupModel, error) {
	result := make([]CdnFrontDoorRuleRouteConfigurationOverrideOriginGroupModel, 0)
	if input == nil {
		return result, nil
	}

	v := CdnFrontDoorRuleRouteConfigurationOverrideOriginGroupModel{
		ForwardingProtocol: pointer.FromEnum(input.ForwardingProtocol),
	}

	if input.OriginGroup != nil && input.OriginGroup.Id != nil {
		// Azure can return IDs with inconsistently cased static segments; normalize them to avoid state drift and refresh failures (#32953).
		originGroupID, err := afdorigins.ParseOriginGroupIDInsensitively(*input.OriginGroup.Id)
		if err != nil {
			return result, err
		}
		v.CdnFrontdoorOriginGroupID = originGroupID.ID()
	}

	return append(result, v), nil
}

func expandCdnFrontDoorBatchRuleSetConditions(input []CdnFrontDoorRuleConditionsModel) ([]rulesets.DeliveryRuleCondition, error) {
	results := make([]rulesets.DeliveryRuleCondition, 0)
	if len(input) == 0 {
		return results, nil
	}

	conditions := input[0]

	expanded, err := expandCdnFrontDoorBatchRuleSetConditionBaseModel(conditions.RemoteAddress, "remote_address", expandCdnFrontDoorBatchRuleSetRemoteAddressCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorBatchRuleSetConditionWithTransformsModel(conditions.QueryString, "query_string", expandCdnFrontDoorBatchRuleSetQueryStringCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorBatchRuleSetConditionWithTransformsModel(conditions.RequestURL, "request_url", expandCdnFrontDoorBatchRuleSetRequestURLCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorBatchRuleSetConditionWithTransformsModel(conditions.RequestBody, "request_body", expandCdnFrontDoorBatchRuleSetRequestBodyCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorBatchRuleSetConditionWithTransformsModel(conditions.RequestPath, "request_path", expandCdnFrontDoorBatchRuleSetRequestPathCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorBatchRuleSetConditionWithTransformsModel(conditions.RequestFileExtension, "request_file_extension", expandCdnFrontDoorBatchRuleSetRequestFileExtensionCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorBatchRuleSetConditionWithTransformsModel(conditions.RequestFilename, "request_filename", expandCdnFrontDoorBatchRuleSetRequestFilenameCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorBatchRuleSetConditionBaseModel(conditions.SocketAddress, "socket_address", expandCdnFrontDoorBatchRuleSetSocketAddressCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorBatchRuleSetConditionBaseModel(conditions.ClientPort, "client_port", expandCdnFrontDoorBatchRuleSetClientPortCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorBatchRuleSetConditionBaseModel(conditions.ServerPort, "server_port", expandCdnFrontDoorBatchRuleSetServerPortCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorBatchRuleSetConditionWithTransformsModel(conditions.HostName, "host_name", expandCdnFrontDoorBatchRuleSetHostNameCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorBatchRuleSetConditionBaseModel(conditions.RequestMethod, "request_method", expandCdnFrontDoorBatchRuleSetRequestMethodCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorBatchRuleSetConditionWithNameAndTransformsModel(conditions.PostArgs, "post_argument", expandCdnFrontDoorBatchRuleSetPostArgsCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorBatchRuleSetConditionWithNameAndTransformsModel(conditions.RequestHeader, "request_header", expandCdnFrontDoorBatchRuleSetRequestHeaderCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorBatchRuleSetConditionBaseModel(conditions.RequestScheme, "request_scheme", expandCdnFrontDoorBatchRuleSetRequestSchemeCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorBatchRuleSetConditionBaseModel(conditions.HTTPVersion, "http_version", expandCdnFrontDoorBatchRuleSetHTTPVersionCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorBatchRuleSetConditionWithNameAndTransformsModel(conditions.RequestCookies, "request_cookies", expandCdnFrontDoorBatchRuleSetCookiesCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorBatchRuleSetConditionBaseModel(conditions.DeviceType, "device_type", expandCdnFrontDoorBatchRuleSetDeviceTypeCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	expanded, err = expandCdnFrontDoorBatchRuleSetConditionBaseModel(conditions.SSLProtocol, "ssl_protocol", expandCdnFrontDoorBatchRuleSetSSLProtocolCondition)
	if err != nil {
		return nil, err
	}
	results = append(results, expanded...)

	if len(results) > 10 {
		return nil, fmt.Errorf("the `conditions` block may only contain up to 10 match conditions, got %d", len(results))
	}

	return results, nil
}

func flattenCdnFrontDoorBatchRuleSetConditions(input *[]rulesets.DeliveryRuleCondition) ([]CdnFrontDoorRuleConditionsModel, error) {
	if input == nil || len(*input) == 0 {
		return []CdnFrontDoorRuleConditionsModel{}, nil
	}

	conditions := CdnFrontDoorRuleConditionsModel{}
	for _, condition := range *input {
		switch c := condition.(type) {
		case rulesets.DeliveryRuleRemoteAddressCondition:
			conditions.RemoteAddress = append(conditions.RemoteAddress, flattenCdnFrontDoorRuleConditionBaseModel(string(c.Parameters.Operator), pointer.From(c.Parameters.MatchValues), c.Parameters.NegateCondition))
		case rulesets.DeliveryRuleRequestMethodCondition:
			conditions.RequestMethod = append(conditions.RequestMethod, flattenCdnFrontDoorRuleConditionBaseModel(string(c.Parameters.Operator), pointer.FromEnumSlice(c.Parameters.MatchValues), c.Parameters.NegateCondition))
		case rulesets.DeliveryRuleQueryStringCondition:
			conditions.QueryString = append(conditions.QueryString, flattenCdnFrontDoorBatchRuleSetConditionWithTransformsModel(string(c.Parameters.Operator), pointer.From(c.Parameters.MatchValues), pointer.FromEnumSlice(c.Parameters.Transforms), c.Parameters.NegateCondition))
		case rulesets.DeliveryRulePostArgsCondition:
			conditions.PostArgs = append(conditions.PostArgs, flattenCdnFrontDoorRuleConditionWithNameAndTransformsModel(string(c.Parameters.Operator), c.Parameters.Selector, pointer.From(c.Parameters.MatchValues), pointer.FromEnumSlice(c.Parameters.Transforms), c.Parameters.NegateCondition))
		case rulesets.DeliveryRuleRequestUriCondition:
			conditions.RequestURL = append(conditions.RequestURL, flattenCdnFrontDoorBatchRuleSetConditionWithTransformsModel(string(c.Parameters.Operator), pointer.From(c.Parameters.MatchValues), pointer.FromEnumSlice(c.Parameters.Transforms), c.Parameters.NegateCondition))
		case rulesets.DeliveryRuleRequestHeaderCondition:
			conditions.RequestHeader = append(conditions.RequestHeader, flattenCdnFrontDoorRuleConditionWithNameAndTransformsModel(string(c.Parameters.Operator), c.Parameters.Selector, pointer.From(c.Parameters.MatchValues), pointer.FromEnumSlice(c.Parameters.Transforms), c.Parameters.NegateCondition))
		case rulesets.DeliveryRuleRequestBodyCondition:
			conditions.RequestBody = append(conditions.RequestBody, flattenCdnFrontDoorBatchRuleSetConditionWithTransformsModel(string(c.Parameters.Operator), pointer.From(c.Parameters.MatchValues), pointer.FromEnumSlice(c.Parameters.Transforms), c.Parameters.NegateCondition))
		case rulesets.DeliveryRuleRequestSchemeCondition:
			conditions.RequestScheme = append(conditions.RequestScheme, flattenCdnFrontDoorRuleConditionBaseModel(string(c.Parameters.Operator), pointer.FromEnumSlice(c.Parameters.MatchValues), c.Parameters.NegateCondition))
		case rulesets.DeliveryRuleURLPathCondition:
			conditions.RequestPath = append(conditions.RequestPath, flattenCdnFrontDoorBatchRuleSetConditionWithTransformsModel(string(c.Parameters.Operator), pointer.From(c.Parameters.MatchValues), pointer.FromEnumSlice(c.Parameters.Transforms), c.Parameters.NegateCondition))
		case rulesets.DeliveryRuleURLFileExtensionCondition:
			conditions.RequestFileExtension = append(conditions.RequestFileExtension, flattenCdnFrontDoorBatchRuleSetConditionWithTransformsModel(string(c.Parameters.Operator), pointer.From(c.Parameters.MatchValues), pointer.FromEnumSlice(c.Parameters.Transforms), c.Parameters.NegateCondition))
		case rulesets.DeliveryRuleURLFileNameCondition:
			conditions.RequestFilename = append(conditions.RequestFilename, flattenCdnFrontDoorBatchRuleSetConditionWithTransformsModel(string(c.Parameters.Operator), pointer.From(c.Parameters.MatchValues), pointer.FromEnumSlice(c.Parameters.Transforms), c.Parameters.NegateCondition))
		case rulesets.DeliveryRuleHTTPVersionCondition:
			conditions.HTTPVersion = append(conditions.HTTPVersion, flattenCdnFrontDoorRuleConditionBaseModel(string(c.Parameters.Operator), pointer.FromEnumSlice(c.Parameters.MatchValues), c.Parameters.NegateCondition))
		case rulesets.DeliveryRuleCookiesCondition:
			conditions.RequestCookies = append(conditions.RequestCookies, flattenCdnFrontDoorRuleConditionWithNameAndTransformsModel(string(c.Parameters.Operator), c.Parameters.Selector, pointer.From(c.Parameters.MatchValues), pointer.FromEnumSlice(c.Parameters.Transforms), c.Parameters.NegateCondition))
		case rulesets.DeliveryRuleIsDeviceCondition:
			conditions.DeviceType = append(conditions.DeviceType, flattenCdnFrontDoorRuleConditionBaseModel(string(c.Parameters.Operator), pointer.FromEnumSlice(c.Parameters.MatchValues), c.Parameters.NegateCondition))
		case rulesets.DeliveryRuleSocketAddrCondition:
			conditions.SocketAddress = append(conditions.SocketAddress, flattenCdnFrontDoorRuleConditionBaseModel(string(c.Parameters.Operator), pointer.From(c.Parameters.MatchValues), c.Parameters.NegateCondition))
		case rulesets.DeliveryRuleClientPortCondition:
			conditions.ClientPort = append(conditions.ClientPort, flattenCdnFrontDoorRuleConditionBaseModel(string(c.Parameters.Operator), pointer.From(c.Parameters.MatchValues), c.Parameters.NegateCondition))
		case rulesets.DeliveryRuleServerPortCondition:
			conditions.ServerPort = append(conditions.ServerPort, flattenCdnFrontDoorRuleConditionBaseModel(string(c.Parameters.Operator), pointer.From(c.Parameters.MatchValues), c.Parameters.NegateCondition))
		case rulesets.DeliveryRuleHostNameCondition:
			conditions.HostName = append(conditions.HostName, flattenCdnFrontDoorBatchRuleSetConditionWithTransformsModel(string(c.Parameters.Operator), pointer.From(c.Parameters.MatchValues), pointer.FromEnumSlice(c.Parameters.Transforms), c.Parameters.NegateCondition))
		case rulesets.DeliveryRuleSslProtocolCondition:
			conditions.SSLProtocol = append(conditions.SSLProtocol, flattenCdnFrontDoorRuleConditionBaseModel(string(c.Parameters.Operator), pointer.FromEnumSlice(c.Parameters.MatchValues), c.Parameters.NegateCondition))
		default:
			return []CdnFrontDoorRuleConditionsModel{}, fmt.Errorf("unsupported condition (`%s`) encountered", condition.DeliveryRuleCondition().Name)
		}
	}

	return []CdnFrontDoorRuleConditionsModel{conditions}, nil
}

func expandCdnFrontDoorBatchRuleSetClientPortCondition(input CdnFrontDoorRuleConditionBaseModel) (rulesets.DeliveryRuleCondition, error) {
	operator, negated := expandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := validateCdnFrontDoorRuleConditionValues("client_port", operator, input.Values); err != nil {
		return nil, err
	}

	return rulesets.DeliveryRuleClientPortCondition{
		Name: rulesets.MatchVariableClientPort,
		Parameters: rulesets.ClientPortMatchConditionParameters{
			TypeName:        rulesets.DeliveryRuleConditionParametersTypeDeliveryRuleClientPortConditionParameters,
			Operator:        rulesets.ClientPortOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
		},
	}, nil
}

func expandCdnFrontDoorBatchRuleSetCookiesCondition(input CdnFrontDoorRuleConditionWithNameAndTransformsModel) (rulesets.DeliveryRuleCondition, error) {
	operator, negated := expandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := validateCdnFrontDoorRuleConditionValues("request_cookies", operator, input.Values); err != nil {
		return nil, err
	}

	return rulesets.DeliveryRuleCookiesCondition{
		Name: rulesets.MatchVariableCookies,
		Parameters: rulesets.CookiesMatchConditionParameters{
			TypeName:        rulesets.DeliveryRuleConditionParametersTypeDeliveryRuleCookiesConditionParameters,
			Selector:        pointer.To(input.Name),
			Operator:        rulesets.CookiesOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
			Transforms:      pointer.ToEnumSlice[rulesets.Transform](input.Transforms),
		},
	}, nil
}

func expandCdnFrontDoorBatchRuleSetHostNameCondition(input CdnFrontDoorRuleConditionWithTransformsModel) (rulesets.DeliveryRuleCondition, error) {
	operator, negated := expandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := validateCdnFrontDoorRuleConditionValues("host_name", operator, input.Values); err != nil {
		return nil, err
	}

	return rulesets.DeliveryRuleHostNameCondition{
		Name: rulesets.MatchVariableHostName,
		Parameters: rulesets.HostNameMatchConditionParameters{
			TypeName:        rulesets.DeliveryRuleConditionParametersTypeDeliveryRuleHostNameConditionParameters,
			Operator:        rulesets.HostNameOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
			Transforms:      pointer.ToEnumSlice[rulesets.Transform](input.Transforms),
		},
	}, nil
}

func expandCdnFrontDoorBatchRuleSetHTTPVersionCondition(input CdnFrontDoorRuleConditionBaseModel) (rulesets.DeliveryRuleCondition, error) {
	operator, negated := expandCdnFrontDoorRuleConditionOperator(input.Operator)

	return rulesets.DeliveryRuleHTTPVersionCondition{
		Name: rulesets.MatchVariableHTTPVersion,
		Parameters: rulesets.HTTPVersionMatchConditionParameters{
			TypeName:        rulesets.DeliveryRuleConditionParametersTypeDeliveryRuleHTTPVersionConditionParameters,
			Operator:        rulesets.HTTPVersionOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
		},
	}, nil
}

func expandCdnFrontDoorBatchRuleSetDeviceTypeCondition(input CdnFrontDoorRuleConditionBaseModel) (rulesets.DeliveryRuleCondition, error) {
	operator, negated := expandCdnFrontDoorRuleConditionOperator(input.Operator)

	return rulesets.DeliveryRuleIsDeviceCondition{
		Name: rulesets.MatchVariableIsDevice,
		Parameters: rulesets.IsDeviceMatchConditionParameters{
			TypeName:        rulesets.DeliveryRuleConditionParametersTypeDeliveryRuleIsDeviceConditionParameters,
			Operator:        rulesets.IsDeviceOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.ToEnumSlice[rulesets.IsDeviceMatchValue](input.Values),
		},
	}, nil
}

func expandCdnFrontDoorBatchRuleSetPostArgsCondition(input CdnFrontDoorRuleConditionWithNameAndTransformsModel) (rulesets.DeliveryRuleCondition, error) {
	operator, negated := expandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := validateCdnFrontDoorRuleConditionValues("post_argument", operator, input.Values); err != nil {
		return nil, err
	}

	return rulesets.DeliveryRulePostArgsCondition{
		Name: rulesets.MatchVariablePostArgs,
		Parameters: rulesets.PostArgsMatchConditionParameters{
			TypeName:        rulesets.DeliveryRuleConditionParametersTypeDeliveryRulePostArgsConditionParameters,
			Selector:        pointer.To(input.Name),
			Operator:        rulesets.PostArgsOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
			Transforms:      pointer.ToEnumSlice[rulesets.Transform](input.Transforms),
		},
	}, nil
}

func expandCdnFrontDoorBatchRuleSetQueryStringCondition(input CdnFrontDoorRuleConditionWithTransformsModel) (rulesets.DeliveryRuleCondition, error) {
	operator, negated := expandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := validateCdnFrontDoorRuleConditionValues("query_string", operator, input.Values); err != nil {
		return nil, err
	}

	return rulesets.DeliveryRuleQueryStringCondition{
		Name: rulesets.MatchVariableQueryString,
		Parameters: rulesets.QueryStringMatchConditionParameters{
			TypeName:        rulesets.DeliveryRuleConditionParametersTypeDeliveryRuleQueryStringConditionParameters,
			Operator:        rulesets.QueryStringOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
			Transforms:      pointer.ToEnumSlice[rulesets.Transform](input.Transforms),
		},
	}, nil
}

func expandCdnFrontDoorBatchRuleSetRemoteAddressCondition(input CdnFrontDoorRuleConditionBaseModel) (rulesets.DeliveryRuleCondition, error) {
	operator, negated := expandCdnFrontDoorRuleConditionOperator(input.Operator)

	switch rulesets.RemoteAddressOperator(operator) {
	case rulesets.RemoteAddressOperatorGeoMatch:
		for _, v := range input.Values {
			if ok, _ := helperValidate.RegExHelper(v, "values", `^[A-Z]{2}$`); !ok {
				return nil, fmt.Errorf("when `conditions.remote_address.operator` is `%s` the values in `conditions.remote_address.values` must be valid country codes consisting of 2 uppercase characters, got `%s`", input.Operator, v)
			}
		}
	case rulesets.RemoteAddressOperatorIPMatch:
		values := make([]interface{}, 0, len(input.Values))
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

	return rulesets.DeliveryRuleRemoteAddressCondition{
		Name: rulesets.MatchVariableRemoteAddress,
		Parameters: rulesets.RemoteAddressMatchConditionParameters{
			TypeName:        rulesets.DeliveryRuleConditionParametersTypeDeliveryRuleRemoteAddressConditionParameters,
			Operator:        rulesets.RemoteAddressOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
		},
	}, nil
}

func expandCdnFrontDoorBatchRuleSetRequestBodyCondition(input CdnFrontDoorRuleConditionWithTransformsModel) (rulesets.DeliveryRuleCondition, error) {
	operator, negated := expandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := validateCdnFrontDoorRuleConditionValues("request_body", operator, input.Values); err != nil {
		return nil, err
	}

	return rulesets.DeliveryRuleRequestBodyCondition{
		Name: rulesets.MatchVariableRequestBody,
		Parameters: rulesets.RequestBodyMatchConditionParameters{
			TypeName:        rulesets.DeliveryRuleConditionParametersTypeDeliveryRuleRequestBodyConditionParameters,
			Operator:        rulesets.RequestBodyOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
			Transforms:      pointer.ToEnumSlice[rulesets.Transform](input.Transforms),
		},
	}, nil
}

func expandCdnFrontDoorBatchRuleSetRequestHeaderCondition(input CdnFrontDoorRuleConditionWithNameAndTransformsModel) (rulesets.DeliveryRuleCondition, error) {
	operator, negated := expandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := validateCdnFrontDoorRuleConditionValues("request_header", operator, input.Values); err != nil {
		return nil, err
	}

	return rulesets.DeliveryRuleRequestHeaderCondition{
		Name: rulesets.MatchVariableRequestHeader,
		Parameters: rulesets.RequestHeaderMatchConditionParameters{
			TypeName:        rulesets.DeliveryRuleConditionParametersTypeDeliveryRuleRequestHeaderConditionParameters,
			Selector:        pointer.To(input.Name),
			Operator:        rulesets.RequestHeaderOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
			Transforms:      pointer.ToEnumSlice[rulesets.Transform](input.Transforms),
		},
	}, nil
}

func expandCdnFrontDoorBatchRuleSetRequestMethodCondition(input CdnFrontDoorRuleConditionBaseModel) (rulesets.DeliveryRuleCondition, error) {
	operator, negated := expandCdnFrontDoorRuleConditionOperator(input.Operator)

	return rulesets.DeliveryRuleRequestMethodCondition{
		Name: rulesets.MatchVariableRequestMethod,
		Parameters: rulesets.RequestMethodMatchConditionParameters{
			TypeName:        rulesets.DeliveryRuleConditionParametersTypeDeliveryRuleRequestMethodConditionParameters,
			Operator:        rulesets.RequestMethodOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.ToEnumSlice[rulesets.RequestMethodMatchValue](input.Values),
		},
	}, nil
}

func expandCdnFrontDoorBatchRuleSetRequestSchemeCondition(input CdnFrontDoorRuleConditionBaseModel) (rulesets.DeliveryRuleCondition, error) {
	operator, negated := expandCdnFrontDoorRuleConditionOperator(input.Operator)

	return rulesets.DeliveryRuleRequestSchemeCondition{
		Name: rulesets.MatchVariableRequestScheme,
		Parameters: rulesets.RequestSchemeMatchConditionParameters{
			TypeName:        rulesets.DeliveryRuleConditionParametersTypeDeliveryRuleRequestSchemeConditionParameters,
			Operator:        rulesets.RequestSchemeMatchConditionParametersOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.ToEnumSlice[rulesets.RequestSchemeMatchValue](input.Values),
		},
	}, nil
}

func expandCdnFrontDoorBatchRuleSetRequestURLCondition(input CdnFrontDoorRuleConditionWithTransformsModel) (rulesets.DeliveryRuleCondition, error) {
	operator, negated := expandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := validateCdnFrontDoorRuleConditionValues("request_url", operator, input.Values); err != nil {
		return nil, err
	}

	return rulesets.DeliveryRuleRequestUriCondition{
		Name: rulesets.MatchVariableRequestUri,
		Parameters: rulesets.RequestUriMatchConditionParameters{
			TypeName:        rulesets.DeliveryRuleConditionParametersTypeDeliveryRuleRequestUriConditionParameters,
			Operator:        rulesets.RequestUriOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
			Transforms:      pointer.ToEnumSlice[rulesets.Transform](input.Transforms),
		},
	}, nil
}

func expandCdnFrontDoorBatchRuleSetServerPortCondition(input CdnFrontDoorRuleConditionBaseModel) (rulesets.DeliveryRuleCondition, error) {
	operator, negated := expandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := validateCdnFrontDoorRuleConditionValues("server_port", operator, input.Values); err != nil {
		return nil, err
	}

	return rulesets.DeliveryRuleServerPortCondition{
		Name: rulesets.MatchVariableServerPort,
		Parameters: rulesets.ServerPortMatchConditionParameters{
			TypeName:        rulesets.DeliveryRuleConditionParametersTypeDeliveryRuleServerPortConditionParameters,
			Operator:        rulesets.ServerPortOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
		},
	}, nil
}

func expandCdnFrontDoorBatchRuleSetSocketAddressCondition(input CdnFrontDoorRuleConditionBaseModel) (rulesets.DeliveryRuleCondition, error) {
	operator, negated := expandCdnFrontDoorRuleConditionOperator(input.Operator)

	return rulesets.DeliveryRuleSocketAddrCondition{
		Name: rulesets.MatchVariableSocketAddr,
		Parameters: rulesets.SocketAddrMatchConditionParameters{
			TypeName:        rulesets.DeliveryRuleConditionParametersTypeDeliveryRuleSocketAddrConditionParameters,
			Operator:        rulesets.SocketAddrOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
		},
	}, nil
}

func expandCdnFrontDoorBatchRuleSetSSLProtocolCondition(input CdnFrontDoorRuleConditionBaseModel) (rulesets.DeliveryRuleCondition, error) {
	operator, negated := expandCdnFrontDoorRuleConditionOperator(input.Operator)

	return rulesets.DeliveryRuleSslProtocolCondition{
		Name: rulesets.MatchVariableSslProtocol,
		Parameters: rulesets.SslProtocolMatchConditionParameters{
			TypeName:        rulesets.DeliveryRuleConditionParametersTypeDeliveryRuleSslProtocolConditionParameters,
			Operator:        rulesets.SslProtocolOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.ToEnumSlice[rulesets.SslProtocol](input.Values),
		},
	}, nil
}

func expandCdnFrontDoorBatchRuleSetRequestFileExtensionCondition(input CdnFrontDoorRuleConditionWithTransformsModel) (rulesets.DeliveryRuleCondition, error) {
	operator, negated := expandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := validateCdnFrontDoorRuleConditionValues("request_file_extension", operator, input.Values); err != nil {
		return nil, err
	}

	return rulesets.DeliveryRuleURLFileExtensionCondition{
		Name: rulesets.MatchVariableURLFileExtension,
		Parameters: rulesets.URLFileExtensionMatchConditionParameters{
			TypeName:        rulesets.DeliveryRuleConditionParametersTypeDeliveryRuleURLFileExtensionMatchConditionParameters,
			Operator:        rulesets.URLFileExtensionOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
			Transforms:      pointer.ToEnumSlice[rulesets.Transform](input.Transforms),
		},
	}, nil
}

func expandCdnFrontDoorBatchRuleSetRequestFilenameCondition(input CdnFrontDoorRuleConditionWithTransformsModel) (rulesets.DeliveryRuleCondition, error) {
	operator, negated := expandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := validateCdnFrontDoorRuleConditionValues("request_filename", operator, input.Values); err != nil {
		return nil, err
	}

	return rulesets.DeliveryRuleURLFileNameCondition{
		Name: rulesets.MatchVariableURLFileName,
		Parameters: rulesets.URLFileNameMatchConditionParameters{
			TypeName:        rulesets.DeliveryRuleConditionParametersTypeDeliveryRuleURLFilenameConditionParameters,
			Operator:        rulesets.URLFileNameOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
			Transforms:      pointer.ToEnumSlice[rulesets.Transform](input.Transforms),
		},
	}, nil
}

func expandCdnFrontDoorBatchRuleSetRequestPathCondition(input CdnFrontDoorRuleConditionWithTransformsModel) (rulesets.DeliveryRuleCondition, error) {
	operator, negated := expandCdnFrontDoorRuleConditionOperator(input.Operator)
	if err := validateCdnFrontDoorRuleConditionValues("request_path", operator, input.Values); err != nil {
		return nil, err
	}

	return rulesets.DeliveryRuleURLPathCondition{
		Name: rulesets.MatchVariableURLPath,
		Parameters: rulesets.URLPathMatchConditionParameters{
			TypeName:        rulesets.DeliveryRuleConditionParametersTypeDeliveryRuleURLPathMatchConditionParameters,
			Operator:        rulesets.URLPathOperator(operator),
			NegateCondition: pointer.To(negated),
			MatchValues:     pointer.To(input.Values),
			Transforms:      pointer.ToEnumSlice[rulesets.Transform](input.Transforms),
		},
	}, nil
}

func expandCdnFrontDoorBatchRuleSetConditionBaseModel(input []CdnFrontDoorRuleConditionBaseModel, key string, expand func(CdnFrontDoorRuleConditionBaseModel) (rulesets.DeliveryRuleCondition, error)) ([]rulesets.DeliveryRuleCondition, error) {
	result := make([]rulesets.DeliveryRuleCondition, 0, len(input))
	for _, c := range input {
		expandedCondition, err := expand(c)
		if err != nil {
			return nil, fmt.Errorf("expanding `%s`: %+v", key, err)
		}
		result = append(result, expandedCondition)
	}
	return result, nil
}

func expandCdnFrontDoorBatchRuleSetConditionWithNameAndTransformsModel(input []CdnFrontDoorRuleConditionWithNameAndTransformsModel, key string, expand func(CdnFrontDoorRuleConditionWithNameAndTransformsModel) (rulesets.DeliveryRuleCondition, error)) ([]rulesets.DeliveryRuleCondition, error) {
	result := make([]rulesets.DeliveryRuleCondition, 0, len(input))
	for _, c := range input {
		expandedCondition, err := expand(c)
		if err != nil {
			return nil, fmt.Errorf("expanding `%s`: %+v", key, err)
		}
		result = append(result, expandedCondition)
	}
	return result, nil
}

func expandCdnFrontDoorBatchRuleSetConditionWithTransformsModel(input []CdnFrontDoorRuleConditionWithTransformsModel, key string, expand func(CdnFrontDoorRuleConditionWithTransformsModel) (rulesets.DeliveryRuleCondition, error)) ([]rulesets.DeliveryRuleCondition, error) {
	result := make([]rulesets.DeliveryRuleCondition, 0, len(input))
	for _, c := range input {
		expandedCondition, err := expand(c)
		if err != nil {
			return nil, fmt.Errorf("expanding `%s`: %+v", key, err)
		}
		result = append(result, expandedCondition)
	}
	return result, nil
}

// Validation

func validateCdnFrontDoorBatchRules(input []CdnFrontDoorBatchRuleSetRuleModel) error {
	names := make(map[string]struct{}, len(input))
	orders := make(map[int64]struct{}, len(input))
	lastOrder := int64(math.MinInt64)

	for _, item := range input {
		if _, exists := names[item.Name]; exists {
			return fmt.Errorf("the `rule` blocks must have unique `name` values, got duplicate `%s`", item.Name)
		}
		names[item.Name] = struct{}{}

		if _, exists := orders[item.Order]; exists {
			return fmt.Errorf("the `rule` blocks must have unique `order` values, got duplicate `%d`", item.Order)
		}
		orders[item.Order] = struct{}{}

		if item.Order < lastOrder {
			return fmt.Errorf("the `rule` blocks must be declared in ascending `order`, got `%d` before `%d`", lastOrder, item.Order)
		}

		if err := validateCdnFrontDoorRuleActionCounts(countCdnFrontDoorBatchRuleActions(item.Actions)); err != nil {
			return err
		}

		lastOrder = item.Order
	}

	return nil
}
