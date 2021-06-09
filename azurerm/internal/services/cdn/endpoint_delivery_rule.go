package cdn

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/deliveryruleactions"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/deliveryruleconditions"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func endpointDeliveryRule() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validate.EndpointDeliveryRuleName(),
				},

				"order": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntAtLeast(1),
				},

				"cookies_condition": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem:     deliveryruleconditions.Cookies(),
				},

				"http_version_condition": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem:     deliveryruleconditions.HTTPVersion(),
				},

				"device_condition": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem:     deliveryruleconditions.Device(),
				},

				"post_arg_condition": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem:     deliveryruleconditions.PostArg(),
				},

				"query_string_condition": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem:     deliveryruleconditions.QueryString(),
				},

				"remote_address_condition": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem:     deliveryruleconditions.RemoteAddress(),
				},

				"request_body_condition": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem:     deliveryruleconditions.RequestBody(),
				},

				"request_header_condition": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem:     deliveryruleconditions.RequestHeader(),
				},

				"request_method_condition": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem:     deliveryruleconditions.RequestMethod(),
				},

				"request_scheme_condition": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem:     deliveryruleconditions.RequestScheme(),
				},

				"request_uri_condition": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem:     deliveryruleconditions.RequestURI(),
				},

				"url_file_extension_condition": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem:     deliveryruleconditions.URLFileExtension(),
				},

				"url_file_name_condition": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem:     deliveryruleconditions.URLFileName(),
				},

				"url_path_condition": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem:     deliveryruleconditions.URLPath(),
				},

				"cache_expiration_action": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem:     deliveryruleactions.CacheExpiration(),
				},

				"cache_key_query_string_action": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem:     deliveryruleactions.CacheKeyQueryString(),
				},

				"modify_request_header_action": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem:     deliveryruleactions.ModifyRequestHeader(),
				},

				"modify_response_header_action": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem:     deliveryruleactions.ModifyResponseHeader(),
				},

				"url_redirect_action": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem:     deliveryruleactions.URLRedirect(),
				},

				"url_rewrite_action": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem:     deliveryruleactions.URLRewrite(),
				},
			},
		},
	}
}

func expandArmCdnEndpointDeliveryRule(rule map[string]interface{}) (*cdn.DeliveryRule, error) {
	deliveryRule := cdn.DeliveryRule{
		Name:  utils.String(rule["name"].(string)),
		Order: utils.Int32(int32(rule["order"].(int))),
	}

	deliveryRule.Conditions = expandDeliveryRuleConditions(rule)

	actions, err := expandDeliveryRuleActions(rule)
	if err != nil {
		return nil, err
	}
	deliveryRule.Actions = &actions

	return &deliveryRule, nil
}

func expandDeliveryRuleConditions(input map[string]interface{}) *[]cdn.BasicDeliveryRuleCondition {
	conditions := make([]cdn.BasicDeliveryRuleCondition, 0)

	// @tombuildsstuff: we'd generally avoid over generalization, but this is /very/ repetitive so makes sense
	type expandFunc func(input []interface{}) []cdn.BasicDeliveryRuleCondition
	conditionTypes := map[string]expandFunc{
		"cookies_condition":            deliveryruleconditions.ExpandArmCdnEndpointConditionCookies,
		"device_condition":             deliveryruleconditions.ExpandArmCdnEndpointConditionDevice,
		"http_version_condition":       deliveryruleconditions.ExpandArmCdnEndpointConditionHTTPVersion,
		"query_string_condition":       deliveryruleconditions.ExpandArmCdnEndpointConditionQueryString,
		"post_arg_condition":           deliveryruleconditions.ExpandArmCdnEndpointConditionPostArg,
		"request_body_condition":       deliveryruleconditions.ExpandArmCdnEndpointConditionRequestBody,
		"request_header_condition":     deliveryruleconditions.ExpandArmCdnEndpointConditionRequestHeader,
		"request_method_condition":     deliveryruleconditions.ExpandArmCdnEndpointConditionRequestMethod,
		"remote_address_condition":     deliveryruleconditions.ExpandArmCdnEndpointConditionRemoteAddress,
		"request_scheme_condition":     deliveryruleconditions.ExpandArmCdnEndpointConditionRequestScheme,
		"request_uri_condition":        deliveryruleconditions.ExpandArmCdnEndpointConditionRequestURI,
		"url_file_extension_condition": deliveryruleconditions.ExpandArmCdnEndpointConditionURLFileExtension,
		"url_file_name_condition":      deliveryruleconditions.ExpandArmCdnEndpointConditionURLFileName,
		"url_path_condition":           deliveryruleconditions.ExpandArmCdnEndpointConditionURLPath,
	}

	for schemaKey, expandFunc := range conditionTypes {
		raw := input[schemaKey].([]interface{})
		expanded := expandFunc(raw)
		conditions = append(conditions, expanded...)
	}

	return &conditions
}

func expandDeliveryRuleActions(input map[string]interface{}) ([]cdn.BasicDeliveryRuleAction, error) {
	actions := make([]cdn.BasicDeliveryRuleAction, 0)

	type expandFunc func(input []interface{}) (*[]cdn.BasicDeliveryRuleAction, error)
	actionTypes := map[string]expandFunc{
		"cache_expiration_action":       deliveryruleactions.ExpandArmCdnEndpointActionCacheExpiration,
		"cache_key_query_string_action": deliveryruleactions.ExpandArmCdnEndpointActionCacheKeyQueryString,
		"modify_request_header_action":  deliveryruleactions.ExpandArmCdnEndpointActionModifyRequestHeader,
		"modify_response_header_action": deliveryruleactions.ExpandArmCdnEndpointActionModifyResponseHeader,
		"url_redirect_action":           deliveryruleactions.ExpandArmCdnEndpointActionUrlRedirect,
		"url_rewrite_action":            deliveryruleactions.ExpandArmCdnEndpointActionURLRewrite,
	}

	for schemaKey, expandFunc := range actionTypes {
		raw := input[schemaKey].([]interface{})
		expanded, err := expandFunc(raw)
		if err != nil {
			return nil, err
		}

		actions = append(actions, *expanded...)
	}

	return actions, nil
}

func flattenArmCdnEndpointDeliveryRule(deliveryRule cdn.DeliveryRule) (*map[string]interface{}, error) {
	name := ""
	if deliveryRule.Name != nil {
		name = *deliveryRule.Name
	}

	order := -1
	if deliveryRule.Order != nil {
		order = int(*deliveryRule.Order)
	}

	output := map[string]interface{}{
		"name":  name,
		"order": order,
	}

	conditions, err := flattenDeliveryRuleConditions(deliveryRule.Conditions)
	if err != nil {
		return nil, err
	}

	for key, value := range *conditions {
		output[key] = value
	}

	actions, err := flattenDeliveryRuleActions(deliveryRule.Actions)
	if err != nil {
		return nil, err
	}

	for key, value := range *actions {
		output[key] = value
	}

	return &output, nil
}

func flattenDeliveryRuleActions(actions *[]cdn.BasicDeliveryRuleAction) (*map[string][]interface{}, error) {
	type flattenFunc = func(input cdn.BasicDeliveryRuleAction) (*map[string]interface{}, error)
	type validateFunc = func(input cdn.BasicDeliveryRuleAction) bool

	actionTypes := map[string]struct {
		flattenFunc  flattenFunc
		validateFunc validateFunc
	}{
		"cache_expiration_action": {
			flattenFunc: deliveryruleactions.FlattenArmCdnEndpointActionCacheExpiration,
			validateFunc: func(action cdn.BasicDeliveryRuleAction) bool {
				_, ok := action.AsDeliveryRuleCacheExpirationAction()
				return ok
			},
		},
		"cache_key_query_string_action": {
			flattenFunc: deliveryruleactions.FlattenArmCdnEndpointActionCacheKeyQueryString,
			validateFunc: func(action cdn.BasicDeliveryRuleAction) bool {
				_, ok := action.AsDeliveryRuleCacheKeyQueryStringAction()
				return ok
			},
		},
		"modify_request_header_action": {
			flattenFunc: deliveryruleactions.FlattenArmCdnEndpointActionModifyRequestHeader,
			validateFunc: func(action cdn.BasicDeliveryRuleAction) bool {
				_, ok := action.AsDeliveryRuleRequestHeaderAction()
				return ok
			},
		},
		"modify_response_header_action": {
			flattenFunc: deliveryruleactions.FlattenArmCdnEndpointActionModifyResponseHeader,
			validateFunc: func(action cdn.BasicDeliveryRuleAction) bool {
				_, ok := action.AsDeliveryRuleResponseHeaderAction()
				return ok
			},
		},
		"url_redirect_action": {
			flattenFunc: deliveryruleactions.FlattenArmCdnEndpointActionUrlRedirect,
			validateFunc: func(action cdn.BasicDeliveryRuleAction) bool {
				_, ok := action.AsURLRedirectAction()
				return ok
			},
		},
		"url_rewrite_action": {
			flattenFunc: deliveryruleactions.FlattenArmCdnEndpointActionURLRewrite,
			validateFunc: func(action cdn.BasicDeliveryRuleAction) bool {
				_, ok := action.AsURLRewriteAction()
				return ok
			},
		},
	}

	// first ensure there's a map for all of the keys
	output := make(map[string][]interface{})
	for schemaKey := range actionTypes {
		output[schemaKey] = make([]interface{}, 0)
	}

	// intentionally bail here now we have defaults populated
	if actions == nil {
		return &output, nil
	}

	// then iterate over all the actions and map them as necessary
	for _, action := range *actions {
		for schemaKey, actionType := range actionTypes {
			suitable := actionType.validateFunc(action)
			if !suitable {
				continue
			}

			mapped, err := actionType.flattenFunc(action)
			if err != nil {
				return nil, err
			}

			output[schemaKey] = append(output[schemaKey], mapped)
			break
		}
	}

	return &output, nil
}

func flattenDeliveryRuleConditions(conditions *[]cdn.BasicDeliveryRuleCondition) (*map[string][]interface{}, error) {
	type flattenFunc = func(input cdn.BasicDeliveryRuleCondition) (*map[string]interface{}, error)
	type validateFunc = func(input cdn.BasicDeliveryRuleCondition) bool

	conditionTypes := map[string]struct {
		flattenFunc  flattenFunc
		validateFunc validateFunc
	}{
		"cookies_condition": {
			flattenFunc: deliveryruleconditions.FlattenArmCdnEndpointConditionCookies,
			validateFunc: func(condition cdn.BasicDeliveryRuleCondition) bool {
				_, ok := condition.AsDeliveryRuleCookiesCondition()
				return ok
			},
		},
		"device_condition": {
			flattenFunc: deliveryruleconditions.FlattenArmCdnEndpointConditionDevice,
			validateFunc: func(condition cdn.BasicDeliveryRuleCondition) bool {
				_, ok := condition.AsDeliveryRuleIsDeviceCondition()
				return ok
			},
		},
		"http_version_condition": {
			flattenFunc: deliveryruleconditions.FlattenArmCdnEndpointConditionHTTPVersion,
			validateFunc: func(condition cdn.BasicDeliveryRuleCondition) bool {
				_, ok := condition.AsDeliveryRuleHTTPVersionCondition()
				return ok
			},
		},
		"query_string_condition": {
			flattenFunc: deliveryruleconditions.FlattenArmCdnEndpointConditionQueryString,
			validateFunc: func(condition cdn.BasicDeliveryRuleCondition) bool {
				_, ok := condition.AsDeliveryRuleQueryStringCondition()
				return ok
			},
		},
		"post_arg_condition": {
			flattenFunc: deliveryruleconditions.FlattenArmCdnEndpointConditionPostArg,
			validateFunc: func(condition cdn.BasicDeliveryRuleCondition) bool {
				_, ok := condition.AsDeliveryRulePostArgsCondition()
				return ok
			},
		},
		"remote_address_condition": {
			flattenFunc: deliveryruleconditions.FlattenArmCdnEndpointConditionRemoteAddress,
			validateFunc: func(condition cdn.BasicDeliveryRuleCondition) bool {
				_, ok := condition.AsDeliveryRuleRemoteAddressCondition()
				return ok
			},
		},
		"request_body_condition": {
			flattenFunc: deliveryruleconditions.FlattenArmCdnEndpointConditionRequestBody,
			validateFunc: func(condition cdn.BasicDeliveryRuleCondition) bool {
				_, ok := condition.AsDeliveryRuleRequestBodyCondition()
				return ok
			},
		},
		"request_header_condition": {
			flattenFunc: deliveryruleconditions.FlattenArmCdnEndpointConditionRequestHeader,
			validateFunc: func(condition cdn.BasicDeliveryRuleCondition) bool {
				_, ok := condition.AsDeliveryRuleRequestHeaderCondition()
				return ok
			},
		},
		"request_method_condition": {
			flattenFunc: deliveryruleconditions.FlattenArmCdnEndpointConditionRequestMethod,
			validateFunc: func(condition cdn.BasicDeliveryRuleCondition) bool {
				_, ok := condition.AsDeliveryRuleRequestMethodCondition()
				return ok
			},
		},
		"request_scheme_condition": {
			flattenFunc: deliveryruleconditions.FlattenArmCdnEndpointConditionRequestScheme,
			validateFunc: func(condition cdn.BasicDeliveryRuleCondition) bool {
				_, ok := condition.AsDeliveryRuleRequestSchemeCondition()
				return ok
			},
		},
		"request_uri_condition": {
			flattenFunc: deliveryruleconditions.FlattenArmCdnEndpointConditionRequestURI,
			validateFunc: func(condition cdn.BasicDeliveryRuleCondition) bool {
				_, ok := condition.AsDeliveryRuleRequestURICondition()
				return ok
			},
		},
		"url_file_extension_condition": {
			flattenFunc: deliveryruleconditions.FlattenArmCdnEndpointConditionURLFileExtension,
			validateFunc: func(condition cdn.BasicDeliveryRuleCondition) bool {
				_, ok := condition.AsDeliveryRuleURLFileExtensionCondition()
				return ok
			},
		},
		"url_file_name_condition": {
			flattenFunc: deliveryruleconditions.FlattenArmCdnEndpointConditionURLFileName,
			validateFunc: func(condition cdn.BasicDeliveryRuleCondition) bool {
				_, ok := condition.AsDeliveryRuleURLFileNameCondition()
				return ok
			},
		},
		"url_path_condition": {
			flattenFunc: deliveryruleconditions.FlattenArmCdnEndpointConditionURLPath,
			validateFunc: func(condition cdn.BasicDeliveryRuleCondition) bool {
				_, ok := condition.AsDeliveryRuleURLPathCondition()
				return ok
			},
		},
	}

	// first ensure there's a map for all of the keys
	output := make(map[string][]interface{})
	for schemaKey := range conditionTypes {
		output[schemaKey] = make([]interface{}, 0)
	}

	// intentionally bail here now we have defaults populated
	if conditions == nil {
		return &output, nil
	}

	// then iterate over all the conditions and map them as necessary
	for _, condition := range *conditions {
		for schemaKey, conditionType := range conditionTypes {
			suitable := conditionType.validateFunc(condition)
			if !suitable {
				continue
			}

			mapped, err := conditionType.flattenFunc(condition)
			if err != nil {
				return nil, err
			}

			output[schemaKey] = append(output[schemaKey], mapped)
			break
		}
	}

	return &output, nil
}
