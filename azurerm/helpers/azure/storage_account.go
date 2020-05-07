package azure

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func SchemaStorageAccountCorsRule() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 5,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"allowed_origins": {
					Type:     schema.TypeList,
					Required: true,
					MaxItems: 64,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
				"exposed_headers": {
					Type:     schema.TypeList,
					Required: true,
					MaxItems: 64,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
				"allowed_headers": {
					Type:     schema.TypeList,
					Required: true,
					MaxItems: 64,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
				"allowed_methods": {
					Type:     schema.TypeList,
					Required: true,
					MaxItems: 64,
					Elem: &schema.Schema{
						Type: schema.TypeString,
						ValidateFunc: validation.StringInSlice([]string{
							"DELETE",
							"GET",
							"HEAD",
							"MERGE",
							"POST",
							"OPTIONS",
							"PUT"}, false),
					},
				},
				"max_age_in_seconds": {
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validation.IntBetween(1, 2000000000),
				},
			},
		},
	}
}
