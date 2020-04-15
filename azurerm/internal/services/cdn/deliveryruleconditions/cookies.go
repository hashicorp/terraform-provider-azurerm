package deliveryruleconditions

import (
	"fmt"

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
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(cdn.Lowercase),
						string(cdn.Uppercase),
					}, false),
				},
			},
		},
	}
}

func ExpandArmCdnEndpointConditionCookies(input []interface{}) []cdn.BasicDeliveryRuleCondition {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)
	for _, v := range input {
		item := v.(map[string]interface{})
		cookiesCondition := cdn.DeliveryRuleCookiesCondition{
			Name: cdn.NameCookies,
			Parameters: &cdn.CookiesMatchConditionParameters{
				OdataType:       utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleCookiesConditionParameters"),
				Selector:        utils.String(item["selector"].(string)),
				Operator:        cdn.CookiesOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].(*schema.Set).List()),
			},
		}

		if rawTransforms := item["transforms"].([]interface{}); len(rawTransforms) != 0 {
			transforms := make([]cdn.Transform, 0)
			for _, t := range rawTransforms {
				transforms = append(transforms, cdn.Transform(t.(string)))
			}
			cookiesCondition.Parameters.Transforms = &transforms
		}

		output = append(output, cookiesCondition)
	}

	return output
}

func FlattenArmCdnEndpointConditionCookies(input cdn.BasicDeliveryRuleCondition) (*map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleCookiesCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule cookie condition")
	}

	selector := ""
	operator := ""
	negateCondition := false
	matchValues := make([]interface{}, 0)
	transforms := make([]string, 0)

	if params := condition.Parameters; params != nil {
		if params.Selector != nil {
			selector = *params.Selector
		}

		operator = string(params.Operator)

		if params.NegateCondition != nil {
			negateCondition = *params.NegateCondition
		}

		if params.MatchValues != nil {
			matchValues = utils.FlattenStringSlice(params.MatchValues)
		}

		if params.Transforms != nil {
			for _, transform := range *params.Transforms {
				transforms = append(transforms, string(transform))
			}
		}
	}

	return &map[string]interface{}{
		"match_values":     schema.NewSet(schema.HashString, matchValues),
		"negate_condition": negateCondition,
		"operator":         operator,
		"selector":         selector,
		"transforms":       transforms,
	}, nil
}
