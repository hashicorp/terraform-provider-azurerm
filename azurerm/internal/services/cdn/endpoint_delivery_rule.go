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
				Elem:     delivery_rule_conditions.RuleConditionRequestScheme(),
			},

			"url_redirect_action": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     delivery_rule_actions.RuleActionUrlRedirect(),
			},
		},
	}
}

func expandArmCdnEndpointDeliveryRule(rule map[string]interface{}) cdn.DeliveryRule {
	deliveryRule := cdn.DeliveryRule{
		Name:  utils.String(rule["name"].(string)),
		Order: utils.Int32(rule["order"].(int32)),
	}

	conditions := make([]cdn.BasicDeliveryRuleCondition, 0)

	if rsc, ok := rule["request_scheme_condition"]; ok {
		conditions = append(conditions, *delivery_rule_conditions.ExpandArmCdnEndpointConditionRequestScheme(rsc.([]interface{})[0].(map[string]interface{})))
	}

	deliveryRule.Conditions = &conditions

	actions := make([]cdn.BasicDeliveryRuleAction, 0)

	if ura, ok := rule["url_redirect_action"]; ok {
		actions = append(actions, *delivery_rule_actions.ExpandArmCdnEndpointActionUrlRedirect(ura.([]interface{})[0].(map[string]interface{})))
	}

	deliveryRule.Actions = &actions

	return deliveryRule
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
				res["request_scheme_condition"] = delivery_rule_conditions.FlattenArmCdnEndpointConditionRequestScheme(condition)
				continue
			}
		}
	}

	if deliveryRule.Actions != nil {
		for _, basicDeliveryRuleAction := range *deliveryRule.Actions {
			if action, isURLRedirectAction := basicDeliveryRuleAction.AsURLRedirectAction(); isURLRedirectAction {
				res["url_redirect_action"] = delivery_rule_actions.FlattenArmCdnEndpointActionUrlRedirect(action)
			}
		}
	}

	return res
}
