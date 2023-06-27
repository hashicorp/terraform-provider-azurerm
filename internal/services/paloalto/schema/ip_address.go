package schema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type IPAddress struct {
	Address     string `tfschema:"address"`
	IPAddressID string `tfschema:"ip_address_id"` // TODO - What is this? `ResourceID` in the API?
}

func IPAddressSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"address": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: validation.IsIPv4Address,
				},

				"ip_address_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,                              // Required?
					ValidateFunc: networkValidate.PublicIpAddressID, // TODO?
				},
			},
		},
	}
}
