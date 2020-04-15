package deliveryruleactions

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/cdn/mgmt/2019-04-15/cdn"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func ModifyResponseHeader() *schema.Resource {
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

func ExpandArmCdnEndpointActionModifyResponseHeader(input []interface{}) (*[]cdn.BasicDeliveryRuleAction, error) {
	output := make([]cdn.BasicDeliveryRuleAction, 0)

	for _, v := range input {
		item := v.(map[string]interface{})

		requestHeaderAction := cdn.DeliveryRuleResponseHeaderAction{
			Name: cdn.NameModifyResponseHeader,
			Parameters: &cdn.HeaderActionParameters{
				OdataType:    utils.String("Microsoft.Azure.Cdn.Models.DeliveryRuleHeaderActionParameters"),
				HeaderAction: cdn.HeaderAction(item["action"].(string)),
				HeaderName:   utils.String(item["name"].(string)),
			},
		}

		if value := item["value"].(string); value != "" {
			requestHeaderAction.Parameters.Value = utils.String(value)
		}

		output = append(output, requestHeaderAction)
	}

	return &output, nil
}

func FlattenArmCdnEndpointActionModifyResponseHeader(input cdn.BasicDeliveryRuleAction) (*map[string]interface{}, error) {
	action, ok := input.AsDeliveryRuleResponseHeaderAction()
	if !ok {
		return nil, fmt.Errorf("expected a delivery rule response header action!")
	}

	headerAction := ""
	headerName := ""
	value := ""
	if params := action.Parameters; params != nil {
		headerAction = string(params.HeaderAction)
		if params.HeaderName != nil {
			headerName = *params.HeaderName
		}

		if params.Value != nil {
			value = *params.Value
		}
	}

	return &map[string]interface{}{
		"action": headerAction,
		"name":   headerName,
		"value":  value,
	}, nil
}
