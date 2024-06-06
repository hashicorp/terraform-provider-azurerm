// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkinterfaces"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type networkInterfaceUpdateInformation struct {
	applicationGatewayBackendAddressPoolIDs []string
	applicationSecurityGroupIDs             []string
	loadBalancerBackendAddressPoolIDs       map[string]string
	loadBalancerInboundNatRuleIDs           []string
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
	applicationGatewayBackendAddressPoolIds := make(map[string]struct{})
	loadBalancerBackendAddressPoolIds := make(map[string]string)
	loadBalancerInboundNatRuleIds := make(map[string]struct{})

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

			if props.ApplicationGatewayBackendAddressPools != nil {
				for _, pool := range *props.ApplicationGatewayBackendAddressPools {
					if pool.Id != nil {
						applicationGatewayBackendAddressPoolIds[*pool.Id] = struct{}{}
					}
				}
			}

			if props.LoadBalancerBackendAddressPools != nil {
				for _, pool := range *props.LoadBalancerBackendAddressPools {
					if v.Name != nil && pool.Id != nil {
						loadBalancerBackendAddressPoolIds[*pool.Id] = *v.Name
					}
				}
			}

			if props.LoadBalancerInboundNatRules != nil {
				for _, rule := range *props.LoadBalancerInboundNatRules {
					if rule.Id != nil {
						loadBalancerInboundNatRuleIds[*rule.Id] = struct{}{}
					}
				}
			}
		}
	}

	return networkInterfaceUpdateInformation{
		applicationGatewayBackendAddressPoolIDs: mapToSlice(applicationGatewayBackendAddressPoolIds),
		applicationSecurityGroupIDs:             mapToSlice(applicationSecurityGroupIds),
		loadBalancerBackendAddressPoolIDs:       loadBalancerBackendAddressPoolIds,
		loadBalancerInboundNatRuleIDs:           mapToSlice(loadBalancerInboundNatRuleIds),
		networkSecurityGroupID:                  networkSecurityGroupId,
	}
}

func mapFieldsToNetworkInterface(input *[]networkinterfaces.NetworkInterfaceIPConfiguration, info networkInterfaceUpdateInformation) *[]networkinterfaces.NetworkInterfaceIPConfiguration {
	output := input

	applicationSecurityGroups := make([]networkinterfaces.ApplicationSecurityGroup, 0)
	for _, id := range info.applicationSecurityGroupIDs {
		applicationSecurityGroups = append(applicationSecurityGroups, networkinterfaces.ApplicationSecurityGroup{
			Id: utils.String(id),
		})
	}

	applicationGatewayBackendAddressPools := make([]networkinterfaces.ApplicationGatewayBackendAddressPool, 0)
	for _, id := range info.applicationGatewayBackendAddressPoolIDs {
		applicationGatewayBackendAddressPools = append(applicationGatewayBackendAddressPools, networkinterfaces.ApplicationGatewayBackendAddressPool{
			Id: utils.String(id),
		})
	}

	loadBalancerInboundNatRules := make([]networkinterfaces.InboundNatRule, 0)
	for _, id := range info.loadBalancerInboundNatRuleIDs {
		loadBalancerInboundNatRules = append(loadBalancerInboundNatRules, networkinterfaces.InboundNatRule{
			Id: utils.String(id),
		})
	}

	for _, config := range *output {
		if config.Properties == nil || config.Properties.PrivateIPAddressVersion == nil || *config.Properties.PrivateIPAddressVersion != networkinterfaces.IPVersionIPvFour {
			continue
		}

		loadBalancerBackendAddressPools := make([]networkinterfaces.BackendAddressPool, 0)
		for id, name := range info.loadBalancerBackendAddressPoolIDs {
			if config.Name != nil && *config.Name == name {
				loadBalancerBackendAddressPools = append(loadBalancerBackendAddressPools, networkinterfaces.BackendAddressPool{
					Id: utils.String(id),
				})
			}

			config.Properties.LoadBalancerBackendAddressPools = &loadBalancerBackendAddressPools
		}

		config.Properties.ApplicationSecurityGroups = &applicationSecurityGroups
		config.Properties.ApplicationGatewayBackendAddressPools = &applicationGatewayBackendAddressPools
		config.Properties.LoadBalancerInboundNatRules = &loadBalancerInboundNatRules
	}

	return output
}
