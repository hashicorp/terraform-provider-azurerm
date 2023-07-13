package schema

import (
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type VWanConfiguration struct {
	VHubID string `tfschema:"virtual_hub_id"`

	IpOfTrust       string `tfschema:"ip_of_trust_for_udr"`
	TrustedSubnet   string `tfschema:"trusted_subnet"`
	UnTrustedSubnet string `tfschema:"untrusted_subnet"`
	ApplianceID     string `tfschema:"virtual_network_appliance_id"`
}

func VWanConfigurationSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Optional: true,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"virtual_hub_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ValidateFunc: networkValidate.VirtualHubID,
				},

				"trusted_subnet_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"untrusted_subnet_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"ip_of_trust_for_udr": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"virtual_network_appliance_id": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},
			},
		},
	}
}
