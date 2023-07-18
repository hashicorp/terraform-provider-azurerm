package schema

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-02-01/networkvirtualappliances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/firewalls"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetworkProfileVnet struct {
	// Required
	PublicIPIDs []string `tfschema:"public_ip_ids"`

	// Optional
	EgressNatIPIDs    []string            `tfschema:"egress_nat_ip_ids"`
	VnetConfiguration []VnetConfiguration `tfschema:"vnet_configuration"`

	// Computed
	PublicIPs   []string `tfschema:"public_ip"`
	EgressNatIP []string `tfschema:"egress_nat_ips"`
	// Inferred
	// NetworkType string      `tfschema:"network_type"`
	// EnableEgressNat bool
}

type NetworkProfileVHub struct {
	VHubID      string   `tfschema:"virtual_hub_id"`
	PublicIPIDs []string `tfschema:"public_ip_ids"`

	// Optional
	EgressNatIPIDs []string `tfschema:"egress_nat_ip_ids"`

	// Computed
	PublicIPs       []string `tfschema:"public_ip"`
	EgressNatIP     []string `tfschema:"egress_nat_ip_ids"`
	IpOfTrust       string   `tfschema:"ip_of_trust_for_udr"`
	TrustedSubnet   string   `tfschema:"trusted_subnet_id"`
	UnTrustedSubnet string   `tfschema:"untrusted_subnet_id"`
	ApplianceID     string   `tfschema:"network_virtual_appliance_id"`
}

func VnetNetworkProfileSchema() *pluginsdk.Schema {
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

func VHubNetworkProfileSchema() *pluginsdk.Schema {
	return &pluginsdk.Schema{
		Type:     pluginsdk.TypeList,
		Required: true,
		MaxItems: 1,
		Elem: &pluginsdk.Resource{
			Schema: map[string]*pluginsdk.Schema{
				"virtual_hub_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: networkValidate.VirtualHubID,
				},

				"network_virtual_appliance_id": {
					Type:         pluginsdk.TypeString,
					Required:     true,
					ForceNew:     true,
					ValidateFunc: networkvirtualappliances.ValidateNetworkVirtualApplianceID,
				},

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

func ExpandNetworkProfileVnet(input []NetworkProfileVnet) firewalls.NetworkProfile {
	result := firewalls.NetworkProfile{
		EnableEgressNat: firewalls.EgressNatDISABLED,
		EgressNatIP:     &[]firewalls.IPAddress{},
	}
	if len(input) == 0 {
		return result
	}

	profile := input[0]

	if len(profile.PublicIPIDs) > 0 {
		ipIDs := make([]firewalls.IPAddress, 0)
		for _, v := range profile.PublicIPIDs {
			ipIDs = append(ipIDs, firewalls.IPAddress{
				ResourceId: pointer.To(v),
			})
		}
		result.PublicIPs = ipIDs
	}

	if len(profile.EgressNatIPIDs) > 0 {
		result.EnableEgressNat = firewalls.EgressNatENABLED
		egressNatIPs := make([]firewalls.IPAddress, 0)
		for _, v := range profile.EgressNatIPIDs {
			egressNatIPs = append(egressNatIPs, firewalls.IPAddress{
				ResourceId: pointer.To(v),
			})
		}
		result.EgressNatIP = pointer.To(egressNatIPs)
	}

	result.NetworkType = firewalls.NetworkTypeVNET
	vnet := profile.VnetConfiguration[0]
	result.VnetConfiguration = &firewalls.VnetConfiguration{
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

	return result
}

func ExpandNetworkProfileVHub(input []NetworkProfileVHub) firewalls.NetworkProfile {
	result := firewalls.NetworkProfile{
		EnableEgressNat: firewalls.EgressNatDISABLED,
		EgressNatIP:     &[]firewalls.IPAddress{},
	}
	if len(input) == 0 {
		return result
	}

	profile := input[0]

	if len(profile.PublicIPIDs) > 0 {
		ipIDs := make([]firewalls.IPAddress, 0)
		for _, v := range profile.PublicIPIDs {
			ipIDs = append(ipIDs, firewalls.IPAddress{
				ResourceId: pointer.To(v),
			})
		}
		result.PublicIPs = ipIDs
	}

	if len(profile.EgressNatIPIDs) > 0 {
		result.EnableEgressNat = firewalls.EgressNatENABLED
		egressNatIPs := make([]firewalls.IPAddress, 0)
		for _, v := range profile.EgressNatIPIDs {
			egressNatIPs = append(egressNatIPs, firewalls.IPAddress{
				ResourceId: pointer.To(v),
			})
		}
	}

	result.NetworkType = firewalls.NetworkTypeVWAN

	result.VwanConfiguration = &firewalls.VwanConfiguration{
		VHub: firewalls.IPAddressSpace{
			ResourceId: pointer.To(profile.VHubID),
		},
		NetworkVirtualApplianceId: pointer.To(profile.ApplianceID),
	}

	return result
}

func FlattenNetworkProfileVnet(input firewalls.NetworkProfile) NetworkProfileVnet {
	result := NetworkProfileVnet{}

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

func FlattenNetworkProfileVHub(input firewalls.NetworkProfile) (*NetworkProfileVHub, error) {
	result := NetworkProfileVHub{}

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

		result.VHubID = pointer.From(v.VHub.ResourceId)
		applianceID, err := networkvirtualappliances.ParseNetworkVirtualApplianceID(pointer.From(v.NetworkVirtualApplianceId))
		if err != nil {
			return nil, err
		}
		result.ApplianceID = applianceID.ID()

		if v.TrustSubnet != nil {
			result.TrustedSubnet = pointer.From(v.TrustSubnet.ResourceId)
		}

		if v.UnTrustSubnet != nil {
			result.UnTrustedSubnet = pointer.From(v.UnTrustSubnet.ResourceId)
		}

		if v.IPOfTrustSubnetForUdr != nil {
			result.IpOfTrust = pointer.From(v.IPOfTrustSubnetForUdr.Address)
		}
	}

	return pointer.To(result), nil
}
