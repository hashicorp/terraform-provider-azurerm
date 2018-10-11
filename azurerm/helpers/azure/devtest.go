package azure

import (
	"github.com/Azure/azure-sdk-for-go/services/devtestlabs/mgmt/2016-05-15/dtl"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func SchemaDevTestVirtualMachineInboundNatRule() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		// since these aren't returned from the API
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"protocol": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(dtl.TCP),
						string(dtl.UDP),
					}, false),
				},

				"backend_port": {
					Type:         schema.TypeInt,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: validate.PortNumber,
				},

				"frontend_port": {
					Type:     schema.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func ExpandDevTestLabVirtualMachineNatRules(input *schema.Set) []dtl.InboundNatRule {
	rules := make([]dtl.InboundNatRule, 0)
	if input == nil {
		return rules
	}

	for _, val := range input.List() {
		v := val.(map[string]interface{})
		backendPort := v["backend_port"].(int)
		protocol := v["protocol"].(string)

		rule := dtl.InboundNatRule{
			TransportProtocol: dtl.TransportProtocol(protocol),
			BackendPort:       utils.Int32(int32(backendPort)),
		}

		rules = append(rules, rule)
	}

	return rules
}
