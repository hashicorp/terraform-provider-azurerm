package cdnfrontdoorruleconditions

import (
	"fmt"
	"net"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2021-06-01/cdn"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type CdnFrontdoorConditionParameters struct {
	Name       cdn.Name
	TypeName   string
	ConfigName string
}

type CdnFrontdoorCondtionsMappings struct {
	ClientPort       CdnFrontdoorConditionParameters
	Cookies          CdnFrontdoorConditionParameters
	HostName         CdnFrontdoorConditionParameters
	HttpVersion      CdnFrontdoorConditionParameters
	IsDevice         CdnFrontdoorConditionParameters
	PostArgs         CdnFrontdoorConditionParameters
	QueryString      CdnFrontdoorConditionParameters
	RemoteAddress    CdnFrontdoorConditionParameters
	RequestBody      CdnFrontdoorConditionParameters
	RequestHeader    CdnFrontdoorConditionParameters
	RequestMethod    CdnFrontdoorConditionParameters
	RequestScheme    CdnFrontdoorConditionParameters
	RequestUri       CdnFrontdoorConditionParameters
	ServerPort       CdnFrontdoorConditionParameters
	SocketAddress    CdnFrontdoorConditionParameters
	SslProtocol      CdnFrontdoorConditionParameters
	UrlFileExtension CdnFrontdoorConditionParameters
	UrlFilename      CdnFrontdoorConditionParameters
	UrlPath          CdnFrontdoorConditionParameters
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

func InitializeCdnFrontdoorConditionMappings() *CdnFrontdoorCondtionsMappings {
	m := new(CdnFrontdoorCondtionsMappings)

	m.ClientPort = CdnFrontdoorConditionParameters{
		Name:       cdn.NameClientPort,
		TypeName:   "DeliveryRuleClientPortConditionParameters",
		ConfigName: "client_port_condition",
	}

	m.Cookies = CdnFrontdoorConditionParameters{
		Name:       cdn.NameCookies,
		TypeName:   "DeliveryRuleCookiesConditionParameters",
		ConfigName: "cookies_condition",
	}

	m.HostName = CdnFrontdoorConditionParameters{
		Name:       cdn.NameHostName,
		TypeName:   "DeliveryRuleHostNameConditionParameters",
		ConfigName: "host_name_condition",
	}

	m.HttpVersion = CdnFrontdoorConditionParameters{
		Name:       cdn.NameHTTPVersion,
		TypeName:   "DeliveryRuleHttpVersionConditionParameters",
		ConfigName: "http_version_condition",
	}

	m.IsDevice = CdnFrontdoorConditionParameters{
		Name:       cdn.NameIsDevice,
		TypeName:   "DeliveryRuleIsDeviceConditionParameters",
		ConfigName: "is_device_condition",
	}

	m.PostArgs = CdnFrontdoorConditionParameters{
		Name:       cdn.NamePostArgs,
		TypeName:   "DeliveryRulePostArgsConditionParameters",
		ConfigName: "post_args_condition",
	}

	m.QueryString = CdnFrontdoorConditionParameters{
		Name:       cdn.NameQueryString,
		TypeName:   "DeliveryRuleQueryStringConditionParameters",
		ConfigName: "query_string_condition",
	}

	m.RemoteAddress = CdnFrontdoorConditionParameters{
		Name:       cdn.NameRemoteAddress,
		TypeName:   "DeliveryRuleRemoteAddressConditionParameters",
		ConfigName: "remote_address_condition",
	}

	m.RequestBody = CdnFrontdoorConditionParameters{
		Name:       cdn.NameRequestBody,
		TypeName:   "DeliveryRuleRequestBodyConditionParameters",
		ConfigName: "request_body_condition",
	}

	m.RequestHeader = CdnFrontdoorConditionParameters{
		Name:       cdn.NameRequestHeader,
		TypeName:   "DeliveryRuleRequestHeaderConditionParameters",
		ConfigName: "request_header_condition",
	}

	m.RequestMethod = CdnFrontdoorConditionParameters{
		Name:       cdn.NameRequestMethod,
		TypeName:   "DeliveryRuleRequestMethodConditionParameters",
		ConfigName: "request_method_condition",
	}

	m.RequestScheme = CdnFrontdoorConditionParameters{
		Name:       cdn.NameRequestScheme,
		TypeName:   "DeliveryRuleRequestSchemeConditionParameters",
		ConfigName: "request_scheme_condition",
	}

	m.RequestUri = CdnFrontdoorConditionParameters{
		Name:       cdn.NameRequestURI,
		TypeName:   "DeliveryRuleRequestUriConditionParameters",
		ConfigName: "request_uri_condition",
	}

	m.ServerPort = CdnFrontdoorConditionParameters{
		Name:       cdn.NameServerPort,
		TypeName:   "DeliveryRuleServerPortConditionParameters",
		ConfigName: "server_port_condition",
	}

	m.SocketAddress = CdnFrontdoorConditionParameters{
		Name:       cdn.NameSocketAddr,
		TypeName:   "DeliveryRuleSocketAddrConditionParameters",
		ConfigName: "socket_address_condition",
	}

	m.SslProtocol = CdnFrontdoorConditionParameters{
		Name:       cdn.NameSslProtocol,
		TypeName:   "DeliveryRuleSslProtocolConditionParameters",
		ConfigName: "ssl_protocol_condition",
	}

	m.UrlFileExtension = CdnFrontdoorConditionParameters{
		Name:       cdn.NameURLFileExtension,
		TypeName:   "DeliveryRuleUrlFileExtensionMatchConditionParameters",
		ConfigName: "url_file_extension_condition",
	}

	m.UrlFilename = CdnFrontdoorConditionParameters{
		Name:       cdn.NameURLFileName,
		TypeName:   "DeliveryRuleUrlFilenameConditionParameters",
		ConfigName: "url_filename_condition",
	}

	m.UrlPath = CdnFrontdoorConditionParameters{
		Name:       cdn.NameURLPath,
		TypeName:   "DeliveryRuleUrlPathMatchConditionParameters",
		ConfigName: "url_path_condition",
	}

	return m
}

func checkForDuplicateCIDRs(input []interface{}) error {
	if len(input) <= 1 {
		return nil
	}

	tmp := make(map[string]bool)
	for _, CIDR := range input {
		if _, value := tmp[CIDR.(string)]; !value {
			tmp[CIDR.(string)] = true
		} else {
			return fmt.Errorf("%q CIDRs must be unique, there is a duplicate entry for CIDR %q in the %[1]q field. Please remove the duplicate entry and re-apply", "match_values", CIDR)
		}
	}

	return nil
}

func checkForCIDROverlap(matchValues []interface{}) error {
	// verify there are no duplicates in the CIDRs
	err := checkForDuplicateCIDRs(matchValues)
	if err != nil {
		return err
	}

	// separate the CIDRs into IPv6 and IPv4 variants
	IPv4CIDRs := make([]string, 0)
	IPv6CIDRs := make([]string, 0)

	for _, matchValue := range matchValues {
		if matchValue != nil {
			CIDR := matchValue.(string)

			// if CIDR is colon-hexadecimal it's IPv6
			if strings.Contains(CIDR, ":") {
				IPv6CIDRs = append(IPv6CIDRs, CIDR)
			} else {
				IPv4CIDRs = append(IPv4CIDRs, CIDR)
			}
		}
	}

	// check to see if the CIDR address ranges overlap based on the type of CIDR
	if len(IPv4CIDRs) > 1 {
		for _, sourceCIDR := range IPv4CIDRs {
			for _, checkCIDR := range IPv4CIDRs {
				if sourceCIDR == checkCIDR {
					continue
				}

				cidrOverlaps, err := validateCIDROverlap(sourceCIDR, checkCIDR)
				if err != nil {
					return err
				}

				if cidrOverlaps {
					return fmt.Errorf("the IPv4 %q CIDR %q address range overlaps with %q IPv4 CIDR address range", "match_values", sourceCIDR, checkCIDR)
				}
			}
		}
	}

	if len(IPv6CIDRs) > 1 {
		for _, sourceCIDR := range IPv6CIDRs {
			for _, checkCIDR := range IPv6CIDRs {
				if sourceCIDR == checkCIDR {
					continue
				}

				cidrOverlaps, err := validateCIDROverlap(sourceCIDR, checkCIDR)
				if err != nil {
					return fmt.Errorf("unable to validate IPv6 CIDR address ranges overlap: %+v", err)
				}

				if cidrOverlaps {
					return fmt.Errorf("the %q IPv6 CIDR %q address range overlaps with %q IPv6 CIDR address range", "match_values", sourceCIDR, checkCIDR)
				}
			}
		}
	}

	return nil
}

// evaluates if the passed CIDR is a valid IPv4 or IPv6 CIDR or not.
func isValidCidr(cidr interface{}) bool {
	if strings.Contains(cidr.(string), ":") {
		// evaluates if the passed CIDR is a valid IPv6 CIDR or not.
		ok, _ := validate.RegExHelper(cidr, "match_values", `^s*((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]d|1dd|[1-9]?d)(.(25[0-5]|2[0-4]d|1dd|[1-9]?d)){3}))|:)))(%.+)?s*(\/([0-9]|[1-9][0-9]|1[0-1][0-9]|12[0-8]))?$`)
		return ok
	}

	// evaluates if the passed CIDR is a valid IPv4 CIDR or not.
	ok, _ := validate.RegExHelper(cidr, "match_values", `^([0-9]{1,3}\.){3}[0-9]{1,3}(/([0-9]|[1-2][0-9]|3[0-2]))?$`)
	return ok

}

func validateCIDROverlap(sourceCIDR string, checkCIDR string) (bool, error) {
	_, sourceNetwork, err := net.ParseCIDR(sourceCIDR)
	if err != nil {
		return false, err
	}

	sourceOnes, sourceBits := sourceNetwork.Mask.Size()
	if sourceOnes == 0 && sourceBits == 0 {
		return false, fmt.Errorf("%q CIDR must be in its canonical form", sourceCIDR)
	}

	_, checkNetwork, err := net.ParseCIDR(checkCIDR)
	if err != nil {
		return false, err
	}

	checkOnes, checkBits := checkNetwork.Mask.Size()
	if checkOnes == 0 && checkBits == 0 {
		return false, fmt.Errorf("%q CIDR must be in its canonical form", checkCIDR)
	}

	ipStr := checkNetwork.IP.String()
	checkIp := net.ParseIP(ipStr)
	if checkIp == nil {
		return false, fmt.Errorf("unable to parse %q, invalid IP address", ipStr)
	}

	ipStr = sourceNetwork.IP.String()
	sourceIp := net.ParseIP(ipStr)
	if sourceIp == nil {
		return false, fmt.Errorf("unable to parse %q, invalid IP address", ipStr)
	}

	// swap the check values depending on which CIDR is more specific
	// So much time and so little to do. Wait a minute.
	// Strike that. Reverse it.
	if sourceOnes > checkOnes {
		sourceNetwork = checkNetwork
		checkIp = sourceIp
	}

	// validate that the passed CIDRs don't overlap
	if !sourceNetwork.Contains(checkIp) {
		return false, nil
	}

	// CIDR overlap was detected
	return true, nil
}

func expandNormalizeCdnFrontdoorTransforms(input []interface{}) []cdn.Transform {
	transforms := make([]cdn.Transform, 0)
	if len(input) == 0 {
		return transforms
	}

	for _, t := range input {
		transforms = append(transforms, cdn.Transform(t.(string)))
	}

	return transforms
}

func flattenAndValidateCdnFrontdoorNormalizedCondition(condition normalizedCondition) map[string]interface{} {
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

func validateCdnFrontdoorExpandConditionOperatorValues(operator string, matchValues *[]string, m CdnFrontdoorConditionParameters) error {
	if operator == "" {
		return fmt.Errorf("%q is invalid: no %q value has been set, got %q", m.ConfigName, "operator", operator)
	}

	if operator == string(cdn.OperatorAny) && len(*matchValues) > 0 {
		return fmt.Errorf("%q is invalid: the %q field must not be set if the conditions %q is set to %q", m.ConfigName, "match_values", "operator", cdn.OperatorAny)
	}

	return nil
}

func ExpandCdnFrontdoorRemoteAddressCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontdoorConditionMappings()
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
						return nil, fmt.Errorf("%q is invalid: when the %q is set to %q the value must be a valid country code consisting of 2 uppercase characters, got %q", conditionMapping.ConfigName, "operator", cdn.RemoteAddressOperatorGeoMatch, matchValue)
					}
				}
			}
		}

		if condition.Parameters.Operator == cdn.RemoteAddressOperatorIPMatch {
			for _, matchValue := range item["match_values"].([]interface{}) {
				address := ""
				if matchValue != nil {
					address = matchValue.(string)
				}

				if !isValidCidr(address) {
					return nil, fmt.Errorf("%q is invalid: when the %q is set to %q the value must be a valid IPv4 or IPv6 CIDR, got %q", conditionMapping.ConfigName, "operator", cdn.RemoteAddressOperatorIPMatch, address)
				}
			}

			// Check for CIDR overlap and CIDR duplicates in the match values
			err := checkForCIDROverlap(item["match_values"].([]interface{}))
			if err != nil {
				return nil, fmt.Errorf("%q is invalid: %+v", conditionMapping.ConfigName, err)
			}
		}

		if err := validateCdnFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontdoorRequestMethodCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontdoorConditionMappings()
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

		if err := validateCdnFrontdoorExpandConditionOperatorValues(*condition.Parameters.Operator, condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontdoorQueryStringCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontdoorConditionMappings()
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
			expanded := expandNormalizeCdnFrontdoorTransforms(transformsRaw)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateCdnFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontdoorPostArgsCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontdoorConditionMappings()
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
			expanded := expandNormalizeCdnFrontdoorTransforms(transformsRaw)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateCdnFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontdoorRequestUriCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontdoorConditionMappings()
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
			expanded := expandNormalizeCdnFrontdoorTransforms(transformsRaw)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateCdnFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontdoorRequestHeaderCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontdoorConditionMappings()
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
			expanded := expandNormalizeCdnFrontdoorTransforms(transformsRaw)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateCdnFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontdoorRequestBodyCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontdoorConditionMappings()
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
			expanded := expandNormalizeCdnFrontdoorTransforms(transformsRaw)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateCdnFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontdoorRequestSchemeCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontdoorConditionMappings()
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

		if err := validateCdnFrontdoorExpandConditionOperatorValues(*condition.Parameters.Operator, condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontdoorUrlPathCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontdoorConditionMappings()
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
			expanded := expandNormalizeCdnFrontdoorTransforms(transformsRaw)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateCdnFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontdoorUrlFileExtensionCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontdoorConditionMappings()
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
			expanded := expandNormalizeCdnFrontdoorTransforms(transformsRaw)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateCdnFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontdoorUrlFileNameCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontdoorConditionMappings()
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
			expanded := expandNormalizeCdnFrontdoorTransforms(transformsRaw)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateCdnFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontdoorHttpVersionCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontdoorConditionMappings()
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

		if err := validateCdnFrontdoorExpandConditionOperatorValues(*condition.Parameters.Operator, condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontdoorCookiesCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontdoorConditionMappings()
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

		if tt := item["transforms"].([]interface{}); len(tt) != 0 {
			transforms := make([]cdn.Transform, 0)
			for _, t := range tt {
				transforms = append(transforms, cdn.Transform(t.(string)))
			}
			condition.Parameters.Transforms = &transforms
		}

		if err := validateCdnFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontdoorIsDeviceCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontdoorConditionMappings()
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

		if err := validateCdnFrontdoorExpandConditionOperatorValues(*condition.Parameters.Operator, condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontdoorSocketAddressCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontdoorConditionMappings()
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
			for _, matchValue := range item["match_values"].([]interface{}) {
				address := ""
				if matchValue != nil {
					address = matchValue.(string)
				}

				if !isValidCidr(address) {
					return nil, fmt.Errorf("%q is invalid: when the %q is set to %q the %q must be a valid IPv4 or IPv6 CIDR, got %q", conditionMapping.ConfigName, "operator", cdn.SocketAddrOperatorIPMatch, "match_values", address)
				}
			}

			// Check for CIDR overlap and CIDR duplicates in the match values
			err := checkForCIDROverlap(item["match_values"].([]interface{}))
			if err != nil {
				return nil, fmt.Errorf("%q is invalid: %+v", conditionMapping.ConfigName, err)
			}
		}

		if err := validateCdnFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontdoorClientPortCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontdoorConditionMappings()
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

		if err := validateCdnFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontdoorServerPortCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontdoorConditionMappings()
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

		if err := validateCdnFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontdoorHostNameCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontdoorConditionMappings()
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
			expanded := expandNormalizeCdnFrontdoorTransforms(transformsRaw)
			condition.Parameters.Transforms = &expanded
		}

		if err := validateCdnFrontdoorExpandConditionOperatorValues(string(condition.Parameters.Operator), condition.Parameters.MatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func ExpandCdnFrontdoorSslProtocolCondition(input []interface{}) (*[]cdn.BasicDeliveryRuleCondition, error) {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	m := InitializeCdnFrontdoorConditionMappings()
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

		if err := validateCdnFrontdoorExpandConditionOperatorValues(*condition.Parameters.Operator, &validationMatchValues, conditionMapping); err != nil {
			return nil, err
		}

		output = append(output, condition)
	}

	return &output, nil
}

func FlattenFrontdoorRemoteAddressCondition(input cdn.BasicDeliveryRuleCondition, m CdnFrontdoorConditionParameters) (map[string]interface{}, error) {
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
			transforms:      nil,
		}

		return flattenAndValidateCdnFrontdoorNormalizedCondition(condition), nil
	}

	return nil, nil
}

func FlattenFrontdoorRequestMethodCondition(input cdn.BasicDeliveryRuleCondition, m CdnFrontdoorConditionParameters) (map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleRequestMethodCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule request method condition")
	}

	if params := condition.Parameters; params != nil {
		condition := normalizedCondition{
			selector:        nil,
			operator:        *params.Operator,
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      nil,
		}

		return flattenAndValidateCdnFrontdoorNormalizedCondition(condition), nil
	}

	return nil, nil
}

func FlattenFrontdoorQueryStringCondition(input cdn.BasicDeliveryRuleCondition, m CdnFrontdoorConditionParameters) (map[string]interface{}, error) {
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

		return flattenAndValidateCdnFrontdoorNormalizedCondition(condition), nil
	}

	return nil, nil
}

func FlattenFrontdoorPostArgsCondition(input cdn.BasicDeliveryRuleCondition, m CdnFrontdoorConditionParameters) (map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRulePostArgsCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule post args condition")
	}

	if params := condition.Parameters; params != nil {
		condition := normalizedCondition{
			selector:        &normalizedSelector{name: utils.String("post_args_name"), value: params.Selector},
			operator:        string(params.Operator),
			negateCondition: params.NegateCondition,
			matchValues:     params.MatchValues,
			transforms:      params.Transforms,
		}

		return flattenAndValidateCdnFrontdoorNormalizedCondition(condition), nil
	}

	return nil, nil
}

func FlattenFrontdoorRequestUriCondition(input cdn.BasicDeliveryRuleCondition, m CdnFrontdoorConditionParameters) (map[string]interface{}, error) {
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

		return flattenAndValidateCdnFrontdoorNormalizedCondition(condition), nil
	}

	return nil, nil
}

func FlattenFrontdoorRequestHeaderCondition(input cdn.BasicDeliveryRuleCondition, m CdnFrontdoorConditionParameters) (map[string]interface{}, error) {
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

		return flattenAndValidateCdnFrontdoorNormalizedCondition(condition), nil
	}

	return nil, nil
}

func FlattenFrontdoorRequestBodyCondition(input cdn.BasicDeliveryRuleCondition, m CdnFrontdoorConditionParameters) (map[string]interface{}, error) {
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

		return flattenAndValidateCdnFrontdoorNormalizedCondition(condition), nil
	}

	return nil, nil
}

func FlattenFrontdoorRequestSchemeCondition(input cdn.BasicDeliveryRuleCondition, m CdnFrontdoorConditionParameters) (map[string]interface{}, error) {
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

		return flattenAndValidateCdnFrontdoorNormalizedCondition(condition), nil
	}

	return nil, nil
}

func FlattenFrontdoorUrlPathCondition(input cdn.BasicDeliveryRuleCondition, m CdnFrontdoorConditionParameters) (map[string]interface{}, error) {
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

		return flattenAndValidateCdnFrontdoorNormalizedCondition(condition), nil
	}

	return nil, nil
}

func FlattenFrontdoorUrlFileExtensionCondition(input cdn.BasicDeliveryRuleCondition, m CdnFrontdoorConditionParameters) (map[string]interface{}, error) {
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

		return flattenAndValidateCdnFrontdoorNormalizedCondition(condition), nil
	}

	return nil, nil
}

func FlattenFrontdoorUrlFileNameCondition(input cdn.BasicDeliveryRuleCondition, m CdnFrontdoorConditionParameters) (map[string]interface{}, error) {
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

		return flattenAndValidateCdnFrontdoorNormalizedCondition(condition), nil
	}

	return nil, nil
}

func FlattenFrontdoorHttpVersionCondition(input cdn.BasicDeliveryRuleCondition, m CdnFrontdoorConditionParameters) (map[string]interface{}, error) {
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

		return flattenAndValidateCdnFrontdoorNormalizedCondition(condition), nil
	}

	return nil, nil
}

func FlattenFrontdoorCookiesCondition(input cdn.BasicDeliveryRuleCondition, m CdnFrontdoorConditionParameters) (map[string]interface{}, error) {
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

		return flattenAndValidateCdnFrontdoorNormalizedCondition(condition), nil
	}

	return nil, nil
}

func FlattenFrontdoorIsDeviceCondition(input cdn.BasicDeliveryRuleCondition, m CdnFrontdoorConditionParameters) (map[string]interface{}, error) {
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

		return flattenAndValidateCdnFrontdoorNormalizedCondition(condition), nil
	}

	return nil, nil
}

func FlattenFrontdoorSocketAddressCondition(input cdn.BasicDeliveryRuleCondition, m CdnFrontdoorConditionParameters) (map[string]interface{}, error) {
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
			transforms:      nil,
		}

		return flattenAndValidateCdnFrontdoorNormalizedCondition(condition), nil
	}

	return nil, nil
}

func FlattenFrontdoorClientPortCondition(input cdn.BasicDeliveryRuleCondition, m CdnFrontdoorConditionParameters) (map[string]interface{}, error) {
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
			transforms:      nil,
		}

		return flattenAndValidateCdnFrontdoorNormalizedCondition(condition), nil
	}

	return nil, nil
}

func FlattenFrontdoorServerPortCondition(input cdn.BasicDeliveryRuleCondition, m CdnFrontdoorConditionParameters) (map[string]interface{}, error) {
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
			transforms:      nil,
		}

		return flattenAndValidateCdnFrontdoorNormalizedCondition(condition), nil
	}

	return nil, nil
}

func FlattenFrontdoorHostNameCondition(input cdn.BasicDeliveryRuleCondition, m CdnFrontdoorConditionParameters) (map[string]interface{}, error) {
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

		return flattenAndValidateCdnFrontdoorNormalizedCondition(condition), nil
	}

	return nil, nil
}

func FlattenFrontdoorSslProtocolCondition(input cdn.BasicDeliveryRuleCondition, m CdnFrontdoorConditionParameters) (map[string]interface{}, error) {
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
			operator:        *params.Operator,
			negateCondition: params.NegateCondition,
			matchValues:     &matchValues,
			transforms:      nil,
		}

		return flattenAndValidateCdnFrontdoorNormalizedCondition(condition), nil
	}

	return nil, nil
}
