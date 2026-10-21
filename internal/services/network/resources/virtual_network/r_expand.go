// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_network

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networksecuritygroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/routetables"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/virtualnetworks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandVirtualNetworkDdosProtectionPlan(input []any) (*virtualnetworks.SubResource, *bool) {
	if len(input) == 0 || input[0] == nil {
		return nil, nil
	}

	var id string
	var enabled bool

	ddosPPlan := input[0].(map[string]any)

	if v, ok := ddosPPlan["id"]; ok {
		id = v.(string)
	}

	if v, ok := ddosPPlan["enable"]; ok {
		enabled = v.(bool)
	}

	return &virtualnetworks.SubResource{
		Id: pointer.To(id),
	}, &enabled
}

func expandVirtualNetworkEncryption(input []any) *virtualnetworks.VirtualNetworkEncryption {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	attr := input[0].(map[string]any)
	return &virtualnetworks.VirtualNetworkEncryption{
		Enabled:     true,
		Enforcement: pointer.ToEnum[virtualnetworks.VirtualNetworkEncryptionEnforcement](attr["enforcement"].(string)),
	}
}

func expandVirtualNetworkSubnets(ctx context.Context, client virtualnetworks.VirtualNetworksClient, input []any, id commonids.VirtualNetworkId) (*[]virtualnetworks.Subnet, *[]string, error) {
	subnets := make([]virtualnetworks.Subnet, 0)
	routeTables := make([]string, 0)

	if len(input) == 0 {
		return &subnets, &routeTables, nil
	}

	for _, subnetRaw := range input {
		if subnetRaw == nil {
			continue
		}
		subnet := subnetRaw.(map[string]any)

		name := subnet["name"].(string)
		log.Printf("[INFO] setting subnets inside vNet, processing %q", name)
		// since subnets can also be created outside of vNet definition (as root objects)
		// do a GET on subnet properties from the server before setting them
		subnetObj, err := getExistingSubnet(ctx, client, id, name)
		if err != nil {
			return nil, nil, err
		}
		log.Printf("[INFO] Completed GET of Subnet props")

		// set the props from config and leave the rest intact
		subnetObj.Name = pointer.To(name)
		if subnetObj.Properties == nil {
			subnetObj.Properties = &virtualnetworks.SubnetPropertiesFormat{}
		}

		addressPrefixes := make([]string, 0)
		for _, prefix := range subnet["address_prefixes"].([]any) {
			addressPrefixes = append(addressPrefixes, prefix.(string))
		}

		subnetObj.Properties.AddressPrefixes = pointer.To(addressPrefixes)
		subnetObj.Properties.AddressPrefix = nil

		privateEndpointNetworkPolicies := virtualnetworks.VirtualNetworkPrivateEndpointNetworkPolicies(subnet["private_endpoint_network_policies"].(string))
		privateLinkServiceNetworkPolicies := virtualnetworks.VirtualNetworkPrivateLinkServiceNetworkPoliciesDisabled
		if subnet["private_link_service_network_policies_enabled"].(bool) {
			privateLinkServiceNetworkPolicies = virtualnetworks.VirtualNetworkPrivateLinkServiceNetworkPoliciesEnabled
		}
		subnetObj.Properties.DefaultOutboundAccess = pointer.To(subnet["default_outbound_access_enabled"].(bool))
		subnetObj.Properties.Delegations = expandVirtualNetworkSubnetDelegation(subnet["delegation"].([]any))
		subnetObj.Properties.PrivateEndpointNetworkPolicies = pointer.To(privateEndpointNetworkPolicies)
		subnetObj.Properties.PrivateLinkServiceNetworkPolicies = pointer.To(privateLinkServiceNetworkPolicies)

		if routeTableId := subnet["route_table_id"].(string); routeTableId != "" {
			id, err := routetables.ParseRouteTableID(routeTableId)
			if err != nil {
				return nil, nil, err
			}

			// Collecting a list of route tables to lock on outside of this function
			routeTables = append(routeTables, id.ID())
			subnetObj.Properties.RouteTable = &virtualnetworks.RouteTable{
				Id: pointer.To(id.ID()),
			}
		} else {
			subnetObj.Properties.RouteTable = nil
		}

		subnetObj.Properties.ServiceEndpointPolicies = expandVirtualNetworkSubnetServiceEndpointPolicies(subnet["service_endpoint_policy_ids"].(*pluginsdk.Set).List())
		subnetObj.Properties.ServiceEndpoints = expandVirtualNetworkSubnetServiceEndpoint(subnet["service_endpoint"].([]any))

		if secGroup := subnet["security_group"].(string); secGroup != "" {
			subnetObj.Properties.NetworkSecurityGroup = &virtualnetworks.NetworkSecurityGroup{
				Id: &secGroup,
			}
		} else {
			subnetObj.Properties.NetworkSecurityGroup = nil
		}

		subnets = append(subnets, *subnetObj)
	}

	return &subnets, &routeTables, nil
}

func expandVirtualNetworkProperties(ctx context.Context, client virtualnetworks.VirtualNetworksClient, id commonids.VirtualNetworkId, d *pluginsdk.ResourceData) (*virtualnetworks.VirtualNetworkPropertiesFormat, *[]string, error) {
	subnets := make([]virtualnetworks.Subnet, 0)
	routeTables := make([]string, 0)
	if subs := d.Get("subnet").(*pluginsdk.Set); subs.Len() > 0 {
		for _, subnet := range subs.List() {
			subnet := subnet.(map[string]any)

			name := subnet["name"].(string)
			log.Printf("[INFO] setting subnets inside vNet, processing %q", name)
			// since subnets can also be created outside of vNet definition (as root objects)
			// do a GET on subnet properties from the server before setting them
			subnetObj, err := getExistingSubnet(ctx, client, id, name)
			if err != nil {
				return nil, nil, err
			}
			log.Printf("[INFO] Completed GET of Subnet props")

			// set the props from config and leave the rest intact
			subnetObj.Name = pointer.To(name)
			if subnetObj.Properties == nil {
				subnetObj.Properties = &virtualnetworks.SubnetPropertiesFormat{}
			}

			addressPrefixes := make([]string, 0)
			for _, prefix := range subnet["address_prefixes"].([]any) {
				addressPrefixes = append(addressPrefixes, prefix.(string))
			}

			subnetObj.Properties.AddressPrefixes = pointer.To(addressPrefixes)

			privateEndpointNetworkPolicies := virtualnetworks.VirtualNetworkPrivateEndpointNetworkPolicies(subnet["private_endpoint_network_policies"].(string))
			privateLinkServiceNetworkPolicies := virtualnetworks.VirtualNetworkPrivateLinkServiceNetworkPoliciesDisabled
			if subnet["private_link_service_network_policies_enabled"].(bool) {
				privateLinkServiceNetworkPolicies = virtualnetworks.VirtualNetworkPrivateLinkServiceNetworkPoliciesEnabled
			}
			subnetObj.Properties.DefaultOutboundAccess = pointer.To(subnet["default_outbound_access_enabled"].(bool))
			subnetObj.Properties.Delegations = expandVirtualNetworkSubnetDelegation(subnet["delegation"].([]any))
			subnetObj.Properties.PrivateEndpointNetworkPolicies = pointer.To(privateEndpointNetworkPolicies)
			subnetObj.Properties.PrivateLinkServiceNetworkPolicies = pointer.To(privateLinkServiceNetworkPolicies)

			if routeTableId := subnet["route_table_id"].(string); routeTableId != "" {
				id, err := routetables.ParseRouteTableID(routeTableId)
				if err != nil {
					return nil, nil, err
				}

				// Collecting a list of route tables to lock on outside of this function
				routeTables = append(routeTables, id.ID())
				subnetObj.Properties.RouteTable = &virtualnetworks.RouteTable{
					Id: pointer.To(id.ID()),
				}
			}

			subnetObj.Properties.ServiceEndpointPolicies = expandVirtualNetworkSubnetServiceEndpointPolicies(subnet["service_endpoint_policy_ids"].(*pluginsdk.Set).List())
			subnetObj.Properties.ServiceEndpoints = expandVirtualNetworkSubnetServiceEndpoint(subnet["service_endpoint"].([]any))

			if secGroup := subnet["security_group"].(string); secGroup != "" {
				subnetObj.Properties.NetworkSecurityGroup = &virtualnetworks.NetworkSecurityGroup{
					Id: &secGroup,
				}
			} else {
				subnetObj.Properties.NetworkSecurityGroup = nil
			}

			subnets = append(subnets, *subnetObj)
		}
	}

	properties := &virtualnetworks.VirtualNetworkPropertiesFormat{
		AddressSpace: &virtualnetworks.AddressSpace{},
		DhcpOptions: &virtualnetworks.DhcpOptions{
			DnsServers: pluginsdk.ExpandStringSlice(d.Get("dns_servers").([]any)),
		},
		PrivateEndpointVNetPolicies: pointer.ToEnum[virtualnetworks.PrivateEndpointVNetPolicies](d.Get("private_endpoint_vnet_policies").(string)),
		Subnets:                     &subnets,
	}

	if v, ok := d.GetOk("address_space"); ok {
		properties.AddressSpace.AddressPrefixes = pluginsdk.ExpandStringSlice(v.(*pluginsdk.Set).List())
	}

	if v, ok := d.GetOk("ddos_protection_plan"); ok {
		rawList := v.([]any)

		var ddosPPlan map[string]any
		if len(rawList) > 0 {
			ddosPPlan = rawList[0].(map[string]any)
		}

		if v, ok := ddosPPlan["id"]; ok {
			properties.DdosProtectionPlan = &virtualnetworks.SubResource{
				Id: pointer.To(v.(string)),
			}
		}

		if v, ok := ddosPPlan["enable"]; ok {
			properties.EnableDdosProtection = pointer.To(v.(bool))
		}
	}

	if v, ok := d.GetOk("encryption"); ok {
		if vList := v.([]any); len(vList) > 0 && vList[0] != nil {
			encryptionConf := vList[0].(map[string]any)
			properties.Encryption = &virtualnetworks.VirtualNetworkEncryption{
				Enabled:     true,
				Enforcement: pointer.ToEnum[virtualnetworks.VirtualNetworkEncryptionEnforcement](encryptionConf["enforcement"].(string)),
			}
		}
	}

	if v, ok := d.GetOk("ip_address_pool"); ok {
		properties.AddressSpace.IPamPoolPrefixAllocations = expandVirtualNetworkIPAddressPool(v.([]any))
	}

	if v, ok := d.GetOk("bgp_community"); ok {
		properties.BgpCommunities = &virtualnetworks.VirtualNetworkBgpCommunities{VirtualNetworkCommunity: v.(string)}
	}

	return properties, &routeTables, nil
}

func expandVirtualNetworkIPAddressPool(input []any) *[]virtualnetworks.IPamPoolPrefixAllocation {
	if len(input) == 0 {
		return nil
	}

	outputs := make([]virtualnetworks.IPamPoolPrefixAllocation, 0)
	for _, v := range input {
		ipPoolRaw := v.(map[string]any)
		output := virtualnetworks.IPamPoolPrefixAllocation{}

		if v, ok := ipPoolRaw["number_of_ip_addresses"]; ok {
			output.NumberOfIPAddresses = pointer.To(v.(string))
		}

		if v, ok := ipPoolRaw["id"]; ok {
			output.Pool = &virtualnetworks.IPamPoolPrefixAllocationPool{
				Id: pointer.To(v.(string)),
			}
		}

		outputs = append(outputs, output)
	}

	return &outputs
}

func expandResourcesForLocking(d *pluginsdk.ResourceData) ([]string, []string, error) {
	nsgIds := make([]string, 0)
	routeTableIds := make([]string, 0)

	if v, ok := d.GetOk("subnet"); ok {
		subnets := v.(*pluginsdk.Set).List()
		for _, subnet := range subnets {
			subnet, ok := subnet.(map[string]any)
			if !ok {
				return nil, nil, fmt.Errorf("[ERROR] Subnet should be a Hash - was '%+v'", subnet)
			}

			networkSecurityGroupId := subnet["security_group"].(string)
			if networkSecurityGroupId != "" {
				parsedNsgID, err := networksecuritygroups.ParseNetworkSecurityGroupID(networkSecurityGroupId)
				if err != nil {
					return nil, nil, err
				}
				nsgIds = append(nsgIds, parsedNsgID.ID())
			}

			routeTableId := subnet["route_table_id"].(string)
			if routeTableId != "" {
				parsedRouteTableID, err := routetables.ParseRouteTableID(routeTableId)
				if err != nil {
					return nil, nil, err
				}
				routeTableIds = append(routeTableIds, parsedRouteTableID.ID())
			}
		}
	}

	return nsgIds, routeTableIds, nil
}

func expandVirtualNetworkSubnetServiceEndpointPolicies(input []any) *[]virtualnetworks.ServiceEndpointPolicy {
	output := make([]virtualnetworks.ServiceEndpointPolicy, 0)
	for _, policy := range input {
		output = append(output, virtualnetworks.ServiceEndpointPolicy{Id: pointer.To(policy.(string))})
	}
	return &output
}

func expandVirtualNetworkSubnetDelegation(input []any) *[]virtualnetworks.Delegation {
	retDelegations := make([]virtualnetworks.Delegation, 0)

	for _, deleValue := range input {
		deleData := deleValue.(map[string]any)
		srvDelegations := deleData["service_delegation"].([]any)
		srvDelegation := srvDelegations[0].(map[string]any)

		srvActions := srvDelegation["actions"].(*pluginsdk.Set).List()

		retSrvActions := make([]string, 0)
		for _, srvAction := range srvActions {
			srvActionData := srvAction.(string)
			retSrvActions = append(retSrvActions, srvActionData)
		}

		retDelegation := virtualnetworks.Delegation{
			Name: pointer.To(deleData["name"].(string)),
			Properties: &virtualnetworks.ServiceDelegationPropertiesFormat{
				ServiceName: pointer.To(srvDelegation["name"].(string)),
				Actions:     &retSrvActions,
			},
		}

		retDelegations = append(retDelegations, retDelegation)
	}

	return &retDelegations
}

func expandVirtualNetworkSubnetServiceEndpoint(input []any) *[]virtualnetworks.ServiceEndpointPropertiesFormat {
	endpoints := make([]virtualnetworks.ServiceEndpointPropertiesFormat, 0)

	for _, item := range input {
		v := item.(map[string]any)
		endpoint := virtualnetworks.ServiceEndpointPropertiesFormat{
			Service: pointer.To(v["service"].(string)),
		}
		if networkIdentifier := v["network_identifier"].(string); networkIdentifier != "" {
			endpoint.NetworkIdentifier = &virtualnetworks.SubResource{
				Id: pointer.To(networkIdentifier),
			}
		}
		endpoints = append(endpoints, endpoint)
	}

	return &endpoints
}
