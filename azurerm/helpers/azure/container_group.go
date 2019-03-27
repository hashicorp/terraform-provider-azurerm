package azure

import (
	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-10-01/containerinstance"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func SchemaContainerGroupProbe() *schema.Schema {
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

				"httpget": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"path": {
								Type:         schema.TypeString,
								Optional:     true,
								ForceNew:     true,
								ValidateFunc: validate.NoEmptyStrings,
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
									string(containerinstance.HTTP),
									string(containerinstance.HTTPS),
								}, true),
							},
						},
					},
				},

				"inital_delay_seconds": {
					Type:     schema.TypeInt,
					Optional: true,
					ForceNew: true,
				},

				"period_seconds": {
					Type:     schema.TypeInt,
					Optional: true,
					ForceNew: true,
				},

				"failure_threashold": {
					Type:     schema.TypeInt,
					Optional: true,
					ForceNew: true,
				},

				"sucess_threashold": {
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
