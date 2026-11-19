// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_virtual_network

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/virtualnetworks"
)

func expandDevTestVirtualNetworkSubnets(input []any, subscriptionId, resourceGroupName, virtualNetworkName string) *[]virtualnetworks.SubnetOverride {
	results := make([]virtualnetworks.SubnetOverride, 0)
	// default found from the Portal
	name := fmt.Sprintf("%sSubnet", virtualNetworkName)
	subnetId := commonids.NewSubnetID(subscriptionId, resourceGroupName, virtualNetworkName, name)
	allow := virtualnetworks.UsagePermissionTypeAllow

	if len(input) == 0 {
		result := virtualnetworks.SubnetOverride{
			ResourceId:                   pointer.To(subnetId.ID()),
			LabSubnetName:                pointer.To(name),
			UsePublicIPAddressPermission: &allow,
			UseInVMCreationPermission:    &allow,
		}
		results = append(results, result)
		return &results
	}

	for _, val := range input {
		v := val.(map[string]any)

		subnet := virtualnetworks.SubnetOverride{
			ResourceId:                         pointer.To(subnetId.ID()),
			LabSubnetName:                      pointer.To(name),
			UsePublicIPAddressPermission:       pointer.ToEnum[virtualnetworks.UsagePermissionType](v["use_public_ip_address"].(string)),
			UseInVMCreationPermission:          pointer.ToEnum[virtualnetworks.UsagePermissionType](v["use_in_virtual_machine_creation"].(string)),
			SharedPublicIPAddressConfiguration: expandDevTestVirtualNetworkSubnetIpAddressConfiguration(v["shared_public_ip_address"].([]any)),
		}
		results = append(results, subnet)
	}

	return &results
}

func expandDevTestVirtualNetworkSubnetIpAddressConfiguration(input []any) *virtualnetworks.SubnetSharedPublicIPAddressConfiguration {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]any)

	return &virtualnetworks.SubnetSharedPublicIPAddressConfiguration{
		AllowedPorts: expandDevTestVirtualNetworkSubnetAllowedPorts(v["allowed_ports"].([]any)),
	}
}

func expandDevTestVirtualNetworkSubnetAllowedPorts(input []any) *[]virtualnetworks.Port {
	results := make([]virtualnetworks.Port, 0)

	for _, val := range input {
		v := val.(map[string]any)

		allowedPort := virtualnetworks.Port{
			BackendPort:       pointer.To(int64(v["backend_port"].(int))),
			TransportProtocol: pointer.ToEnum[virtualnetworks.TransportProtocol](v["transport_protocol"].(string)),
		}
		results = append(results, allowedPort)
	}

	return &results
}
