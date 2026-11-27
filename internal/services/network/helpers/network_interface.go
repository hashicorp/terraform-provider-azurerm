// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/networkinterfaces"
)

type NetworkInterfaceUpdateInformation struct {
	ApplicationGatewayBackendAddressPoolIDs []string
	ApplicationSecurityGroupIDs             []string
	LoadBalancerBackendAddressPoolIDs       map[string]string
	LoadBalancerInboundNatRuleIDs           []string
	NetworkSecurityGroupID                  string
}

func ParseFieldsFromNetworkInterface(input networkinterfaces.NetworkInterfacePropertiesFormat) NetworkInterfaceUpdateInformation {
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

	return NetworkInterfaceUpdateInformation{
		ApplicationGatewayBackendAddressPoolIDs: mapToSlice(applicationGatewayBackendAddressPoolIds),
		ApplicationSecurityGroupIDs:             mapToSlice(applicationSecurityGroupIds),
		LoadBalancerBackendAddressPoolIDs:       loadBalancerBackendAddressPoolIds,
		LoadBalancerInboundNatRuleIDs:           mapToSlice(loadBalancerInboundNatRuleIds),
		NetworkSecurityGroupID:                  networkSecurityGroupId,
	}
}

func MapFieldsToNetworkInterface(input *[]networkinterfaces.NetworkInterfaceIPConfiguration, info NetworkInterfaceUpdateInformation) *[]networkinterfaces.NetworkInterfaceIPConfiguration {
	output := input

	applicationSecurityGroups := make([]networkinterfaces.ApplicationSecurityGroup, 0)
	for _, id := range info.ApplicationSecurityGroupIDs {
		applicationSecurityGroups = append(applicationSecurityGroups, networkinterfaces.ApplicationSecurityGroup{
			Id: pointer.To(id),
		})
	}

	applicationGatewayBackendAddressPools := make([]networkinterfaces.ApplicationGatewayBackendAddressPool, 0)
	for _, id := range info.ApplicationGatewayBackendAddressPoolIDs {
		applicationGatewayBackendAddressPools = append(applicationGatewayBackendAddressPools, networkinterfaces.ApplicationGatewayBackendAddressPool{
			Id: pointer.To(id),
		})
	}

	loadBalancerInboundNatRules := make([]networkinterfaces.InboundNatRule, 0)
	for _, id := range info.LoadBalancerInboundNatRuleIDs {
		loadBalancerInboundNatRules = append(loadBalancerInboundNatRules, networkinterfaces.InboundNatRule{
			Id: pointer.To(id),
		})
	}

	for _, config := range *output {
		if config.Properties == nil || config.Properties.PrivateIPAddressVersion == nil || *config.Properties.PrivateIPAddressVersion != networkinterfaces.IPVersionIPvFour {
			continue
		}

		loadBalancerBackendAddressPools := make([]networkinterfaces.BackendAddressPool, 0)
		for id, name := range info.LoadBalancerBackendAddressPoolIDs {
			if config.Name != nil && *config.Name == name {
				loadBalancerBackendAddressPools = append(loadBalancerBackendAddressPools, networkinterfaces.BackendAddressPool{
					Id: pointer.To(id),
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
