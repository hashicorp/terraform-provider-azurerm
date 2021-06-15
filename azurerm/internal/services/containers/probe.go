package containers

import (
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func SchemaContainerGroupProbe() *pluginsdk.Schema {
	//lintignore:XS003
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ForceNew: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"exec": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: validation.NoZeroValues,
					},
				},

				//lintignore:XS003
				"http_get": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					ForceNew: true,
					Elem: &pluginsdk.Resource{
						Schema: map[string]*pluginsdk.Schema{
							"path": {
								Type:         pluginsdk.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validation.StringIsNotEmpty,
							},
							"port": {
								Type:         pluginsdk.TypeInt,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validate.PortNumber,
							},
							"scheme": {
								Type:     pluginsdk.TypeString,
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
					Type:     pluginsdk.TypeInt,
					Optional: true,
					ForceNew: true,
				},

				"period_seconds": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					ForceNew: true,
				},

				"failure_threshold": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					ForceNew: true,
				},

				"success_threshold": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					ForceNew: true,
				},

				"timeout_seconds": {
					Type:     pluginsdk.TypeInt,
					Optional: true,
					ForceNew: true,
				},
			},
		},
	}
}
