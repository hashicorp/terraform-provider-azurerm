// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package deliveryruleactions

import (
	"errors"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2020-09-01/cdn" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func CacheExpiration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"behavior": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.CacheBehaviorBypassCache),
					string(cdn.CacheBehaviorOverride),
					string(cdn.CacheBehaviorSetIfMissing),
				}, false),
			},

			"duration": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validate.RuleActionCacheExpirationDuration(),
			},
		},
	}
}

func ExpandArmCdnEndpointActionCacheExpiration(input []interface{}) (*[]cdn.BasicDeliveryRuleAction, error) {
	output := make([]cdn.BasicDeliveryRuleAction, 0)

	for _, v := range input {
		item := v.(map[string]interface{})

		cacheExpirationAction := cdn.DeliveryRuleCacheExpirationAction{
			Name: cdn.NameBasicDeliveryRuleActionNameCacheExpiration,
			Parameters: &cdn.CacheExpirationActionParameters{
				OdataType:     pointer.To("Microsoft.Azure.Cdn.Models.DeliveryRuleCacheExpirationActionParameters"),
				CacheBehavior: cdn.CacheBehavior(item["behavior"].(string)),
				CacheType:     pointer.To("All"),
			},
		}

		if duration := item["duration"].(string); duration != "" {
			if cacheExpirationAction.Parameters.CacheBehavior == cdn.CacheBehaviorBypassCache {
				return nil, errors.New("cache expiration duration must not be set when using behavior `BypassCache`")
			}

			cacheExpirationAction.Parameters.CacheDuration = pointer.To(duration)
		}

		output = append(output, cacheExpirationAction)
	}

	return &output, nil
}

func FlattenArmCdnEndpointActionCacheExpiration(input cdn.BasicDeliveryRuleAction) (*map[string]interface{}, error) {
	action, ok := input.AsDeliveryRuleCacheExpirationAction()
	if !ok {
		return nil, errors.New("expected a delivery rule cache expiration action")
	}

	behaviour := ""
	duration := ""
	if params := action.Parameters; params != nil {
		behaviour = string(params.CacheBehavior)

		if params.CacheDuration != nil {
			duration = *params.CacheDuration
		}
	}

	return &map[string]interface{}{
		"behavior": behaviour,
		"duration": duration,
	}, nil
}
