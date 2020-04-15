package deliveryruleconditions

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func HTTPVersion() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"operator": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Equal",
				ValidateFunc: validation.StringInSlice([]string{
					"Equal",
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
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"0.9",
						"1.0",
						"1.1",
						"2.0",
					}, false),
				},
			},
		},
	}
}

func ExpandArmCdnEndpointConditionHTTPVersion(input []interface{}) []cdn.BasicDeliveryRuleCondition {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)

	for _, v := range input {
		item := v.(map[string]interface{})
		output = append(output, cdn.DeliveryRuleHTTPVersionCondition{
			Name: cdn.NameHTTPVersion,
			Parameters: &cdn.HTTPVersionMatchConditionParameters{
				OdataType:       utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleHttpVersionConditionParameters"),
				Operator:        utils.String(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].(*schema.Set).List()),
			},
		})
	}

	return output
}

func FlattenArmCdnEndpointConditionHTTPVersion(input cdn.BasicDeliveryRuleCondition) (*map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleHTTPVersionCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule http version condition")
	}

	operator := ""
	matchValues := make([]interface{}, 0)
	negateCondition := false
	if params := condition.Parameters; params != nil {
		if params.Operator != nil {
			operator = *params.Operator
		}

		if params.NegateCondition != nil {
			negateCondition = *params.NegateCondition
		}

		if params.MatchValues != nil {
			matchValues = utils.FlattenStringSlice(params.MatchValues)
		}
	}

	return &map[string]interface{}{
		"operator":         operator,
		"match_values":     schema.NewSet(schema.HashString, matchValues),
		"negate_condition": negateCondition,
	}, nil
}
