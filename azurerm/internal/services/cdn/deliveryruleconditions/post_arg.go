package deliveryruleconditions

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func PostArg() *schema.Resource {
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
					string(cdn.PostArgsOperatorAny),
					string(cdn.PostArgsOperatorBeginsWith),
					string(cdn.PostArgsOperatorContains),
					string(cdn.PostArgsOperatorEndsWith),
					string(cdn.PostArgsOperatorEqual),
					string(cdn.PostArgsOperatorGreaterThan),
					string(cdn.PostArgsOperatorGreaterThanOrEqual),
					string(cdn.PostArgsOperatorLessThan),
					string(cdn.PostArgsOperatorLessThanOrEqual),
				}, false),
			},

			"negate_condition": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"match_values": {
				Type:     schema.TypeSet,
				Optional: true,
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

func ExpandArmCdnEndpointConditionPostArg(input []interface{}) []cdn.BasicDeliveryRuleCondition {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)

	for _, v := range input {
		item := v.(map[string]interface{})

		condition := cdn.DeliveryRulePostArgsCondition{
			Name: cdn.NameCookies,
			Parameters: &cdn.PostArgsMatchConditionParameters{
				OdataType:       utils.String("Microsoft.Azure.Cdn.Models.DeliveryRulePostArgsConditionParameters"),
				Selector:        utils.String(item["selector"].(string)),
				Operator:        cdn.PostArgsOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].(*schema.Set).List()),
			},
		}

		if rawTransforms := item["transforms"].([]interface{}); len(rawTransforms) != 0 {
			transforms := make([]cdn.Transform, 0)
			for _, t := range rawTransforms {
				transforms = append(transforms, cdn.Transform(t.(string)))
			}
			condition.Parameters.Transforms = &transforms
		}

		output = append(output, condition)
	}

	return output
}

func FlattenArmCdnEndpointConditionPostArg(input cdn.BasicDeliveryRuleCondition) (*map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRulePostArgsCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule post args condition")
	}

	operator := ""
	matchValues := make([]interface{}, 0)
	negateCondition := false
	selector := ""
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
		"operator":         operator,
		"match_values":     schema.NewSet(schema.HashString, matchValues),
		"negate_condition": negateCondition,
		"selector":         selector,
		"transforms":       transforms,
	}, nil
}
