package deliveryruleactions

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func CacheKeyQueryString() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"behavior": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.Exclude),
					string(cdn.ExcludeAll),
					string(cdn.Include),
					string(cdn.IncludeAll),
				}, false),
			},

			"parameters": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func ExpandArmCdnEndpointActionCacheKeyQueryString(ckqsa map[string]interface{}) (*cdn.DeliveryRuleCacheKeyQueryStringAction, error) {
	cacheKeyQueryStringAction := cdn.DeliveryRuleCacheKeyQueryStringAction{
		Name: cdn.NameCacheKeyQueryString,
		Parameters: &cdn.CacheKeyQueryStringActionParameters{
			OdataType:           utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleCacheKeyQueryStringBehaviorActionParameters"),
			QueryStringBehavior: cdn.QueryStringBehavior(ckqsa["behavior"].(string)),
		},
	}

	if parameters := ckqsa["parameters"].(string); parameters == "" {
		if behavior := cacheKeyQueryStringAction.Parameters.QueryStringBehavior; behavior == cdn.Include || behavior == cdn.Exclude {
			return nil, fmt.Errorf("Parameters can not be empty if the behavior is either Include or Exclude.")
		}
	} else {
		cacheKeyQueryStringAction.Parameters.QueryParameters = utils.String(parameters)
	}

	return &cacheKeyQueryStringAction, nil
}

func FlattenArmCdnEndpointActionCacheKeyQueryString(ckqsa *cdn.DeliveryRuleCacheKeyQueryStringAction) map[string]interface{} {
	res := make(map[string]interface{}, 1)

	if params := ckqsa.Parameters; params != nil {
		res["behavior"] = string(params.QueryStringBehavior)

		if params.QueryParameters != nil {
			res["parameters"] = *params.QueryParameters
		}
	}

	return res
}
