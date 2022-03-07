package frontdoorruleconditions

import (
	"fmt"

	track1 "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

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

func flattenNormalizedCondition(condition normalizedCondition) map[string]interface{} {
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

	return v
}

func ExpandFrontdoorRemoteAddressCondition(input []interface{}) []track1.BasicDeliveryRuleCondition {
	output := make([]track1.BasicDeliveryRuleCondition, 0)

	for _, v := range input {
		item := v.(map[string]interface{})
		remoteAddrCondition := track1.DeliveryRuleRemoteAddressCondition{
			Name: track1.NameRemoteAddress,
			Parameters: &track1.RemoteAddressMatchConditionParameters{
				TypeName:        utils.String("DeliveryRuleRemoteAddressMatchConditionParameters"),
				Operator:        track1.RemoteAddressOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		output = append(output, remoteAddrCondition)
	}

	return output
}

func ExpandFrontdoorRequestMethodCondition(input []interface{}) []track1.BasicDeliveryRuleCondition {
	output := make([]track1.BasicDeliveryRuleCondition, 0)

	for _, v := range input {
		item := v.(map[string]interface{})
		requestMethodCondition := track1.DeliveryRuleRequestMethodCondition{
			Name: track1.NameRequestMethod,
			Parameters: &track1.RequestMethodMatchConditionParameters{
				TypeName:        utils.String("DeliveryRuleRequestMethodMatchConditionParameters"),
				Operator:        utils.String(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		output = append(output, requestMethodCondition)
	}

	return output
}

func ExpandFrontdoorQueryStringCondition(input []interface{}) []track1.BasicDeliveryRuleCondition {
	output := make([]track1.BasicDeliveryRuleCondition, 0)

	for _, v := range input {
		item := v.(map[string]interface{})
		queryStringCondition := track1.DeliveryRuleQueryStringCondition{
			Name: track1.NameQueryString,
			Parameters: &track1.QueryStringMatchConditionParameters{
				TypeName:        utils.String("DeliveryRuleQueryStringMatchConditionParameters"),
				Operator:        track1.QueryStringOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if tt := item["transforms"].([]interface{}); len(tt) != 0 {
			transforms := make([]track1.Transform, 0)
			for _, t := range tt {
				transforms = append(transforms, track1.Transform(t.(string)))
			}
			queryStringCondition.Parameters.Transforms = &transforms
		}

		output = append(output, queryStringCondition)
	}

	return output
}

func ExpandFrontdoorPostArgsCondition(input []interface{}) []track1.BasicDeliveryRuleCondition {
	output := make([]track1.BasicDeliveryRuleCondition, 0)

	for _, v := range input {
		item := v.(map[string]interface{})
		cookiesCondition := track1.DeliveryRulePostArgsCondition{
			Name: track1.NamePostArgs,
			Parameters: &track1.PostArgsMatchConditionParameters{
				TypeName:        utils.String("DeliveryRulePostArgsMatchConditionParameters"),
				Selector:        utils.String(item["postargs_name"].(string)),
				Operator:        track1.PostArgsOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if tt := item["transforms"].([]interface{}); len(tt) != 0 {
			transforms := make([]track1.Transform, 0)
			for _, t := range tt {
				transforms = append(transforms, track1.Transform(t.(string)))
			}
			cookiesCondition.Parameters.Transforms = &transforms
		}

		output = append(output, cookiesCondition)
	}

	return output
}

func ExpandFrontdoorRequestUriCondition(input []interface{}) []track1.BasicDeliveryRuleCondition {
	output := make([]track1.BasicDeliveryRuleCondition, 0)

	for _, v := range input {
		item := v.(map[string]interface{})
		requestUriCondition := track1.DeliveryRuleRequestURICondition{
			Name: track1.NameRequestURI,
			Parameters: &track1.RequestURIMatchConditionParameters{
				TypeName:        utils.String("DeliveryRuleRequestURIMatchConditionParameters"),
				Operator:        track1.RequestURIOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if tt := item["transforms"].([]interface{}); len(tt) != 0 {
			transforms := make([]track1.Transform, 0)
			for _, t := range tt {
				transforms = append(transforms, track1.Transform(t.(string)))
			}
			requestUriCondition.Parameters.Transforms = &transforms
		}

		output = append(output, requestUriCondition)
	}

	return output
}

func ExpandFrontdoorRequestHeaderCondition(input []interface{}) []track1.BasicDeliveryRuleCondition {
	output := make([]track1.BasicDeliveryRuleCondition, 0)

	for _, v := range input {
		item := v.(map[string]interface{})
		requestHeaderCondition := track1.DeliveryRuleRequestHeaderCondition{
			Name: track1.NameRequestHeader,
			Parameters: &track1.RequestHeaderMatchConditionParameters{
				TypeName:        utils.String("DeliveryRuleRequestHeaderConditionParameters"),
				Selector:        utils.String(item["header_name"].(string)),
				Operator:        track1.RequestHeaderOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if tt := item["transforms"].([]interface{}); len(tt) != 0 {
			transforms := make([]track1.Transform, 0)
			for _, t := range tt {
				transforms = append(transforms, track1.Transform(t.(string)))
			}
			requestHeaderCondition.Parameters.Transforms = &transforms
		}

		output = append(output, requestHeaderCondition)
	}

	return output
}

func ExpandFrontdoorRequestBodyCondition(input []interface{}) []track1.BasicDeliveryRuleCondition {
	output := make([]track1.BasicDeliveryRuleCondition, 0)

	for _, v := range input {
		item := v.(map[string]interface{})
		requestBodyCondition := track1.DeliveryRuleRequestBodyCondition{
			Name: track1.NameRequestBody,
			Parameters: &track1.RequestBodyMatchConditionParameters{
				TypeName:        utils.String("DeliveryRuleRequestBodyMatchConditionParameters"),
				Operator:        track1.RequestBodyOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if tt := item["transforms"].([]interface{}); len(tt) != 0 {
			transforms := make([]track1.Transform, 0)
			for _, t := range tt {
				transforms = append(transforms, track1.Transform(t.(string)))
			}
			requestBodyCondition.Parameters.Transforms = &transforms
		}

		output = append(output, requestBodyCondition)
	}

	return output
}

func ExpandFrontdoorRequestSchemeCondition(input []interface{}) []track1.BasicDeliveryRuleCondition {
	output := make([]track1.BasicDeliveryRuleCondition, 0)

	for _, v := range input {
		item := v.(map[string]interface{})
		requestBodyCondition := track1.DeliveryRuleRequestSchemeCondition{
			Name: track1.NameRequestScheme,
			Parameters: &track1.RequestSchemeMatchConditionParameters{
				TypeName:        utils.String("DeliveryRuleRequestSchemeMatchConditionParameters"),
				Operator:        utils.String(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		output = append(output, requestBodyCondition)
	}

	return output
}

func ExpandFrontdoorUrlPathCondition(input []interface{}) []track1.BasicDeliveryRuleCondition {
	output := make([]track1.BasicDeliveryRuleCondition, 0)

	for _, v := range input {
		item := v.(map[string]interface{})
		urlPathCondition := track1.DeliveryRuleURLPathCondition{
			Name: track1.NameURLPath,
			Parameters: &track1.URLPathMatchConditionParameters{
				TypeName:        utils.String("DeliveryRuleURLPathMatchConditionParameters"),
				Operator:        track1.URLPathOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if tt := item["transforms"].([]interface{}); len(tt) != 0 {
			transforms := make([]track1.Transform, 0)
			for _, t := range tt {
				transforms = append(transforms, track1.Transform(t.(string)))
			}
			urlPathCondition.Parameters.Transforms = &transforms
		}

		output = append(output, urlPathCondition)
	}

	return output
}

func ExpandFrontdoorUrlFileExtensionCondition(input []interface{}) []track1.BasicDeliveryRuleCondition {
	output := make([]track1.BasicDeliveryRuleCondition, 0)

	for _, v := range input {
		item := v.(map[string]interface{})
		urlFileExtensionCondition := track1.DeliveryRuleURLFileExtensionCondition{
			Name: track1.NameURLFileExtension,
			Parameters: &track1.URLFileExtensionMatchConditionParameters{
				TypeName:        utils.String("DeliveryRuleURLFileExtensionMatchConditionParameters"),
				Operator:        track1.URLFileExtensionOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if tt := item["transforms"].([]interface{}); len(tt) != 0 {
			transforms := make([]track1.Transform, 0)
			for _, t := range tt {
				transforms = append(transforms, track1.Transform(t.(string)))
			}
			urlFileExtensionCondition.Parameters.Transforms = &transforms
		}

		output = append(output, urlFileExtensionCondition)
	}

	return output
}

func ExpandFrontdoorUrlFileNameCondition(input []interface{}) []track1.BasicDeliveryRuleCondition {
	output := make([]track1.BasicDeliveryRuleCondition, 0)

	for _, v := range input {
		item := v.(map[string]interface{})
		urlFileNameCondition := track1.DeliveryRuleURLFileNameCondition{
			Name: track1.NameURLFileName,
			Parameters: &track1.URLFileNameMatchConditionParameters{
				TypeName:        utils.String("DeliveryRuleURLFileExtensionMatchConditionParameters"),
				Operator:        track1.URLFileNameOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		if tt := item["transforms"].([]interface{}); len(tt) != 0 {
			transforms := make([]track1.Transform, 0)
			for _, t := range tt {
				transforms = append(transforms, track1.Transform(t.(string)))
			}
			urlFileNameCondition.Parameters.Transforms = &transforms
		}

		output = append(output, urlFileNameCondition)
	}

	return output
}

func ExpandFrontdoorHttpVersionCondition(input []interface{}) []track1.BasicDeliveryRuleCondition {
	output := make([]track1.BasicDeliveryRuleCondition, 0)

	for _, v := range input {
		item := v.(map[string]interface{})
		httpVersionCondition := track1.DeliveryRuleHTTPVersionCondition{
			Name: track1.NameHTTPVersion,
			Parameters: &track1.HTTPVersionMatchConditionParameters{
				TypeName:        utils.String("DeliveryRuleURLFileExtensionMatchConditionParameters"),
				Operator:        utils.String(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		output = append(output, httpVersionCondition)
	}

	return output
}

func ExpandFrontdoorCookiesCondition(input []interface{}) []track1.BasicDeliveryRuleCondition {
	output := make([]track1.BasicDeliveryRuleCondition, 0)
	for _, v := range input {
		item := v.(map[string]interface{})
		cookiesCondition := track1.DeliveryRuleCookiesCondition{
			Name: track1.NameCookies,
			Parameters: &track1.CookiesMatchConditionParameters{
				TypeName:        utils.String("DeliveryRuleCookiesConditionParameters"),
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
			cookiesCondition.Parameters.Transforms = &transforms
		}

		output = append(output, cookiesCondition)
	}

	return output
}

func ExpandFrontdoorIsDeviceCondition(input []interface{}) []track1.BasicDeliveryRuleCondition {
	output := make([]track1.BasicDeliveryRuleCondition, 0)

	for _, v := range input {
		item := v.(map[string]interface{})
		IsDeviceCondition := track1.DeliveryRuleIsDeviceCondition{
			Name: track1.NameIsDevice,
			Parameters: &track1.IsDeviceMatchConditionParameters{
				TypeName:        utils.String("DeliveryRuleIsDeviceMatchConditionParameters"),
				Operator:        utils.String(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].([]interface{})),
			},
		}

		output = append(output, IsDeviceCondition)
	}

	return output
}

func FlattenFrontdoorRemoteAddressCondition(input track1.BasicDeliveryRuleCondition) ([]interface{}, error) {
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

		return []interface{}{flattenNormalizedCondition(condition)}, nil
	}

	return nil, nil
}

func FlattenFrontdoorRequestMethodCondition(input track1.BasicDeliveryRuleCondition) ([]interface{}, error) {
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

		return []interface{}{flattenNormalizedCondition(condition)}, nil
	}

	return nil, nil
}

func FlattenFrontdoorQueryStringCondition(input track1.BasicDeliveryRuleCondition) ([]interface{}, error) {
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

		return []interface{}{flattenNormalizedCondition(condition)}, nil
	}

	return nil, nil
}

func FlattenFrontdoorPostArgsCondition(input track1.BasicDeliveryRuleCondition) ([]interface{}, error) {
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

		return []interface{}{flattenNormalizedCondition(condition)}, nil
	}

	return nil, nil
}

func FlattenFrontdoorRequestUriCondition(input track1.BasicDeliveryRuleCondition) ([]interface{}, error) {
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

		return []interface{}{flattenNormalizedCondition(condition)}, nil
	}

	return nil, nil
}

func FlattenFrontdoorRequestHeaderCondition(input track1.BasicDeliveryRuleCondition) ([]interface{}, error) {
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

		return []interface{}{flattenNormalizedCondition(condition)}, nil
	}

	return nil, nil
}

func FlattenFrontdoorRequestBodyCondition(input track1.BasicDeliveryRuleCondition) ([]interface{}, error) {
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

		return []interface{}{flattenNormalizedCondition(condition)}, nil
	}

	return nil, nil
}

func FlattenFrontdoorRequestSchemeCondition(input track1.BasicDeliveryRuleCondition) ([]interface{}, error) {
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

		return []interface{}{flattenNormalizedCondition(condition)}, nil
	}

	return nil, nil
}

func FlattenFrontdoorUrlPathCondition(input track1.BasicDeliveryRuleCondition) ([]interface{}, error) {
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

		return []interface{}{flattenNormalizedCondition(condition)}, nil
	}

	return nil, nil
}

func FlattenFrontdoorUrlFileExtensionCondition(input track1.BasicDeliveryRuleCondition) ([]interface{}, error) {
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

		return []interface{}{flattenNormalizedCondition(condition)}, nil
	}

	return nil, nil
}

func FlattenFrontdoorUrlFileNameCondition(input track1.BasicDeliveryRuleCondition) ([]interface{}, error) {
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

		return []interface{}{flattenNormalizedCondition(condition)}, nil
	}

	return nil, nil
}

func FlattenFrontdoorHttpVersionCondition(input track1.BasicDeliveryRuleCondition) ([]interface{}, error) {
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

		return []interface{}{flattenNormalizedCondition(condition)}, nil
	}

	return nil, nil
}

func FlattenFrontdoorCookiesCondition(input track1.BasicDeliveryRuleCondition) ([]interface{}, error) {
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

		return []interface{}{flattenNormalizedCondition(condition)}, nil
	}

	return nil, nil
}

func FlattenFrontdoorIsDeviceCondition(input track1.BasicDeliveryRuleCondition) ([]interface{}, error) {
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

		return []interface{}{flattenNormalizedCondition(condition)}, nil
	}

	return nil, nil
}
