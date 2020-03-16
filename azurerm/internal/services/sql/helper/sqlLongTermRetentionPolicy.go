package helper

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func SqlLongTermRententionPolicy() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				// "retention_in_days": {
				// 	Type:         schema.TypeInt,
				// 	Optional:     true,
				// 	ValidateFunc: validation.IntBetween(0, 3285),
				// },
			},
		},
	}
}
