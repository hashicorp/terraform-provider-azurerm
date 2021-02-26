package helper

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func StorageSettingSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"luns": {
					Type:     schema.TypeList,
					Required: true,
					Elem: &schema.Schema{
						Type: schema.TypeInt,
					},
				},
				"default_file_path": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}
