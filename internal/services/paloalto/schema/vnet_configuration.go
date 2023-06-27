package schema

import (
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type VnetConfiguration struct {
	VNetID            string `tfschema:"virtual_network_id"`
	TrustedSubnetID   string `tfschema:"trusted_subnet_id"`
	UntrustedSubnetID string `tfschema:"untrusted_subnet_id"`
	IpOfTrust         string `tfschema:"ip_of_trust_for_udr"` // TODO - What is this?
}

func VnetConfigurationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		MaxItems: 1,
		ExactlyOneOf: []string{
			"network_profile.0.vnet_configuration",
			"network_profile.0.vwan_configuration",
		},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"virtual_network_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: networkValidate.VirtualNetworkID,
				},

				"trusted_subnet_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: networkValidate.SubnetID,
				},

				"untrusted_subnet_id": {
					Type:         pluginsdk.TypeString,
					Optional:     true,
					ValidateFunc: networkValidate.SubnetID,
				},

				"ip_of_trust_for_udr": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}
