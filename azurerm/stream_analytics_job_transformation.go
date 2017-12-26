package azurerm

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func streamAnalyticsTransformationSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		MaxItems: 1,
		MinItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": &schema.Schema{
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringLenBetween(3, 64),
				},
				"streaming_units": &schema.Schema{
					Type:     schema.TypeInt,
					Optional: true,
					Default:  1,
				},
				"query": &schema.Schema{
					Type:     schema.TypeString,
					Required: true,
				},
			},
		},
	}
}
