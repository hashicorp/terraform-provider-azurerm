package cdn

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/delivery_rule_actions"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/delivery_rule_conditions"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func EndpointDeliveryRule() *schema.Resource {
	return &schema.Resource{
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
				Elem:     delivery_rule_conditions.RequestScheme(),
			},

			"url_redirect_action": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     delivery_rule_actions.URLRedirect(),
			},

			"cache_expiration_action": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     delivery_rule_actions.CacheExpiration(),
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
		conditions = append(conditions, *delivery_rule_conditions.ExpandArmCdnEndpointConditionRequestScheme(rsc.([]interface{})[0].(map[string]interface{})))
	}

	deliveryRule.Conditions = &conditions

	actions := make([]cdn.BasicDeliveryRuleAction, 0)

	if ura := rule["url_redirect_action"]; len(ura.([]interface{})) > 0 {
		actions = append(actions, *delivery_rule_actions.ExpandArmCdnEndpointActionUrlRedirect(ura.([]interface{})[0].(map[string]interface{})))
	}

	if cea := rule["cache_expiration_action"]; len(cea.([]interface{})) > 0 {
		action, err := delivery_rule_actions.ExpandArmCdnEndpointActionCacheExpiration(cea.([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		actions = append(actions, *action)
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
				res["request_scheme_condition"] = []interface{}{delivery_rule_conditions.FlattenArmCdnEndpointConditionRequestScheme(condition)}
				continue
			}
		}
	}

	if deliveryRule.Actions != nil {
		for _, basicDeliveryRuleAction := range *deliveryRule.Actions {
			if action, isCacheExpirationAction := basicDeliveryRuleAction.AsDeliveryRuleCacheExpirationAction(); isCacheExpirationAction {
				res["cache_expiration_action"] = []interface{}{delivery_rule_actions.FlattenArmCdnEndpointActionCacheExpiration(action)}
				continue
			}

			if action, isURLRedirectAction := basicDeliveryRuleAction.AsURLRedirectAction(); isURLRedirectAction {
				res["url_redirect_action"] = []interface{}{delivery_rule_actions.FlattenArmCdnEndpointActionUrlRedirect(action)}
				continue
			}
		}
	}

	return res
}
