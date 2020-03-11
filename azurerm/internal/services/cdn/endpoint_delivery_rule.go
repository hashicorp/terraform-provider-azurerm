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

	if rsc := rule["request_scheme_condition"]; len(rsc.([]interface{})) > 0 {
		conditions = append(conditions, *deliveryruleconditions.ExpandArmCdnEndpointConditionRequestScheme(rsc.([]interface{})[0].(map[string]interface{})))
	}

	deliveryRule.Conditions = &conditions

	actions := make([]cdn.BasicDeliveryRuleAction, 0)

	if cea := rule["cache_expiration_action"]; len(cea.([]interface{})) > 0 {
		action, err := deliveryruleactions.ExpandArmCdnEndpointActionCacheExpiration(cea.([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		actions = append(actions, *action)
	}

	if ura := rule["url_redirect_action"]; len(ura.([]interface{})) > 0 {
		actions = append(actions, *deliveryruleactions.ExpandArmCdnEndpointActionUrlRedirect(ura.([]interface{})[0].(map[string]interface{})))
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

			if action, isURLRedirectAction := basicDeliveryRuleAction.AsURLRedirectAction(); isURLRedirectAction {
				res["url_redirect_action"] = []interface{}{deliveryruleactions.FlattenArmCdnEndpointActionUrlRedirect(action)}
				continue
			}
		}
	}

	return res
}
