package cdn

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/deliveryruleactions"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func endpointGlobalDeliveryRule() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
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

func expandArmCdnEndpointGlobalDeliveryRule(rule map[string]interface{}) (*cdn.DeliveryRule, error) {
	deliveryRule := cdn.DeliveryRule{
		Name:  utils.String("Global"),
		Order: utils.Int32(0),
	}

	actions, err := expandDeliveryRuleActions(rule)
	if err != nil {
		return nil, err
	}
	deliveryRule.Actions = &actions

	return &deliveryRule, nil
}

func flattenArmCdnEndpointGlobalDeliveryRule(deliveryRule cdn.DeliveryRule) (*map[string]interface{}, error) {
	actions, err := flattenDeliveryRuleActions(deliveryRule.Actions)
	if err != nil {
		return nil, err
	}

	output := make(map[string]interface{})
	for key, value := range *actions {
		output[key] = value
	}
	return &output, nil
}
