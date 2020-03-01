package cdn

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func EndpointDeliveryPolicy() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: false,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"description": {
					Type:     schema.TypeString,
					Required: false,
				},

				"rule": {
					Type:     schema.TypeList,
					Required: true,
					MinItems: 1,
					MaxItems: 4,
					Elem:     EndpointDeliveryRule(),
				},
			},
		},
	}
}
