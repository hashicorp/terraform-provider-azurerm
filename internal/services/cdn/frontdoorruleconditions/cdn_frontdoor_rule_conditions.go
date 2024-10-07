// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package CdnFrontDoorruleconditions

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	cdnValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontDoorConditionParameters struct {
	Name       cdn.Name
	TypeName   string
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
	transforms      *[]cdn.Transform
}

func InitializeCdnFrontDoorConditionMappings() *CdnFrontDoorCondtionsMappings {
	m := CdnFrontDoorCondtionsMappings{}

	m.ClientPort = CdnFrontDoorConditionParameters{
		Name:       cdn.NameClientPort,
		TypeName:   "DeliveryRuleClientPortConditionParameters",
		ConfigName: "client_port_condition",
	}

	m.Cookies = CdnFrontDoorConditionParameters{
		Name:       cdn.NameCookies,
		TypeName:   "DeliveryRuleCookiesConditionParameters",
		ConfigName: "cookies_condition",
	}

	m.HostName = CdnFrontDoorConditionParameters{
		Name:       cdn.NameHostName,
		TypeName:   "DeliveryRuleHostNameConditionParameters",
		ConfigName: "host_name_condition",
	}

	m.HttpVersion = CdnFrontDoorConditionParameters{
		Name:       cdn.NameHTTPVersion,
		TypeName:   "DeliveryRuleHttpVersionConditionParameters",
		ConfigName: "http_version_condition",
	}

	m.IsDevice = CdnFrontDoorConditionParameters{
		Name:       cdn.NameIsDevice,
		TypeName:   "DeliveryRuleIsDeviceConditionParameters",
		ConfigName: "is_device_condition",
	}

	m.PostArgs = CdnFrontDoorConditionParameters{
		Name:       cdn.NamePostArgs,
		TypeName:   "DeliveryRulePostArgsConditionParameters",
		ConfigName: "post_args_condition",
	}

	m.QueryString = CdnFrontDoorConditionParameters{
		Name:       cdn.NameQueryString,
		TypeName:   "DeliveryRuleQueryStringConditionParameters",
		ConfigName: "query_string_condition",
	}

	m.RemoteAddress = CdnFrontDoorConditionParameters{
		Name:       cdn.NameRemoteAddress,
		TypeName:   "DeliveryRuleRemoteAddressConditionParameters",
		ConfigName: "remote_address_condition",
	}

	m.RequestBody = CdnFrontDoorConditionParameters{
		Name:       cdn.NameRequestBody,
		TypeName:   "DeliveryRuleRequestBodyConditionParameters",
		ConfigName: "request_body_condition",
	}

	m.RequestHeader = CdnFrontDoorConditionParameters{
		Name:       cdn.NameRequestHeader,
		TypeName:   "DeliveryRuleRequestHeaderConditionParameters",
		ConfigName: "request_header_condition",
	}

	m.RequestMethod = CdnFrontDoorConditionParameters{
		Name:       cdn.NameRequestMethod,
		TypeName:   "DeliveryRuleRequestMethodConditionParameters",
		ConfigName: "request_method_condition",
	}

	m.RequestScheme = CdnFrontDoorConditionParameters{
		Name:       cdn.NameRequestScheme,
		TypeName:   "DeliveryRuleRequestSchemeConditionParameters",
		ConfigName: "request_scheme_condition",
	}

	m.RequestUri = CdnFrontDoorConditionParameters{
		Name:       cdn.NameRequestURI,
		TypeName:   "DeliveryRuleRequestUriConditionParameters",
		ConfigName: "request_uri_condition",
	}

	m.ServerPort = CdnFrontDoorConditionParameters{
		Name:       cdn.NameServerPort,
		TypeName:   "DeliveryRuleServerPortConditionParameters",
		ConfigName: "server_port_condition",
	}

	m.SocketAddress = CdnFrontDoorConditionParameters{
		Name:       cdn.NameSocketAddr,
		TypeName:   "DeliveryRuleSocketAddrConditionParameters",
		ConfigName: "socket_address_condition",
	}

	m.SslProtocol = CdnFrontDoorConditionParameters{
		Name:       cdn.NameSslProtocol,
		TypeName:   "DeliveryRuleSslProtocolConditionParameters",
		ConfigName: "ssl_protocol_condition",
	}

	m.UrlFileExtension = CdnFrontDoorConditionParameters{
		Name:       cdn.NameURLFileExtension,
		TypeName:   "DeliveryRuleUrlFileExtensionMatchConditionParameters",
		ConfigName: "url_file_extension_condition",
	}

	m.UrlFilename = CdnFrontDoorConditionParameters{
		Name:       cdn.NameURLFileName,
		TypeName:   "DeliveryRuleUrlFilenameConditionParameters",
		ConfigName: "url_filename_condition",
	}

	m.UrlPath = CdnFrontDoorConditionParameters{
		Name:       cdn.NameURLPath,
		TypeName:   "DeliveryRuleUrlPathMatchConditionParameters",
		ConfigName: "url_path_condition",
	}

	return &m
}

func expandNormalizeCdnFrontDoorTransforms(input []interface{}) []cdn.Transform {
	transforms := make([]cdn.Transform, 0)
	if len(input) == 0 {
		return transforms
	}

	for _, t := range input {
		transforms = append(transforms, cdn.Transform(t.(string)))
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

	if operator == string(cdn.OperatorAny) && len(*matchValues) > 0 {
		return fmt.Errorf("%q is invalid: the 'match_values' field must not be set if the conditions 'operator' is set to 'Any'", m.ConfigName)
	}

	// make the 'match_values' field required if the operator is not set to 'Any'...
	if operator != string(cdn.OperatorAny) && len(*matchValues) == 0 {
		return fmt.Errorf("%q is invalid: the 'match_values' field must be set if the conditions 'operator' is not set to 'Any'", m.ConfigName)
	}

	return nil
}

func ExpandCdnFrontDoorRemoteAddressCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.RemoteAddress

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := cdn.DeliveryRuleRemoteAddressCondition{
			Name: conditionMapping.Name,
			Parameters: &cdn.RemoteAddressMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        cdn.RemoteAddressOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if condition.Parameters.Operator == cdn.RemoteAddressOperatorGeoMatch {
			for _, matchValue := range item["match_values"].([]interface{}) {
				if matchValue != nil {
					if ok, _ := validate.RegExHelper(matchValue, "match_values", `^[A-Z]{2}$`); !ok {
						return nil, fmt.Errorf("%q is invalid: when the 'operator' is set to 'GeoMatch' the value must be a valid country code consisting of 2 uppercase characters, got %q", conditionMapping.ConfigName, matchValue)
					}
				}
			}
		}

		if condition.Parameters.Operator == cdn.RemoteAddressOperatorIPMatch {
			// make sure all of the passed CIDRs are valid
			for _, matchValue := range item["match_values"].([]interface{}) {
				if _, err := cdnValidate.FrontDoorRuleCidrIsValid(matchValue, "match_values"); err != nil {
					return nil, fmt.Errorf("%q is invalid: when the 'operator' is set to 'IPMatch' the 'match_values' must be a valid IPv4 or IPv6 CIDR, got %q", conditionMapping.ConfigName, matchValue.(string))
				}
			}

			// Check for CIDR overlap and CIDR duplicates in the match values
			_, err := cdnValidate.FrontDoorRuleCidrOverlap(item["match_values"].([]interface{}), "match_values")
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

func ExpandCdnFrontDoorRequestMethodCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.RequestMethod

	for _, v := range input {
		item := v.(map[string]interface{})
		matchValuesRaw := item["match_values"].(*pluginsdk.Set).List()

		condition := cdn.DeliveryRuleRequestMethodCondition{
			Name: conditionMapping.Name,
			Parameters: &cdn.RequestMethodMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        utils.String(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(matchValuesRaw),
			},
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(*condition.Parameters.Operator, condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontDoorQueryStringCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.QueryString

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := cdn.DeliveryRuleQueryStringCondition{
			Name: conditionMapping.Name,
			Parameters: &cdn.QueryStringMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        cdn.QueryStringOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
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

func ExpandCdnFrontDoorPostArgsCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.PostArgs

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := cdn.DeliveryRulePostArgsCondition{
			Name: conditionMapping.Name,
			Parameters: &cdn.PostArgsMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Selector:        utils.String(item["post_args_name"].(string)),
				Operator:        cdn.PostArgsOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
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

func ExpandCdnFrontDoorRequestUriCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.RequestUri

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := cdn.DeliveryRuleRequestURICondition{
			Name: conditionMapping.Name,
			Parameters: &cdn.RequestURIMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        cdn.RequestURIOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
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

func ExpandCdnFrontDoorRequestHeaderCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.RequestHeader

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := cdn.DeliveryRuleRequestHeaderCondition{
			Name: conditionMapping.Name,
			Parameters: &cdn.RequestHeaderMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Selector:        utils.String(item["header_name"].(string)),
				Operator:        cdn.RequestHeaderOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
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

func ExpandCdnFrontDoorRequestBodyCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.RequestBody

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := cdn.DeliveryRuleRequestBodyCondition{
			Name: conditionMapping.Name,
			Parameters: &cdn.RequestBodyMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        cdn.RequestBodyOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
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

func ExpandCdnFrontDoorRequestSchemeCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.RequestScheme

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := cdn.DeliveryRuleRequestSchemeCondition{
			Name: conditionMapping.Name,
			Parameters: &cdn.RequestSchemeMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        utils.String(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(*condition.Parameters.Operator, condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontDoorUrlPathCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.UrlPath

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := cdn.DeliveryRuleURLPathCondition{
			Name: conditionMapping.Name,
			Parameters: &cdn.URLPathMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        cdn.URLPathOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
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

func ExpandCdnFrontDoorUrlFileExtensionCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.UrlFileExtension

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := cdn.DeliveryRuleURLFileExtensionCondition{
			Name: conditionMapping.Name,
			Parameters: &cdn.URLFileExtensionMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        cdn.URLFileExtensionOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
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

func ExpandCdnFrontDoorUrlFileNameCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.UrlFilename

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := cdn.DeliveryRuleURLFileNameCondition{
			Name: conditionMapping.Name,
			Parameters: &cdn.URLFileNameMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        cdn.URLFileNameOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
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

func ExpandCdnFrontDoorHttpVersionCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.HttpVersion

	for _, v := range input {
		item := v.(map[string]interface{})
		matchValuesRaw := item["match_values"].(*pluginsdk.Set).List()

		condition := cdn.DeliveryRuleHTTPVersionCondition{
			Name: conditionMapping.Name,
			Parameters: &cdn.HTTPVersionMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        utils.String(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(matchValuesRaw),
			},
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(*condition.Parameters.Operator, condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontDoorCookiesCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.Cookies

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := cdn.DeliveryRuleCookiesCondition{
			Name: conditionMapping.Name,
			Parameters: &cdn.CookiesMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Selector:        utils.String(item["cookie_name"].(string)),
				Operator:        cdn.CookiesOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
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

func ExpandCdnFrontDoorIsDeviceCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.IsDevice

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := cdn.DeliveryRuleIsDeviceCondition{
			Name: conditionMapping.Name,
			Parameters: &cdn.IsDeviceMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        utils.String(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(*condition.Parameters.Operator, condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontDoorSocketAddressCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.SocketAddress

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := cdn.DeliveryRuleSocketAddrCondition{
			Name: conditionMapping.Name,
			Parameters: &cdn.SocketAddrMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        cdn.SocketAddrOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if condition.Parameters.Operator == cdn.SocketAddrOperatorIPMatch {
			// make sure all of the passed CIDRs are valid
			for _, matchValue := range item["match_values"].([]interface{}) {
				if _, err := cdnValidate.FrontDoorRuleCidrIsValid(matchValue, "match_values"); err != nil {
					return nil, fmt.Errorf("%q is invalid: when the 'operator' is set to 'IPMatch' the 'match_values' must be a valid IPv4 or IPv6 CIDR, got %q", conditionMapping.ConfigName, matchValue.(string))
				}
			}

			// Check for CIDR overlap and CIDR duplicates in the match values
			_, err := cdnValidate.FrontDoorRuleCidrOverlap(item["match_values"].([]interface{}), "match_values")
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

func ExpandCdnFrontDoorClientPortCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.ClientPort

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := cdn.DeliveryRuleClientPortCondition{
			Name: conditionMapping.Name,
			Parameters: &cdn.ClientPortMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        cdn.ClientPortOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
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

func ExpandCdnFrontDoorServerPortCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.ServerPort

	for _, v := range input {
		item := v.(map[string]interface{})
		matchValuesRaw := item["match_values"].(*pluginsdk.Set).List()

		condition := cdn.DeliveryRuleServerPortCondition{
			Name: conditionMapping.Name,
			Parameters: &cdn.ServerPortMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        cdn.ServerPortOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
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

func ExpandCdnFrontDoorHostNameCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.HostName

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := cdn.DeliveryRuleHostNameCondition{
			Name: conditionMapping.Name,
			Parameters: &cdn.HostNameMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        cdn.HostNameOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
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

func ExpandCdnFrontDoorSslProtocolCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontDoorConditionMappings()
	conditionMapping := m.SslProtocol

	for _, v := range input {
		item := v.(map[string]interface{})

		matchValues := make([]cdn.SslProtocol, 0)
		validationMatchValues := make([]string, 0)
		matchValuesRaw := item["match_values"].(*pluginsdk.Set).List()

		for _, value := range matchValuesRaw {
			matchValues = append(matchValues, cdn.SslProtocol(value.(string)))
			validationMatchValues = append(validationMatchValues, value.(string))
		}

		condition := cdn.DeliveryRuleSslProtocolCondition{
			Name: conditionMapping.Name,
			Parameters: &cdn.SslProtocolMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        utils.String(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     &matchValues,
			},
		}

		if err := validateCdnFrontDoorExpandConditionOperatorValues(*condition.Parameters.Operator, &validationMatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func FlattenFrontdoorRemoteAddressCondition(input cdn.BasicDeliveryRuleCondition) (map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleRemoteAddressCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule remote address condition")
	}

	normalized := createCdnFrontDoorNormalizedConditionStub()

	if params := condition.Parameters; params != nil {
		normalized = normalizedCondition{
			selector:        nil,
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      nil,
		}
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorRequestMethodCondition(input cdn.BasicDeliveryRuleCondition) (map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleRequestMethodCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule request method condition")
	}

	normalized := createCdnFrontDoorNormalizedConditionStub()

	if params := condition.Parameters; params != nil {
		normalized = normalizedCondition{
			selector:        nil,
			operator:        *params.Operator,
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      nil,
		}
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorQueryStringCondition(input cdn.BasicDeliveryRuleCondition) (map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleQueryStringCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule query string condition")
	}

	normalized := createCdnFrontDoorNormalizedConditionStub()

	if params := condition.Parameters; params != nil {
		normalized = normalizedCondition{
			selector:        nil,
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorPostArgsCondition(input cdn.BasicDeliveryRuleCondition) (map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRulePostArgsCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule post args condition")
	}

	normalized := createCdnFrontDoorNormalizedConditionStub()

	if params := condition.Parameters; params != nil {
		normalized = normalizedCondition{
			selector:        &normalizedSelector{name: utils.String("post_args_name"), value: params.Selector},
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorRequestUriCondition(input cdn.BasicDeliveryRuleCondition) (map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleRequestURICondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule request URI condition")
	}

	normalized := createCdnFrontDoorNormalizedConditionStub()

	if params := condition.Parameters; params != nil {
		normalized = normalizedCondition{
			selector:        nil,
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorRequestHeaderCondition(input cdn.BasicDeliveryRuleCondition) (map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleRequestHeaderCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule request header condition")
	}

	normalized := createCdnFrontDoorNormalizedConditionStub()

	if params := condition.Parameters; params != nil {
		normalized = normalizedCondition{
			selector:        &normalizedSelector{name: utils.String("header_name"), value: params.Selector},
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorRequestBodyCondition(input cdn.BasicDeliveryRuleCondition) (map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleRequestBodyCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule request body condition")
	}

	normalized := createCdnFrontDoorNormalizedConditionStub()

	if params := condition.Parameters; params != nil {
		normalized = normalizedCondition{
			selector:        nil,
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorRequestSchemeCondition(input cdn.BasicDeliveryRuleCondition) (map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleRequestSchemeCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule request scheme condition")
	}

	normalized := createCdnFrontDoorNormalizedConditionStub()

	if params := condition.Parameters; params != nil {
		normalized = normalizedCondition{
			selector:        nil,
			operator:        *params.Operator,
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      nil,
		}
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorUrlPathCondition(input cdn.BasicDeliveryRuleCondition) (map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleURLPathCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule url path condition")
	}

	normalized := createCdnFrontDoorNormalizedConditionStub()

	if params := condition.Parameters; params != nil {
		normalized = normalizedCondition{
			selector:        nil,
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorUrlFileExtensionCondition(input cdn.BasicDeliveryRuleCondition) (map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleURLFileExtensionCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule url file extension condition")
	}

	normalized := createCdnFrontDoorNormalizedConditionStub()

	if params := condition.Parameters; params != nil {
		normalized = normalizedCondition{
			selector:        nil,
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorUrlFileNameCondition(input cdn.BasicDeliveryRuleCondition) (map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleURLFileNameCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule url file name condition")
	}

	normalized := createCdnFrontDoorNormalizedConditionStub()

	if params := condition.Parameters; params != nil {
		normalized = normalizedCondition{
			selector:        nil,
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorHttpVersionCondition(input cdn.BasicDeliveryRuleCondition) (map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleHTTPVersionCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule http version condition")
	}

	normalized := createCdnFrontDoorNormalizedConditionStub()

	if params := condition.Parameters; params != nil {
		normalized = normalizedCondition{
			selector:        nil,
			operator:        *params.Operator,
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      nil,
		}
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorCookiesCondition(input cdn.BasicDeliveryRuleCondition) (map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleCookiesCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule cookie condition")
	}

	// create a normalized condition stub first with empty values so we can push that to the state file
	normalized := createCdnFrontDoorNormalizedConditionStub()

	// if this has values override the stub values for the actual values
	if params := condition.Parameters; params != nil {
		normalized = normalizedCondition{
			selector:        &normalizedSelector{name: utils.String("cookie_name"), value: params.Selector},
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}
	}

	// flatten the normalized condition regardless if it is a stub or not...
	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorIsDeviceCondition(input cdn.BasicDeliveryRuleCondition) (map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleIsDeviceCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule is device condition")
	}

	normalized := createCdnFrontDoorNormalizedConditionStub()

	if params := condition.Parameters; params != nil {
		normalized = normalizedCondition{
			selector:        nil,
			operator:        *params.Operator,
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      nil,
		}
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorSocketAddressCondition(input cdn.BasicDeliveryRuleCondition) (map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleSocketAddrCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule socket address condition")
	}

	normalized := createCdnFrontDoorNormalizedConditionStub()

	if params := condition.Parameters; params != nil {
		normalized = normalizedCondition{
			selector:        nil,
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      nil,
		}
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorClientPortCondition(input cdn.BasicDeliveryRuleCondition) (map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleClientPortCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule client port condition")
	}

	normalized := createCdnFrontDoorNormalizedConditionStub()

	if params := condition.Parameters; params != nil {
		normalized = normalizedCondition{
			selector:        nil,
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      nil,
		}
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorServerPortCondition(input cdn.BasicDeliveryRuleCondition) (map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleServerPortCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule server port condition")
	}

	normalized := createCdnFrontDoorNormalizedConditionStub()

	if params := condition.Parameters; params != nil {
		normalized = normalizedCondition{
			selector:        nil,
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      nil,
		}
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorHostNameCondition(input cdn.BasicDeliveryRuleCondition) (map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleHostNameCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule host name condition")
	}

	normalized := createCdnFrontDoorNormalizedConditionStub()

	if params := condition.Parameters; params != nil {
		normalized = normalizedCondition{
			selector:        nil,
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func FlattenFrontdoorSslProtocolCondition(input cdn.BasicDeliveryRuleCondition) (map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleSslProtocolCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule ssl protocol condition")
	}

	normalized := createCdnFrontDoorNormalizedConditionStub()

	if params := condition.Parameters; params != nil {
		matchValues := make([]string, 0)
		for _, value := range *params.MatchValues {
			matchValues = append(matchValues, string(value))
		}

		normalized = normalizedCondition{
			selector:        nil,
			operator:        *params.Operator,
			negateCondition: params.NegateCondition,
			matchValues:     &matchValues,
			transforms:      nil,
		}
	}

	return flattenCdnFrontDoorNormalizedCondition(normalized), nil
}

func createCdnFrontDoorNormalizedConditionStub() normalizedCondition {
	matchValues := make([]string, 0)

	stub := normalizedCondition{
		selector:        nil,
		operator:        "",
		negateCondition: utils.Bool(false),
		matchValues:     &matchValues,
		transforms:      nil,
	}

	return stub
}
