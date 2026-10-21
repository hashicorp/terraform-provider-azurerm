// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network_interface

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceNetworkInterfaceFlatten(d *pluginsdk.ResourceData, id *commonids.NetworkInterfaceId, ni *networkinterfaces.NetworkInterface) error {
	d.Set("name", id.NetworkInterfaceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if ni != nil {
		d.Set("location", location.NormalizeNilable(ni.Location))
		d.Set("edge_zone", helpers.FlattenEdgeZoneModel(ni.ExtendedLocation))

		if props := ni.Properties; props != nil {
			primaryPrivateIPAddress, privateIPAddresses := flattenNetworkInterfacePrivateIPAddresses(props.IPConfigurations)
			d.Set("private_ip_address", primaryPrivateIPAddress)
			if err := d.Set("private_ip_addresses", privateIPAddresses); err != nil {
				return fmt.Errorf("setting `private_ip_addresses`: %+v", err)
			}

			appliedDNSServers := make([]string, 0)
			dnsServers := make([]string, 0)
			internalDnsNameLabel := ""
			internalDomainNameSuffix := ""
			if dnsSettings := props.DnsSettings; dnsSettings != nil {
				appliedDNSServers = flattenNetworkInterfaceDnsServers(dnsSettings.AppliedDnsServers)
				dnsServers = flattenNetworkInterfaceDnsServers(dnsSettings.DnsServers)
				internalDnsNameLabel = pointer.From(dnsSettings.InternalDnsNameLabel)
				internalDomainNameSuffix = pointer.From(dnsSettings.InternalDomainNameSuffix)
			}

			if err := d.Set("applied_dns_servers", appliedDNSServers); err != nil {
				return fmt.Errorf("setting `applied_dns_servers`: %+v", err)
			}

			if err := d.Set("dns_servers", dnsServers); err != nil {
				return fmt.Errorf("setting `applied_dns_servers`: %+v", err)
			}
			d.Set("internal_dns_name_label", internalDnsNameLabel)
			d.Set("internal_domain_name_suffix", internalDomainNameSuffix)

			virtualMachineId := ""
			if props.VirtualMachine != nil && props.VirtualMachine.Id != nil {
				virtualMachineId = *props.VirtualMachine.Id
			}
			d.Set("virtual_machine_id", virtualMachineId)

			auxiliaryMode := ""
			if props.AuxiliaryMode != nil && *props.AuxiliaryMode != networkinterfaces.NetworkInterfaceAuxiliaryModeNone {
				auxiliaryMode = string(*props.AuxiliaryMode)
			}
			d.Set("auxiliary_mode", auxiliaryMode)

			auxiliarySku := ""
			if props.AuxiliarySku != nil && *props.AuxiliarySku != networkinterfaces.NetworkInterfaceAuxiliarySkuNone {
				auxiliarySku = string(*props.AuxiliarySku)
			}
			d.Set("auxiliary_sku", auxiliarySku)

			d.Set("ip_forwarding_enabled", props.EnableIPForwarding)
			d.Set("accelerated_networking_enabled", props.EnableAcceleratedNetworking)
			d.Set("mac_address", props.MacAddress)

			if err := d.Set("ip_configuration", flattenNetworkInterfaceIPConfigurations(props.IPConfigurations)); err != nil {
				return fmt.Errorf("setting `ip_configuration`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, ni.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenNetworkInterfaceIPConfigurations(input *[]networkinterfaces.NetworkInterfaceIPConfiguration) []any {
	if input == nil {
		return []any{}
	}

	result := make([]any, 0)
	for _, ipConfig := range *input {
		props := ipConfig.Properties

		subnetId := ""
		if props.Subnet != nil && props.Subnet.Id != nil {
			subnetId = *props.Subnet.Id
		}

		privateIPAllocationMethod := ""
		if props.PrivateIPAllocationMethod != nil {
			privateIPAllocationMethod = string(*props.PrivateIPAllocationMethod)
		}

		privateIPAddressVersion := ""
		if props.PrivateIPAddressVersion != nil {
			privateIPAddressVersion = string(*props.PrivateIPAddressVersion)
		}

		publicIPAddressId := ""
		if props.PublicIPAddress != nil && props.PublicIPAddress.Id != nil {
			publicIPAddressId = *props.PublicIPAddress.Id
		}

		gatewayLBFrontendIPConfigId := ""
		if props.GatewayLoadBalancer != nil && props.GatewayLoadBalancer.Id != nil {
			gatewayLBFrontendIPConfigId = *props.GatewayLoadBalancer.Id
		}

		result = append(result, map[string]any{
			"name":                          pointer.From(ipConfig.Name),
			"primary":                       pointer.From(props.Primary),
			"private_ip_address":            pointer.From(props.PrivateIPAddress),
			"private_ip_address_allocation": privateIPAllocationMethod,
			"private_ip_address_version":    privateIPAddressVersion,
			"public_ip_address_id":          publicIPAddressId,
			"subnet_id":                     subnetId,
			"gateway_load_balancer_frontend_ip_configuration_id": gatewayLBFrontendIPConfigId,
		})
	}
	return result
}

func flattenNetworkInterfaceDnsServers(input *[]string) []string {
	if input == nil {
		return make([]string, 0)
	}

	return *input
}

func flattenNetworkInterfacePrivateIPAddresses(input *[]networkinterfaces.NetworkInterfaceIPConfiguration) (string, []any) {
	primary := ""
	result := make([]any, 0)

	if input == nil {
		return primary, result
	}

	for idx, config := range *input {
		if props := config.Properties; props != nil && props.PrivateIPAddress != nil {
			privateIP := *props.PrivateIPAddress
			if idx == 0 {
				primary = privateIP
			}

			result = append(result, privateIP)
		}
	}

	return primary, result
}
