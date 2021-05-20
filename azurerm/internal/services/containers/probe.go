package containers

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func SchemaContainerGroupProbe() *schema.Schema {
	//lintignore:XS003
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"exec": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validation.NoZeroValues,
					},
				},

				//lintignore:XS003
				"http_get": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"path": {
								Type:         schema.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"port": {
								Type:         schema.TypeInt,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validate.PortNumber,
							},
							"scheme": {
								Type:     schema.TypeString,
								Optional: true,
								ForceNew: true,
								ValidateFunc: validation.StringInSlice([]string{
									"Http",
									"Https",
								}, false),
							},
						},
					},
				},

				"initial_delay_seconds": {
					Type:     schema.TypeInt,
					Optional: true,
					ForceNew: true,
				},

				"period_seconds": {
					Type:     schema.TypeInt,
					Optional: true,
					ForceNew: true,
				},

				"failure_threshold": {
					Type:     schema.TypeInt,
					Optional: true,
					ForceNew: true,
				},

				"success_threshold": {
					Type:     schema.TypeInt,
					Optional: true,
					ForceNew: true,
				},

				"timeout_seconds": {
					Type:     schema.TypeInt,
					Optional: true,
					ForceNew: true,
				},
			},
		},
	}
}
