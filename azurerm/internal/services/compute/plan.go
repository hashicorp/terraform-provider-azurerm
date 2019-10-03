package compute

import (
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func PlanSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},

				"publisher": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},

				"product": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
			},
		},
	}
}

func ExpandPlan(input []interface{}) *compute.Plan {
	if len(input) == 0 {
		return nil
	}

	raw := input[0].(map[string]interface{})

	return &compute.Plan{
		Name:      utils.String(raw["name"].(string)),
		Product:   utils.String(raw["product"].(string)),
		Publisher: utils.String(raw["publisher"].(string)),
	}
}

func FlattenPlan(input *compute.Plan) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	name := ""
	if input.Name != nil {
		name = *input.Name
	}

	product := ""
	if input.Product != nil {
		product = *input.Product
	}

	publisher := ""
	if input.Publisher != nil {
		publisher = *input.Publisher
	}

	return []interface{}{
		map[string]interface{}{
			"name":      name,
			"product":   product,
			"publisher": publisher,
		},
	}
}
