// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package deliveryruleactions

import (
	"errors"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func CacheKeyQueryString() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"behavior": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.QueryStringBehaviorExclude),
					string(cdn.QueryStringBehaviorExcludeAll),
					string(cdn.QueryStringBehaviorInclude),
					string(cdn.QueryStringBehaviorIncludeAll),
				}, false),
			},

			"parameters": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},
		},
	}
}

func ExpandArmCdnEndpointActionCacheKeyQueryString(input []interface{}) (*[]cdn.BasicDeliveryRuleAction, error) {
	output := make([]cdn.BasicDeliveryRuleAction, 0)

	for _, v := range input {
		item := v.(map[string]interface{})

		cacheKeyQueryStringAction := cdn.DeliveryRuleCacheKeyQueryStringAction{
			Name: cdn.NameBasicDeliveryRuleActionNameCacheKeyQueryString,
			Parameters: &cdn.CacheKeyQueryStringActionParameters{
				OdataType:           pointer.To("Microsoft.Azure.Cdn.Models.DeliveryRuleCacheKeyQueryStringBehaviorActionParameters"),
				QueryStringBehavior: cdn.QueryStringBehavior(item["behavior"].(string)),
			},
		}

		if parameters := item["parameters"].(string); parameters == "" {
			if behavior := cacheKeyQueryStringAction.Parameters.QueryStringBehavior; behavior == cdn.QueryStringBehaviorInclude || behavior == cdn.QueryStringBehaviorExclude {
				return nil, errors.New("parameters can not be empty if the behaviour is either Include or Exclude")
			}
		} else {
			cacheKeyQueryStringAction.Parameters.QueryParameters = pointer.To(parameters)
		}

		output = append(output, cacheKeyQueryStringAction)
	}

	return &output, nil
}

func FlattenArmCdnEndpointActionCacheKeyQueryString(input cdn.BasicDeliveryRuleAction) (*map[string]interface{}, error) {
	action, ok := input.AsDeliveryRuleCacheKeyQueryStringAction()
	if !ok {
		return nil, errors.New("expected a delivery rule cache key query string action")
	}

	behaviour := ""
	parameters := ""
	if params := action.Parameters; params != nil {
		behaviour = string(params.QueryStringBehavior)

		if params.QueryParameters != nil {
			parameters = *params.QueryParameters
		}
	}

	return &map[string]interface{}{
		"behavior":   behaviour,
		"parameters": parameters,
	}, nil
}
