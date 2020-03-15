package deliveryruleconditions

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func Device() *schema.Resource {
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
						"Desktop",
						"Mobile",
					}, false),
				},
			},
		},
	}
}

func ExpandArmCdnEndpointConditionDevice(dc map[string]interface{}) *cdn.DeliveryRuleIsDeviceCondition {
	return &cdn.DeliveryRuleIsDeviceCondition{
		Name: cdn.NameHTTPVersion,
		Parameters: &cdn.IsDeviceMatchConditionParameters{
			OdataType:       utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleIsDeviceConditionParameters"),
			Operator:        utils.String(dc["operator"].(string)),
			NegateCondition: utils.Bool(dc["negate_condition"].(bool)),
			MatchValues:     utils.ExpandStringSlice(dc["match_values"].(*schema.Set).List()),
		},
	}
}

func FlattenArmCdnEndpointConditionDevice(dc *cdn.DeliveryRuleIsDeviceCondition) map[string]interface{} {
	res := make(map[string]interface{}, 1)

	if params := dc.Parameters; params != nil {
		if params.Operator != nil {
			res["operator"] = *params.Operator
		}

		if params.NegateCondition != nil {
			res["negate_condition"] = *params.NegateCondition
		}

		if params.MatchValues != nil {
			res["match_values"] = schema.NewSet(schema.HashString, utils.FlattenStringSlice(params.MatchValues))
		}
	}

	return res
}
