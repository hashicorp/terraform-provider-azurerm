package deliveryruleconditions

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func Cookies() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"selector": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},

			"operator": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.Any),
					string(cdn.BeginsWith),
					string(cdn.Contains),
					string(cdn.EndsWith),
					string(cdn.Equal),
					string(cdn.GreaterThan),
					string(cdn.GreaterThanOrEqual),
					string(cdn.LessThan),
					string(cdn.LessThanOrEqual),
				}, false),
			},

			"negate_condition": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"match_values": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
			},

			"transforms": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.Lowercase),
					string(cdn.Uppercase),
				}, false),
			},
		},
	}
}

func ExpandArmCdnEndpointConditionCookies(cc map[string]interface{}) *cdn.DeliveryRuleCookiesCondition {
	cookiesCondition := cdn.DeliveryRuleCookiesCondition{
		Name: cdn.NameCookies,
		Parameters: &cdn.CookiesMatchConditionParameters{
			OdataType:       utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleCookiesConditionParameters"),
			Selector:        utils.String(cc["selector"].(string)),
			Operator:        cdn.CookiesOperator(cc["operator"].(string)),
			NegateCondition: utils.Bool(cc["negate_condition"].(bool)),
			MatchValues:     utils.ExpandStringSlice(cc["match_values"].(*schema.Set).List()),
		},
	}

	if transform := cc["transforms"].(string); transform != "" {
		transforms := []cdn.Transform{cdn.Transform(transform)}
		cookiesCondition.Parameters.Transforms = &transforms
	}

	return &cookiesCondition
}

func FlattenArmCdnEndpointConditionCookies(cc *cdn.DeliveryRuleCookiesCondition) map[string]interface{} {
	res := make(map[string]interface{}, 1)

	if params := cc.Parameters; params != nil {
		if params.Selector != nil {
			res["selector"] = *params.Selector
		}

		res["operator"] = params.Operator

		if params.NegateCondition != nil {
			res["negate_condition"] = *params.NegateCondition
		}

		if params.MatchValues != nil {
			res["match_values"] = schema.NewSet(schema.HashString, utils.FlattenStringSlice(params.MatchValues))
		}

		if params.Transforms != nil && len(*params.Transforms) > 0 {
			transforms := make([]string, 0)
			for _, transform := range *params.Transforms {
				transforms = append(transforms, string(transform))
			}
			res["transforms"] = &transforms
		}
	}

	return res
}
