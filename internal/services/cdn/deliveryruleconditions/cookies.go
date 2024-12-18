// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package deliveryruleconditions

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2024-02-01/rules"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func Cookies() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"selector": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},

			"operator": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(rules.PossibleValuesForCookiesOperator(),
					false),
			},

			"negate_condition": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"match_values": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotWhiteSpace,
				},
			},

			"transforms": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						string(rules.TransformLowercase),
						string(rules.TransformUppercase),
					}, false),
				},
			},
		},
	}
}

func ExpandArmCdnEndpointConditionCookies(input []interface{}) []rules.DeliveryRuleCondition {
	output := make([]rules.DeliveryRuleCondition, 0)
	for _, v := range input {
		item := v.(map[string]interface{})

		cookiesCondition := rules.DeliveryRuleCookiesCondition{
			Name: rules.MatchVariableCookies,
			Parameters: rules.CookiesMatchConditionParameters{
				TypeName:        rules.DeliveryRuleConditionParametersTypeDeliveryRuleCookiesConditionParameters,
				Selector:        pointer.To(item["selector"].(string)),
				Operator:        rules.CookiesOperator(item["operator"].(string)),
				NegateCondition: pointer.To(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].(*pluginsdk.Set).List()),
				Transforms:      expandTransforms(item["transforms"].([]interface{})),
			},
		}

		output = append(output, cookiesCondition)
	}

	return output
}

func FlattenArmCdnEndpointConditionCookies(input rules.DeliveryRuleCondition) (*map[string]interface{}, error) {
	condition, ok := AsDeliveryRuleCookiesCondition(input)
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule cookie condition")
	}

	selector := ""
	operator := ""
	negateCondition := false
	matchValues := make([]interface{}, 0)
	transforms := make([]string, 0)

	if params := condition; params != nil {
		operator = string(params.Operator)

		if params.Selector != nil {
			selector = *params.Selector
		}

		if params.NegateCondition != nil {
			negateCondition = *params.NegateCondition
		}

		if params.MatchValues != nil {
			matchValues = utils.FlattenStringSlice(params.MatchValues)
		}

		if params.Transforms != nil {
			transforms = flattenTransforms(params.Transforms)
		}
	}

	return &map[string]interface{}{
		"match_values":     pluginsdk.NewSet(pluginsdk.HashString, matchValues),
		"negate_condition": negateCondition,
		"operator":         operator,
		"selector":         selector,
		"transforms":       transforms,
	}, nil
}
