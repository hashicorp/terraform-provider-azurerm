package deliveryruleactions

import (
	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func ModifyRequestHeader() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(cdn.Append),
					string(cdn.Delete),
					string(cdn.Overwrite),
				}, false),
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func ExpandArmCdnEndpointActionModifyRequestHeader(mrha map[string]interface{}) *cdn.DeliveryRuleRequestHeaderAction {
	requestHeaderAction := cdn.DeliveryRuleRequestHeaderAction{
		Name: cdn.NameModifyRequestHeader,
		Parameters: &cdn.HeaderActionParameters{
			OdataType:    utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleHeaderActionParameters"),
			HeaderAction: cdn.HeaderAction(mrha["action"].(string)),
			HeaderName:   utils.String(mrha["name"].(string)),
		},
	}

	if value := mrha["value"].(string); value != "" {
		requestHeaderAction.Parameters.Value = utils.String(value)
	}

	return &requestHeaderAction
}

func FlattenArmCdnEndpointActionModifyRequestHeader(mrha *cdn.DeliveryRuleRequestHeaderAction) map[string]interface{} {
	res := make(map[string]interface{}, 1)

	if params := mrha.Parameters; params != nil {
		res["action"] = string(params.HeaderAction)
		res["name"] = *params.HeaderName

		if params.Value != nil {
			res["value"] = *params.Value
		}
	}

	return res
}
