// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/networkinterfaces"
)

type networkInterfaceUpdateInformation struct {
	applicationGatewayBackendAddressPoolIDs map[string]string
	applicationSecurityGroupIDs             []string
	loadBalancerBackendAddressPoolIDs       map[string]string
	loadBalancerInboundNatRuleIDs           map[string]string
	networkSecurityGroupID                  string
}

func parseFieldsFromNetworkInterface(input networkinterfaces.NetworkInterfacePropertiesFormat) networkInterfaceUpdateInformation {
	networkSecurityGroupId := ""
	if input.NetworkSecurityGroup != nil && input.NetworkSecurityGroup.Id != nil {
		networkSecurityGroupId = *input.NetworkSecurityGroup.Id
	}

	mapToSlice := func(input map[string]struct{}) []string {
		output := make([]string, 0)

		for id := range input {
			output = append(output, id)
		}

		return output
	}

	applicationSecurityGroupIds := make(map[string]struct{})
	applicationGatewayBackendAddressPoolIds := make(map[string]string)
	loadBalancerBackendAddressPoolIds := make(map[string]string)
	loadBalancerInboundNatRuleIds := make(map[string]string)

	if input.IPConfigurations != nil {
		for _, v := range *input.IPConfigurations {
			if v.Properties == nil {
				continue
			}

			props := *v.Properties
			if props.ApplicationSecurityGroups != nil {
				for _, asg := range *props.ApplicationSecurityGroups {
					if asg.Id != nil {
						applicationSecurityGroupIds[*asg.Id] = struct{}{}
					}
				}
			}

			if v.Name == nil {
				continue
			}

			if props.ApplicationGatewayBackendAddressPools != nil {
				for _, pool := range *props.ApplicationGatewayBackendAddressPools {
					if pool.Id != nil {
						applicationGatewayBackendAddressPoolIds[*pool.Id] = *v.Name
					}
				}
			}

			if props.LoadBalancerBackendAddressPools != nil {
				for _, pool := range *props.LoadBalancerBackendAddressPools {
					if pool.Id != nil {
						loadBalancerBackendAddressPoolIds[*pool.Id] = *v.Name
					}
				}
			}

			if props.LoadBalancerInboundNatRules != nil {
				for _, rule := range *props.LoadBalancerInboundNatRules {
					if rule.Id != nil {
						loadBalancerInboundNatRuleIds[*rule.Id] = *v.Name
					}
				}
			}
		}
	}

	return networkInterfaceUpdateInformation{
		applicationGatewayBackendAddressPoolIDs: applicationGatewayBackendAddressPoolIds,
		applicationSecurityGroupIDs:             mapToSlice(applicationSecurityGroupIds),
		loadBalancerBackendAddressPoolIDs:       loadBalancerBackendAddressPoolIds,
		loadBalancerInboundNatRuleIDs:           loadBalancerInboundNatRuleIds,
		networkSecurityGroupID:                  networkSecurityGroupId,
	}
}

func mapFieldsToNetworkInterface(input *[]networkinterfaces.NetworkInterfaceIPConfiguration, info networkInterfaceUpdateInformation) *[]networkinterfaces.NetworkInterfaceIPConfiguration {
	output := input

	applicationSecurityGroups := make([]networkinterfaces.ApplicationSecurityGroup, 0)
	for _, id := range info.applicationSecurityGroupIDs {
		applicationSecurityGroups = append(applicationSecurityGroups, networkinterfaces.ApplicationSecurityGroup{
			Id: pointer.To(id),
		})
	}

	for _, config := range *output {
		if config.Properties == nil {
			continue
		}

		config.Properties.ApplicationSecurityGroups = &applicationSecurityGroups

		loadBalancerBackendAddressPools := make([]networkinterfaces.BackendAddressPool, 0)
		for id, name := range info.loadBalancerBackendAddressPoolIDs {
			if config.Name != nil && *config.Name == name {
				loadBalancerBackendAddressPools = append(loadBalancerBackendAddressPools, networkinterfaces.BackendAddressPool{
					Id: pointer.To(id),
				})
			}
		}
		config.Properties.LoadBalancerBackendAddressPools = &loadBalancerBackendAddressPools

		loadBalancerInboundNatRules := make([]networkinterfaces.InboundNatRule, 0)
		for id, name := range info.loadBalancerInboundNatRuleIDs {
			if config.Name != nil && *config.Name == name {
				loadBalancerInboundNatRules = append(loadBalancerInboundNatRules, networkinterfaces.InboundNatRule{
					Id: pointer.To(id),
				})
			}
		}
		config.Properties.LoadBalancerInboundNatRules = &loadBalancerInboundNatRules

		applicationGatewayBackendAddressPools := make([]networkinterfaces.ApplicationGatewayBackendAddressPool, 0)
		for id, name := range info.applicationGatewayBackendAddressPoolIDs {
			if config.Name != nil && *config.Name == name {
				applicationGatewayBackendAddressPools = append(applicationGatewayBackendAddressPools, networkinterfaces.ApplicationGatewayBackendAddressPool{
					Id: pointer.To(id),
				})
			}
		}
		config.Properties.ApplicationGatewayBackendAddressPools = &applicationGatewayBackendAddressPools
	}

	return output
}
