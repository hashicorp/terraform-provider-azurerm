package schema

import (
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type VWanConfiguration struct {
	VHubID string `tfschema:"virtual_hub_id"`

	ApplianceID     string `tfschema:"virtual_network_appliance_id"`
	TrustedSubnet   string `tfschema:"trusted_subnet"`
	UnTrustedSubnet string `tfschema:"untrusted_subnet"`
	IpOfTrust       string `tfschema:"ip_of_trust_for_udr"`
}

func VWanConfigurationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		ExactlyOneOf: []string{
			"network_profile.0.vwan_configuration",
			"network_profile.0.vnet_configuration",
		},
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"virtual_hub_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: networkValidate.VirtualHubID,
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
