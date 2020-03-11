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
					ValidateFunc: validate.CdnEndpointDeliveryPolicyRuleName(),
				},

				"order": {
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntAtLeast(1),
				},

				"request_scheme_condition": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem:     deliveryruleconditions.RequestScheme(),
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

				"url_redirect_action": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem:     deliveryruleactions.URLRedirect(),
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

	conditions := make([]cdn.BasicDeliveryRuleCondition, 0)

	if rsc := rule["request_scheme_condition"].([]interface{}); len(rsc) > 0 {
		conditions = append(conditions, *deliveryruleconditions.ExpandArmCdnEndpointConditionRequestScheme(rsc[0].(map[string]interface{})))
	}

	deliveryRule.Conditions = &conditions

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

	if ura := rule["url_redirect_action"].([]interface{}); len(ura) > 0 {
		actions = append(actions, *deliveryruleactions.ExpandArmCdnEndpointActionUrlRedirect(ura[0].(map[string]interface{})))
	}

	deliveryRule.Actions = &actions

	return &deliveryRule, nil
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
			if condition, isRequestSchemeCondition := basicDeliveryRuleCondition.AsDeliveryRuleRequestSchemeCondition(); isRequestSchemeCondition {
				res["request_scheme_condition"] = []interface{}{deliveryruleconditions.FlattenArmCdnEndpointConditionRequestScheme(condition)}
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

			if action, isURLRedirectAction := basicDeliveryRuleAction.AsURLRedirectAction(); isURLRedirectAction {
				res["url_redirect_action"] = []interface{}{deliveryruleactions.FlattenArmCdnEndpointActionUrlRedirect(action)}
				continue
			}
		}
	}

	return res
}
