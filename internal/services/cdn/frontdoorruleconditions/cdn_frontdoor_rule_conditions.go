// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cdnfrontdoorruleconditions

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-09-01/rules"
	helperValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontDoorConditionParameters struct {
	Name       rules.MatchVariable
	TypeName   rules.DeliveryRuleConditionParametersType
	ConfigName string
}

type CdnFrontDoorCondtionsMappings struct {
	ClientPort       CdnFrontDoorConditionParameters
	Cookies          CdnFrontDoorConditionParameters
	HostName         CdnFrontDoorConditionParameters
	HttpVersion      CdnFrontDoorConditionParameters
	IsDevice         CdnFrontDoorConditionParameters
	PostArgs         CdnFrontDoorConditionParameters
	QueryString      CdnFrontDoorConditionParameters
	RemoteAddress    CdnFrontDoorConditionParameters
	RequestBody      CdnFrontDoorConditionParameters
	RequestHeader    CdnFrontDoorConditionParameters
	RequestMethod    CdnFrontDoorConditionParameters
	RequestScheme    CdnFrontDoorConditionParameters
	RequestUri       CdnFrontDoorConditionParameters
	ServerPort       CdnFrontDoorConditionParameters
	SocketAddress    CdnFrontDoorConditionParameters
	SslProtocol      CdnFrontDoorConditionParameters
	UrlFileExtension CdnFrontDoorConditionParameters
	UrlFilename      CdnFrontDoorConditionParameters
	UrlPath          CdnFrontDoorConditionParameters
}

type normalizedSelector struct {
	name  *string
	value *string
}

type normalizedCondition struct {
	selector        *normalizedSelector
	operator        string
	negateCondition *bool
	matchValues     *[]string
	transforms      *[]rules.Transform
}

func InitializeCdnFrontDoorConditionMappings() *CdnFrontDoorCondtionsMappings {
	m := CdnFrontDoorCondtionsMappings{}

	m.ClientPort = CdnFrontDoorConditionParameters{
		Name:       rules.MatchVariableClientPort,
		TypeName:   rules.DeliveryRuleConditionParametersTypeDeliveryRuleClientPortConditionParameters,
		ConfigName: "client_port_condition",
	}

	m.Cookies = CdnFrontDoorConditionParameters{
		Name:       rules.MatchVariableCookies,
		TypeName:   rules.DeliveryRuleConditionParametersTypeDeliveryRuleCookiesConditionParameters,
		ConfigName: "cookies_condition",
	}

	m.HostName = CdnFrontDoorConditionParameters{
		Name:       rules.MatchVariableHostName,
		TypeName:   rules.DeliveryRuleConditionParametersTypeDeliveryRuleHostNameConditionParameters,
		ConfigName: "host_name_condition",
	}

	m.HttpVersion = CdnFrontDoorConditionParameters{
		Name:       rules.MatchVariableHTTPVersion,
		TypeName:   rules.DeliveryRuleConditionParametersTypeDeliveryRuleHTTPVersionConditionParameters,
		ConfigName: "http_version_condition",
	}

	m.IsDevice = CdnFrontDoorConditionParameters{
		Name:       rules.MatchVariableIsDevice,
		TypeName:   rules.DeliveryRuleConditionParametersTypeDeliveryRuleIsDeviceConditionParameters,
		ConfigName: "is_device_condition",
	}

	m.PostArgs = CdnFrontDoorConditionParameters{
		Name:       rules.MatchVariablePostArgs,
		TypeName:   rules.DeliveryRuleConditionParametersTypeDeliveryRulePostArgsConditionParameters,
		ConfigName: "post_args_condition",
	}

	m.QueryString = CdnFrontDoorConditionParameters{
		Name:       rules.MatchVariableQueryString,
		TypeName:   rules.DeliveryRuleConditionParametersTypeDeliveryRuleQueryStringConditionParameters,
		ConfigName: "query_string_condition",
	}

	m.RemoteAddress = CdnFrontDoorConditionParameters{
		Name:       rules.MatchVariableRemoteAddress,
		TypeName:   rules.DeliveryRuleConditionParametersTypeDeliveryRuleRemoteAddressConditionParameters,
		ConfigName: "remote_address_condition",
	}

	m.RequestBody = CdnFrontDoorConditionParameters{
		Name:       rules.MatchVariableRequestBody,
		TypeName:   rules.DeliveryRuleConditionParametersTypeDeliveryRuleRequestBodyConditionParameters,
		ConfigName: "request_body_condition",
	}

	m.RequestHeader = CdnFrontDoorConditionParameters{
		Name:       rules.MatchVariableRequestHeader,
		TypeName:   rules.DeliveryRuleConditionParametersTypeDeliveryRuleRequestHeaderConditionParameters,
		ConfigName: "request_header_condition",
	}

	m.RequestMethod = CdnFrontDoorConditionParameters{
		Name:       rules.MatchVariableRequestMethod,
		TypeName:   rules.DeliveryRuleConditionParametersTypeDeliveryRuleRequestMethodConditionParameters,
		ConfigName: "request_method_condition",
	}

	m.RequestScheme = CdnFrontDoorConditionParameters{
		Name:       rules.MatchVariableRequestScheme,
		TypeName:   rules.DeliveryRuleConditionParametersTypeDeliveryRuleRequestSchemeConditionParameters,
		ConfigName: "request_scheme_condition",
	}

	m.RequestUri = CdnFrontDoorConditionParameters{
		Name:       rules.MatchVariableRequestUri,
		TypeName:   rules.DeliveryRuleConditionParametersTypeDeliveryRuleRequestUriConditionParameters,
		ConfigName: "request_uri_condition",
	}

	m.ServerPort = CdnFrontDoorConditionParameters{
		Name:       rules.MatchVariableServerPort,
		TypeName:   rules.DeliveryRuleConditionParametersTypeDeliveryRuleServerPortConditionParameters,
		ConfigName: "server_port_condition",
	}

	m.SocketAddress = CdnFrontDoorConditionParameters{
		Name:       rules.MatchVariableSocketAddr,
		TypeName:   rules.DeliveryRuleConditionParametersTypeDeliveryRuleSocketAddrConditionParameters,
		ConfigName: "socket_address_condition",
	}

	m.SslProtocol = CdnFrontDoorConditionParameters{
		Name:       rules.MatchVariableSslProtocol,
		TypeName:   rules.DeliveryRuleConditionParametersTypeDeliveryRuleSslProtocolConditionParameters,
		ConfigName: "ssl_protocol_condition",
	}

	m.UrlFileExtension = CdnFrontDoorConditionParameters{
		Name:       rules.MatchVariableURLFileExtension,
		TypeName:   rules.DeliveryRuleConditionParametersTypeDeliveryRuleURLFileExtensionMatchConditionParameters,
		ConfigName: "url_file_extension_condition",
	}

	m.UrlFilename = CdnFrontDoorConditionParameters{
		Name:       rules.MatchVariableURLFileName,
		TypeName:   rules.DeliveryRuleConditionParametersTypeDeliveryRuleURLFilenameConditionParameters,
		ConfigName: "url_filename_condition",
	}

	m.UrlPath = CdnFrontDoorConditionParameters{
		Name:       rules.MatchVariableURLPath,
		TypeName:   rules.DeliveryRuleConditionParametersTypeDeliveryRuleURLPathMatchConditionParameters,
		ConfigName: "url_path_condition",
	}

	return &m
}

func expandNormalizeCdnFrontDoorTransforms(input []interface{}) []rules.Transform {
	transforms := make([]rules.Transform, 0)
	if len(input) == 0 {
		return transforms
	}

	for _, t := range input {
		transforms = append(transforms, rules.Transform(t.(string)))
	}

	return transforms
}

func flattenCdnFrontDoorNormalizedCondition(condition normalizedCondition) map[string]interface{} {
	operator := ""
	negateCondition := false
	matchValues := make([]interface{}, 0)
	conditionTransforms := make([]interface{}, 0)

	if condition.operator != "" {
		operator = condition.operator
	}

	if condition.negateCondition != nil {
		negateCondition = *condition.negateCondition
	}

	if condition.matchValues != nil {
		matchValues = utils.FlattenStringSlice(condition.matchValues)
	}

	flattened := map[string]interface{}{
		"operator":         operator,
		"negate_condition": negateCondition,
		"match_values":     matchValues,
	}

	if condition.selector != nil {
		flattened[*condition.selector.name] = *condition.selector.value
	}

	if condition.transforms != nil {
		for _, transform := range *condition.transforms {
			conditionTransforms = append(conditionTransforms, string(transform))
		}

		flattened["transforms"] = conditionTransforms
	}

	return flattened
}

func validateCdnFrontDoorExpandConditionOperatorValues(operator string, matchValues *[]string, m CdnFrontDoorConditionParameters) error {
	if operator == "" {
		return fmt.Errorf("%q is invalid: no 'operator' value has been set, got %q", m.ConfigName, operator)
	}

	// NOTE: There are now 14 different 'Any' operators in the new API, however they are all the same value so I am just hardcoding the value evaluation here,
	//       the multiple values appears to be by design as they have exposed an "Any" field for each condition type (e.g., 'ClientPortOperatorAny',
	//       'CookiesOperatorAny', 'HostNameOperatorAny', etc.) due to the swagger issue raised by TomBuildsStuff...
	if operator == "Any" && len(*matchValues) > 0 {
		return fmt.Errorf("%q is invalid: the 'match_values' field must not be set if the conditions 'operator' is set to 'Any'", m.ConfigName)
	}

	// make the 'match_values' field required if the operator is not set to 'Any'...
	if operator != "Any" && len(*matchValues) == 0 {
		return fmt.Errorf("%q is invalid: the 'match_values' field must be set if the conditions 'operator' is not set to 'Any'", m.ConfigName)
	}

	return nil
}

func ExpandCdnFrontDoorRemoteAddressCondition(input []interface{}) (*[]rules.DeliveryRuleCondition, error) {
	output := make([]rules.DeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.RemoteAddress

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := rules.DeliveryRuleRemoteAddressCondition{
			Name: conditionMapping.Name,
			Parameters: rules.RemoteAddressMatchConditionParameters{
				TypeName:        conditionMapping.TypeName,
				Operator:        rules.RemoteAddressOperator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if condition.Parameters.Operator == rules.RemoteAddressOperatorGeoMatch {
			for _, matchValue := range item["match_values"].([]interface{}) {
				if matchValue != nil {
					if ok, _ := helperValidate.RegExHelper(matchValue, "match_values", `^[A-Z]{2}$`); !ok {
						return nil, fmt.Errorf("%q is invalid: when the 'operator' is set to 'GeoMatch' the value must be a valid country code consisting of 2 uppercase characters, got %q", conditionMapping.ConfigName, matchValue)
					}
				}
			}
		}

		if condition.Parameters.Operator == rules.RemoteAddressOperatorIPMatch {
			// make sure all of the passed CIDRs are valid
			for _, matchValue := range item["match_values"].([]interface{}) {
				if _, err := validate.FrontDoorRuleCidrIsValid(matchValue, "match_values"); err != nil {
					return nil, fmt.Errorf("%q is invalid: when the 'operator' is set to 'IPMatch' the 'match_values' must be a valid IPv4 or IPv6 CIDR, got %q", conditionMapping.ConfigName, matchValue.(string))
				}
			}

			// Check for CIDR overlap and CIDR duplicates in the match values
			_, err := validate.FrontDoorRuleCidrOverlap(item["match_values"].([]interface{}), "match_values")
			if err != nil {
				return nil, fmt.Errorf("%q is invalid: %+v", conditionMapping.ConfigName, err)
			}
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontDoorRequestMethodCondition(input []interface{}) (*[]rules.DeliveryRuleCondition, error) {
	output := make([]rules.DeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.RequestMethod

	for _, v := range input {
		item := v.(map[string]interface{})
		matchValuesRaw := item["match_values"].(*pluginsdk.Set).List()

		condition := rules.DeliveryRuleRequestMethodCondition{
			Name: conditionMapping.Name,
			Parameters: rules.RequestMethodMatchConditionParameters{
				TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleRequestMethodConditionParameters,
				Operator:        rules.RequestMethodOperator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     expandRequestMethodMatchValues(matchValuesRaw),
			},
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(string(condition.Parameters.Operator), utils.ExpandStringSlice(matchValuesRaw), conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontDoorQueryStringCondition(input []interface{}) (*[]rules.DeliveryRuleCondition, error) {
	output := make([]rules.DeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.QueryString

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := rules.DeliveryRuleQueryStringCondition{
			Name: conditionMapping.Name,
			Parameters: rules.QueryStringMatchConditionParameters{
				TypeName:        conditionMapping.TypeName,
				Operator:        rules.QueryStringOperator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		transformsRaw := item["transforms"].(*pluginsdk.Set).List()
		if len(transformsRaw) != 0 {
			expanded := expandNormalizeCdnFrontDoorTransforms(transformsRaw)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontDoorPostArgsCondition(input []interface{}) (*[]rules.DeliveryRuleCondition, error) {
	output := make([]rules.DeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.PostArgs

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := rules.DeliveryRulePostArgsCondition{
			Name: conditionMapping.Name,
			Parameters: rules.PostArgsMatchConditionParameters{
				TypeName:        conditionMapping.TypeName,
				Selector:        pointer.To(item["post_args_name"].(string)),
				Operator:        rules.PostArgsOperator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		transformsRaw := item["transforms"].(*pluginsdk.Set).List()
		if len(transformsRaw) != 0 {
			expanded := expandNormalizeCdnFrontDoorTransforms(transformsRaw)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontDoorRequestUriCondition(input []interface{}) (*[]rules.DeliveryRuleCondition, error) {
	output := make([]rules.DeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.RequestUri

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := rules.DeliveryRuleRequestUriCondition{
			Name: conditionMapping.Name,
			Parameters: rules.RequestUriMatchConditionParameters{
				TypeName:        conditionMapping.TypeName,
				Operator:        rules.RequestUriOperator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		transformsRaw := item["transforms"].(*pluginsdk.Set).List()
		if len(transformsRaw) != 0 {
			expanded := expandNormalizeCdnFrontDoorTransforms(transformsRaw)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontDoorRequestHeaderCondition(input []interface{}) (*[]rules.DeliveryRuleCondition, error) {
	output := make([]rules.DeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.RequestHeader

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := rules.DeliveryRuleRequestHeaderCondition{
			Name: conditionMapping.Name,
			Parameters: rules.RequestHeaderMatchConditionParameters{
				TypeName:        conditionMapping.TypeName,
				Selector:        pointer.To(item["header_name"].(string)),
				Operator:        rules.RequestHeaderOperator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		transformsRaw := item["transforms"].(*pluginsdk.Set).List()
		if len(transformsRaw) != 0 {
			expanded := expandNormalizeCdnFrontDoorTransforms(transformsRaw)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontDoorRequestBodyCondition(input []interface{}) (*[]rules.DeliveryRuleCondition, error) {
	output := make([]rules.DeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.RequestBody

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := rules.DeliveryRuleRequestBodyCondition{
			Name: conditionMapping.Name,
			Parameters: rules.RequestBodyMatchConditionParameters{
				TypeName:        conditionMapping.TypeName,
				Operator:        rules.RequestBodyOperator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		transformsRaw := item["transforms"].(*pluginsdk.Set).List()
		if len(transformsRaw) != 0 {
			expanded := expandNormalizeCdnFrontDoorTransforms(transformsRaw)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontDoorRequestSchemeCondition(input []interface{}) (*[]rules.DeliveryRuleCondition, error) {
	output := make([]rules.DeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.RequestScheme

	for _, v := range input {
		item := v.(map[string]interface{})
		matchValuesRaw := item["match_values"].([]interface{})

		condition := rules.DeliveryRuleRequestSchemeCondition{
			Name: conditionMapping.Name,
			Parameters: rules.RequestSchemeMatchConditionParameters{
				TypeName:        conditionMapping.TypeName,
				Operator:        rules.Operator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     expandRequestSchemeMatchValues(matchValuesRaw),
			},
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(string(condition.Parameters.Operator), utils.ExpandStringSlice(matchValuesRaw), conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontDoorUrlPathCondition(input []interface{}) (*[]rules.DeliveryRuleCondition, error) {
	output := make([]rules.DeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.UrlPath

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := rules.DeliveryRuleURLPathCondition{
			Name: conditionMapping.Name,
			Parameters: rules.URLPathMatchConditionParameters{
				TypeName:        conditionMapping.TypeName,
				Operator:        rules.URLPathOperator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		transformsRaw := item["transforms"].(*pluginsdk.Set).List()
		if len(transformsRaw) != 0 {
			expanded := expandNormalizeCdnFrontDoorTransforms(transformsRaw)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontDoorUrlFileExtensionCondition(input []interface{}) (*[]rules.DeliveryRuleCondition, error) {
	output := make([]rules.DeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.UrlFileExtension

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := rules.DeliveryRuleURLFileExtensionCondition{
			Name: conditionMapping.Name,
			Parameters: rules.URLFileExtensionMatchConditionParameters{
				TypeName:        conditionMapping.TypeName,
				Operator:        rules.URLFileExtensionOperator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		transformsRaw := item["transforms"].(*pluginsdk.Set).List()
		if len(transformsRaw) != 0 {
			expanded := expandNormalizeCdnFrontDoorTransforms(transformsRaw)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontDoorUrlFileNameCondition(input []interface{}) (*[]rules.DeliveryRuleCondition, error) {
	output := make([]rules.DeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.UrlFilename

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := rules.DeliveryRuleURLFileNameCondition{
			Name: conditionMapping.Name,
			Parameters: rules.URLFileNameMatchConditionParameters{
				TypeName:        conditionMapping.TypeName,
				Operator:        rules.URLFileNameOperator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		transformsRaw := item["transforms"].(*pluginsdk.Set).List()
		if len(transformsRaw) != 0 {
			expanded := expandNormalizeCdnFrontDoorTransforms(transformsRaw)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontDoorHttpVersionCondition(input []interface{}) (*[]rules.DeliveryRuleCondition, error) {
	output := make([]rules.DeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.HttpVersion

	for _, v := range input {
		item := v.(map[string]interface{})
		matchValuesRaw := item["match_values"].(*pluginsdk.Set).List()

		condition := rules.DeliveryRuleHTTPVersionCondition{
			Name: conditionMapping.Name,
			Parameters: rules.HTTPVersionMatchConditionParameters{
				TypeName:        conditionMapping.TypeName,
				Operator:        rules.HTTPVersionOperator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(matchValuesRaw),
			},
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontDoorCookiesCondition(input []interface{}) (*[]rules.DeliveryRuleCondition, error) {
	output := make([]rules.DeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.Cookies

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := rules.DeliveryRuleCookiesCondition{
			Name: conditionMapping.Name,
			Parameters: rules.CookiesMatchConditionParameters{
				TypeName:        conditionMapping.TypeName,
				Selector:        pointer.To(item["cookie_name"].(string)),
				Operator:        rules.CookiesOperator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		transformsRaw := item["transforms"].(*pluginsdk.Set).List()
		if len(transformsRaw) != 0 {
			expanded := expandNormalizeCdnFrontDoorTransforms(transformsRaw)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontDoorIsDeviceCondition(input []interface{}) (*[]rules.DeliveryRuleCondition, error) {
	output := make([]rules.DeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.IsDevice

	for _, v := range input {
		item := v.(map[string]interface{})
		matchValuesRaw := item["match_values"].([]interface{})

		condition := rules.DeliveryRuleIsDeviceCondition{
			Name: conditionMapping.Name,
			Parameters: rules.IsDeviceMatchConditionParameters{
				TypeName:        conditionMapping.TypeName,
				Operator:        rules.IsDeviceOperator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     expandIsDeviceMatchValues(matchValuesRaw),
			},
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(string(condition.Parameters.Operator), utils.ExpandStringSlice(matchValuesRaw), conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontDoorSocketAddressCondition(input []interface{}) (*[]rules.DeliveryRuleCondition, error) {
	output := make([]rules.DeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.SocketAddress

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := rules.DeliveryRuleSocketAddrCondition{
			Name: conditionMapping.Name,
			Parameters: rules.SocketAddrMatchConditionParameters{
				TypeName:        conditionMapping.TypeName,
				Operator:        rules.SocketAddrOperator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if condition.Parameters.Operator == rules.SocketAddrOperatorIPMatch {
			// make sure all of the passed CIDRs are valid
			for _, matchValue := range item["match_values"].([]interface{}) {
				if _, err := validate.FrontDoorRuleCidrIsValid(matchValue, "match_values"); err != nil {
					return nil, fmt.Errorf("%q is invalid: when the 'operator' is set to 'IPMatch' the 'match_values' must be a valid IPv4 or IPv6 CIDR, got %q", conditionMapping.ConfigName, matchValue.(string))
				}
			}

			// Check for CIDR overlap and CIDR duplicates in the match values
			_, err := validate.FrontDoorRuleCidrOverlap(item["match_values"].([]interface{}), "match_values")
			if err != nil {
				return nil, fmt.Errorf("%q is invalid: %+v", conditionMapping.ConfigName, err)
			}
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontDoorClientPortCondition(input []interface{}) (*[]rules.DeliveryRuleCondition, error) {
	output := make([]rules.DeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.ClientPort

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := rules.DeliveryRuleClientPortCondition{
			Name: conditionMapping.Name,
			Parameters: rules.ClientPortMatchConditionParameters{
				TypeName:        conditionMapping.TypeName,
				Operator:        rules.ClientPortOperator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontDoorServerPortCondition(input []interface{}) (*[]rules.DeliveryRuleCondition, error) {
	output := make([]rules.DeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.ServerPort

	for _, v := range input {
		item := v.(map[string]interface{})
		matchValuesRaw := item["match_values"].(*pluginsdk.Set).List()

		condition := rules.DeliveryRuleServerPortCondition{
			Name: conditionMapping.Name,
			Parameters: rules.ServerPortMatchConditionParameters{
				TypeName:        conditionMapping.TypeName,
				Operator:        rules.ServerPortOperator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(matchValuesRaw),
			},
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontDoorHostNameCondition(input []interface{}) (*[]rules.DeliveryRuleCondition, error) {
	output := make([]rules.DeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.HostName

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := rules.DeliveryRuleHostNameCondition{
			Name: conditionMapping.Name,
			Parameters: rules.HostNameMatchConditionParameters{
				TypeName:        conditionMapping.TypeName,
				Operator:        rules.HostNameOperator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		transformsRaw := item["transforms"].(*pluginsdk.Set).List()
		if len(transformsRaw) != 0 {
			expanded := expandNormalizeCdnFrontDoorTransforms(transformsRaw)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontDoorSslProtocolCondition(input []interface{}) (*[]rules.DeliveryRuleCondition, error) {
	output := make([]rules.DeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.SslProtocol

	for _, v := range input {
		item := v.(map[string]interface{})

		matchValues := make([]rules.SslProtocol, 0)
		validationMatchValues := make([]string, 0)
		matchValuesRaw := item["match_values"].(*pluginsdk.Set).List()

		for _, value := range matchValuesRaw {
			matchValues = append(matchValues, rules.SslProtocol(value.(string)))
			validationMatchValues = append(validationMatchValues, value.(string))
		}

		condition := rules.DeliveryRuleSslProtocolCondition{
			Name: conditionMapping.Name,
			Parameters: rules.SslProtocolMatchConditionParameters{
				TypeName:        conditionMapping.TypeName,
				Operator:        rules.SslProtocolOperator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     &matchValues,
			},
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(string(condition.Parameters.Operator), &validationMatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func FlattenFrontdoorRemoteAddressCondition(input rules.DeliveryRuleRemoteAddressCondition) (map[string]interface{}, error) {
	params := input.Parameters

	normalized := normalizedCondition{
		selector:        nil,
		operator:        string(params.Operator),
		negateCondition: params.NegateCondition,
		matchValues:     params.MatchValues,
		transforms:      nil,
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorRequestMethodCondition(input rules.DeliveryRuleRequestMethodCondition) (map[string]interface{}, error) {
	params := input.Parameters

	normalized := normalizedCondition{
		selector:        nil,
		operator:        string(params.Operator),
		negateCondition: params.NegateCondition,
		matchValues:     flattenRequestMethodMatchValues(params.MatchValues),
		transforms:      nil,
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorQueryStringCondition(input rules.DeliveryRuleQueryStringCondition) (map[string]interface{}, error) {
	params := input.Parameters

	normalized := normalizedCondition{
		selector:        nil,
		operator:        string(params.Operator),
		negateCondition: params.NegateCondition,
		matchValues:     params.MatchValues,
		transforms:      params.Transforms,
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorPostArgsCondition(input rules.DeliveryRulePostArgsCondition) (map[string]interface{}, error) {
	params := input.Parameters

	normalized := normalizedCondition{
		selector:        &normalizedSelector{name: pointer.To("post_args_name"), value: params.Selector},
		operator:        string(params.Operator),
		negateCondition: params.NegateCondition,
		matchValues:     params.MatchValues,
		transforms:      params.Transforms,
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorRequestUriCondition(input rules.DeliveryRuleRequestUriCondition) (map[string]interface{}, error) {
	params := input.Parameters

	normalized := normalizedCondition{
		selector:        nil,
		operator:        string(params.Operator),
		negateCondition: params.NegateCondition,
		matchValues:     params.MatchValues,
		transforms:      params.Transforms,
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorRequestHeaderCondition(input rules.DeliveryRuleRequestHeaderCondition) (map[string]interface{}, error) {
	params := input.Parameters

	normalized := normalizedCondition{
		selector:        &normalizedSelector{name: pointer.To("header_name"), value: params.Selector},
		operator:        string(params.Operator),
		negateCondition: params.NegateCondition,
		matchValues:     params.MatchValues,
		transforms:      params.Transforms,
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorRequestBodyCondition(input rules.DeliveryRuleRequestBodyCondition) (map[string]interface{}, error) {
	params := input.Parameters

	normalized := normalizedCondition{
		selector:        nil,
		operator:        string(params.Operator),
		negateCondition: params.NegateCondition,
		matchValues:     params.MatchValues,
		transforms:      params.Transforms,
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorRequestSchemeCondition(input rules.DeliveryRuleRequestSchemeCondition) (map[string]interface{}, error) {
	params := input.Parameters

	normalized := normalizedCondition{
		selector:        nil,
		operator:        string(params.Operator),
		negateCondition: params.NegateCondition,
		matchValues:     flattenRequestSchemeMatchValues(params.MatchValues),
		transforms:      nil,
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorUrlPathCondition(input rules.DeliveryRuleURLPathCondition) (map[string]interface{}, error) {
	params := input.Parameters

	normalized := normalizedCondition{
		selector:        nil,
		operator:        string(params.Operator),
		negateCondition: params.NegateCondition,
		matchValues:     params.MatchValues,
		transforms:      params.Transforms,
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorUrlFileExtensionCondition(input rules.DeliveryRuleURLFileExtensionCondition) (map[string]interface{}, error) {
	params := input.Parameters

	normalized := normalizedCondition{
		selector:        nil,
		operator:        string(params.Operator),
		negateCondition: params.NegateCondition,
		matchValues:     params.MatchValues,
		transforms:      params.Transforms,
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorUrlFileNameCondition(input rules.DeliveryRuleURLFileNameCondition) (map[string]interface{}, error) {
	params := input.Parameters

	normalized := normalizedCondition{
		selector:        nil,
		operator:        string(params.Operator),
		negateCondition: params.NegateCondition,
		matchValues:     params.MatchValues,
		transforms:      params.Transforms,
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorHttpVersionCondition(input rules.DeliveryRuleHTTPVersionCondition) (map[string]interface{}, error) {
	params := input.Parameters

	normalized := normalizedCondition{
		selector:        nil,
		operator:        string(params.Operator),
		negateCondition: params.NegateCondition,
		matchValues:     params.MatchValues,
		transforms:      nil,
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorCookiesCondition(input rules.DeliveryRuleCookiesCondition) (map[string]interface{}, error) {
	params := input.Parameters

	normalized := normalizedCondition{
		selector:        &normalizedSelector{name: pointer.To("cookie_name"), value: params.Selector},
		operator:        string(params.Operator),
		negateCondition: params.NegateCondition,
		matchValues:     params.MatchValues,
		transforms:      params.Transforms,
	}

	// flatten the normalized condition regardless if it is a stub or not...
	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorIsDeviceCondition(input rules.DeliveryRuleIsDeviceCondition) (map[string]interface{}, error) {
	params := input.Parameters

	normalized := normalizedCondition{
		selector:        nil,
		operator:        string(params.Operator),
		negateCondition: params.NegateCondition,
		matchValues:     flattenIsDeviceMatchValues(params.MatchValues),
		transforms:      nil,
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorSocketAddressCondition(input rules.DeliveryRuleSocketAddrCondition) (map[string]interface{}, error) {
	params := input.Parameters

	normalized := normalizedCondition{
		selector:        nil,
		operator:        string(params.Operator),
		negateCondition: params.NegateCondition,
		matchValues:     params.MatchValues,
		transforms:      nil,
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorClientPortCondition(input rules.DeliveryRuleClientPortCondition) (map[string]interface{}, error) {
	params := input.Parameters

	normalized := normalizedCondition{
		selector:        nil,
		operator:        string(params.Operator),
		negateCondition: params.NegateCondition,
		matchValues:     params.MatchValues,
		transforms:      nil,
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorServerPortCondition(input rules.DeliveryRuleServerPortCondition) (map[string]interface{}, error) {
	params := input.Parameters

	normalized := normalizedCondition{
		selector:        nil,
		operator:        string(params.Operator),
		negateCondition: params.NegateCondition,
		matchValues:     params.MatchValues,
		transforms:      nil,
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorHostNameCondition(input rules.DeliveryRuleHostNameCondition) (map[string]interface{}, error) {
	params := input.Parameters

	normalized := normalizedCondition{
		selector:        nil,
		operator:        string(params.Operator),
		negateCondition: params.NegateCondition,
		matchValues:     params.MatchValues,
		transforms:      params.Transforms,
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorSslProtocolCondition(input rules.DeliveryRuleSslProtocolCondition) (map[string]interface{}, error) {
	params := input.Parameters

	matchValues := make([]string, 0)
	for _, value := range *params.MatchValues {
		matchValues = append(matchValues, string(value))
	}

	normalized := normalizedCondition{
		selector:        nil,
		operator:        string(params.Operator),
		negateCondition: params.NegateCondition,
		matchValues:     &matchValues,
		transforms:      nil,
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func expandRequestMethodMatchValues(input []interface{}) *[]rules.RequestMethodMatchValue {
	result := make([]rules.RequestMethodMatchValue, 0)
	if len(input) == 0 {
		return nil
	}

	for _, v := range input {
		matchValue := v.(string)
		result = append(result, rules.RequestMethodMatchValue(matchValue))
	}

	return &result
}

func flattenRequestMethodMatchValues(input *[]rules.RequestMethodMatchValue) *[]string {
	result := make([]string, 0)
	if input == nil {
		return nil
	}

	for _, v := range *input {
		matchValue := string(v)
		result = append(result, matchValue)
	}

	return &result
}

func expandRequestSchemeMatchValues(input []interface{}) *[]rules.RequestSchemeMatchValue {
	result := make([]rules.RequestSchemeMatchValue, 0)
	if len(input) == 0 {
		return nil
	}

	for _, v := range input {
		matchValue := v.(string)
		result = append(result, rules.RequestSchemeMatchValue(matchValue))
	}

	return &result
}

func flattenRequestSchemeMatchValues(input *[]rules.RequestSchemeMatchValue) *[]string {
	result := make([]string, 0)
	if input == nil {
		return nil
	}

	for _, v := range *input {
		matchValue := string(v)
		result = append(result, matchValue)
	}

	return &result
}

func expandIsDeviceMatchValues(input []interface{}) *[]rules.IsDeviceMatchValue {
	result := make([]rules.IsDeviceMatchValue, 0)
	if len(input) == 0 {
		return nil
	}

	for _, v := range input {
		matchValue := v.(string)
		result = append(result, rules.IsDeviceMatchValue(matchValue))
	}

	return &result
}

func flattenIsDeviceMatchValues(input *[]rules.IsDeviceMatchValue) *[]string {
	result := make([]string, 0)
	if input == nil {
		return nil
	}

	for _, v := range *input {
		matchValue := string(v)
		result = append(result, matchValue)
	}

	return &result
}
