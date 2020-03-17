package deliveryruleconditions

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func URLFileName() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"operator": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.URLFileNameOperatorAny),
					string(cdn.URLFileNameOperatorBeginsWith),
					string(cdn.URLFileNameOperatorContains),
					string(cdn.URLFileNameOperatorEndsWith),
					string(cdn.URLFileNameOperatorEqual),
					string(cdn.URLFileNameOperatorGreaterThan),
					string(cdn.URLFileNameOperatorGreaterThanOrEqual),
					string(cdn.URLFileNameOperatorLessThan),
					string(cdn.URLFileNameOperatorLessThanOrEqual),
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

func ExpandArmCdnEndpointConditionURLFileName(ufnc map[string]interface{}) *cdn.DeliveryRuleURLFileNameCondition {
	requestURICondition := cdn.DeliveryRuleURLFileNameCondition{
		Name: cdn.NameURLFileName,
		Parameters: &cdn.URLFileNameMatchConditionParameters{
			OdataType:       utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleUrlFilenameConditionParameters"),
			Operator:        cdn.URLFileNameOperator(ufnc["operator"].(string)),
			NegateCondition: utils.Bool(ufnc["negate_condition"].(bool)),
			MatchValues:     utils.ExpandStringSlice(ufnc["match_values"].(*schema.Set).List()),
		},
	}

	if rawTransforms := ufnc["transforms"].([]interface{}); len(rawTransforms) != 0 {
		transforms := make([]cdn.Transform, 0)
		for _, t := range rawTransforms {
			transforms = append(transforms, cdn.Transform(t.(string)))
		}
		requestURICondition.Parameters.Transforms = &transforms
	}

	return &requestURICondition
}

func FlattenArmCdnEndpointConditionURLFileName(ufnc *cdn.DeliveryRuleURLFileNameCondition) map[string]interface{} {
	res := make(map[string]interface{}, 1)

	if params := ufnc.Parameters; params != nil {
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
