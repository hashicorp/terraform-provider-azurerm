package schema

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type IPAddressSpace struct {
	AddressSpace string `tfschema:"address"`
	IPAddressID  string `tfschema:"ip_address_id"` // TODO - What is this? `ResourceID` in the API?
}

func IPAddressSpaceSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"address_space": {
					Type:         pluginsdk.TypeString,
					Optional:     true, // TODO - Should this be Computed Only from the ID?
					ValidateFunc: validate.CIDR,
				},

				"ip_address_space_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,                             // Required?
					ValidateFunc: networkValidate.PublicIpPrefixID, // TODO - Should this be IPPrefix?
				},
			},
		},
	}
}
