// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_endpoint

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/deliveryruleactions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func endpointGlobalDeliveryRule() *pluginsdk.Schema {
	// lintignore:XS003
	return &pluginsdk.Schema{
		Description: "Actions that are valid for all resources regardless of any conditions.",
		Type:        pluginsdk.TypeList,
		Optional:    true,
		MaxItems:    1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"cache_expiration_action": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					AtLeastOneOf: []string{
						"global_delivery_rule.0.cache_expiration_action",
						"global_delivery_rule.0.cache_key_query_string_action",
						"global_delivery_rule.0.modify_request_header_action",
						"global_delivery_rule.0.modify_response_header_action",
						"global_delivery_rule.0.url_redirect_action",
						"global_delivery_rule.0.url_rewrite_action",
					},
					Elem: deliveryruleactions.CacheExpiration(),
				},

				"cache_key_query_string_action": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					AtLeastOneOf: []string{
						"global_delivery_rule.0.cache_expiration_action",
						"global_delivery_rule.0.cache_key_query_string_action",
						"global_delivery_rule.0.modify_request_header_action",
						"global_delivery_rule.0.modify_response_header_action",
						"global_delivery_rule.0.url_redirect_action",
						"global_delivery_rule.0.url_rewrite_action",
					},
					Elem: deliveryruleactions.CacheKeyQueryString(),
				},

				"modify_request_header_action": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					AtLeastOneOf: []string{
						"global_delivery_rule.0.cache_expiration_action",
						"global_delivery_rule.0.cache_key_query_string_action",
						"global_delivery_rule.0.modify_request_header_action",
						"global_delivery_rule.0.modify_response_header_action",
						"global_delivery_rule.0.url_redirect_action",
						"global_delivery_rule.0.url_rewrite_action",
					},
					Elem: deliveryruleactions.ModifyRequestHeader(),
				},

				"modify_response_header_action": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					AtLeastOneOf: []string{
						"global_delivery_rule.0.cache_expiration_action",
						"global_delivery_rule.0.cache_key_query_string_action",
						"global_delivery_rule.0.modify_request_header_action",
						"global_delivery_rule.0.modify_response_header_action",
						"global_delivery_rule.0.url_redirect_action",
						"global_delivery_rule.0.url_rewrite_action",
					},
					Elem: deliveryruleactions.ModifyResponseHeader(),
				},

				"url_redirect_action": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					AtLeastOneOf: []string{
						"global_delivery_rule.0.cache_expiration_action",
						"global_delivery_rule.0.cache_key_query_string_action",
						"global_delivery_rule.0.modify_request_header_action",
						"global_delivery_rule.0.modify_response_header_action",
						"global_delivery_rule.0.url_redirect_action",
						"global_delivery_rule.0.url_rewrite_action",
					},
					Elem: deliveryruleactions.URLRedirect(),
				},

				"url_rewrite_action": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MaxItems: 1,
					AtLeastOneOf: []string{
						"global_delivery_rule.0.cache_expiration_action",
						"global_delivery_rule.0.cache_key_query_string_action",
						"global_delivery_rule.0.modify_request_header_action",
						"global_delivery_rule.0.modify_response_header_action",
						"global_delivery_rule.0.url_redirect_action",
						"global_delivery_rule.0.url_rewrite_action",
					},
					Elem: deliveryruleactions.URLRewrite(),
				},
			},
		},
	}
}

func expandArmCdnEndpointGlobalDeliveryRule(rule map[string]any) (*cdn.DeliveryRule, error) {
	deliveryRule := cdn.DeliveryRule{
		Name:  pointer.To("Global"),
		Order: pointer.To(int32(0)),
	}

	actions, err := helpers.ExpandDeliveryRuleActions(rule)
	if err != nil {
		return nil, err
	}
	deliveryRule.Actions = &actions

	return &deliveryRule, nil
}

func flattenArmCdnEndpointGlobalDeliveryRule(deliveryRule cdn.DeliveryRule) (*map[string]any, error) {
	actions, err := helpers.FlattenDeliveryRuleActions(deliveryRule.Actions)
	if err != nil {
		return nil, err
	}

	output := make(map[string]any)
	for key, value := range *actions {
		output[key] = value
	}
	return &output, nil
}
