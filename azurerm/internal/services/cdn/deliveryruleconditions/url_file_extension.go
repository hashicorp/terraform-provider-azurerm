package deliveryruleconditions

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func URLFileExtension() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"operator": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.URLFileExtensionOperatorAny),
					string(cdn.URLFileExtensionOperatorBeginsWith),
					string(cdn.URLFileExtensionOperatorContains),
					string(cdn.URLFileExtensionOperatorEndsWith),
					string(cdn.URLFileExtensionOperatorEqual),
					string(cdn.URLFileExtensionOperatorGreaterThan),
					string(cdn.URLFileExtensionOperatorGreaterThanOrEqual),
					string(cdn.URLFileExtensionOperatorLessThan),
					string(cdn.URLFileExtensionOperatorLessThanOrEqual),
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

func ExpandArmCdnEndpointConditionURLFileExtension(input []interface{}) []cdn.BasicDeliveryRuleCondition {
	output := make([]cdn.BasicDeliveryRuleCondition, 0)

	for _, v := range input {
		item := v.(map[string]interface{})

		requestURICondition := cdn.DeliveryRuleURLFileExtensionCondition{
			Name: cdn.NameURLFileExtension,
			Parameters: &cdn.URLFileExtensionMatchConditionParameters{
				OdataType:       utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleUrlFileExtensionMatchConditionParameters"),
				Operator:        cdn.URLFileExtensionOperator(item["operator"].(string)),
				NegateCondition: utils.Bool(item["negate_condition"].(bool)),
				MatchValues:     utils.ExpandStringSlice(item["match_values"].(*schema.Set).List()),
			},
		}

		if rawTransforms := item["transforms"].([]interface{}); len(rawTransforms) != 0 {
			transforms := make([]cdn.Transform, 0)
			for _, t := range rawTransforms {
				transforms = append(transforms, cdn.Transform(t.(string)))
			}
			requestURICondition.Parameters.Transforms = &transforms
		}

		output = append(output, requestURICondition)
	}

	return output
}

func FlattenArmCdnEndpointConditionURLFileExtension(input cdn.BasicDeliveryRuleCondition) (*map[string]interface{}, error) {
	condition, ok := input.AsDeliveryRuleURLFileExtensionCondition()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule url file extension condition")
	}

	matchValues := make([]interface{}, 0)
	negateCondition := false
	operator := ""
	transforms := make([]string, 0)
	if params := condition.Parameters; params != nil {
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
		"transforms":       transforms,
	}, nil
}
