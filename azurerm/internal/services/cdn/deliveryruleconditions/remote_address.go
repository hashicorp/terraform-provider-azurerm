package deliveryruleconditions

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func RemoteAddress() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"operator": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.RemoteAddressOperatorAny),
					string(cdn.RemoteAddressOperatorGeoMatch),
					string(cdn.RemoteAddressOperatorIPMatch),
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
		},
	}
}

func ExpandArmCdnEndpointConditionRemoteAddress(qsc map[string]interface{}) *cdn.DeliveryRuleRemoteAddressCondition {
	remoteAddressCondition := cdn.DeliveryRuleRemoteAddressCondition{
		Name: cdn.NameRemoteAddress,
		Parameters: &cdn.RemoteAddressMatchConditionParameters{
			OdataType:       utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleRemoteAddressConditionParameters"),
			Operator:        cdn.RemoteAddressOperator(qsc["operator"].(string)),
			NegateCondition: utils.Bool(qsc["negate_condition"].(bool)),
			MatchValues:     utils.ExpandStringSlice(qsc["match_values"].(*schema.Set).List()),
		},
	}

	return &remoteAddressCondition
}

func FlattenArmCdnEndpointConditionRemoteAddress(qsc *cdn.DeliveryRuleRemoteAddressCondition) map[string]interface{} {
	res := make(map[string]interface{}, 1)

	if params := qsc.Parameters; params != nil {
		res["operator"] = string(params.Operator)

		if params.NegateCondition != nil {
			res["negate_condition"] = *params.NegateCondition
		}

		if params.MatchValues != nil {
			res["match_values"] = schema.NewSet(schema.HashString, utils.FlattenStringSlice(params.MatchValues))
		}
	}

	return res
}
