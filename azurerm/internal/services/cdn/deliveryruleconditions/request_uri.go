package deliveryruleconditions

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func RequestURI() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"operator": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.RequestURIOperatorAny),
					string(cdn.RequestURIOperatorBeginsWith),
					string(cdn.RequestURIOperatorContains),
					string(cdn.RequestURIOperatorEndsWith),
					string(cdn.RequestURIOperatorEqual),
					string(cdn.RequestURIOperatorGreaterThan),
					string(cdn.RequestURIOperatorGreaterThanOrEqual),
					string(cdn.RequestURIOperatorLessThan),
					string(cdn.RequestURIOperatorLessThanOrEqual),
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

func ExpandArmCdnEndpointConditionRequestURI(ruc map[string]interface{}) *cdn.DeliveryRuleRequestURICondition {
	requestURICondition := cdn.DeliveryRuleRequestURICondition{
		Name: cdn.NameRequestURI,
		Parameters: &cdn.RequestURIMatchConditionParameters{
			OdataType:       utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleRequestUriConditionParameters"),
			Operator:        cdn.RequestURIOperator(ruc["operator"].(string)),
			NegateCondition: utils.Bool(ruc["negate_condition"].(bool)),
			MatchValues:     utils.ExpandStringSlice(ruc["match_values"].(*schema.Set).List()),
		},
	}

	if rawTransforms := ruc["transforms"].([]interface{}); len(rawTransforms) != 0 {
		transforms := make([]cdn.Transform, 0)
		for _, t := range rawTransforms {
			transforms = append(transforms, cdn.Transform(t.(string)))
		}
		requestURICondition.Parameters.Transforms = &transforms
	}

	return &requestURICondition
}

func FlattenArmCdnEndpointConditionRequestURI(ruc *cdn.DeliveryRuleRequestURICondition) map[string]interface{} {
	res := make(map[string]interface{}, 1)

	if params := ruc.Parameters; params != nil {
		res["operator"] = string(params.Operator)

		if params.NegateCondition != nil {
			res["negate_condition"] = *params.NegateCondition
		}

		if params.MatchValues != nil {
			res["match_values"] = schema.NewSet(schema.HashString, utils.FlattenStringSlice(params.MatchValues))
		}

		if params.Transforms != nil {
			transforms := make([]string, 0)
			for _, transform := range *params.Transforms {
				transforms = append(transforms, string(transform))
			}
			res["transforms"] = &transforms
		}
	}

	return res
}
