package cdn

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/deliveryruleactions"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/deliveryruleconditions"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func EndpointDeliveryRule() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validate.EndpointDeliveryPolicyRuleName(),
				},

				"order": {
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntAtLeast(1),
				},

				"cookies_condition": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     deliveryruleconditions.Cookies(),
				},

				"http_version_condition": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     deliveryruleconditions.HTTPVersion(),
				},

				"device_condition": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem:     deliveryruleconditions.Device(),
				},

				"post_arg_condition": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     deliveryruleconditions.PostArg(),
				},

				"query_string_condition": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     deliveryruleconditions.QueryString(),
				},

				"remote_address_condition": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     deliveryruleconditions.RemoteAddress(),
				},

				"request_body_condition": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     deliveryruleconditions.RequestBody(),
				},

				"request_header_condition": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     deliveryruleconditions.RequestHeader(),
				},

				"request_method_condition": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem:     deliveryruleconditions.RequestMethod(),
				},

				"request_scheme_condition": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem:     deliveryruleconditions.RequestScheme(),
				},

				"request_uri_condition": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     deliveryruleconditions.RequestURI(),
				},

				"url_file_extension_condition": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     deliveryruleconditions.URLFileExtension(),
				},

				"cache_expiration_action": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem:     deliveryruleactions.CacheExpiration(),
				},

				"cache_key_query_string_action": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem:     deliveryruleactions.CacheKeyQueryString(),
				},

				"modify_request_header_action": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     deliveryruleactions.ModifyRequestHeader(),
				},

				"modify_response_header_action": {
					Type:     schema.TypeList,
					Optional: true,
					Elem:     deliveryruleactions.ModifyResponseHeader(),
				},

				"url_redirect_action": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem:     deliveryruleactions.URLRedirect(),
				},

				"url_rewrite_action": {
					Type:     schema.TypeList,
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

	conditions := expandDeliveryRuleConditions(rule)
	deliveryRule.Conditions = &conditions

	actions, err := expandDeliveryRuleActions(rule)
	if err != nil {
		return nil, err
	}
	deliveryRule.Actions = &actions

	return &deliveryRule, nil
}

func expandDeliveryRuleConditions(rule map[string]interface{}) []cdn.BasicDeliveryRuleCondition {
	conditions := make([]cdn.BasicDeliveryRuleCondition, 0)

	if ccs := rule["cookies_condition"].([]interface{}); len(ccs) > 0 {
		for _, cc := range ccs {
			conditions = append(conditions, *deliveryruleconditions.ExpandArmCdnEndpointConditionCookies(cc.(map[string]interface{})))
		}
	}

	if hvcs := rule["http_version_condition"].([]interface{}); len(hvcs) > 0 {
		for _, hvc := range hvcs {
			conditions = append(conditions, *deliveryruleconditions.ExpandArmCdnEndpointConditionHTTPVersion(hvc.(map[string]interface{})))
		}
	}

	if rsc := rule["device_condition"].([]interface{}); len(rsc) > 0 {
		conditions = append(conditions, *deliveryruleconditions.ExpandArmCdnEndpointConditionDevice(rsc[0].(map[string]interface{})))
	}

	if pacs := rule["post_arg_condition"].([]interface{}); len(pacs) > 0 {
		for _, pac := range pacs {
			conditions = append(conditions, *deliveryruleconditions.ExpandArmCdnEndpointConditionPostArg(pac.(map[string]interface{})))
		}
	}

	if qscs := rule["query_string_condition"].([]interface{}); len(qscs) > 0 {
		for _, qsc := range qscs {
			conditions = append(conditions, *deliveryruleconditions.ExpandArmCdnEndpointConditionQueryString(qsc.(map[string]interface{})))
		}
	}

	if racs := rule["remote_address_condition"].([]interface{}); len(racs) > 0 {
		for _, rac := range racs {
			conditions = append(conditions, *deliveryruleconditions.ExpandArmCdnEndpointConditionRemoteAddress(rac.(map[string]interface{})))
		}
	}

	if rbcs := rule["request_body_condition"].([]interface{}); len(rbcs) > 0 {
		for _, rbc := range rbcs {
			conditions = append(conditions, *deliveryruleconditions.ExpandArmCdnEndpointConditionRequestBody(rbc.(map[string]interface{})))
		}
	}

	if rhcs := rule["request_header_condition"].([]interface{}); len(rhcs) > 0 {
		for _, rhc := range rhcs {
			conditions = append(conditions, *deliveryruleconditions.ExpandArmCdnEndpointConditionRequestHeader(rhc.(map[string]interface{})))
		}
	}

	if rsc := rule["request_method_condition"].([]interface{}); len(rsc) > 0 {
		conditions = append(conditions, *deliveryruleconditions.ExpandArmCdnEndpointConditionRequestMethod(rsc[0].(map[string]interface{})))
	}

	if rsc := rule["request_scheme_condition"].([]interface{}); len(rsc) > 0 {
		conditions = append(conditions, *deliveryruleconditions.ExpandArmCdnEndpointConditionRequestScheme(rsc[0].(map[string]interface{})))
	}

	if rucs := rule["request_uri_condition"].([]interface{}); len(rucs) > 0 {
		for _, ruc := range rucs {
			conditions = append(conditions, *deliveryruleconditions.ExpandArmCdnEndpointConditionRequestURI(ruc.(map[string]interface{})))
		}
	}

	if ufecs := rule["url_file_extension_condition"].([]interface{}); len(ufecs) > 0 {
		for _, ufec := range ufecs {
			conditions = append(conditions, *deliveryruleconditions.ExpandArmCdnEndpointConditionURLFileExtension(ufec.(map[string]interface{})))
		}
	}

	return conditions
}

func expandDeliveryRuleActions(rule map[string]interface{}) ([]cdn.BasicDeliveryRuleAction, error) {
	actions := make([]cdn.BasicDeliveryRuleAction, 0)

	if cea := rule["cache_expiration_action"].([]interface{}); len(cea) > 0 {
		action, err := deliveryruleactions.ExpandArmCdnEndpointActionCacheExpiration(cea[0].(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		actions = append(actions, *action)
	}

	if ckqsa := rule["cache_key_query_string_action"].([]interface{}); len(ckqsa) > 0 {
		action, err := deliveryruleactions.ExpandArmCdnEndpointActionCacheKeyQueryString(ckqsa[0].(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		actions = append(actions, *action)
	}

	if mrha := rule["modify_request_header_action"].([]interface{}); len(mrha) > 0 {
		for _, rawAction := range mrha {
			actions = append(actions, *deliveryruleactions.ExpandArmCdnEndpointActionModifyRequestHeader(rawAction.(map[string]interface{})))
		}
	}

	if mrha := rule["modify_response_header_action"].([]interface{}); len(mrha) > 0 {
		for _, rawAction := range mrha {
			actions = append(actions, *deliveryruleactions.ExpandArmCdnEndpointActionModifyResponseHeader(rawAction.(map[string]interface{})))
		}
	}

	if ura := rule["url_redirect_action"].([]interface{}); len(ura) > 0 {
		actions = append(actions, *deliveryruleactions.ExpandArmCdnEndpointActionUrlRedirect(ura[0].(map[string]interface{})))
	}

	if ura := rule["url_rewrite_action"].([]interface{}); len(ura) > 0 {
		actions = append(actions, *deliveryruleactions.ExpandArmCdnEndpointActionURLRewrite(ura[0].(map[string]interface{})))
	}

	return actions, nil
}

func flattenArmCdnEndpointDeliveryRule(deliveryRule *cdn.DeliveryRule) map[string]interface{} {
	res := make(map[string]interface{}, 0)

	if deliveryRule == nil {
		return res
	}

	if deliveryRule.Name != nil {
		res["name"] = *deliveryRule.Name
	}

	if deliveryRule.Order != nil {
		res["order"] = *deliveryRule.Order
	}

	if deliveryRule.Conditions != nil {
		for _, basicDeliveryRuleCondition := range *deliveryRule.Conditions {
			if condition, isCookiesCondition := basicDeliveryRuleCondition.AsDeliveryRuleCookiesCondition(); isCookiesCondition {
				if _, ok := res["cookies_condition"]; !ok {
					res["cookies_condition"] = []map[string]interface{}{}
				}

				res["cookies_condition"] = append(res["cookies_condition"].([]map[string]interface{}), deliveryruleconditions.FlattenArmCdnEndpointConditionCookies(condition))
				continue
			}

			if condition, isHTTPVersionCondition := basicDeliveryRuleCondition.AsDeliveryRuleHTTPVersionCondition(); isHTTPVersionCondition {
				if _, ok := res["http_version_condition"]; !ok {
					res["http_version_condition"] = []map[string]interface{}{}
				}

				res["http_version_condition"] = append(res["http_version_condition"].([]map[string]interface{}), deliveryruleconditions.FlattenArmCdnEndpointConditionHTTPVersion(condition))
				continue
			}

			if condition, isDeviceCondition := basicDeliveryRuleCondition.AsDeliveryRuleIsDeviceCondition(); isDeviceCondition {
				res["device_condition"] = []interface{}{deliveryruleconditions.FlattenArmCdnEndpointConditionDevice(condition)}
				continue
			}

			if condition, isPostArgCondition := basicDeliveryRuleCondition.AsDeliveryRulePostArgsCondition(); isPostArgCondition {
				if _, ok := res["post_arg_condition"]; !ok {
					res["post_arg_condition"] = []map[string]interface{}{}
				}

				res["post_arg_condition"] = append(res["post_arg_condition"].([]map[string]interface{}), deliveryruleconditions.FlattenArmCdnEndpointConditionPostArg(condition))
				continue
			}

			if condition, isQueryStringCondition := basicDeliveryRuleCondition.AsDeliveryRuleQueryStringCondition(); isQueryStringCondition {
				if _, ok := res["query_string_condition"]; !ok {
					res["query_string_condition"] = []map[string]interface{}{}
				}

				res["query_string_condition"] = append(res["query_string_condition"].([]map[string]interface{}), deliveryruleconditions.FlattenArmCdnEndpointConditionQueryString(condition))
				continue
			}

			if condition, isRemoteAddressCondition := basicDeliveryRuleCondition.AsDeliveryRuleRemoteAddressCondition(); isRemoteAddressCondition {
				if _, ok := res["remote_address_condition"]; !ok {
					res["remote_address_condition"] = []map[string]interface{}{}
				}

				res["remote_address_condition"] = append(res["remote_address_condition"].([]map[string]interface{}), deliveryruleconditions.FlattenArmCdnEndpointConditionRemoteAddress(condition))
				continue
			}

			if condition, isRequestBodyCondition := basicDeliveryRuleCondition.AsDeliveryRuleRequestBodyCondition(); isRequestBodyCondition {
				if _, ok := res["request_body_condition"]; !ok {
					res["request_body_condition"] = []map[string]interface{}{}
				}

				res["request_body_condition"] = append(res["request_body_condition"].([]map[string]interface{}), deliveryruleconditions.FlattenArmCdnEndpointConditionRequestBody(condition))
				continue
			}

			if condition, isRequestMethodCondition := basicDeliveryRuleCondition.AsDeliveryRuleRequestMethodCondition(); isRequestMethodCondition {
				res["request_method_condition"] = []interface{}{deliveryruleconditions.FlattenArmCdnEndpointConditionRequestMethod(condition)}
				continue
			}

			if condition, isRequestSchemeCondition := basicDeliveryRuleCondition.AsDeliveryRuleRequestSchemeCondition(); isRequestSchemeCondition {
				res["request_scheme_condition"] = []interface{}{deliveryruleconditions.FlattenArmCdnEndpointConditionRequestScheme(condition)}
				continue
			}

			if condition, isRequestURICondition := basicDeliveryRuleCondition.AsDeliveryRuleRequestURICondition(); isRequestURICondition {
				if _, ok := res["request_uri_condition"]; !ok {
					res["request_uri_condition"] = []map[string]interface{}{}
				}

				res["request_uri_condition"] = append(res["request_uri_condition"].([]map[string]interface{}), deliveryruleconditions.FlattenArmCdnEndpointConditionRequestURI(condition))
				continue
			}

			if condition, isRequestURICondition := basicDeliveryRuleCondition.AsDeliveryRuleURLFileExtensionCondition(); isRequestURICondition {
				if _, ok := res["url_file_extension_condition"]; !ok {
					res["url_file_extension_condition"] = []map[string]interface{}{}
				}

				res["url_file_extension_condition"] = append(res["url_file_extension_condition"].([]map[string]interface{}), deliveryruleconditions.FlattenArmCdnEndpointConditionURLFileExtension(condition))
				continue
			}
		}
	}

	if deliveryRule.Actions != nil {
		for _, basicDeliveryRuleAction := range *deliveryRule.Actions {
			if action, isCacheExpirationAction := basicDeliveryRuleAction.AsDeliveryRuleCacheExpirationAction(); isCacheExpirationAction {
				res["cache_expiration_action"] = []interface{}{deliveryruleactions.FlattenArmCdnEndpointActionCacheExpiration(action)}
				continue
			}

			if action, isCacheKeyQueryStringAction := basicDeliveryRuleAction.AsDeliveryRuleCacheKeyQueryStringAction(); isCacheKeyQueryStringAction {
				res["cache_key_query_string_action"] = []interface{}{deliveryruleactions.FlattenArmCdnEndpointActionCacheKeyQueryString(action)}
				continue
			}

			if action, isModifyRequestHeaderAction := basicDeliveryRuleAction.AsDeliveryRuleRequestHeaderAction(); isModifyRequestHeaderAction {
				res["modify_request_header_action"] = []interface{}{deliveryruleactions.FlattenArmCdnEndpointActionModifyRequestHeader(action)}
				continue
			}

			if action, isModifyResponseHeaderAction := basicDeliveryRuleAction.AsDeliveryRuleResponseHeaderAction(); isModifyResponseHeaderAction {
				res["modify_response_header_action"] = []interface{}{deliveryruleactions.FlattenArmCdnEndpointActionModifyResponseHeader(action)}
				continue
			}

			if action, isURLRedirectAction := basicDeliveryRuleAction.AsURLRedirectAction(); isURLRedirectAction {
				res["url_redirect_action"] = []interface{}{deliveryruleactions.FlattenArmCdnEndpointActionUrlRedirect(action)}
				continue
			}

			if action, isURLRewriteAction := basicDeliveryRuleAction.AsURLRewriteAction(); isURLRewriteAction {
				res["url_rewrite_action"] = []interface{}{deliveryruleactions.FlattenArmCdnEndpointActionURLRewrite(action)}
				continue
			}
		}
	}

	return res
}
