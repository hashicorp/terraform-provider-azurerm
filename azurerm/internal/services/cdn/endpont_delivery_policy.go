package cdn

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func EndpointDeliveryPolicy() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: false,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"description": {
					Type:     schema.TypeString,
					Required: false,
				},

				"rule": {
					Type:     schema.TypeList,
					Required: true,
					MinItems: 1,
					MaxItems: 4,
					Elem:     EndpointDeliveryRule(),
				},
			},
		},
	}
}

func expandArmCdnEndpointDeliveryPolicy(d *schema.ResourceData) *cdn.EndpointPropertiesUpdateParametersDeliveryPolicy {
	policies := d.Get("delivery_policy").([]interface{})
	if len(policies) == 0 {
		return nil
	}

	deliveryPolicy := cdn.EndpointPropertiesUpdateParametersDeliveryPolicy{}

	policy := policies[0].(map[string]interface{})
	if descr, ok := policy["description"]; ok {
		deliveryPolicy.Description = utils.String(descr.(string))
	}

	rules := policy["rule"].([]interface{})
	deliveryRules := make([]cdn.DeliveryRule, len(rules))
	for i, rule := range rules {
		deliveryRules[i] = expandArmCdnEndpointDeliveryRule(rule.(map[string]interface{}))
	}
	deliveryPolicy.Rules = &deliveryRules

	return &deliveryPolicy
}

func flattenArmCdnEndpointDeliveryPolicy(deliveryPolicy *cdn.EndpointPropertiesUpdateParametersDeliveryPolicy) []interface{} {
	deliveryPolicies := make([]interface{}, 0)

	if deliveryPolicy == nil {
		return deliveryPolicies
	}

	dp := make(map[string]interface{})
	if deliveryPolicy.Description != nil {
		dp["description"] = *deliveryPolicy.Description
	}

	if deliveryPolicy.Rules != nil && len(*deliveryPolicy.Rules) > 0 {
		rules := make([]map[string]interface{}, len(*deliveryPolicy.Rules))
		for i, rule := range *deliveryPolicy.Rules {
			rules[i] = flattenArmCdnEndpointDeliveryRule(&rule)
		}
	}

	return deliveryPolicies
}
