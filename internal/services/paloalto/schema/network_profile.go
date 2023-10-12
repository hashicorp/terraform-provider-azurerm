package schema

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/networkvirtualappliances"
	"github.com/hashicorp/go-azure-sdk/resource-manager/paloaltonetworks/2022-08-29/firewalls"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type NetworkProfileVnet struct {
	// Required
	PublicIPIDs []string `tfschema:"public_ip_address_ids"`

	// Optional
	EgressNatIPIDs    []string            `tfschema:"egress_nat_ip_address_ids"`
	VnetConfiguration []VnetConfiguration `tfschema:"vnet_configuration"`

	// Computed
	PublicIPs   []string `tfschema:"public_ip_addresses"`
	EgressNatIP []string `tfschema:"egress_nat_ip_addresses"`
}

type NetworkProfileVHub struct {
	VHubID      string   `tfschema:"virtual_hub_id"`
	PublicIPIDs []string `tfschema:"public_ip_address_ids"`

	// Optional
	EgressNatIPIDs []string `tfschema:"egress_nat_ip_address_ids"`

	// Computed
	PublicIPs       []string `tfschema:"public_ip_addresses"`
	EgressNatIP     []string `tfschema:"egress_nat_ip_addresses"`
	IpOfTrust       string   `tfschema:"ip_of_trust_for_user_defined_routes"`
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
				"public_ip_address_ids": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MinItems: 1,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: networkValidate.PublicIpAddressID,
					},
				},

				"egress_nat_ip_address_ids": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MinItems: 1,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: networkValidate.PublicIpAddressID,
					},
				},

				"vnet_configuration": VnetConfigurationSchema(),

				// Computed

				"public_ip_addresses": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"egress_nat_ip_addresses": {
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
		NetworkType:     firewalls.NetworkTypeVNET,
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

func FlattenNetworkProfileVnet(input firewalls.NetworkProfile) []NetworkProfileVnet {
	result := NetworkProfileVnet{}

	publicIPIDs := make([]string, 0)
	publicIPs := make([]string, 0)
	for _, v := range input.PublicIPs {
		if id := pointer.From(v.ResourceId); id != "" {
			publicIPIDs = append(publicIPIDs, id)
		}
		if ip := pointer.From(v.Address); ip != "" {
			publicIPs = append(publicIPs, ip)
		}
	}
	result.PublicIPIDs = publicIPIDs
	result.PublicIPs = publicIPs

	egressIds := make([]string, 0)
	egressIPs := make([]string, 0)
	if input.EgressNatIP != nil {
		for _, v := range *input.EgressNatIP {
			if id := pointer.From(v.ResourceId); id != "" {
				egressIds = append(egressIds, id)
			}
			if ip := pointer.From(v.Address); ip != "" {
				egressIPs = append(egressIPs, ip)
			}
		}
	}
	result.EgressNatIPIDs = egressIds
	result.EgressNatIP = egressIPs

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

	return []NetworkProfileVnet{result}
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

				"public_ip_address_ids": {
					Type:     pluginsdk.TypeList,
					Required: true,
					MinItems: 1,
					Elem: &pluginsdk.Schema{
						Type:         pluginsdk.TypeString,
						ValidateFunc: networkValidate.PublicIpAddressID,
					},
				},

				"egress_nat_ip_address_ids": {
					Type:     pluginsdk.TypeList,
					Optional: true,
					MinItems: 1,
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

				"ip_of_trust_for_user_defined_routes": {
					Type:     pluginsdk.TypeString,
					Computed: true,
				},

				"public_ip_addresses": {
					Type:     pluginsdk.TypeList,
					Computed: true,
					Elem: &pluginsdk.Schema{
						Type: pluginsdk.TypeString,
					},
				},

				"egress_nat_ip_addresses": {
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

		result.EgressNatIP = pointer.To(egressNatIPs)
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

func FlattenNetworkProfileVHub(input firewalls.NetworkProfile) (*NetworkProfileVHub, error) {
	result := NetworkProfileVHub{}

	publicIPIDs := make([]string, 0)
	publicIPs := make([]string, 0)
	for _, v := range input.PublicIPs {
		if id := pointer.From(v.ResourceId); id != "" {
			publicIPIDs = append(publicIPIDs, id)
		}
		if ip := pointer.From(v.Address); ip != "" {
			publicIPs = append(publicIPs, ip)
		}
	}
	result.PublicIPIDs = publicIPIDs
	result.PublicIPs = publicIPs

	egressIds := make([]string, 0)
	egressIPs := make([]string, 0)
	if input.EgressNatIP != nil {
		for _, v := range *input.EgressNatIP {
			if id := pointer.From(v.ResourceId); id != "" {
				egressIds = append(egressIds, id)
			}
			if ip := pointer.From(v.Address); ip != "" {
				egressIPs = append(egressIPs, ip)
			}
		}
	}
	result.EgressNatIPIDs = egressIds
	result.EgressNatIP = egressIPs

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
