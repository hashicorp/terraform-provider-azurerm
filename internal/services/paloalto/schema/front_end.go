package schema

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/firewalls"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type FrontEnd struct {
	Name                  string                  `tfschema:"name"`
	Protocol              string                  `tfschema:"protocol"`
	FrontendConfiguration []EndpointConfiguration `tfschema:"front_end_config"`
	BackendConfiguration  []EndpointConfiguration `tfschema:"back_end_config"`
}

type EndpointConfiguration struct {
	PublicIPID string `tfschema:"public_ip_address_id"`
	Port       int    `tfschema:"port"`
}

// FrontEndSchema returns the schema for a Palo Alto NGFW Front End Settings
func FrontEndSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"name": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: nil, // TODO - Validation needed
				},

				"protocol": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(firewalls.PossibleValuesForProtocolType(), false),
				},

				"back_end_config": EndpointSchema(),

				"front_end_config": EndpointSchema(),
			},
		},
	}
}

func EndpointSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		MaxItems: 1,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"public_ip_address_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: networkValidate.PublicIpAddressID,
				},

				"port": {
					Type:         pluginsdk.TypeInt,
					Required:     true,
					ValidateFunc: nil, // TODO - Need a atoi validation func for 1 - 65535
				},
			},
		},
	}
}
