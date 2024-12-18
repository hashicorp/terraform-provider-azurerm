// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package deliveryruleconditions

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func Device() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"operator": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(rules.IsDeviceOperatorEqual),
				ValidateFunc: validation.StringInSlice(rules.PossibleValuesForIsDeviceOperator(),
					false),
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
					ValidateFunc: validation.StringInSlice(rules.PossibleValuesForIsDeviceMatchValue(),
						false),
				},
			},
		},
	}
}

func ExpandArmCdnEndpointConditionDevice(input []interface{}) []rules.DeliveryRuleCondition {
	output := make([]rules.DeliveryRuleCondition, 0)

	for _, v := range input {
		item := v.(map[string]interface{})

		output = append(output, rules.DeliveryRuleIsDeviceCondition{
			Name: rules.MatchVariableIsDevice,
			Parameters: rules.IsDeviceMatchConditionParameters{
				TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleIsDeviceConditionParameters,
				Operator:        rules.IsDeviceOperator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     expandIsDeviceMatchValue(item["match_values"].(*pluginsdk.Set).List()),
			},
		})
	}

	return output
}

func FlattenArmCdnEndpointConditionDevice(input rules.DeliveryRuleCondition) (*map[string]interface{}, error) {
	condition, ok := AsDeliveryRuleIsDeviceCondition(input)
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule device condition")
	}

	operator := ""
	matchValues := make([]interface{}, 0)
	negateCondition := false

	if params := condition; params != nil {
		operator = string(params.Operator)

		if params.NegateCondition != nil {
			negateCondition = *params.NegateCondition
		}

		if params.MatchValues != nil {
			matchValues = flattenIsDeviceMatchValue(params.MatchValues)
		}
	}

	return &map[string]interface{}{
		"operator":         operator,
		"match_values":     pluginsdk.NewSet(pluginsdk.HashString, matchValues),
		"negate_condition": negateCondition,
	}, nil
}
