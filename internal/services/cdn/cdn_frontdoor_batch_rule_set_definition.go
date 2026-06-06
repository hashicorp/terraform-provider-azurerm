// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cdn

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/profiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2025-12-01/rules"
	helperValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/azuresdkhacks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type CdnFrontDoorBatchRuleSetModel struct {
	Name                  string                           `tfschema:"name"`
	CdnFrontDoorProfileID string                           `tfschema:"cdn_frontdoor_profile_id"`
	Rules                 []CdnFrontDoorBatchRuleRuleModel `tfschema:"rules"`
	ID                    string                           `tfschema:"id"`
}

type CdnFrontDoorBatchRuleRuleModel struct {
	Name            string                                 `tfschema:"name"`
	BehaviorOnMatch string                                 `tfschema:"behavior_on_match"`
	Order           int64                                  `tfschema:"order"`
	Actions         []CdnFrontDoorBatchRuleActionsModel    `tfschema:"actions"`
	Conditions      []CdnFrontDoorBatchRuleConditionsModel `tfschema:"conditions"`
}

type CdnFrontDoorBatchRuleActionsModel struct {
	URLRedirectAction                []CdnFrontDoorBatchRuleURLRedirectActionModel                `tfschema:"url_redirect_action"`
	URLRewriteAction                 []CdnFrontDoorBatchRuleURLRewriteActionModel                 `tfschema:"url_rewrite_action"`
	RequestHeaderAction              []CdnFrontDoorBatchRuleHeaderActionModel                     `tfschema:"request_header_action"`
	ResponseHeaderAction             []CdnFrontDoorBatchRuleHeaderActionModel                     `tfschema:"response_header_action"`
	RouteConfigurationOverrideAction []CdnFrontDoorBatchRuleRouteConfigurationOverrideActionModel `tfschema:"route_configuration_override_action"`
}

type CdnFrontDoorBatchRuleURLRedirectActionModel struct {
	RedirectType        string `tfschema:"redirect_type"`
	RedirectProtocol    string `tfschema:"redirect_protocol"`
	DestinationPath     string `tfschema:"destination_path"`
	DestinationHostname string `tfschema:"destination_hostname"`
	QueryString         string `tfschema:"query_string"`
	DestinationFragment string `tfschema:"destination_fragment"`
}

type CdnFrontDoorBatchRuleURLRewriteActionModel struct {
	SourcePattern         string `tfschema:"source_pattern"`
	Destination           string `tfschema:"destination"`
	PreserveUnmatchedPath bool   `tfschema:"preserve_unmatched_path"`
}

type CdnFrontDoorBatchRuleHeaderActionModel struct {
	HeaderAction string `tfschema:"header_action"`
	HeaderName   string `tfschema:"header_name"`
	Value        string `tfschema:"value"`
}

type CdnFrontDoorBatchRuleRouteConfigurationOverrideActionModel struct {
	CdnFrontDoorOriginGroupID  string   `tfschema:"cdn_frontdoor_origin_group_id"`
	ForwardingProtocol         string   `tfschema:"forwarding_protocol"`
	QueryStringCachingBehavior string   `tfschema:"query_string_caching_behavior"`
	QueryStringParameters      []string `tfschema:"query_string_parameters"`
	CompressionEnabled         bool     `tfschema:"compression_enabled"`
	CacheBehavior              string   `tfschema:"cache_behavior"`
	CacheDuration              string   `tfschema:"cache_duration"`
}

type CdnFrontDoorBatchRuleConditionsModel struct {
	RemoteAddressCondition    []CdnFrontDoorBatchRuleStringConditionModel        `tfschema:"remote_address_condition"`
	RequestMethodCondition    []CdnFrontDoorBatchRuleRequestMethodConditionModel `tfschema:"request_method_condition"`
	QueryStringCondition      []CdnFrontDoorBatchRuleStringConditionModel        `tfschema:"query_string_condition"`
	PostArgsCondition         []CdnFrontDoorBatchRulePostArgsConditionModel      `tfschema:"post_args_condition"`
	RequestURICondition       []CdnFrontDoorBatchRuleStringConditionModel        `tfschema:"request_uri_condition"`
	RequestHeaderCondition    []CdnFrontDoorBatchRuleRequestHeaderConditionModel `tfschema:"request_header_condition"`
	RequestBodyCondition      []CdnFrontDoorBatchRuleStringConditionModel        `tfschema:"request_body_condition"`
	RequestSchemeCondition    []CdnFrontDoorBatchRuleRequestSchemeConditionModel `tfschema:"request_scheme_condition"`
	URLPathCondition          []CdnFrontDoorBatchRuleStringConditionModel        `tfschema:"url_path_condition"`
	URLFileExtensionCondition []CdnFrontDoorBatchRuleStringConditionModel        `tfschema:"url_file_extension_condition"`
	URLFilenameCondition      []CdnFrontDoorBatchRuleStringConditionModel        `tfschema:"url_filename_condition"`
	HTTPVersionCondition      []CdnFrontDoorBatchRuleHTTPVersionConditionModel   `tfschema:"http_version_condition"`
	CookiesCondition          []CdnFrontDoorBatchRuleCookiesConditionModel       `tfschema:"cookies_condition"`
	IsDeviceCondition         []CdnFrontDoorBatchRuleIsDeviceConditionModel      `tfschema:"is_device_condition"`
	SocketAddressCondition    []CdnFrontDoorBatchRuleStringConditionModel        `tfschema:"socket_address_condition"`
	ClientPortCondition       []CdnFrontDoorBatchRuleStringConditionModel        `tfschema:"client_port_condition"`
	ServerPortCondition       []CdnFrontDoorBatchRuleStringConditionModel        `tfschema:"server_port_condition"`
	HostNameCondition         []CdnFrontDoorBatchRuleStringConditionModel        `tfschema:"host_name_condition"`
	SSLProtocolCondition      []CdnFrontDoorBatchRuleSSLProtocolConditionModel   `tfschema:"ssl_protocol_condition"`
}

type CdnFrontDoorBatchRuleStringConditionModel struct {
	Operator        string   `tfschema:"operator"`
	NegateCondition bool     `tfschema:"negate_condition"`
	MatchValues     []string `tfschema:"match_values"`
	Transforms      []string `tfschema:"transforms"`
}

type CdnFrontDoorBatchRulePostArgsConditionModel struct {
	PostArgsName    string   `tfschema:"post_args_name"`
	Operator        string   `tfschema:"operator"`
	NegateCondition bool     `tfschema:"negate_condition"`
	MatchValues     []string `tfschema:"match_values"`
	Transforms      []string `tfschema:"transforms"`
}

type CdnFrontDoorBatchRuleRequestHeaderConditionModel struct {
	HeaderName      string   `tfschema:"header_name"`
	Operator        string   `tfschema:"operator"`
	NegateCondition bool     `tfschema:"negate_condition"`
	MatchValues     []string `tfschema:"match_values"`
	Transforms      []string `tfschema:"transforms"`
}

type CdnFrontDoorBatchRuleCookiesConditionModel struct {
	CookieName      string   `tfschema:"cookie_name"`
	Operator        string   `tfschema:"operator"`
	NegateCondition bool     `tfschema:"negate_condition"`
	MatchValues     []string `tfschema:"match_values"`
	Transforms      []string `tfschema:"transforms"`
}

type CdnFrontDoorBatchRuleRequestMethodConditionModel struct {
	Operator        string   `tfschema:"operator"`
	NegateCondition bool     `tfschema:"negate_condition"`
	MatchValues     []string `tfschema:"match_values"`
}

type CdnFrontDoorBatchRuleRequestSchemeConditionModel struct {
	Operator        string   `tfschema:"operator"`
	NegateCondition bool     `tfschema:"negate_condition"`
	MatchValues     []string `tfschema:"match_values"`
}

type CdnFrontDoorBatchRuleHTTPVersionConditionModel struct {
	Operator        string   `tfschema:"operator"`
	NegateCondition bool     `tfschema:"negate_condition"`
	MatchValues     []string `tfschema:"match_values"`
}

type CdnFrontDoorBatchRuleIsDeviceConditionModel struct {
	Operator        string   `tfschema:"operator"`
	NegateCondition bool     `tfschema:"negate_condition"`
	MatchValues     []string `tfschema:"match_values"`
}

type CdnFrontDoorBatchRuleSSLProtocolConditionModel struct {
	Operator        string   `tfschema:"operator"`
	NegateCondition bool     `tfschema:"negate_condition"`
	MatchValues     []string `tfschema:"match_values"`
}

func cdnFrontDoorBatchRuleSetArguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FrontDoorRuleSetName,
		},
		"cdn_frontdoor_profile_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.FrontDoorProfileID,
		},
		"rules": cdnFrontDoorBatchRulesSchema(),
	}
}

func cdnFrontDoorBatchRulesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MinItems: 1,
		Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.CdnFrontDoorRuleName,
			},
			"behavior_on_match": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      string(rules.MatchProcessingBehaviorContinue),
				ValidateFunc: validation.StringInSlice(rules.PossibleValuesForMatchProcessingBehavior(), false),
			},
			"order": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},
			"actions":    cdnFrontDoorBatchRuleActionsSchema(),
			"conditions": cdnFrontDoorBatchRuleConditionsSchema(),
		}},
	}
}

// lintignore:AZSD002 The `rules` block is repeatable, so Plugin SDK cannot use schema-level AtLeastOneOf paths here. The requirement that each rule define at least one action is enforced per rule in provider validation, and the batch rules API accepts an empty `route_configuration_override_action` block.
func cdnFrontDoorBatchRuleActionsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"url_redirect_action": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
					"redirect_type":        {Type: pluginsdk.TypeString, Required: true, ValidateFunc: validation.StringInSlice(rules.PossibleValuesForRedirectType(), false)},
					"redirect_protocol":    {Type: pluginsdk.TypeString, Optional: true, Default: string(rules.DestinationProtocolMatchRequest), ValidateFunc: validation.StringInSlice(rules.PossibleValuesForDestinationProtocol(), false)},
					"destination_path":     {Type: pluginsdk.TypeString, Optional: true, Default: "", ValidateFunc: validate.CdnFrontDoorUrlRedirectActionDestinationPath},
					"destination_hostname": {Type: pluginsdk.TypeString, Required: true, ValidateFunc: validation.StringLenBetween(0, 2048)},
					"query_string":         {Type: pluginsdk.TypeString, Optional: true, Default: "", ValidateFunc: validate.CdnFrontDoorUrlRedirectActionQueryString},
					"destination_fragment": {Type: pluginsdk.TypeString, Optional: true, Default: "", ValidateFunc: validation.StringLenBetween(0, 1024)},
				}},
			},
			"url_rewrite_action": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
					"source_pattern":          {Type: pluginsdk.TypeString, Required: true, ValidateFunc: validation.StringIsNotEmpty},
					"destination":             {Type: pluginsdk.TypeString, Required: true, ValidateFunc: validation.StringIsNotEmpty},
					"preserve_unmatched_path": {Type: pluginsdk.TypeBool, Optional: true, Default: false},
				}},
			},
			"request_header_action": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
					"header_action": {Type: pluginsdk.TypeString, Required: true, ValidateFunc: validation.StringInSlice(rules.PossibleValuesForHeaderAction(), false)},
					"header_name":   {Type: pluginsdk.TypeString, Required: true, ValidateFunc: validation.StringIsNotEmpty},
					"value":         {Type: pluginsdk.TypeString, Optional: true, ValidateFunc: validation.StringIsNotEmpty},
				}},
			},
			"response_header_action": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
					"header_action": {Type: pluginsdk.TypeString, Required: true, ValidateFunc: validation.StringInSlice(rules.PossibleValuesForHeaderAction(), false)},
					"header_name":   {Type: pluginsdk.TypeString, Required: true, ValidateFunc: validation.StringIsNotEmpty},
					"value":         {Type: pluginsdk.TypeString, Optional: true, ValidateFunc: validation.StringIsNotEmpty},
				}},
			},
			"route_configuration_override_action": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
					"cdn_frontdoor_origin_group_id": {Type: pluginsdk.TypeString, Optional: true, ValidateFunc: validate.FrontDoorOriginGroupID},
					"forwarding_protocol":           {Type: pluginsdk.TypeString, Optional: true, ValidateFunc: validation.StringInSlice(rules.PossibleValuesForForwardingProtocol(), false)},
					"query_string_caching_behavior": {Type: pluginsdk.TypeString, Optional: true, ValidateFunc: validation.StringInSlice(rules.PossibleValuesForRuleQueryStringCachingBehavior(), false)},
					"query_string_parameters":       {Type: pluginsdk.TypeList, Optional: true, MaxItems: 100, Elem: &pluginsdk.Schema{Type: pluginsdk.TypeString}},
					"compression_enabled":           {Type: pluginsdk.TypeBool, Optional: true},
					"cache_behavior":                {Type: pluginsdk.TypeString, Optional: true, ValidateFunc: validation.StringInSlice(PossibleValuesForRuleCacheBehavior(), false)},
					"cache_duration":                {Type: pluginsdk.TypeString, Optional: true, ValidateFunc: validate.CdnFrontDoorCacheDuration},
				}},
			},
		}},
	}
}

func cdnFrontDoorBatchRuleConditionsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{Schema: map[string]*pluginsdk.Schema{
			"remote_address_condition":     batchConditionListSchema(batchOperatorRemoteAddressSchema(), batchMatchValuesSchema(), false, ""),
			"request_method_condition":     batchConditionListSchema(batchEqualOnlyOperatorSchema(), batchRequestMethodMatchValuesSchema(), false, ""),
			"query_string_condition":       batchConditionListSchema(batchOperatorSchema(), batchMatchValuesSchema(), true, ""),
			"post_args_condition":          batchConditionListSchema(batchOperatorSchema(), batchMatchValuesSchema(), true, "post_args_name"),
			"request_uri_condition":        batchConditionListSchema(batchOperatorSchema(), batchMatchValuesSchema(), true, ""),
			"request_header_condition":     batchConditionListSchema(batchOperatorSchema(), batchMatchValuesSchema(), true, "header_name"),
			"request_body_condition":       batchConditionListSchema(batchOperatorSchema(), batchMatchValuesRequiredSchema(), true, ""),
			"request_scheme_condition":     batchConditionListSchema(batchEqualOnlyOperatorSchema(), batchProtocolMatchValuesSchema(), false, ""),
			"url_path_condition":           batchConditionListSchema(batchURLPathOperatorSchema(), batchURLPathMatchValuesSchema(), true, ""),
			"url_file_extension_condition": batchConditionListSchema(batchOperatorSchema(), batchMatchValuesRequiredSchema(), true, ""),
			"url_filename_condition":       batchConditionListSchema(batchOperatorSchema(), batchMatchValuesSchema(), true, ""),
			"http_version_condition":       batchConditionListSchema(batchEqualOnlyOperatorSchema(), batchHTTPVersionMatchValuesSchema(), false, ""),
			"cookies_condition":            batchConditionListSchema(batchOperatorSchema(), batchMatchValuesSchema(), true, "cookie_name"),
			"is_device_condition":          batchConditionListSchema(batchEqualOnlyOperatorSchema(), batchIsDeviceMatchValuesSchema(), false, ""),
			"socket_address_condition":     batchConditionListSchema(batchOperatorSocketAddressSchema(), batchMatchValuesSchema(), false, ""),
			"client_port_condition":        batchConditionListSchema(batchOperatorSchema(), batchMatchValuesSchema(), false, ""),
			"server_port_condition":        batchConditionListSchema(batchOperatorSchema(), batchServerPortMatchValuesSchema(), false, ""),
			"host_name_condition":          batchConditionListSchema(batchOperatorSchema(), batchMatchValuesSchema(), true, ""),
			"ssl_protocol_condition":       batchConditionListSchema(batchEqualOnlyOperatorSchema(), batchSSLProtocolMatchValuesSchema(), false, ""),
		}},
	}
}

func batchConditionListSchema(operatorSchema, matchValuesSchema *pluginsdk.Schema, includeTransforms bool, selectorField string) *pluginsdk.Schema {
	schema := map[string]*pluginsdk.Schema{
		"operator":         operatorSchema,
		"negate_condition": batchNegateConditionSchema(),
		"match_values":     matchValuesSchema,
	}
	if includeTransforms {
		schema["transforms"] = batchTransformsSchema()
	}
	if selectorField != "" {
		schema[selectorField] = &pluginsdk.Schema{Type: pluginsdk.TypeString, Required: true, ValidateFunc: validation.StringIsNotEmpty}
	}
	return &pluginsdk.Schema{Type: pluginsdk.TypeList, Optional: true, Elem: &pluginsdk.Resource{Schema: schema}}
}

func batchOperatorSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{Type: pluginsdk.TypeString, Required: true, ValidateFunc: validation.StringInSlice(rules.PossibleValuesForQueryStringOperator(), false)}
}

func batchURLPathOperatorSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{Type: pluginsdk.TypeString, Required: true, ValidateFunc: validation.StringInSlice(rules.PossibleValuesForURLPathOperator(), false)}
}

func batchEqualOnlyOperatorSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{Type: pluginsdk.TypeString, Optional: true, Default: string(rules.RequestMethodOperatorEqual), ValidateFunc: validation.StringInSlice(rules.PossibleValuesForRequestMethodOperator(), false)}
}

func batchOperatorRemoteAddressSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{Type: pluginsdk.TypeString, Optional: true, Default: string(rules.RemoteAddressOperatorIPMatch), ValidateFunc: validation.StringInSlice(rules.PossibleValuesForRemoteAddressOperator(), false)}
}

func batchOperatorSocketAddressSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{Type: pluginsdk.TypeString, Optional: true, Default: string(rules.SocketAddrOperatorIPMatch), ValidateFunc: validation.StringInSlice(rules.PossibleValuesForSocketAddrOperator(), false)}
}

func batchNegateConditionSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{Type: pluginsdk.TypeBool, Optional: true, Default: false}
}

func batchMatchValuesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{Type: pluginsdk.TypeList, Optional: true, MaxItems: 25, Elem: &pluginsdk.Schema{Type: pluginsdk.TypeString}}
}

func batchMatchValuesRequiredSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{Type: pluginsdk.TypeList, Required: true, MaxItems: 25, Elem: &pluginsdk.Schema{Type: pluginsdk.TypeString, ValidateFunc: validation.StringIsNotEmpty}}
}

func batchServerPortMatchValuesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{Type: pluginsdk.TypeSet, Required: true, MaxItems: 2, Elem: &pluginsdk.Schema{Type: pluginsdk.TypeString, ValidateFunc: validation.StringInSlice([]string{"80", "443"}, false)}}
}

func batchSSLProtocolMatchValuesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{Type: pluginsdk.TypeSet, Required: true, MaxItems: 3, Elem: &pluginsdk.Schema{Type: pluginsdk.TypeString, ValidateFunc: validation.StringInSlice(rules.PossibleValuesForSslProtocol(), false)}}
}

func batchURLPathMatchValuesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{Type: pluginsdk.TypeList, Optional: true, MaxItems: 25, Elem: &pluginsdk.Schema{Type: pluginsdk.TypeString, ValidateFunc: validate.CdnFrontDoorUrlPathConditionMatchValue}}
}

func batchRequestMethodMatchValuesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{Type: pluginsdk.TypeSet, Required: true, MaxItems: 7, Elem: &pluginsdk.Schema{Type: pluginsdk.TypeString, ValidateFunc: validation.StringInSlice(rules.PossibleValuesForRequestMethodMatchValue(), false)}}
}

func batchProtocolMatchValuesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{Type: pluginsdk.TypeList, Optional: true, MaxItems: 1, Elem: &pluginsdk.Schema{Type: pluginsdk.TypeString, Default: rules.RequestSchemeMatchValueHTTP, ValidateFunc: validation.StringInSlice(rules.PossibleValuesForRequestSchemeMatchValue(), false)}}
}

func batchIsDeviceMatchValuesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{Type: pluginsdk.TypeList, Optional: true, MaxItems: 1, Elem: &pluginsdk.Schema{Type: pluginsdk.TypeString, ValidateFunc: validation.StringInSlice(rules.PossibleValuesForIsDeviceMatchValue(), false)}}
}

func batchHTTPVersionMatchValuesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{Type: pluginsdk.TypeSet, Required: true, MaxItems: 4, Elem: &pluginsdk.Schema{Type: pluginsdk.TypeString, ValidateFunc: validation.StringInSlice([]string{"2.0", "1.1", "1.0", "0.9"}, false)}}
}

func batchTransformsSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{Type: pluginsdk.TypeSet, Optional: true, MaxItems: 4, Elem: &pluginsdk.Schema{Type: pluginsdk.TypeString, ValidateFunc: validation.StringInSlice(rules.PossibleValuesForTransform(), false)}}
}

func (r CdnFrontDoorBatchRuleSetResource) flatten(metadata sdk.ResourceMetaData, id *rules.RuleSetId, model *azuresdkhacks.BatchRuleSetResource) error {
	state := CdnFrontDoorBatchRuleSetModel{
		ID:                    id.ID(),
		Name:                  id.RuleSetName,
		CdnFrontDoorProfileID: profiles.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName).ID(),
	}

	if model != nil && model.Properties != nil && model.Properties.Rules != nil {
		rulesState, err := flattenCdnFrontDoorBatchRules(pointer.From(model.Properties.Rules))
		if err != nil {
			return fmt.Errorf("flattening `rules`: %+v", err)
		}
		state.Rules = rulesState
	}

	if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, id); err != nil {
		return err
	}

	return metadata.Encode(&state)
}

func flattenCdnFrontDoorBatchRules(input []azuresdkhacks.BatchRuleProperties) ([]CdnFrontDoorBatchRuleRuleModel, error) {
	results := make([]CdnFrontDoorBatchRuleRuleModel, 0, len(input))
	sorted := append([]azuresdkhacks.BatchRuleProperties(nil), input...)
	sort.SliceStable(sorted, func(i, j int) bool {
		return batchRuleOrder(sorted[i]) < batchRuleOrder(sorted[j])
	})

	for _, item := range sorted {
		ruleState, err := flattenCdnFrontDoorBatchRule(item)
		if err != nil {
			return []CdnFrontDoorBatchRuleRuleModel{}, err
		}
		results = append(results, ruleState)
	}

	return results, nil
}

func flattenCdnFrontDoorBatchRule(input azuresdkhacks.BatchRuleProperties) (CdnFrontDoorBatchRuleRuleModel, error) {
	state := CdnFrontDoorBatchRuleRuleModel{
		Name: pointer.From(input.Name),
	}

	state.BehaviorOnMatch = string(pointer.From(input.MatchProcessingBehavior))
	state.Order = pointer.From(input.Order)

	actions, err := flattenCdnFrontDoorBatchRuleActions(input.Actions)
	if err != nil {
		return CdnFrontDoorBatchRuleRuleModel{}, fmt.Errorf("flattening `actions`: %+v", err)
	}
	state.Actions = actions

	conditions, err := flattenCdnFrontDoorBatchRuleConditions(input.Conditions)
	if err != nil {
		return CdnFrontDoorBatchRuleRuleModel{}, fmt.Errorf("flattening `conditions`: %+v", err)
	}
	state.Conditions = conditions

	return state, nil
}

func expandCdnFrontDoorBatchRuleSetPayload(batchMode bool, model CdnFrontDoorBatchRuleSetModel) (azuresdkhacks.BatchRuleSetResource, error) {
	expandedRules, err := expandCdnFrontDoorBatchRules(model.Rules)
	if err != nil {
		return azuresdkhacks.BatchRuleSetResource{}, err
	}

	return azuresdkhacks.BatchRuleSetResource{
		Properties: &azuresdkhacks.BatchRuleSetProperties{
			BatchMode: pointer.To(batchMode),
			Rules:     pointer.To(expandedRules),
		},
	}, nil
}

func expandCdnFrontDoorBatchRules(input []CdnFrontDoorBatchRuleRuleModel) ([]azuresdkhacks.BatchRuleProperties, error) {
	if err := validateCdnFrontDoorBatchRules(input); err != nil {
		return nil, err
	}

	results := make([]azuresdkhacks.BatchRuleProperties, 0, len(input))
	for _, item := range input {
		actions, err := expandCdnFrontDoorBatchRuleActions(item.Actions)
		if err != nil {
			return nil, err
		}

		conditions, err := expandCdnFrontDoorBatchRuleConditions(item.Conditions)
		if err != nil {
			return nil, err
		}

		results = append(results, azuresdkhacks.BatchRuleProperties{
			Actions:                 pointer.To(actions),
			Conditions:              pointer.To(conditions),
			MatchProcessingBehavior: pointer.ToEnum[rules.MatchProcessingBehavior](item.BehaviorOnMatch),
			Order:                   pointer.To(item.Order),
			RuleName:                pointer.To(item.Name),
		})
	}

	sort.SliceStable(results, func(i, j int) bool {
		return batchRuleOrder(results[i]) < batchRuleOrder(results[j])
	})

	return results, nil
}

func validateCdnFrontDoorBatchRules(input []CdnFrontDoorBatchRuleRuleModel) error {
	names := make(map[string]struct{}, len(input))
	orders := make(map[int64]struct{}, len(input))
	var lastOrder int64
	haveLastOrder := false

	for _, item := range input {
		if _, exists := names[item.Name]; exists {
			return fmt.Errorf("the `rules` blocks must have unique `name` values, got duplicate `%s`", item.Name)
		}
		names[item.Name] = struct{}{}

		if _, exists := orders[item.Order]; exists {
			return fmt.Errorf("the `rules` blocks must have unique `order` values, got duplicate `%d`", item.Order)
		}
		orders[item.Order] = struct{}{}

		if haveLastOrder && item.Order < lastOrder {
			return fmt.Errorf("the `rules` blocks must be declared in ascending `order`, got `%d` before `%d`", lastOrder, item.Order)
		}

		if err := validate.CdnFrontDoorValidateActionDefinitions(batchRuleActionCounts(item.Actions)); err != nil {
			return err
		}

		lastOrder = item.Order
		haveLastOrder = true
	}

	return nil
}

func batchRuleOrder(input azuresdkhacks.BatchRuleProperties) int64 {
	return pointer.From(input.Order)
}

func batchRuleActionCounts(input []CdnFrontDoorBatchRuleActionsModel) (int, int, int, int) {
	if len(input) == 0 {
		return 0, 0, 0, 0
	}

	actions := input[0]
	urlRewriteCount := len(actions.URLRewriteAction)
	urlRedirectCount := len(actions.URLRedirectAction)
	routeConfigurationOverrideCount := len(actions.RouteConfigurationOverrideAction)
	totalCount := urlRewriteCount + urlRedirectCount + len(actions.RequestHeaderAction) + len(actions.ResponseHeaderAction) + routeConfigurationOverrideCount

	return urlRewriteCount, urlRedirectCount, routeConfigurationOverrideCount, totalCount
}

func expandCdnFrontDoorBatchRuleActions(input []CdnFrontDoorBatchRuleActionsModel) ([]rules.DeliveryRuleAction, error) {
	results := make([]rules.DeliveryRuleAction, 0)
	if len(input) == 0 {
		return results, nil
	}

	actions := input[0]
	urlRewriteCount, urlRedirectCount, routeConfigurationOverrideCount, totalCount := batchRuleActionCounts(input)
	if err := validate.CdnFrontDoorValidateActionDefinitions(urlRewriteCount, urlRedirectCount, routeConfigurationOverrideCount, totalCount); err != nil {
		return nil, err
	}

	for _, item := range actions.URLRedirectAction {
		results = append(results, rules.URLRedirectAction{
			Name: rules.DeliveryRuleActionNameURLRedirect,
			Parameters: rules.URLRedirectActionParameters{
				TypeName:            rules.DeliveryRuleActionParametersTypeDeliveryRuleURLRedirectActionParameters,
				RedirectType:        rules.RedirectType(item.RedirectType),
				DestinationProtocol: pointerOrNil(rules.DestinationProtocol(item.RedirectProtocol)),
				CustomPath:          stringPointerOrNil(item.DestinationPath),
				CustomHostname:      stringPointerOrNil(item.DestinationHostname),
				CustomQueryString:   stringPointerOrNil(item.QueryString),
				CustomFragment:      stringPointerOrNil(item.DestinationFragment),
			},
		})
	}

	for _, item := range actions.URLRewriteAction {
		results = append(results, rules.URLRewriteAction{
			Name: rules.DeliveryRuleActionNameURLRewrite,
			Parameters: rules.URLRewriteActionParameters{
				TypeName:              rules.DeliveryRuleActionParametersTypeDeliveryRuleURLRewriteActionParameters,
				SourcePattern:         item.SourcePattern,
				Destination:           item.Destination,
				PreserveUnmatchedPath: pointer.To(item.PreserveUnmatchedPath),
			},
		})
	}

	for _, item := range actions.RequestHeaderAction {
		if err := validate.CdnFrontDoorValidateHeaderAction("request_header_action", item.HeaderAction, item.Value); err != nil {
			return nil, err
		}
		results = append(results, rules.DeliveryRuleRequestHeaderAction{
			Name: rules.DeliveryRuleActionNameModifyRequestHeader,
			Parameters: rules.HeaderActionParameters{
				TypeName:     rules.DeliveryRuleActionParametersTypeDeliveryRuleHeaderActionParameters,
				HeaderAction: rules.HeaderAction(item.HeaderAction),
				HeaderName:   item.HeaderName,
				Value:        stringPointerOrNil(item.Value),
			},
		})
	}

	for _, item := range actions.ResponseHeaderAction {
		if err := validate.CdnFrontDoorValidateHeaderAction("response_header_action", item.HeaderAction, item.Value); err != nil {
			return nil, err
		}
		results = append(results, rules.DeliveryRuleResponseHeaderAction{
			Name: rules.DeliveryRuleActionNameModifyResponseHeader,
			Parameters: rules.HeaderActionParameters{
				TypeName:     rules.DeliveryRuleActionParametersTypeDeliveryRuleHeaderActionParameters,
				HeaderAction: rules.HeaderAction(item.HeaderAction),
				HeaderName:   item.HeaderName,
				Value:        stringPointerOrNil(item.Value),
			},
		})
	}

	for _, item := range actions.RouteConfigurationOverrideAction {
		originGroupOverride, cacheConfiguration, err := expandRouteConfigurationOverrideAction(item)
		if err != nil {
			return nil, err
		}
		results = append(results, rules.DeliveryRuleRouteConfigurationOverrideAction{
			Name: rules.DeliveryRuleActionNameRouteConfigurationOverride,
			Parameters: rules.RouteConfigurationOverrideActionParameters{
				TypeName:            rules.DeliveryRuleActionParametersTypeDeliveryRuleRouteConfigurationOverrideActionParameters,
				OriginGroupOverride: originGroupOverride,
				CacheConfiguration:  cacheConfiguration,
			},
		})
	}

	return results, nil
}

func flattenCdnFrontDoorBatchRuleActions(input *[]rules.DeliveryRuleAction) ([]CdnFrontDoorBatchRuleActionsModel, error) {
	if input == nil {
		return []CdnFrontDoorBatchRuleActionsModel{}, nil
	}

	actions := CdnFrontDoorBatchRuleActionsModel{}
	for _, action := range *input {
		switch action.DeliveryRuleAction().Name {
		case rules.DeliveryRuleActionNameURLRedirect:
			item := action.(rules.URLRedirectAction)
			actions.URLRedirectAction = append(actions.URLRedirectAction, CdnFrontDoorBatchRuleURLRedirectActionModel{
				RedirectType:        string(item.Parameters.RedirectType),
				RedirectProtocol:    string(pointer.From(item.Parameters.DestinationProtocol)),
				DestinationPath:     pointer.From(item.Parameters.CustomPath),
				DestinationHostname: pointer.From(item.Parameters.CustomHostname),
				QueryString:         pointer.From(item.Parameters.CustomQueryString),
				DestinationFragment: pointer.From(item.Parameters.CustomFragment),
			})
		case rules.DeliveryRuleActionNameURLRewrite:
			item := action.(rules.URLRewriteAction)
			actions.URLRewriteAction = append(actions.URLRewriteAction, CdnFrontDoorBatchRuleURLRewriteActionModel{
				SourcePattern:         item.Parameters.SourcePattern,
				Destination:           item.Parameters.Destination,
				PreserveUnmatchedPath: pointer.From(item.Parameters.PreserveUnmatchedPath),
			})
		case rules.DeliveryRuleActionNameModifyRequestHeader:
			item := action.(rules.DeliveryRuleRequestHeaderAction)
			actions.RequestHeaderAction = append(actions.RequestHeaderAction, CdnFrontDoorBatchRuleHeaderActionModel{
				HeaderAction: string(item.Parameters.HeaderAction),
				HeaderName:   item.Parameters.HeaderName,
				Value:        pointer.From(item.Parameters.Value),
			})
		case rules.DeliveryRuleActionNameModifyResponseHeader:
			item := action.(rules.DeliveryRuleResponseHeaderAction)
			actions.ResponseHeaderAction = append(actions.ResponseHeaderAction, CdnFrontDoorBatchRuleHeaderActionModel{
				HeaderAction: string(item.Parameters.HeaderAction),
				HeaderName:   item.Parameters.HeaderName,
				Value:        pointer.From(item.Parameters.Value),
			})
		case rules.DeliveryRuleActionNameRouteConfigurationOverride:
			item := action.(rules.DeliveryRuleRouteConfigurationOverrideAction)
			flattened, err := flattenRouteConfigurationOverrideAction(item.Parameters)
			if err != nil {
				return []CdnFrontDoorBatchRuleActionsModel{}, err
			}
			actions.RouteConfigurationOverrideAction = append(actions.RouteConfigurationOverrideAction, flattened)
		default:
			return []CdnFrontDoorBatchRuleActionsModel{}, fmt.Errorf("unsupported batch rule action %q encountered", action.DeliveryRuleAction().Name)
		}
	}

	return []CdnFrontDoorBatchRuleActionsModel{actions}, nil
}

func expandRouteConfigurationOverrideAction(input CdnFrontDoorBatchRuleRouteConfigurationOverrideActionModel) (*rules.OriginGroupOverride, *rules.CacheConfiguration, error) {
	var originGroupOverride *rules.OriginGroupOverride
	if err := validate.CdnFrontDoorValidateRouteConfigurationOverrideAction(validate.CdnFrontDoorRouteConfigurationOverrideInput{
		OriginGroupID:              input.CdnFrontDoorOriginGroupID,
		ForwardingProtocol:         input.ForwardingProtocol,
		QueryStringCachingBehavior: input.QueryStringCachingBehavior,
		QueryStringParameters:      input.QueryStringParameters,
		CompressionEnabled:         input.CompressionEnabled,
		CacheBehavior:              input.CacheBehavior,
		CacheDuration:              input.CacheDuration,
	}); err != nil {
		return nil, nil, err
	}

	if input.CdnFrontDoorOriginGroupID != "" {
		originGroupOverride = &rules.OriginGroupOverride{
			OriginGroup:        &rules.ResourceReference{Id: pointer.To(input.CdnFrontDoorOriginGroupID)},
			ForwardingProtocol: pointer.ToEnum[rules.ForwardingProtocol](input.ForwardingProtocol),
		}
	}

	var cacheConfiguration *rules.CacheConfiguration
	if input.CacheBehavior != "" || input.QueryStringCachingBehavior != "" || len(input.QueryStringParameters) > 0 || input.CacheDuration != "" || input.CompressionEnabled {
		cacheBehavior := input.CacheBehavior
		if cacheBehavior != string(rules.RuleIsCompressionEnabledDisabled) {
			compressionEnabled := rules.RuleIsCompressionEnabledDisabled
			if input.CompressionEnabled {
				compressionEnabled = rules.RuleIsCompressionEnabledEnabled
			}

			cacheConfiguration = &rules.CacheConfiguration{
				CacheBehavior:              pointer.ToEnum[rules.RuleCacheBehavior](cacheBehavior),
				CacheDuration:              stringPointerOrNil(input.CacheDuration),
				IsCompressionEnabled:       pointer.To(compressionEnabled),
				QueryParameters:            stringPointerOrNil(strings.Join(input.QueryStringParameters, ",")),
				QueryStringCachingBehavior: pointerOrNil(rules.RuleQueryStringCachingBehavior(input.QueryStringCachingBehavior)),
			}
		}
	}

	return originGroupOverride, cacheConfiguration, nil
}

func flattenRouteConfigurationOverrideAction(input rules.RouteConfigurationOverrideActionParameters) (CdnFrontDoorBatchRuleRouteConfigurationOverrideActionModel, error) {
	output := CdnFrontDoorBatchRuleRouteConfigurationOverrideActionModel{}
	if input.OriginGroupOverride != nil {
		originGroup, err := parse.FrontDoorOriginGroupIDInsensitively(pointer.From(input.OriginGroupOverride.OriginGroup.Id))
		if err != nil {
			return CdnFrontDoorBatchRuleRouteConfigurationOverrideActionModel{}, err
		}
		output.CdnFrontDoorOriginGroupID = originGroup.ID()
		output.ForwardingProtocol = string(pointer.From(input.OriginGroupOverride.ForwardingProtocol))
	}
	if input.CacheConfiguration != nil {
		output.QueryStringParameters = splitCSV(pointer.From(input.CacheConfiguration.QueryParameters))
		output.QueryStringCachingBehavior = string(pointer.From(input.CacheConfiguration.QueryStringCachingBehavior))
		output.CacheBehavior = string(pointer.From(input.CacheConfiguration.CacheBehavior))
		output.CacheDuration = pointer.From(input.CacheConfiguration.CacheDuration)
		output.CompressionEnabled = pointer.From(input.CacheConfiguration.IsCompressionEnabled) == rules.RuleIsCompressionEnabledEnabled
	}
	return output, nil
}

func expandCdnFrontDoorBatchRuleConditions(input []CdnFrontDoorBatchRuleConditionsModel) ([]rules.DeliveryRuleCondition, error) {
	results := make([]rules.DeliveryRuleCondition, 0)
	if len(input) == 0 {
		return results, nil
	}

	conditions := input[0]

	appendStringCondition := func(configName string, items []CdnFrontDoorBatchRuleStringConditionModel, expand func(CdnFrontDoorBatchRuleStringConditionModel) (rules.DeliveryRuleCondition, error)) error {
		for _, item := range items {
			condition, err := expand(item)
			if err != nil {
				return fmt.Errorf("expanding `%s`: %+v", configName, err)
			}
			results = append(results, condition)
		}
		return nil
	}

	if err := appendStringCondition("remote_address_condition", conditions.RemoteAddressCondition, expandRemoteAddressCondition); err != nil {
		return nil, err
	}
	if err := appendStringCondition("query_string_condition", conditions.QueryStringCondition, expandQueryStringCondition); err != nil {
		return nil, err
	}
	if err := appendStringCondition("request_uri_condition", conditions.RequestURICondition, expandRequestURICondition); err != nil {
		return nil, err
	}
	if err := appendStringCondition("request_body_condition", conditions.RequestBodyCondition, expandRequestBodyCondition); err != nil {
		return nil, err
	}
	if err := appendStringCondition("url_path_condition", conditions.URLPathCondition, expandURLPathCondition); err != nil {
		return nil, err
	}
	if err := appendStringCondition("url_file_extension_condition", conditions.URLFileExtensionCondition, expandURLFileExtensionCondition); err != nil {
		return nil, err
	}
	if err := appendStringCondition("url_filename_condition", conditions.URLFilenameCondition, expandURLFilenameCondition); err != nil {
		return nil, err
	}
	if err := appendStringCondition("socket_address_condition", conditions.SocketAddressCondition, expandSocketAddressCondition); err != nil {
		return nil, err
	}
	if err := appendStringCondition("client_port_condition", conditions.ClientPortCondition, expandClientPortCondition); err != nil {
		return nil, err
	}
	if err := appendStringCondition("server_port_condition", conditions.ServerPortCondition, expandServerPortCondition); err != nil {
		return nil, err
	}
	if err := appendStringCondition("host_name_condition", conditions.HostNameCondition, expandHostNameCondition); err != nil {
		return nil, err
	}

	for _, item := range conditions.RequestMethodCondition {
		condition, err := expandRequestMethodCondition(item)
		if err != nil {
			return nil, fmt.Errorf("expanding `request_method_condition`: %+v", err)
		}
		results = append(results, condition)
	}
	for _, item := range conditions.PostArgsCondition {
		condition, err := expandPostArgsCondition(item)
		if err != nil {
			return nil, fmt.Errorf("expanding `post_args_condition`: %+v", err)
		}
		results = append(results, condition)
	}
	for _, item := range conditions.RequestHeaderCondition {
		condition, err := expandRequestHeaderCondition(item)
		if err != nil {
			return nil, fmt.Errorf("expanding `request_header_condition`: %+v", err)
		}
		results = append(results, condition)
	}
	for _, item := range conditions.RequestSchemeCondition {
		condition, err := expandRequestSchemeCondition(item)
		if err != nil {
			return nil, fmt.Errorf("expanding `request_scheme_condition`: %+v", err)
		}
		results = append(results, condition)
	}
	for _, item := range conditions.HTTPVersionCondition {
		condition, err := expandHTTPVersionCondition(item)
		if err != nil {
			return nil, fmt.Errorf("expanding `http_version_condition`: %+v", err)
		}
		results = append(results, condition)
	}
	for _, item := range conditions.CookiesCondition {
		condition, err := expandCookiesCondition(item)
		if err != nil {
			return nil, fmt.Errorf("expanding `cookies_condition`: %+v", err)
		}
		results = append(results, condition)
	}
	for _, item := range conditions.IsDeviceCondition {
		condition, err := expandIsDeviceCondition(item)
		if err != nil {
			return nil, fmt.Errorf("expanding `is_device_condition`: %+v", err)
		}
		results = append(results, condition)
	}
	for _, item := range conditions.SSLProtocolCondition {
		condition, err := expandSSLProtocolCondition(item)
		if err != nil {
			return nil, fmt.Errorf("expanding `ssl_protocol_condition`: %+v", err)
		}
		results = append(results, condition)
	}

	if len(results) > 10 {
		return nil, fmt.Errorf("the `conditions` block may only contain up to 10 match conditions, got %d", len(results))
	}

	return results, nil
}

func flattenCdnFrontDoorBatchRuleConditions(input *[]rules.DeliveryRuleCondition) ([]CdnFrontDoorBatchRuleConditionsModel, error) {
	if input == nil {
		return []CdnFrontDoorBatchRuleConditionsModel{}, nil
	}

	conditions := CdnFrontDoorBatchRuleConditionsModel{}
	for _, condition := range *input {
		switch condition.DeliveryRuleCondition().Name {
		case rules.MatchVariableRemoteAddress:
			item := condition.(rules.DeliveryRuleRemoteAddressCondition)
			conditions.RemoteAddressCondition = append(conditions.RemoteAddressCondition, flattenStringCondition(string(item.Parameters.Operator), pointer.From(item.Parameters.NegateCondition), pointer.From(item.Parameters.MatchValues), transformsToStrings(item.Parameters.Transforms)))
		case rules.MatchVariableRequestMethod:
			item := condition.(rules.DeliveryRuleRequestMethodCondition)
			conditions.RequestMethodCondition = append(conditions.RequestMethodCondition, CdnFrontDoorBatchRuleRequestMethodConditionModel{Operator: string(item.Parameters.Operator), NegateCondition: pointer.From(item.Parameters.NegateCondition), MatchValues: requestMethodValuesToStrings(item.Parameters.MatchValues)})
		case rules.MatchVariableQueryString:
			item := condition.(rules.DeliveryRuleQueryStringCondition)
			conditions.QueryStringCondition = append(conditions.QueryStringCondition, flattenStringCondition(string(item.Parameters.Operator), pointer.From(item.Parameters.NegateCondition), pointer.From(item.Parameters.MatchValues), transformsToStrings(item.Parameters.Transforms)))
		case rules.MatchVariablePostArgs:
			item := condition.(rules.DeliveryRulePostArgsCondition)
			conditions.PostArgsCondition = append(conditions.PostArgsCondition, CdnFrontDoorBatchRulePostArgsConditionModel{PostArgsName: pointer.From(item.Parameters.Selector), Operator: string(item.Parameters.Operator), NegateCondition: pointer.From(item.Parameters.NegateCondition), MatchValues: pointer.From(item.Parameters.MatchValues), Transforms: transformsToStrings(item.Parameters.Transforms)})
		case rules.MatchVariableRequestUri:
			item := condition.(rules.DeliveryRuleRequestUriCondition)
			conditions.RequestURICondition = append(conditions.RequestURICondition, flattenStringCondition(string(item.Parameters.Operator), pointer.From(item.Parameters.NegateCondition), pointer.From(item.Parameters.MatchValues), transformsToStrings(item.Parameters.Transforms)))
		case rules.MatchVariableRequestHeader:
			item := condition.(rules.DeliveryRuleRequestHeaderCondition)
			conditions.RequestHeaderCondition = append(conditions.RequestHeaderCondition, CdnFrontDoorBatchRuleRequestHeaderConditionModel{HeaderName: pointer.From(item.Parameters.Selector), Operator: string(item.Parameters.Operator), NegateCondition: pointer.From(item.Parameters.NegateCondition), MatchValues: pointer.From(item.Parameters.MatchValues), Transforms: transformsToStrings(item.Parameters.Transforms)})
		case rules.MatchVariableRequestBody:
			item := condition.(rules.DeliveryRuleRequestBodyCondition)
			conditions.RequestBodyCondition = append(conditions.RequestBodyCondition, flattenStringCondition(string(item.Parameters.Operator), pointer.From(item.Parameters.NegateCondition), pointer.From(item.Parameters.MatchValues), transformsToStrings(item.Parameters.Transforms)))
		case rules.MatchVariableRequestScheme:
			item := condition.(rules.DeliveryRuleRequestSchemeCondition)
			conditions.RequestSchemeCondition = append(conditions.RequestSchemeCondition, CdnFrontDoorBatchRuleRequestSchemeConditionModel{Operator: string(item.Parameters.Operator), NegateCondition: pointer.From(item.Parameters.NegateCondition), MatchValues: requestSchemeValuesToStrings(item.Parameters.MatchValues)})
		case rules.MatchVariableURLPath:
			item := condition.(rules.DeliveryRuleURLPathCondition)
			conditions.URLPathCondition = append(conditions.URLPathCondition, flattenStringCondition(string(item.Parameters.Operator), pointer.From(item.Parameters.NegateCondition), pointer.From(item.Parameters.MatchValues), transformsToStrings(item.Parameters.Transforms)))
		case rules.MatchVariableURLFileExtension:
			item := condition.(rules.DeliveryRuleURLFileExtensionCondition)
			conditions.URLFileExtensionCondition = append(conditions.URLFileExtensionCondition, flattenStringCondition(string(item.Parameters.Operator), pointer.From(item.Parameters.NegateCondition), pointer.From(item.Parameters.MatchValues), transformsToStrings(item.Parameters.Transforms)))
		case rules.MatchVariableURLFileName:
			item := condition.(rules.DeliveryRuleURLFileNameCondition)
			conditions.URLFilenameCondition = append(conditions.URLFilenameCondition, flattenStringCondition(string(item.Parameters.Operator), pointer.From(item.Parameters.NegateCondition), pointer.From(item.Parameters.MatchValues), transformsToStrings(item.Parameters.Transforms)))
		case rules.MatchVariableHTTPVersion:
			item := condition.(rules.DeliveryRuleHTTPVersionCondition)
			conditions.HTTPVersionCondition = append(conditions.HTTPVersionCondition, CdnFrontDoorBatchRuleHTTPVersionConditionModel{Operator: string(item.Parameters.Operator), NegateCondition: pointer.From(item.Parameters.NegateCondition), MatchValues: pointer.From(item.Parameters.MatchValues)})
		case rules.MatchVariableCookies:
			item := condition.(rules.DeliveryRuleCookiesCondition)
			conditions.CookiesCondition = append(conditions.CookiesCondition, CdnFrontDoorBatchRuleCookiesConditionModel{CookieName: pointer.From(item.Parameters.Selector), Operator: string(item.Parameters.Operator), NegateCondition: pointer.From(item.Parameters.NegateCondition), MatchValues: pointer.From(item.Parameters.MatchValues), Transforms: transformsToStrings(item.Parameters.Transforms)})
		case rules.MatchVariableIsDevice:
			item := condition.(rules.DeliveryRuleIsDeviceCondition)
			conditions.IsDeviceCondition = append(conditions.IsDeviceCondition, CdnFrontDoorBatchRuleIsDeviceConditionModel{Operator: string(item.Parameters.Operator), NegateCondition: pointer.From(item.Parameters.NegateCondition), MatchValues: isDeviceValuesToStrings(item.Parameters.MatchValues)})
		case rules.MatchVariableSocketAddr:
			item := condition.(rules.DeliveryRuleSocketAddrCondition)
			conditions.SocketAddressCondition = append(conditions.SocketAddressCondition, flattenStringCondition(string(item.Parameters.Operator), pointer.From(item.Parameters.NegateCondition), pointer.From(item.Parameters.MatchValues), nil))
		case rules.MatchVariableClientPort:
			item := condition.(rules.DeliveryRuleClientPortCondition)
			conditions.ClientPortCondition = append(conditions.ClientPortCondition, flattenStringCondition(string(item.Parameters.Operator), pointer.From(item.Parameters.NegateCondition), pointer.From(item.Parameters.MatchValues), nil))
		case rules.MatchVariableServerPort:
			item := condition.(rules.DeliveryRuleServerPortCondition)
			conditions.ServerPortCondition = append(conditions.ServerPortCondition, flattenStringCondition(string(item.Parameters.Operator), pointer.From(item.Parameters.NegateCondition), pointer.From(item.Parameters.MatchValues), nil))
		case rules.MatchVariableHostName:
			item := condition.(rules.DeliveryRuleHostNameCondition)
			conditions.HostNameCondition = append(conditions.HostNameCondition, flattenStringCondition(string(item.Parameters.Operator), pointer.From(item.Parameters.NegateCondition), pointer.From(item.Parameters.MatchValues), transformsToStrings(item.Parameters.Transforms)))
		case rules.MatchVariableSslProtocol:
			item := condition.(rules.DeliveryRuleSslProtocolCondition)
			conditions.SSLProtocolCondition = append(conditions.SSLProtocolCondition, CdnFrontDoorBatchRuleSSLProtocolConditionModel{Operator: string(item.Parameters.Operator), NegateCondition: pointer.From(item.Parameters.NegateCondition), MatchValues: sslProtocolValuesToStrings(item.Parameters.MatchValues)})
		default:
			return []CdnFrontDoorBatchRuleConditionsModel{}, fmt.Errorf("unsupported batch rule condition %q encountered", condition.DeliveryRuleCondition().Name)
		}
	}

	if !batchRuleHasConditions(conditions) {
		return []CdnFrontDoorBatchRuleConditionsModel{}, nil
	}

	return []CdnFrontDoorBatchRuleConditionsModel{conditions}, nil
}

func expandRemoteAddressCondition(input CdnFrontDoorBatchRuleStringConditionModel) (rules.DeliveryRuleCondition, error) {
	if err := validateAnyCondition("remote_address_condition", input.Operator, input.MatchValues); err != nil {
		return nil, err
	}

	if strings.EqualFold(input.Operator, string(rules.RemoteAddressOperatorGeoMatch)) {
		for _, matchValue := range input.MatchValues {
			if ok, _ := helperValidate.RegExHelper(matchValue, "match_values", `^[A-Z]{2}$`); !ok {
				return nil, fmt.Errorf("`remote_address_condition` is invalid: when `operator` is `GeoMatch` the values in `match_values` must be valid country codes consisting of 2 uppercase characters, got %q", matchValue)
			}
		}
	}

	if strings.EqualFold(input.Operator, string(rules.RemoteAddressOperatorIPMatch)) {
		matchValues := make([]interface{}, 0, len(input.MatchValues))
		for _, matchValue := range input.MatchValues {
			matchValues = append(matchValues, matchValue)
			if _, errs := validate.FrontDoorRuleCidrIsValid(matchValue, "match_values"); len(errs) > 0 {
				return nil, fmt.Errorf("`remote_address_condition` is invalid: when `operator` is `IPMatch` the values in `match_values` must be valid IPv4 or IPv6 CIDRs, got %q", matchValue)
			}
		}

		if _, errs := validate.FrontDoorRuleCidrOverlap(matchValues, "match_values"); len(errs) > 0 {
			return nil, fmt.Errorf("`remote_address_condition` is invalid: %+v", errs[0])
		}
	}

	return rules.DeliveryRuleRemoteAddressCondition{Name: rules.MatchVariableRemoteAddress, Parameters: rules.RemoteAddressMatchConditionParameters{TypeName: rules.DeliveryRuleConditionParametersTypeDeliveryRuleRemoteAddressConditionParameters, Operator: rules.RemoteAddressOperator(input.Operator), NegateCondition: pointer.To(input.NegateCondition), MatchValues: stringSlicePointer(input.MatchValues)}}, nil
}

func expandQueryStringCondition(input CdnFrontDoorBatchRuleStringConditionModel) (rules.DeliveryRuleCondition, error) {
	if err := validateStandardCondition("query_string_condition", input.Operator, input.MatchValues); err != nil {
		return nil, err
	}
	return rules.DeliveryRuleQueryStringCondition{Name: rules.MatchVariableQueryString, Parameters: rules.QueryStringMatchConditionParameters{TypeName: rules.DeliveryRuleConditionParametersTypeDeliveryRuleQueryStringConditionParameters, Operator: rules.QueryStringOperator(input.Operator), NegateCondition: pointer.To(input.NegateCondition), MatchValues: stringSlicePointer(input.MatchValues), Transforms: transformsPointer(input.Transforms)}}, nil
}

func expandPostArgsCondition(input CdnFrontDoorBatchRulePostArgsConditionModel) (rules.DeliveryRuleCondition, error) {
	if err := validateStandardCondition("post_args_condition", input.Operator, input.MatchValues); err != nil {
		return nil, err
	}
	return rules.DeliveryRulePostArgsCondition{Name: rules.MatchVariablePostArgs, Parameters: rules.PostArgsMatchConditionParameters{TypeName: rules.DeliveryRuleConditionParametersTypeDeliveryRulePostArgsConditionParameters, Selector: pointer.To(input.PostArgsName), Operator: rules.PostArgsOperator(input.Operator), NegateCondition: pointer.To(input.NegateCondition), MatchValues: stringSlicePointer(input.MatchValues), Transforms: transformsPointer(input.Transforms)}}, nil
}

func expandRequestURICondition(input CdnFrontDoorBatchRuleStringConditionModel) (rules.DeliveryRuleCondition, error) {
	if err := validateStandardCondition("request_uri_condition", input.Operator, input.MatchValues); err != nil {
		return nil, err
	}
	return rules.DeliveryRuleRequestUriCondition{Name: rules.MatchVariableRequestUri, Parameters: rules.RequestUriMatchConditionParameters{TypeName: rules.DeliveryRuleConditionParametersTypeDeliveryRuleRequestUriConditionParameters, Operator: rules.RequestUriOperator(input.Operator), NegateCondition: pointer.To(input.NegateCondition), MatchValues: stringSlicePointer(input.MatchValues), Transforms: transformsPointer(input.Transforms)}}, nil
}

func expandRequestHeaderCondition(input CdnFrontDoorBatchRuleRequestHeaderConditionModel) (rules.DeliveryRuleCondition, error) {
	if err := validateStandardCondition("request_header_condition", input.Operator, input.MatchValues); err != nil {
		return nil, err
	}
	return rules.DeliveryRuleRequestHeaderCondition{Name: rules.MatchVariableRequestHeader, Parameters: rules.RequestHeaderMatchConditionParameters{TypeName: rules.DeliveryRuleConditionParametersTypeDeliveryRuleRequestHeaderConditionParameters, Selector: pointer.To(input.HeaderName), Operator: rules.RequestHeaderOperator(input.Operator), NegateCondition: pointer.To(input.NegateCondition), MatchValues: stringSlicePointer(input.MatchValues), Transforms: transformsPointer(input.Transforms)}}, nil
}

func expandRequestBodyCondition(input CdnFrontDoorBatchRuleStringConditionModel) (rules.DeliveryRuleCondition, error) {
	if err := validateStandardCondition("request_body_condition", input.Operator, input.MatchValues); err != nil {
		return nil, err
	}
	return rules.DeliveryRuleRequestBodyCondition{Name: rules.MatchVariableRequestBody, Parameters: rules.RequestBodyMatchConditionParameters{TypeName: rules.DeliveryRuleConditionParametersTypeDeliveryRuleRequestBodyConditionParameters, Operator: rules.RequestBodyOperator(input.Operator), NegateCondition: pointer.To(input.NegateCondition), MatchValues: stringSlicePointer(input.MatchValues), Transforms: transformsPointer(input.Transforms)}}, nil
}

func expandRequestMethodCondition(input CdnFrontDoorBatchRuleRequestMethodConditionModel) (rules.DeliveryRuleCondition, error) {
	if err := validateStandardCondition("request_method_condition", input.Operator, input.MatchValues); err != nil {
		return nil, err
	}
	return rules.DeliveryRuleRequestMethodCondition{Name: rules.MatchVariableRequestMethod, Parameters: rules.RequestMethodMatchConditionParameters{TypeName: rules.DeliveryRuleConditionParametersTypeDeliveryRuleRequestMethodConditionParameters, Operator: rules.RequestMethodOperator(input.Operator), NegateCondition: pointer.To(input.NegateCondition), MatchValues: requestMethodValuesPointer(input.MatchValues)}}, nil
}

func expandRequestSchemeCondition(input CdnFrontDoorBatchRuleRequestSchemeConditionModel) (rules.DeliveryRuleCondition, error) {
	if err := validateStandardCondition("request_scheme_condition", input.Operator, input.MatchValues); err != nil {
		return nil, err
	}
	return rules.DeliveryRuleRequestSchemeCondition{Name: rules.MatchVariableRequestScheme, Parameters: rules.RequestSchemeMatchConditionParameters{TypeName: rules.DeliveryRuleConditionParametersTypeDeliveryRuleRequestSchemeConditionParameters, Operator: rules.RequestSchemeMatchConditionParametersOperator(input.Operator), NegateCondition: pointer.To(input.NegateCondition), MatchValues: requestSchemeValuesPointer(input.MatchValues)}}, nil
}

func expandURLPathCondition(input CdnFrontDoorBatchRuleStringConditionModel) (rules.DeliveryRuleCondition, error) {
	if err := validateStandardCondition("url_path_condition", input.Operator, input.MatchValues); err != nil {
		return nil, err
	}
	return rules.DeliveryRuleURLPathCondition{Name: rules.MatchVariableURLPath, Parameters: rules.URLPathMatchConditionParameters{TypeName: rules.DeliveryRuleConditionParametersTypeDeliveryRuleURLPathMatchConditionParameters, Operator: rules.URLPathOperator(input.Operator), NegateCondition: pointer.To(input.NegateCondition), MatchValues: stringSlicePointer(input.MatchValues), Transforms: transformsPointer(input.Transforms)}}, nil
}

func expandURLFileExtensionCondition(input CdnFrontDoorBatchRuleStringConditionModel) (rules.DeliveryRuleCondition, error) {
	if err := validateStandardCondition("url_file_extension_condition", input.Operator, input.MatchValues); err != nil {
		return nil, err
	}
	return rules.DeliveryRuleURLFileExtensionCondition{Name: rules.MatchVariableURLFileExtension, Parameters: rules.URLFileExtensionMatchConditionParameters{TypeName: rules.DeliveryRuleConditionParametersTypeDeliveryRuleURLFileExtensionMatchConditionParameters, Operator: rules.URLFileExtensionOperator(input.Operator), NegateCondition: pointer.To(input.NegateCondition), MatchValues: stringSlicePointer(input.MatchValues), Transforms: transformsPointer(input.Transforms)}}, nil
}

func expandURLFilenameCondition(input CdnFrontDoorBatchRuleStringConditionModel) (rules.DeliveryRuleCondition, error) {
	if err := validateStandardCondition("url_filename_condition", input.Operator, input.MatchValues); err != nil {
		return nil, err
	}
	return rules.DeliveryRuleURLFileNameCondition{Name: rules.MatchVariableURLFileName, Parameters: rules.URLFileNameMatchConditionParameters{TypeName: rules.DeliveryRuleConditionParametersTypeDeliveryRuleURLFilenameConditionParameters, Operator: rules.URLFileNameOperator(input.Operator), NegateCondition: pointer.To(input.NegateCondition), MatchValues: stringSlicePointer(input.MatchValues), Transforms: transformsPointer(input.Transforms)}}, nil
}

func expandHTTPVersionCondition(input CdnFrontDoorBatchRuleHTTPVersionConditionModel) (rules.DeliveryRuleCondition, error) {
	if err := validateStandardCondition("http_version_condition", input.Operator, input.MatchValues); err != nil {
		return nil, err
	}
	return rules.DeliveryRuleHTTPVersionCondition{Name: rules.MatchVariableHTTPVersion, Parameters: rules.HTTPVersionMatchConditionParameters{TypeName: rules.DeliveryRuleConditionParametersTypeDeliveryRuleHTTPVersionConditionParameters, Operator: rules.HTTPVersionOperator(input.Operator), NegateCondition: pointer.To(input.NegateCondition), MatchValues: stringSlicePointer(input.MatchValues)}}, nil
}

func expandCookiesCondition(input CdnFrontDoorBatchRuleCookiesConditionModel) (rules.DeliveryRuleCondition, error) {
	if err := validateStandardCondition("cookies_condition", input.Operator, input.MatchValues); err != nil {
		return nil, err
	}
	return rules.DeliveryRuleCookiesCondition{Name: rules.MatchVariableCookies, Parameters: rules.CookiesMatchConditionParameters{TypeName: rules.DeliveryRuleConditionParametersTypeDeliveryRuleCookiesConditionParameters, Selector: pointer.To(input.CookieName), Operator: rules.CookiesOperator(input.Operator), NegateCondition: pointer.To(input.NegateCondition), MatchValues: stringSlicePointer(input.MatchValues), Transforms: transformsPointer(input.Transforms)}}, nil
}

func expandIsDeviceCondition(input CdnFrontDoorBatchRuleIsDeviceConditionModel) (rules.DeliveryRuleCondition, error) {
	if err := validateStandardCondition("is_device_condition", input.Operator, input.MatchValues); err != nil {
		return nil, err
	}
	return rules.DeliveryRuleIsDeviceCondition{Name: rules.MatchVariableIsDevice, Parameters: rules.IsDeviceMatchConditionParameters{TypeName: rules.DeliveryRuleConditionParametersTypeDeliveryRuleIsDeviceConditionParameters, Operator: rules.IsDeviceOperator(input.Operator), NegateCondition: pointer.To(input.NegateCondition), MatchValues: isDeviceValuesPointer(input.MatchValues)}}, nil
}

func expandSocketAddressCondition(input CdnFrontDoorBatchRuleStringConditionModel) (rules.DeliveryRuleCondition, error) {
	if err := validateAnyCondition("socket_address_condition", input.Operator, input.MatchValues); err != nil {
		return nil, err
	}

	if strings.EqualFold(input.Operator, string(rules.SocketAddrOperatorIPMatch)) {
		matchValues := make([]interface{}, 0, len(input.MatchValues))
		for _, matchValue := range input.MatchValues {
			matchValues = append(matchValues, matchValue)
			if _, errs := validate.FrontDoorRuleCidrIsValid(matchValue, "match_values"); len(errs) > 0 {
				return nil, fmt.Errorf("`socket_address_condition` is invalid: when `operator` is `IPMatch` the values in `match_values` must be valid IPv4 or IPv6 CIDRs, got %q", matchValue)
			}
		}

		if _, errs := validate.FrontDoorRuleCidrOverlap(matchValues, "match_values"); len(errs) > 0 {
			return nil, fmt.Errorf("`socket_address_condition` is invalid: %+v", errs[0])
		}
	}

	return rules.DeliveryRuleSocketAddrCondition{Name: rules.MatchVariableSocketAddr, Parameters: rules.SocketAddrMatchConditionParameters{TypeName: rules.DeliveryRuleConditionParametersTypeDeliveryRuleSocketAddrConditionParameters, Operator: rules.SocketAddrOperator(input.Operator), NegateCondition: pointer.To(input.NegateCondition), MatchValues: stringSlicePointer(input.MatchValues)}}, nil
}

func expandClientPortCondition(input CdnFrontDoorBatchRuleStringConditionModel) (rules.DeliveryRuleCondition, error) {
	if err := validateAnyCondition("client_port_condition", input.Operator, input.MatchValues); err != nil {
		return nil, err
	}
	return rules.DeliveryRuleClientPortCondition{Name: rules.MatchVariableClientPort, Parameters: rules.ClientPortMatchConditionParameters{TypeName: rules.DeliveryRuleConditionParametersTypeDeliveryRuleClientPortConditionParameters, Operator: rules.ClientPortOperator(input.Operator), NegateCondition: pointer.To(input.NegateCondition), MatchValues: stringSlicePointer(input.MatchValues)}}, nil
}

func expandServerPortCondition(input CdnFrontDoorBatchRuleStringConditionModel) (rules.DeliveryRuleCondition, error) {
	if err := validateStandardCondition("server_port_condition", input.Operator, input.MatchValues); err != nil {
		return nil, err
	}
	return rules.DeliveryRuleServerPortCondition{Name: rules.MatchVariableServerPort, Parameters: rules.ServerPortMatchConditionParameters{TypeName: rules.DeliveryRuleConditionParametersTypeDeliveryRuleServerPortConditionParameters, Operator: rules.ServerPortOperator(input.Operator), NegateCondition: pointer.To(input.NegateCondition), MatchValues: stringSlicePointer(input.MatchValues)}}, nil
}

func expandHostNameCondition(input CdnFrontDoorBatchRuleStringConditionModel) (rules.DeliveryRuleCondition, error) {
	if err := validateStandardCondition("host_name_condition", input.Operator, input.MatchValues); err != nil {
		return nil, err
	}
	return rules.DeliveryRuleHostNameCondition{Name: rules.MatchVariableHostName, Parameters: rules.HostNameMatchConditionParameters{TypeName: rules.DeliveryRuleConditionParametersTypeDeliveryRuleHostNameConditionParameters, Operator: rules.HostNameOperator(input.Operator), NegateCondition: pointer.To(input.NegateCondition), MatchValues: stringSlicePointer(input.MatchValues), Transforms: transformsPointer(input.Transforms)}}, nil
}

func expandSSLProtocolCondition(input CdnFrontDoorBatchRuleSSLProtocolConditionModel) (rules.DeliveryRuleCondition, error) {
	if err := validateStandardCondition("ssl_protocol_condition", input.Operator, input.MatchValues); err != nil {
		return nil, err
	}
	return rules.DeliveryRuleSslProtocolCondition{Name: rules.MatchVariableSslProtocol, Parameters: rules.SslProtocolMatchConditionParameters{TypeName: rules.DeliveryRuleConditionParametersTypeDeliveryRuleSslProtocolConditionParameters, Operator: rules.SslProtocolOperator(input.Operator), NegateCondition: pointer.To(input.NegateCondition), MatchValues: sslProtocolValuesPointer(input.MatchValues)}}, nil
}

func validateStandardCondition(configName, operator string, matchValues []string) error {
	return validate.CdnFrontDoorValidateConditionMatchValues(configName, operator, matchValues)
}

func validateAnyCondition(configName, operator string, matchValues []string) error {
	return validate.CdnFrontDoorValidateConditionMatchValues(configName, operator, matchValues)
}

func batchRuleHasConditions(input CdnFrontDoorBatchRuleConditionsModel) bool {
	return len(input.RemoteAddressCondition) > 0 ||
		len(input.RequestMethodCondition) > 0 ||
		len(input.QueryStringCondition) > 0 ||
		len(input.PostArgsCondition) > 0 ||
		len(input.RequestURICondition) > 0 ||
		len(input.RequestHeaderCondition) > 0 ||
		len(input.RequestBodyCondition) > 0 ||
		len(input.RequestSchemeCondition) > 0 ||
		len(input.URLPathCondition) > 0 ||
		len(input.URLFileExtensionCondition) > 0 ||
		len(input.URLFilenameCondition) > 0 ||
		len(input.HTTPVersionCondition) > 0 ||
		len(input.CookiesCondition) > 0 ||
		len(input.IsDeviceCondition) > 0 ||
		len(input.SocketAddressCondition) > 0 ||
		len(input.ClientPortCondition) > 0 ||
		len(input.ServerPortCondition) > 0 ||
		len(input.HostNameCondition) > 0 ||
		len(input.SSLProtocolCondition) > 0
}

func flattenStringCondition(operator string, negate bool, matchValues []string, transforms []string) CdnFrontDoorBatchRuleStringConditionModel {
	return CdnFrontDoorBatchRuleStringConditionModel{Operator: operator, NegateCondition: negate, MatchValues: matchValues, Transforms: transforms}
}

func transformsPointer(input []string) *[]rules.Transform {
	if len(input) == 0 {
		return nil
	}
	values := make([]rules.Transform, 0, len(input))
	for _, value := range input {
		values = append(values, rules.Transform(value))
	}
	return &values
}

func transformsToStrings(input *[]rules.Transform) []string {
	if input == nil {
		return nil
	}
	values := make([]string, 0, len(*input))
	for _, value := range *input {
		values = append(values, string(value))
	}
	return values
}

func stringSlicePointer(input []string) *[]string {
	if len(input) == 0 {
		return nil
	}
	return &input
}

func stringPointerOrNil(input string) *string {
	if input == "" {
		return nil
	}
	return pointer.To(input)
}

func pointerOrNil[T ~string](input T) *T {
	if input == "" {
		return nil
	}
	return &input
}

func splitCSV(input string) []string {
	if input == "" {
		return nil
	}
	return strings.Split(input, ",")
}

func requestMethodValuesPointer(input []string) *[]rules.RequestMethodMatchValue {
	if len(input) == 0 {
		return nil
	}
	values := make([]rules.RequestMethodMatchValue, 0, len(input))
	for _, value := range input {
		values = append(values, rules.RequestMethodMatchValue(value))
	}
	return &values
}

func requestMethodValuesToStrings(input *[]rules.RequestMethodMatchValue) []string {
	if input == nil {
		return nil
	}
	values := make([]string, 0, len(*input))
	for _, value := range *input {
		values = append(values, string(value))
	}
	return values
}

func requestSchemeValuesPointer(input []string) *[]rules.RequestSchemeMatchValue {
	if len(input) == 0 {
		return nil
	}
	values := make([]rules.RequestSchemeMatchValue, 0, len(input))
	for _, value := range input {
		values = append(values, rules.RequestSchemeMatchValue(value))
	}
	return &values
}

func requestSchemeValuesToStrings(input *[]rules.RequestSchemeMatchValue) []string {
	if input == nil {
		return nil
	}
	values := make([]string, 0, len(*input))
	for _, value := range *input {
		values = append(values, string(value))
	}
	return values
}

func isDeviceValuesPointer(input []string) *[]rules.IsDeviceMatchValue {
	if len(input) == 0 {
		return nil
	}
	values := make([]rules.IsDeviceMatchValue, 0, len(input))
	for _, value := range input {
		values = append(values, rules.IsDeviceMatchValue(value))
	}
	return &values
}

func isDeviceValuesToStrings(input *[]rules.IsDeviceMatchValue) []string {
	if input == nil {
		return nil
	}
	values := make([]string, 0, len(*input))
	for _, value := range *input {
		values = append(values, string(value))
	}
	return values
}

func sslProtocolValuesPointer(input []string) *[]rules.SslProtocol {
	if len(input) == 0 {
		return nil
	}
	values := make([]rules.SslProtocol, 0, len(input))
	for _, value := range input {
		values = append(values, rules.SslProtocol(value))
	}
	return &values
}

func sslProtocolValuesToStrings(input *[]rules.SslProtocol) []string {
	if input == nil {
		return nil
	}
	values := make([]string, 0, len(*input))
	for _, value := range *input {
		values = append(values, string(value))
	}
	return values
}
