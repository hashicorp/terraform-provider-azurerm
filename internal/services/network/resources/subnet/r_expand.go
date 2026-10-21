// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subnet

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/subnets"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandSubnetNetworkSecurityGroupID(d *pluginsdk.ResourceData) (*subnets.NetworkSecurityGroup, error) {
	wo, err := pluginsdk.GetWriteOnly(d, "network_security_group_id_wo", cty.String)
	if err != nil {
		return nil, err
	}

	if wo.IsNull() {
		return nil, nil // NOTE - Returning `nil` will remove the SG association!
	}

	return &subnets.NetworkSecurityGroup{
		Id: pointer.To(wo.AsString()),
	}, nil
}

func expandSubnetRouteTableID(d *pluginsdk.ResourceData) (*subnets.RouteTable, error) {
	wo, err := pluginsdk.GetWriteOnly(d, "route_table_id_wo", cty.String)
	if err != nil {
		return nil, err
	}

	if wo.IsNull() {
		return nil, nil // NOTE - Returning `nil` will remove the SG association!
	}

	return &subnets.RouteTable{
		Id: pointer.To(wo.AsString()),
	}, nil
}

func expandSubnetServiceEndpoint(input []any) *[]subnets.ServiceEndpointPropertiesFormat {
	endpoints := make([]subnets.ServiceEndpointPropertiesFormat, 0)

	for _, item := range input {
		v := item.(map[string]any)
		endpoint := subnets.ServiceEndpointPropertiesFormat{
			Service: pointer.To(v["service"].(string)),
		}
		if networkIdentifier := v["network_identifier"].(string); networkIdentifier != "" {
			endpoint.NetworkIdentifier = &subnets.SubResource{
				Id: pointer.To(networkIdentifier),
			}
		}
		endpoints = append(endpoints, endpoint)
	}

	return &endpoints
}

func expandSubnetDelegation(input []any) *[]subnets.Delegation {
	retDelegations := make([]subnets.Delegation, 0)

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

		retDelegation := subnets.Delegation{
			Name: pointer.To(deleData["name"].(string)),
			Properties: &subnets.ServiceDelegationPropertiesFormat{
				ServiceName: pointer.To(srvDelegation["name"].(string)),
				Actions:     &retSrvActions,
			},
		}

		retDelegations = append(retDelegations, retDelegation)
	}

	return &retDelegations
}

func expandSubnetNetworkPolicy(enabled bool) *subnets.VirtualNetworkPrivateLinkServiceNetworkPolicies {
	if enabled {
		return pointer.To(subnets.VirtualNetworkPrivateLinkServiceNetworkPoliciesEnabled)
	}
	return pointer.To(subnets.VirtualNetworkPrivateLinkServiceNetworkPoliciesDisabled)
}

func expandSubnetServiceEndpointPolicies(input []any) *[]subnets.ServiceEndpointPolicy {
	output := make([]subnets.ServiceEndpointPolicy, 0)
	for _, policy := range input {
		output = append(output, subnets.ServiceEndpointPolicy{Id: pointer.To(policy.(string))})
	}
	return &output
}

func expandSubnetIPAddressPool(input []any) *[]subnets.IPamPoolPrefixAllocation {
	if len(input) == 0 {
		return nil
	}

	outputs := make([]subnets.IPamPoolPrefixAllocation, 0)
	for _, v := range input {
		ipPoolRaw := v.(map[string]any)
		output := subnets.IPamPoolPrefixAllocation{}

		if v, ok := ipPoolRaw["number_of_ip_addresses"]; ok {
			output.NumberOfIPAddresses = pointer.To(v.(string))
		}

		if v, ok := ipPoolRaw["id"]; ok {
			output.Pool = &subnets.IPamPoolPrefixAllocationPool{
				Id: pointer.To(v.(string)),
			}
		}

		outputs = append(outputs, output)
	}

	return &outputs
}
