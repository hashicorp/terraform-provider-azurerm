package schema

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/firewalls"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetworkProfile struct {
	// Required
	PublicIPIDs []string `tfschema:"public_ip_ids"`

	// Optional
	EgressNatIPIDs    []string            `tfschema:"egress_nat_ip_ids"`
	VnetConfiguration []VnetConfiguration `tfschema:"vnet_configuration"`
	VWanConfiguration []VWanConfiguration `tfschema:"vwan_configuration"`

	// Computed
	PublicIPs   []string `tfschema:"public_ip"`
	EgressNatIP []string `tfschema:"egress_nat_ip_ids"`
	// Inferred
	// NetworkType string      `tfschema:"network_type"`
	// EnableEgressNat bool
}

func NetworkProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"public_ip_ids": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MinItems: 1,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: networkValidate.PublicIpAddressID,
					},
				},

				"egress_nat_ip_ids": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: networkValidate.PublicIpAddressID,
					},
				},

				"vnet_configuration": VnetConfigurationSchema(),

				"vwan_configuration": VWanConfigurationSchema(),

				// Computed

				"public_ip": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"egress_nat_ips": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},
			},
		},
	}
}

func ExpandNetworkProfile(input NetworkProfile) firewalls.NetworkProfile {
	result := firewalls.NetworkProfile{
		EnableEgressNat: firewalls.EgressNatDISABLED,
	}

	if len(input.PublicIPIDs) > 0 {
		ipIDs := make([]firewalls.IPAddress, 0)
		for _, v := range input.PublicIPIDs {
			ipIDs = append(ipIDs, firewalls.IPAddress{
				ResourceId: pointer.To(v),
			})
		}
		result.PublicIPs = ipIDs
	}

	if len(input.EgressNatIPIDs) > 0 {
		result.EnableEgressNat = firewalls.EgressNatENABLED
		egressNatIPs := make([]firewalls.IPAddress, 0)
		for _, v := range input.EgressNatIP {
			egressNatIPs = append(egressNatIPs, firewalls.IPAddress{
				ResourceId: pointer.To(v),
			})
		}
	}

	if len(input.VWanConfiguration) > 0 {
		result.NetworkType = firewalls.NetworkTypeVWAN
		vhub := input.VWanConfiguration[0]
		result.VwanConfiguration = &firewalls.VwanConfiguration{
			IPOfTrustSubnetForUdr:     nil,
			NetworkVirtualApplianceId: nil, // TODO Needs support for networkVirtualAPpliances adding to the provider
			TrustSubnet:               nil,
			UnTrustSubnet:             nil,
			VHub: firewalls.IPAddressSpace{
				ResourceId: pointer.To(vhub.VHubID),
			},
		}
	}

	if len(input.VnetConfiguration) > 0 {
		result.NetworkType = firewalls.NetworkTypeVNET
		vnet := input.VnetConfiguration[0]
		result.VnetConfiguration = &firewalls.VnetConfiguration{
			IPOfTrustSubnetForUdr: nil, // TODO - What is this?
			TrustSubnet: firewalls.IPAddressSpace{
				ResourceId: pointer.To(vnet.TrustedSubnetID),
			},
			UnTrustSubnet: firewalls.IPAddressSpace{
				ResourceId: pointer.To(vnet.UntrustedSubnetID),
			},
			Vnet: firewalls.IPAddressSpace{
				ResourceId: pointer.To(vnet.VNetID),
			},
		}
	}

	return result
}

func FlattenNetworkProfile(input firewalls.NetworkProfile) NetworkProfile {
	result := NetworkProfile{}

	for _, v := range input.PublicIPs {
		result.PublicIPIDs = append(result.PublicIPIDs, pointer.From(v.ResourceId))
		result.PublicIPs = append(result.PublicIPs, pointer.From(v.ResourceId))
	}

	if egressIPS := pointer.From(input.EgressNatIP); len(egressIPS) > 0 {
		for _, v := range egressIPS {
			result.EgressNatIPIDs = append(result.EgressNatIPIDs, pointer.From(v.ResourceId))
			result.EgressNatIP = append(result.EgressNatIP, pointer.From(v.Address))
		}
	}

	if v := input.VwanConfiguration; v != nil {
		vWan := VWanConfiguration{}

		vWan.VHubID = pointer.From(v.VHub.ResourceId)
		vWan.ApplianceID = pointer.From(v.NetworkVirtualApplianceId)

		if v.TrustSubnet != nil {
			vWan.TrustedSubnet = pointer.From(v.TrustSubnet.ResourceId)
		}

		if v.UnTrustSubnet != nil {
			vWan.UnTrustedSubnet = pointer.From(v.UnTrustSubnet.ResourceId)
		}

		if v.IPOfTrustSubnetForUdr != nil {
			vWan.IpOfTrust = pointer.From(v.IPOfTrustSubnetForUdr.Address)
		}

		result.VWanConfiguration = []VWanConfiguration{vWan}
	}

	if v := input.VnetConfiguration; v != nil {
		vNet := VnetConfiguration{}

		vNet.VNetID = pointer.From(v.Vnet.ResourceId)
		vNet.TrustedSubnetID = pointer.From(v.TrustSubnet.ResourceId)
		vNet.UntrustedSubnetID = pointer.From(v.UnTrustSubnet.ResourceId)

		if v.IPOfTrustSubnetForUdr != nil {
			vNet.IpOfTrust = pointer.From(v.IPOfTrustSubnetForUdr.Address)
		}

		result.VnetConfiguration = []VnetConfiguration{vNet}
	}

	return result
}
