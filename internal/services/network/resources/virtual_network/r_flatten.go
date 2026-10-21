// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routetables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceVirtualNetworkFlatten(d *pluginsdk.ResourceData, id commonids.VirtualNetworkId, vnet *virtualnetworks.VirtualNetwork) error {
	d.Set("name", id.VirtualNetworkName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if vnet != nil {
		d.Set("location", location.NormalizeNilable(vnet.Location))
		d.Set("edge_zone", helpers.FlattenEdgeZoneModel(vnet.ExtendedLocation))

		if props := vnet.Properties; props != nil {
			d.Set("guid", props.ResourceGuid)
			d.Set("flow_timeout_in_minutes", props.FlowTimeoutInMinutes)
			d.Set("private_endpoint_vnet_policies", pointer.FromEnum(props.PrivateEndpointVNetPolicies))

			if space := props.AddressSpace; space != nil {
				if err := d.Set("address_space", space.AddressPrefixes); err != nil {
					return fmt.Errorf("setting `address_space`: %+v", err)
				}

				if err := d.Set("ip_address_pool", flattenVirtualNetworkIPAddressPool(space.IPamPoolPrefixAllocations)); err != nil {
					return fmt.Errorf("setting `ip_address_pool`: %+v", err)
				}
			}

			if err := d.Set("ddos_protection_plan", flattenVirtualNetworkDDoSProtectionPlan(props)); err != nil {
				return fmt.Errorf("setting `ddos_protection_plan`: %+v", err)
			}

			if err := d.Set("encryption", flattenVirtualNetworkEncryption(props.Encryption)); err != nil {
				return fmt.Errorf("setting `encryption`: %+v", err)
			}

			subnet, err := flattenVirtualNetworkSubnets(props.Subnets)
			if err != nil {
				return fmt.Errorf("flattening `subnet`: %+v", err)
			}
			if err := d.Set("subnet", subnet); err != nil {
				return fmt.Errorf("setting `subnet`: %+v", err)
			}

			if err := d.Set("dns_servers", helpers.FlattenVirtualNetworkDNSServers(props.DhcpOptions)); err != nil {
				return fmt.Errorf("setting `dns_servers`: %+v", err)
			}

			bgpCommunity := ""
			if p := props.BgpCommunities; p != nil {
				bgpCommunity = p.VirtualNetworkCommunity
			}
			if err := d.Set("bgp_community", bgpCommunity); err != nil {
				return fmt.Errorf("setting `bgp_community`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, vnet.Tags); err != nil {
			return fmt.Errorf("flattening `tags`: %+v", err)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, pointer.To(id))
}

func flattenVirtualNetworkIPAddressPool(input *[]virtualnetworks.IPamPoolPrefixAllocation) []any {
	if input == nil {
		return []any{}
	}

	outputs := make([]any, 0)
	for _, v := range *input {
		output := map[string]any{
			"number_of_ip_addresses":        pointer.From(v.NumberOfIPAddresses),
			"allocated_ip_address_prefixes": pointer.From(v.AllocatedAddressPrefixes),
		}
		if v.Pool != nil {
			output["id"] = pointer.From(v.Pool.Id)
		}
		outputs = append(outputs, output)
	}

	return outputs
}

func flattenVirtualNetworkDDoSProtectionPlan(input *virtualnetworks.VirtualNetworkPropertiesFormat) []any {
	if input == nil {
		return []any{}
	}

	if input.DdosProtectionPlan == nil || input.DdosProtectionPlan.Id == nil || input.EnableDdosProtection == nil {
		return []any{}
	}

	return []any{
		map[string]any{
			"id":     *input.DdosProtectionPlan.Id,
			"enable": *input.EnableDdosProtection,
		},
	}
}

func flattenVirtualNetworkEncryption(encryption *virtualnetworks.VirtualNetworkEncryption) any {
	if encryption == nil || !encryption.Enabled {
		return make([]any, 0)
	}

	return []any{
		map[string]any{
			"enforcement": encryption.Enforcement,
		},
	}
}

func flattenVirtualNetworkSubnets(input *[]virtualnetworks.Subnet) (*pluginsdk.Set, error) {
	results := &pluginsdk.Set{
		F: resourceAzureSubnetHash,
	}

	if subnets := input; subnets != nil {
		for _, subnet := range *input {
			output := map[string]any{}

			if id := subnet.Id; id != nil {
				output["id"] = *id
			}

			if name := subnet.Name; name != nil {
				output["name"] = *name
			}

			if props := subnet.Properties; props != nil {
				if nsg := props.NetworkSecurityGroup; nsg != nil {
					if nsg.Id != nil {
						output["security_group"] = *nsg.Id
					}
				}

				if props.AddressPrefixes == nil {
					if props.AddressPrefix != nil && len(*props.AddressPrefix) > 0 {
						output["address_prefixes"] = []string{*props.AddressPrefix}
					} else {
						output["address_prefixes"] = []string{}
					}
				} else {
					output["address_prefixes"] = props.AddressPrefixes
				}
				output["delegation"] = flattenVirtualNetworkSubnetDelegation(props.Delegations)
				output["default_outbound_access_enabled"] = pointer.From(props.DefaultOutboundAccess)
				output["private_endpoint_network_policies"] = pointer.FromEnum(props.PrivateEndpointNetworkPolicies)
				output["private_link_service_network_policies_enabled"] = strings.EqualFold(pointer.FromEnum(props.PrivateLinkServiceNetworkPolicies), string(virtualnetworks.VirtualNetworkPrivateEndpointNetworkPoliciesEnabled))
				routeTableId := ""
				if props.RouteTable != nil && props.RouteTable.Id != nil {
					id, err := routetables.ParseRouteTableID(*props.RouteTable.Id)
					if err != nil {
						return nil, err
					}
					routeTableId = id.ID()
				}
				output["route_table_id"] = routeTableId
				output["service_endpoint"] = flattenVirtualNetworkSubnetServiceEndpoint(props.ServiceEndpoints)
				output["service_endpoint_policy_ids"] = flattenVirtualNetworkSubnetServiceEndpointPolicies(props.ServiceEndpointPolicies)
			}

			results.Add(output)
		}
	}

	return results, nil
}

func flattenVirtualNetworkSubnetServiceEndpointPolicies(input *[]virtualnetworks.ServiceEndpointPolicy) []any {
	output := make([]any, 0)
	if input == nil {
		return output
	}

	for _, policy := range *input {
		id := pointer.From(policy.Id)
		output = append(output, id)
	}
	return output
}

func flattenVirtualNetworkSubnetServiceEndpoint(serviceEndpoints *[]virtualnetworks.ServiceEndpointPropertiesFormat) []any {
	endpoints := make([]any, 0)

	if serviceEndpoints == nil {
		return endpoints
	}

	for _, endpoint := range *serviceEndpoints {
		item := map[string]any{
			"service": pointer.From(endpoint.Service),
		}
		if endpoint.NetworkIdentifier != nil {
			item["network_identifier"] = pointer.From(endpoint.NetworkIdentifier.Id)
		}
		endpoints = append(endpoints, item)
	}

	return endpoints
}

func flattenVirtualNetworkSubnetDelegation(delegations *[]virtualnetworks.Delegation) []any {
	if delegations == nil {
		return []any{}
	}

	retDeles := make([]any, 0)

	normalizeServiceName := map[string]string{}
	for _, normName := range subnetDelegationServiceNames {
		normalizeServiceName[strings.ToLower(normName)] = normName
	}

	for _, dele := range *delegations {
		retDele := make(map[string]any)
		if v := dele.Name; v != nil {
			retDele["name"] = *v
		}

		svcDeles := make([]any, 0)
		svcDele := make(map[string]any)
		if props := dele.Properties; props != nil {
			if v := props.ServiceName; v != nil {
				name := *v
				if nv, ok := normalizeServiceName[strings.ToLower(name)]; ok {
					name = nv
				}
				svcDele["name"] = name
			}

			if v := props.Actions; v != nil {
				svcDele["actions"] = *v
			}
		}

		svcDeles = append(svcDeles, svcDele)

		retDele["service_delegation"] = svcDeles

		retDeles = append(retDeles, retDele)
	}

	return retDeles
}
