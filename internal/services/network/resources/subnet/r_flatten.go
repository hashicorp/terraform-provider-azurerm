// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subnet

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routetables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceSubnetFlatten(d *pluginsdk.ResourceData, id commonids.SubnetId, subnet *subnets.Subnet) error {
	d.Set("name", id.SubnetName)
	d.Set("virtual_network_name", id.VirtualNetworkName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("network_security_group_id_wo_version", d.Get("network_security_group_id_wo_version").(int))
	d.Set("route_table_id_wo_version", d.Get("route_table_id_wo_version").(int))

	if subnet != nil {
		if props := subnet.Properties; props != nil {
			if props.AddressPrefixes == nil {
				if props.AddressPrefix != nil && len(*props.AddressPrefix) > 0 {
					d.Set("address_prefixes", []string{*props.AddressPrefix})
				} else {
					d.Set("address_prefixes", []string{})
				}
			} else {
				d.Set("address_prefixes", props.AddressPrefixes)
			}

			defaultOutboundAccessEnabled := true
			if props.DefaultOutboundAccess != nil {
				defaultOutboundAccessEnabled = *props.DefaultOutboundAccess
			}
			d.Set("default_outbound_access_enabled", defaultOutboundAccessEnabled)

			if err := d.Set("delegation", flattenSubnetDelegation(props.Delegations)); err != nil {
				return fmt.Errorf("flattening `delegation`: %+v", err)
			}

			if err := d.Set("ip_address_pool", flattenSubnetIPAddressPool(props.IPamPoolPrefixAllocations)); err != nil {
				return fmt.Errorf("setting `ip_address_pool`: %+v", err)
			}

			d.Set("private_endpoint_network_policies", pointer.FromEnum(props.PrivateEndpointNetworkPolicies))
			d.Set("private_link_service_network_policies_enabled", flattenSubnetNetworkPolicy(pointer.FromEnum(props.PrivateLinkServiceNetworkPolicies)))
			d.Set("sharing_scope", pointer.FromEnum(props.SharingScope))

			if err := d.Set("service_endpoint", flattenSubnetServiceEndpoint(props.ServiceEndpoints)); err != nil {
				return fmt.Errorf("setting `service_endpoint`: %+v", err)
			}

			if err := d.Set("service_endpoint_policy_ids", flattenSubnetServiceEndpointPolicies(props.ServiceEndpointPolicies)); err != nil {
				return fmt.Errorf("setting `service_endpoint_policy_ids`: %+v", err)
			}

			var networkSecurityGroupID string
			if props.NetworkSecurityGroup != nil && props.NetworkSecurityGroup.Id != nil {
				if nsgID, err := networksecuritygroups.ParseNetworkSecurityGroupID(*props.NetworkSecurityGroup.Id); err == nil {
					if nsgID != nil {
						networkSecurityGroupID = nsgID.ID()
					}
				}
			}
			d.Set("network_security_group_id", networkSecurityGroupID)

			var routeTableID string
			if props.RouteTable != nil && props.RouteTable.Id != nil {
				if rtID, err := routetables.ParseRouteTableID(*props.RouteTable.Id); err == nil {
					routeTableID = rtID.ID()
				}
			}
			d.Set("route_table_id", routeTableID)
		}
	}

	return pluginsdk.SetResourceIdentityData(d, &id)
}

func flattenSubnetServiceEndpoint(serviceEndpoints *[]subnets.ServiceEndpointPropertiesFormat) []any {
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

func flattenSubnetDelegation(delegations *[]subnets.Delegation) []any {
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

func flattenSubnetNetworkPolicy(input string) bool {
	return strings.EqualFold(input, string(subnets.VirtualNetworkPrivateEndpointNetworkPoliciesEnabled))
}

func flattenSubnetServiceEndpointPolicies(input *[]subnets.ServiceEndpointPolicy) []any {
	if input == nil {
		return []any{}
	}

	output := make([]any, 0, len(*input))
	for _, policy := range *input {
		id := pointer.From(policy.Id)
		output = append(output, id)
	}
	return output
}

func flattenSubnetIPAddressPool(input *[]subnets.IPamPoolPrefixAllocation) []any {
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
