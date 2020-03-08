package delivery_rule_actions

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/cdn/mgmt/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func CacheExpiration() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"behavior": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.BypassCache),
					string(cdn.Override),
					string(cdn.SetIfMissing),
				}, false),
			},

			"duration": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.RuleActionCacheExpirationDuration(),
			},
		},
	}
}

func ExpandArmCdnEndpointActionCacheExpiration(cea map[string]interface{}) *cdn.DeliveryRuleCacheExpirationAction {
	cacheExpirationAction := cdn.DeliveryRuleCacheExpirationAction{
		Name: cdn.NameCacheExpiration,
		Parameters: &cdn.CacheExpirationActionParameters{
			OdataType:     utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleCacheExpirationActionParameters"),
			CacheBehavior: cdn.CacheBehavior(cea["behavior"].(string)),
			CacheType:     utils.String("All"),
		},
	}

	if duration := cea["duration"].(string); duration != "" {
		cacheExpirationAction.Parameters.CacheDuration = utils.String(duration)
	}

	return &cacheExpirationAction
}

func FlattenArmCdnEndpointActionCacheExpiration(cea *cdn.DeliveryRuleCacheExpirationAction) map[string]interface{} {
	res := make(map[string]interface{}, 1)

	if params := cea.Parameters; params != nil {
		res["behavior"] = string(params.CacheBehavior)

		if params.CacheDuration != nil {
			res["duration"] = *params.CacheDuration
		}
	}

	return res
}
