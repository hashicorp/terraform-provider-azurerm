package deliveryruleconditions

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func RequestScheme() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"operator": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "Equal",
				ValidateFunc: validation.StringInSlice([]string{
					"Equal",
				}, false),
			},

			"negate_condition": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"match_values": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"HTTP",
						"HTTPS",
					}, false),
				},
			},
		},
	}
}

func ExpandArmCdnEndpointConditionRequestScheme(input []interface{}) []cdn.BasicDeliveryRuleCondition {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)

	for _, v := range input {
		item := v.(map[string]interface{})

		requestSchemeCondition := cdn.DeliveryRuleRequestSchemeCondition{
			Name: cdn.NameRequestScheme,
			Parameters: &cdn.RequestSchemeMatchConditionParameters{
				OdataType:       utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleRequestSchemeConditionParameters"),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].(*pluginsdk.Set).List()),
			},
		}

		if operator := item["operator"]; operator.(string) != "" {
			requestSchemeCondition.Parameters.Operator = utils.String(operator.(string))
		}

		output = append(output, requestSchemeCondition)
	}

	return output
}

func FlattenArmCdnEndpointConditionRequestScheme(input cdn.BasicDeliveryRuleCondition) (*map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleRequestSchemeCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule request scheme condition")
	}

	operator := ""
	negateCondition := false
	matchValues := make([]interface{}, 0)
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
		"match_values":     pluginsdk.NewSet(pluginsdk.HashString, matchValues),
		"negate_condition": negateCondition,
	}, nil
}
