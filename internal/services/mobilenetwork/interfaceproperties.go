package mobilenetwork

import (
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type InterfacePropertiesModel struct {
	IPv4Address string `tfschema:"ipv4_address"`
	IPv4Gateway string `tfschema:"ipv4_gateway"`
	IPv4Subnet  string `tfschema:"ipv4_subnet"`
	Name        string `tfschema:"name"`
}

func interfacePropertiesSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:     pluginsdk.TypeString,
					Required: true,
				},
				"ipv4_address": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.IsIPv4Address,
				},
				"ipv4_subnet": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validate.CIDR,
				},
				"ipv4_gateway": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.IsIPv6Address,
				},
			},
		},
	}
}
