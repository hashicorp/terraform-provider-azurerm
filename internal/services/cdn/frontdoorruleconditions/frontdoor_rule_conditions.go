package frontdoorruleconditions

import (
	"fmt"

	track1 "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type FrontdoorConditionParameters struct {
	Name       track1.Name
	TypeName   string
	ConfigName string
}

type FrontdoorCondtionsMappings struct {
	ClientPort       FrontdoorConditionParameters
	Cookies          FrontdoorConditionParameters
	HostName         FrontdoorConditionParameters
	HttpVersion      FrontdoorConditionParameters
	IsDevice         FrontdoorConditionParameters
	PostArgs         FrontdoorConditionParameters
	QueryString      FrontdoorConditionParameters
	RemoteAddress    FrontdoorConditionParameters
	RequestBody      FrontdoorConditionParameters
	RequestHeader    FrontdoorConditionParameters
	RequestMethod    FrontdoorConditionParameters
	RequestScheme    FrontdoorConditionParameters
	RequestUri       FrontdoorConditionParameters
	ServerPort       FrontdoorConditionParameters
	SocketAddress    FrontdoorConditionParameters
	SslProtocol      FrontdoorConditionParameters
	UrlFileExtension FrontdoorConditionParameters
	UrlFilename      FrontdoorConditionParameters
	UrlPath          FrontdoorConditionParameters
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
	transforms      *[]track1.Transform
}

func InitializeFrontdoorConditionMappings() *FrontdoorCondtionsMappings {
	m := new(FrontdoorCondtionsMappings)

	m.ClientPort = FrontdoorConditionParameters{
		Name:       track1.NameClientPort,
		TypeName:   "DeliveryRuleClientPortConditionParameters",
		ConfigName: "client_port_condition",
	}

	m.Cookies = FrontdoorConditionParameters{
		Name:       track1.NameCookies,
		TypeName:   "DeliveryRuleCookiesConditionParameters",
		ConfigName: "cookies_condition",
	}

	m.HostName = FrontdoorConditionParameters{
		Name:       track1.NameHostName,
		TypeName:   "DeliveryRuleHostNameConditionParameters",
		ConfigName: "host_name_condition",
	}

	m.HttpVersion = FrontdoorConditionParameters{
		Name:       track1.NameHTTPVersion,
		TypeName:   "DeliveryRuleHttpVersionConditionParameters",
		ConfigName: "http_version_condition",
	}

	m.IsDevice = FrontdoorConditionParameters{
		Name:       track1.NameIsDevice,
		TypeName:   "DeliveryRuleIsDeviceConditionParameters",
		ConfigName: "is_device_condition",
	}

	m.PostArgs = FrontdoorConditionParameters{
		Name:       track1.NamePostArgs,
		TypeName:   "DeliveryRulePostArgsConditionParameters",
		ConfigName: "postargs_condition",
	}

	m.QueryString = FrontdoorConditionParameters{
		Name:       track1.NameQueryString,
		TypeName:   "DeliveryRuleQueryStringConditionParameters",
		ConfigName: "query_string_condition",
	}

	m.RemoteAddress = FrontdoorConditionParameters{
		Name:       track1.NameRemoteAddress,
		TypeName:   "DeliveryRuleRemoteAddressConditionParameters",
		ConfigName: "remote_address_condition",
	}

	m.RequestBody = FrontdoorConditionParameters{
		Name:       track1.NameRequestBody,
		TypeName:   "DeliveryRuleRequestBodyConditionParameters",
		ConfigName: "request_body_condition",
	}

	m.RequestHeader = FrontdoorConditionParameters{
		Name:       track1.NameRequestHeader,
		TypeName:   "DeliveryRuleRequestHeaderConditionParameters",
		ConfigName: "request_header_condition",
	}

	m.RequestMethod = FrontdoorConditionParameters{
		Name:       track1.NameRequestMethod,
		TypeName:   "DeliveryRuleRequestMethodConditionParameters",
		ConfigName: "request_method_condition",
	}

	m.RequestScheme = FrontdoorConditionParameters{
		Name:       track1.NameRequestScheme,
		TypeName:   "DeliveryRuleRequestSchemeConditionParameters",
		ConfigName: "request_scheme_condition",
	}

	m.RequestUri = FrontdoorConditionParameters{
		Name:       track1.NameRequestURI,
		TypeName:   "DeliveryRuleRequestUriConditionParameters",
		ConfigName: "request_uri_condition",
	}

	m.ServerPort = FrontdoorConditionParameters{
		Name:       track1.NameServerPort,
		TypeName:   "DeliveryRuleServerPortConditionParameters",
		ConfigName: "server_port_condition",
	}

	m.SocketAddress = FrontdoorConditionParameters{
		Name:       track1.NameSocketAddr,
		TypeName:   "DeliveryRuleSocketAddrConditionParameters",
		ConfigName: "socket_address_condition",
	}

	m.SslProtocol = FrontdoorConditionParameters{
		Name:       track1.NameSslProtocol,
		TypeName:   "DeliveryRuleSslProtocolConditionParameters",
		ConfigName: "ssl_protocol_condition",
	}

	m.UrlFileExtension = FrontdoorConditionParameters{
		Name:       track1.NameURLFileExtension,
		TypeName:   "DeliveryRuleURLFileExtensionMatchConditionParameters",
		ConfigName: "url_file_extension_condition",
	}

	m.UrlFilename = FrontdoorConditionParameters{
		Name:       track1.NameURLFileName,
		TypeName:   "DeliveryRuleUrlFilenameConditionParameters",
		ConfigName: "url_filename_condition",
	}

	m.UrlPath = FrontdoorConditionParameters{
		Name:       track1.NameURLPath,
		TypeName:   "DeliveryRuleUrlPathMatchConditionParameters",
		ConfigName: "url_path_condition",
	}

	return m
}

func expandNormalizeTransforms(input []interface{}) []track1.Transform {
	transforms := make([]track1.Transform, 0)
	if len(input) == 0 {
		return transforms
	}

	for _, t := range input {
		transforms = append(transforms, track1.Transform(t.(string)))
	}

	return transforms
}

func flattenAndValidateNormalizedCondition(condition normalizedCondition, m FrontdoorConditionParameters) (map[string]interface{}, error) {
	operator := ""
	negateCondition := false
	matchValues := make([]interface{}, 0)
	conditionTransforms := make([]string, 0)

	if condition.operator != "" {
		operator = condition.operator
	}

	if condition.negateCondition != nil {
		negateCondition = *condition.negateCondition
	}

	if condition.matchValues != nil {
		matchValues = utils.FlattenStringSlice(condition.matchValues)
	}

	v := map[string]interface{}{
		"operator":         operator,
		"negate_condition": negateCondition,
		"match_values":     matchValues,
		"transforms":       conditionTransforms,
	}

	if condition.selector != nil {
		v[*condition.selector.name] = *condition.selector.value
	}

	if condition.transforms != nil {
		for _, transform := range *condition.transforms {
			conditionTransforms = append(conditionTransforms, string(transform))
		}

		v["transforms"] = conditionTransforms
	}

	return v, nil
}

func validateFrontdoorExpandConditionOperatorValues(operator string, matchValues *[]string, m FrontdoorConditionParameters) error {
	if operator == "" {
		return fmt.Errorf("%q is invalid: no %q value has been set, got %q", m.ConfigName, "operator", operator)
	}

	if operator == string(track1.OperatorAny) && len(*matchValues) > 0 {
		return fmt.Errorf("%q is invalid: the %q field must not be set if the conditions %q is set to %q", m.ConfigName, "match_values", "operator", track1.OperatorAny)
	}

	return nil
}

func ExpandFrontdoorRemoteAddressCondition(input []interface{}) (*[]track1.BasicDeliveryRuleCondition, error) {
	output := make([]track1.BasicDeliveryRuleCondition, 0)
	m := InitializeFrontdoorConditionMappings()
	conditionMapping := m.RemoteAddress

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := track1.DeliveryRuleRemoteAddressCondition{
			Name: conditionMapping.Name,
			Parameters: &track1.RemoteAddressMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        track1.RemoteAddressOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if err := validateFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandFrontdoorRequestMethodCondition(input []interface{}) (*[]track1.BasicDeliveryRuleCondition, error) {
	output := make([]track1.BasicDeliveryRuleCondition, 0)
	m := InitializeFrontdoorConditionMappings()
	conditionMapping := m.RequestMethod

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := track1.DeliveryRuleRequestMethodCondition{
			Name: conditionMapping.Name,
			Parameters: &track1.RequestMethodMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        utils.String(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if err := validateFrontdoorExpandConditionOperatorValues(string(*condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandFrontdoorQueryStringCondition(input []interface{}) (*[]track1.BasicDeliveryRuleCondition, error) {
	output := make([]track1.BasicDeliveryRuleCondition, 0)
	m := InitializeFrontdoorConditionMappings()
	conditionMapping := m.QueryString

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := track1.DeliveryRuleQueryStringCondition{
			Name: conditionMapping.Name,
			Parameters: &track1.QueryStringMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        track1.QueryStringOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if transforms := item["transforms"].([]interface{}); len(transforms) != 0 {
			expanded := expandNormalizeTransforms(transforms)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandFrontdoorPostArgsCondition(input []interface{}) (*[]track1.BasicDeliveryRuleCondition, error) {
	output := make([]track1.BasicDeliveryRuleCondition, 0)
	m := InitializeFrontdoorConditionMappings()
	conditionMapping := m.PostArgs

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := track1.DeliveryRulePostArgsCondition{
			Name: conditionMapping.Name,
			Parameters: &track1.PostArgsMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Selector:        utils.String(item["postargs_name"].(string)),
				Operator:        track1.PostArgsOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if transforms := item["transforms"].([]interface{}); len(transforms) != 0 {
			expanded := expandNormalizeTransforms(transforms)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandFrontdoorRequestUriCondition(input []interface{}) (*[]track1.BasicDeliveryRuleCondition, error) {
	output := make([]track1.BasicDeliveryRuleCondition, 0)
	m := InitializeFrontdoorConditionMappings()
	conditionMapping := m.RequestUri

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := track1.DeliveryRuleRequestURICondition{
			Name: conditionMapping.Name,
			Parameters: &track1.RequestURIMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        track1.RequestURIOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if transforms := item["transforms"].([]interface{}); len(transforms) != 0 {
			expanded := expandNormalizeTransforms(transforms)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandFrontdoorRequestHeaderCondition(input []interface{}) (*[]track1.BasicDeliveryRuleCondition, error) {
	output := make([]track1.BasicDeliveryRuleCondition, 0)
	m := InitializeFrontdoorConditionMappings()
	conditionMapping := m.RequestHeader

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := track1.DeliveryRuleRequestHeaderCondition{
			Name: conditionMapping.Name,
			Parameters: &track1.RequestHeaderMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Selector:        utils.String(item["header_name"].(string)),
				Operator:        track1.RequestHeaderOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if transforms := item["transforms"].([]interface{}); len(transforms) != 0 {
			expanded := expandNormalizeTransforms(transforms)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandFrontdoorRequestBodyCondition(input []interface{}) (*[]track1.BasicDeliveryRuleCondition, error) {
	output := make([]track1.BasicDeliveryRuleCondition, 0)
	m := InitializeFrontdoorConditionMappings()
	conditionMapping := m.RequestBody

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := track1.DeliveryRuleRequestBodyCondition{
			Name: conditionMapping.Name,
			Parameters: &track1.RequestBodyMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        track1.RequestBodyOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if transforms := item["transforms"].([]interface{}); len(transforms) != 0 {
			expanded := expandNormalizeTransforms(transforms)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandFrontdoorRequestSchemeCondition(input []interface{}) (*[]track1.BasicDeliveryRuleCondition, error) {
	output := make([]track1.BasicDeliveryRuleCondition, 0)
	m := InitializeFrontdoorConditionMappings()
	conditionMapping := m.RequestScheme

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := track1.DeliveryRuleRequestSchemeCondition{
			Name: conditionMapping.Name,
			Parameters: &track1.RequestSchemeMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        utils.String(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if err := validateFrontdoorExpandConditionOperatorValues(*condition.Parameters.Operator, condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandFrontdoorUrlPathCondition(input []interface{}) (*[]track1.BasicDeliveryRuleCondition, error) {
	output := make([]track1.BasicDeliveryRuleCondition, 0)
	m := InitializeFrontdoorConditionMappings()
	conditionMapping := m.UrlPath

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := track1.DeliveryRuleURLPathCondition{
			Name: conditionMapping.Name,
			Parameters: &track1.URLPathMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        track1.URLPathOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if transforms := item["transforms"].([]interface{}); len(transforms) != 0 {
			expanded := expandNormalizeTransforms(transforms)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandFrontdoorUrlFileExtensionCondition(input []interface{}) (*[]track1.BasicDeliveryRuleCondition, error) {
	output := make([]track1.BasicDeliveryRuleCondition, 0)
	m := InitializeFrontdoorConditionMappings()
	conditionMapping := m.UrlFileExtension

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := track1.DeliveryRuleURLFileExtensionCondition{
			Name: conditionMapping.Name,
			Parameters: &track1.URLFileExtensionMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        track1.URLFileExtensionOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if transforms := item["transforms"].([]interface{}); len(transforms) != 0 {
			expanded := expandNormalizeTransforms(transforms)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandFrontdoorUrlFileNameCondition(input []interface{}) (*[]track1.BasicDeliveryRuleCondition, error) {
	output := make([]track1.BasicDeliveryRuleCondition, 0)
	m := InitializeFrontdoorConditionMappings()
	conditionMapping := m.UrlFilename

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := track1.DeliveryRuleURLFileNameCondition{
			Name: conditionMapping.Name,
			Parameters: &track1.URLFileNameMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        track1.URLFileNameOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if transforms := item["transforms"].([]interface{}); len(transforms) != 0 {
			expanded := expandNormalizeTransforms(transforms)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandFrontdoorHttpVersionCondition(input []interface{}) (*[]track1.BasicDeliveryRuleCondition, error) {
	output := make([]track1.BasicDeliveryRuleCondition, 0)
	m := InitializeFrontdoorConditionMappings()
	conditionMapping := m.HttpVersion

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := track1.DeliveryRuleHTTPVersionCondition{
			Name: conditionMapping.Name,
			Parameters: &track1.HTTPVersionMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        utils.String(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if err := validateFrontdoorExpandConditionOperatorValues(*condition.Parameters.Operator, condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandFrontdoorCookiesCondition(input []interface{}) (*[]track1.BasicDeliveryRuleCondition, error) {
	output := make([]track1.BasicDeliveryRuleCondition, 0)
	m := InitializeFrontdoorConditionMappings()
	conditionMapping := m.Cookies

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := track1.DeliveryRuleCookiesCondition{
			Name: conditionMapping.Name,
			Parameters: &track1.CookiesMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Selector:        utils.String(item["cookie_name"].(string)),
				Operator:        track1.CookiesOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if tt := item["transforms"].([]interface{}); len(tt) != 0 {
			transforms := make([]track1.Transform, 0)
			for _, t := range tt {
				transforms = append(transforms, track1.Transform(t.(string)))
			}
			condition.Parameters.Transforms = &transforms
		}

		if err := validateFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandFrontdoorIsDeviceCondition(input []interface{}) (*[]track1.BasicDeliveryRuleCondition, error) {
	output := make([]track1.BasicDeliveryRuleCondition, 0)
	m := InitializeFrontdoorConditionMappings()
	conditionMapping := m.IsDevice

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := track1.DeliveryRuleIsDeviceCondition{
			Name: conditionMapping.Name,
			Parameters: &track1.IsDeviceMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        utils.String(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if err := validateFrontdoorExpandConditionOperatorValues(*condition.Parameters.Operator, condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandFrontdoorSocketAddressCondition(input []interface{}) (*[]track1.BasicDeliveryRuleCondition, error) {
	output := make([]track1.BasicDeliveryRuleCondition, 0)
	m := InitializeFrontdoorConditionMappings()
	conditionMapping := m.SocketAddress

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := track1.DeliveryRuleSocketAddrCondition{
			Name: conditionMapping.Name,
			Parameters: &track1.SocketAddrMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        track1.SocketAddrOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if transforms := item["transforms"].([]interface{}); len(transforms) != 0 {
			expanded := expandNormalizeTransforms(transforms)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandFrontdoorClientPortCondition(input []interface{}) (*[]track1.BasicDeliveryRuleCondition, error) {
	output := make([]track1.BasicDeliveryRuleCondition, 0)
	m := InitializeFrontdoorConditionMappings()
	conditionMapping := m.ClientPort

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := track1.DeliveryRuleClientPortCondition{
			Name: conditionMapping.Name,
			Parameters: &track1.ClientPortMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        track1.ClientPortOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if transforms := item["transforms"].([]interface{}); len(transforms) != 0 {
			expanded := expandNormalizeTransforms(transforms)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandFrontdoorServerPortCondition(input []interface{}) (*[]track1.BasicDeliveryRuleCondition, error) {
	output := make([]track1.BasicDeliveryRuleCondition, 0)
	m := InitializeFrontdoorConditionMappings()
	conditionMapping := m.ServerPort

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := track1.DeliveryRuleServerPortCondition{
			Name: conditionMapping.Name,
			Parameters: &track1.ServerPortMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        track1.ServerPortOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if transforms := item["transforms"].([]interface{}); len(transforms) != 0 {
			expanded := expandNormalizeTransforms(transforms)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandFrontdoorHostNameCondition(input []interface{}) (*[]track1.BasicDeliveryRuleCondition, error) {
	output := make([]track1.BasicDeliveryRuleCondition, 0)
	m := InitializeFrontdoorConditionMappings()
	conditionMapping := m.HostName

	for _, v := range input {
		item := v.(map[string]interface{})
		condition := track1.DeliveryRuleHostNameCondition{
			Name: conditionMapping.Name,
			Parameters: &track1.HostNameMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        track1.HostNameOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if transforms := item["transforms"].([]interface{}); len(transforms) != 0 {
			expanded := expandNormalizeTransforms(transforms)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandFrontdoorSslProtocolCondition(input []interface{}) (*[]track1.BasicDeliveryRuleCondition, error) {
	output := make([]track1.BasicDeliveryRuleCondition, 0)
	m := InitializeFrontdoorConditionMappings()
	conditionMapping := m.SslProtocol

	for _, v := range input {
		item := v.(map[string]interface{})

		matchValues := make([]track1.SslProtocol, 0)
		validationMatchValues := make([]string, 0)
		for _, value := range item["match_values"].([]interface{}) {
			matchValues = append(matchValues, track1.SslProtocol(value.(string)))
			validationMatchValues = append(validationMatchValues, value.(string))
		}

		condition := track1.DeliveryRuleSslProtocolCondition{
			Name: conditionMapping.Name,
			Parameters: &track1.SslProtocolMatchConditionParameters{
				TypeName:        utils.String(conditionMapping.TypeName),
				Operator:        utils.String(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     &matchValues,
			},
		}

		if transforms := item["transforms"].([]interface{}); len(transforms) != 0 {
			expanded := expandNormalizeTransforms(transforms)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateFrontdoorExpandConditionOperatorValues(*condition.Parameters.Operator, &validationMatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func FlattenFrontdoorRemoteAddressCondition(input track1.BasicDeliveryRuleCondition, m FrontdoorConditionParameters) ([]interface{}, error) {
	condition, ok := input.AsDeliveryRuleRemoteAddressCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule remote address condition")
	}

	if params := condition.Parameters; params != nil {
		condition := normalizedCondition{
			selector:        nil,
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}

		if flattened, err := flattenAndValidateNormalizedCondition(condition, m); err != nil {
			return []interface{}{flattened}, err
		} else {
			return []interface{}{flattened}, nil
		}
	}

	return nil, nil
}

func FlattenFrontdoorRequestMethodCondition(input track1.BasicDeliveryRuleCondition, m FrontdoorConditionParameters) ([]interface{}, error) {
	condition, ok := input.AsDeliveryRuleRequestMethodCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule request method condition")
	}

	if params := condition.Parameters; params != nil {
		condition := normalizedCondition{
			selector:        nil,
			operator:        string(*params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}

		if flattened, err := flattenAndValidateNormalizedCondition(condition, m); err != nil {
			return []interface{}{flattened}, err
		} else {
			return []interface{}{flattened}, nil
		}
	}

	return nil, nil
}

func FlattenFrontdoorQueryStringCondition(input track1.BasicDeliveryRuleCondition, m FrontdoorConditionParameters) ([]interface{}, error) {
	condition, ok := input.AsDeliveryRuleQueryStringCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule query string condition")
	}

	if params := condition.Parameters; params != nil {
		condition := normalizedCondition{
			selector:        nil,
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}

		if flattened, err := flattenAndValidateNormalizedCondition(condition, m); err != nil {
			return []interface{}{flattened}, err
		} else {
			return []interface{}{flattened}, nil
		}
	}

	return nil, nil
}

func FlattenFrontdoorPostArgsCondition(input track1.BasicDeliveryRuleCondition, m FrontdoorConditionParameters) ([]interface{}, error) {
	condition, ok := input.AsDeliveryRulePostArgsCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule post args condition")
	}

	if params := condition.Parameters; params != nil {
		condition := normalizedCondition{
			selector:        &normalizedSelector{name: utils.String("postargs_name"), value: params.Selector},
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}

		if flattened, err := flattenAndValidateNormalizedCondition(condition, m); err != nil {
			return []interface{}{flattened}, err
		} else {
			return []interface{}{flattened}, nil
		}
	}

	return nil, nil
}

func FlattenFrontdoorRequestUriCondition(input track1.BasicDeliveryRuleCondition, m FrontdoorConditionParameters) ([]interface{}, error) {
	condition, ok := input.AsDeliveryRuleRequestURICondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule request URI condition")
	}

	if params := condition.Parameters; params != nil {
		condition := normalizedCondition{
			selector:        nil,
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}

		if flattened, err := flattenAndValidateNormalizedCondition(condition, m); err != nil {
			return []interface{}{flattened}, err
		} else {
			return []interface{}{flattened}, nil
		}
	}

	return nil, nil
}

func FlattenFrontdoorRequestHeaderCondition(input track1.BasicDeliveryRuleCondition, m FrontdoorConditionParameters) ([]interface{}, error) {
	condition, ok := input.AsDeliveryRuleRequestHeaderCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule request header condition")
	}

	if params := condition.Parameters; params != nil {
		condition := normalizedCondition{
			selector:        &normalizedSelector{name: utils.String("header_name"), value: params.Selector},
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}

		if flattened, err := flattenAndValidateNormalizedCondition(condition, m); err != nil {
			return []interface{}{flattened}, err
		} else {
			return []interface{}{flattened}, nil
		}
	}

	return nil, nil
}

func FlattenFrontdoorRequestBodyCondition(input track1.BasicDeliveryRuleCondition, m FrontdoorConditionParameters) ([]interface{}, error) {
	condition, ok := input.AsDeliveryRuleRequestBodyCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule request body condition")
	}

	if params := condition.Parameters; params != nil {
		condition := normalizedCondition{
			selector:        nil,
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}

		if flattened, err := flattenAndValidateNormalizedCondition(condition, m); err != nil {
			return []interface{}{flattened}, err
		} else {
			return []interface{}{flattened}, nil
		}
	}

	return nil, nil
}

func FlattenFrontdoorRequestSchemeCondition(input track1.BasicDeliveryRuleCondition, m FrontdoorConditionParameters) ([]interface{}, error) {
	condition, ok := input.AsDeliveryRuleRequestSchemeCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule request scheme condition")
	}

	if params := condition.Parameters; params != nil {
		condition := normalizedCondition{
			selector:        nil,
			operator:        *params.Operator,
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      nil,
		}

		if flattened, err := flattenAndValidateNormalizedCondition(condition, m); err != nil {
			return []interface{}{flattened}, err
		} else {
			return []interface{}{flattened}, nil
		}
	}

	return nil, nil
}

func FlattenFrontdoorUrlPathCondition(input track1.BasicDeliveryRuleCondition, m FrontdoorConditionParameters) ([]interface{}, error) {
	condition, ok := input.AsDeliveryRuleURLPathCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule url path condition")
	}

	if params := condition.Parameters; params != nil {
		condition := normalizedCondition{
			selector:        nil,
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}

		if flattened, err := flattenAndValidateNormalizedCondition(condition, m); err != nil {
			return []interface{}{flattened}, err
		} else {
			return []interface{}{flattened}, nil
		}
	}

	return nil, nil
}

func FlattenFrontdoorUrlFileExtensionCondition(input track1.BasicDeliveryRuleCondition, m FrontdoorConditionParameters) ([]interface{}, error) {
	condition, ok := input.AsDeliveryRuleURLFileExtensionCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule url file extension condition")
	}

	if params := condition.Parameters; params != nil {
		condition := normalizedCondition{
			selector:        nil,
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}

		if flattened, err := flattenAndValidateNormalizedCondition(condition, m); err != nil {
			return []interface{}{flattened}, err
		} else {
			return []interface{}{flattened}, nil
		}
	}

	return nil, nil
}

func FlattenFrontdoorUrlFileNameCondition(input track1.BasicDeliveryRuleCondition, m FrontdoorConditionParameters) ([]interface{}, error) {
	condition, ok := input.AsDeliveryRuleURLFileNameCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule url file name condition")
	}

	if params := condition.Parameters; params != nil {
		condition := normalizedCondition{
			selector:        nil,
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}

		if flattened, err := flattenAndValidateNormalizedCondition(condition, m); err != nil {
			return []interface{}{flattened}, err
		} else {
			return []interface{}{flattened}, nil
		}
	}

	return nil, nil
}

func FlattenFrontdoorHttpVersionCondition(input track1.BasicDeliveryRuleCondition, m FrontdoorConditionParameters) ([]interface{}, error) {
	condition, ok := input.AsDeliveryRuleHTTPVersionCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule http version condition")
	}

	if params := condition.Parameters; params != nil {
		condition := normalizedCondition{
			selector:        nil,
			operator:        *params.Operator,
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      nil,
		}

		if flattened, err := flattenAndValidateNormalizedCondition(condition, m); err != nil {
			return []interface{}{flattened}, err
		} else {
			return []interface{}{flattened}, nil
		}
	}

	return nil, nil
}

func FlattenFrontdoorCookiesCondition(input track1.BasicDeliveryRuleCondition, m FrontdoorConditionParameters) ([]interface{}, error) {
	condition, ok := input.AsDeliveryRuleCookiesCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule cookie condition")
	}

	if params := condition.Parameters; params != nil {
		condition := normalizedCondition{
			selector:        &normalizedSelector{name: utils.String("cookie_name"), value: params.Selector},
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}

		if flattened, err := flattenAndValidateNormalizedCondition(condition, m); err != nil {
			return []interface{}{flattened}, err
		} else {
			return []interface{}{flattened}, nil
		}
	}

	return nil, nil
}

func FlattenFrontdoorIsDeviceCondition(input track1.BasicDeliveryRuleCondition, m FrontdoorConditionParameters) ([]interface{}, error) {
	condition, ok := input.AsDeliveryRuleIsDeviceCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule is device condition")
	}

	if params := condition.Parameters; params != nil {
		condition := normalizedCondition{
			selector:        nil,
			operator:        *params.Operator,
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      nil,
		}

		if flattened, err := flattenAndValidateNormalizedCondition(condition, m); err != nil {
			return []interface{}{flattened}, err
		} else {
			return []interface{}{flattened}, nil
		}
	}

	return nil, nil
}

func FlattenFrontdoorSocketAddressCondition(input track1.BasicDeliveryRuleCondition, m FrontdoorConditionParameters) ([]interface{}, error) {
	condition, ok := input.AsDeliveryRuleSocketAddrCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule socket address condition")
	}

	if params := condition.Parameters; params != nil {
		condition := normalizedCondition{
			selector:        nil,
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}

		if flattened, err := flattenAndValidateNormalizedCondition(condition, m); err != nil {
			return []interface{}{flattened}, err
		} else {
			return []interface{}{flattened}, nil
		}
	}

	return nil, nil
}

func FlattenFrontdoorClientPortCondition(input track1.BasicDeliveryRuleCondition, m FrontdoorConditionParameters) ([]interface{}, error) {
	condition, ok := input.AsDeliveryRuleClientPortCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule client port condition")
	}

	if params := condition.Parameters; params != nil {
		condition := normalizedCondition{
			selector:        nil,
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}

		if flattened, err := flattenAndValidateNormalizedCondition(condition, m); err != nil {
			return []interface{}{flattened}, err
		} else {
			return []interface{}{flattened}, nil
		}
	}

	return nil, nil
}

func FlattenFrontdoorServerPortCondition(input track1.BasicDeliveryRuleCondition, m FrontdoorConditionParameters) ([]interface{}, error) {
	condition, ok := input.AsDeliveryRuleServerPortCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule server port condition")
	}

	if params := condition.Parameters; params != nil {
		condition := normalizedCondition{
			selector:        nil,
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}

		if flattened, err := flattenAndValidateNormalizedCondition(condition, m); err != nil {
			return []interface{}{flattened}, err
		} else {
			return []interface{}{flattened}, nil
		}
	}

	return nil, nil
}

func FlattenFrontdoorHostNameCondition(input track1.BasicDeliveryRuleCondition, m FrontdoorConditionParameters) ([]interface{}, error) {
	condition, ok := input.AsDeliveryRuleHostNameCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule host name condition")
	}

	if params := condition.Parameters; params != nil {
		condition := normalizedCondition{
			selector:        nil,
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}

		if flattened, err := flattenAndValidateNormalizedCondition(condition, m); err != nil {
			return []interface{}{flattened}, err
		} else {
			return []interface{}{flattened}, nil
		}
	}

	return nil, nil
}

func FlattenFrontdoorSslProtocolCondition(input track1.BasicDeliveryRuleCondition, m FrontdoorConditionParameters) ([]interface{}, error) {
	condition, ok := input.AsDeliveryRuleSslProtocolCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule ssl protocol condition")
	}

	if params := condition.Parameters; params != nil {

		matchValues := make([]string, 0)
		for _, value := range *params.MatchValues {
			matchValues = append(matchValues, string(value))
		}

		condition := normalizedCondition{
			selector:        nil,
			operator:        string(*params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     &matchValues,
			transforms:      params.Transforms,
		}

		if flattened, err := flattenAndValidateNormalizedCondition(condition, m); err != nil {
			return []interface{}{flattened}, err
		} else {
			return []interface{}{flattened}, nil
		}
	}

	return nil, nil
}
