package deliveryruleactions

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/cdn/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func CacheExpiration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"behavior": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.BypassCache),
					string(cdn.Override),
					string(cdn.SetIfMissing),
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
			Name: cdn.NameCacheExpiration,
			Parameters: &cdn.CacheExpirationActionParameters{
				OdataType:     utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleCacheExpirationActionParameters"),
				CacheBehavior: cdn.CacheBehavior(item["behavior"].(string)),
				CacheType:     utils.String("All"),
			},
		}

		if duration := item["duration"].(string); duration != "" {
			if cacheExpirationAction.Parameters.CacheBehavior == cdn.BypassCache {
				return nil, fmt.Errorf("Cache expiration duration must not be set when using behavior `BypassCache`")
			}

			cacheExpirationAction.Parameters.CacheDuration = utils.String(duration)
		}

		output = append(output, cacheExpirationAction)
	}

	return &output, nil
}

func FlattenArmCdnEndpointActionCacheExpiration(input cdn.BasicDeliveryRuleAction) (*map[string]interface{}, error) {
	action, ok := input.AsDeliveryRuleCacheExpirationAction()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule cache expiration action!")
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
