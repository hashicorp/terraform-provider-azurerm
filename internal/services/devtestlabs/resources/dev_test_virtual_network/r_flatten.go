// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package dev_test_virtual_network

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/devtestlab/2018-09-15/virtualnetworks"
)

func flattenDevTestVirtualNetworkSubnets(input *[]virtualnetworks.SubnetOverride) []any {
	outputs := make([]any, 0)
	if input == nil {
		return outputs
	}

	for _, v := range *input {
		output := make(map[string]any)
		if v.LabSubnetName != nil {
			output["name"] = *v.LabSubnetName
		}
		output["use_public_ip_address"] = v.UsePublicIPAddressPermission
		output["use_in_virtual_machine_creation"] = v.UseInVMCreationPermission
		output["shared_public_ip_address"] = flattenDevTestVirtualNetworkSubnetIpAddressConfiguration(v.SharedPublicIPAddressConfiguration)

		outputs = append(outputs, output)
	}

	return outputs
}

func flattenDevTestVirtualNetworkSubnetIpAddressConfiguration(input *virtualnetworks.SubnetSharedPublicIPAddressConfiguration) []any {
	outputs := make([]any, 0)

	if input == nil {
		return outputs
	}

	output := make(map[string]any)
	if input.AllowedPorts != nil {
		output["allowed_ports"] = flattenDevTestVirtualNetworkSubnetAllowedPorts(input.AllowedPorts)
	}
	outputs = append(outputs, output)
	return outputs
}

func flattenDevTestVirtualNetworkSubnetAllowedPorts(input *[]virtualnetworks.Port) []any {
	outputs := make([]any, 0)
	if input == nil {
		return outputs
	}

	for _, v := range *input {
		output := make(map[string]any)
		output["backend_port"] = pointer.From(v.BackendPort)
		output["transport_protocol"] = pointer.From(v.TransportProtocol)
		outputs = append(outputs, output)
	}

	return outputs
}
