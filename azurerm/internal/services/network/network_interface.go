package network

import (
	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-07-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

type networkInterfaceUpdateInformation struct {
	applicationGatewayBackendAddressPoolIDs []string
	applicationSecurityGroupIDs             []string
	loadBalancerBackendAddressPoolIDs       []string
	loadBalancerInboundNatRuleIDs           []string
	networkSecurityGroupID                  string
}

func parseFieldsFromNetworkInterface(input network.InterfacePropertiesFormat) networkInterfaceUpdateInformation {
	networkSecurityGroupId := ""
	if input.NetworkSecurityGroup != nil && input.NetworkSecurityGroup.ID != nil {
		networkSecurityGroupId = *input.NetworkSecurityGroup.ID
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
	loadBalancerBackendAddressPoolIds := make(map[string]struct{})
	loadBalancerInboundNatRuleIds := make(map[string]struct{})

	if input.IPConfigurations != nil {
		for _, v := range *input.IPConfigurations {
			if v.InterfaceIPConfigurationPropertiesFormat == nil {
				continue
			}

			props := *v.InterfaceIPConfigurationPropertiesFormat
			if props.ApplicationSecurityGroups != nil {
				for _, asg := range *props.ApplicationSecurityGroups {
					if asg.ID != nil {
						applicationSecurityGroupIds[*asg.ID] = struct{}{}
					}
				}
			}

			if props.ApplicationGatewayBackendAddressPools != nil {
				for _, pool := range *props.ApplicationGatewayBackendAddressPools {
					if pool.ID != nil {
						applicationGatewayBackendAddressPoolIds[*pool.ID] = struct{}{}
					}
				}
			}

			if props.LoadBalancerBackendAddressPools != nil {
				for _, pool := range *props.LoadBalancerBackendAddressPools {
					if pool.ID != nil {
						loadBalancerBackendAddressPoolIds[*pool.ID] = struct{}{}
					}
				}
			}

			if props.LoadBalancerInboundNatRules != nil {
				for _, rule := range *props.LoadBalancerInboundNatRules {
					if rule.ID != nil {
						loadBalancerInboundNatRuleIds[*rule.ID] = struct{}{}
					}
				}
			}
		}
	}

	return networkInterfaceUpdateInformation{
		applicationGatewayBackendAddressPoolIDs: mapToSlice(applicationGatewayBackendAddressPoolIds),
		applicationSecurityGroupIDs:             mapToSlice(applicationSecurityGroupIds),
		loadBalancerBackendAddressPoolIDs:       mapToSlice(loadBalancerBackendAddressPoolIds),
		loadBalancerInboundNatRuleIDs:           mapToSlice(loadBalancerInboundNatRuleIds),
		networkSecurityGroupID:                  networkSecurityGroupId,
	}
}

func mapFieldsToNetworkInterface(input *[]network.InterfaceIPConfiguration, info networkInterfaceUpdateInformation) *[]network.InterfaceIPConfiguration {
	output := input

	applicationSecurityGroups := make([]network.ApplicationSecurityGroup, 0)
	for _, id := range info.applicationSecurityGroupIDs {
		applicationSecurityGroups = append(applicationSecurityGroups, network.ApplicationSecurityGroup{
			ID: utils.String(id),
		})
	}

	applicationGatewayBackendAddressPools := make([]network.ApplicationGatewayBackendAddressPool, 0)
	for _, id := range info.applicationGatewayBackendAddressPoolIDs {
		applicationGatewayBackendAddressPools = append(applicationGatewayBackendAddressPools, network.ApplicationGatewayBackendAddressPool{
			ID: utils.String(id),
		})
	}

	loadBalancerBackendAddressPools := make([]network.BackendAddressPool, 0)
	for _, id := range info.loadBalancerBackendAddressPoolIDs {
		loadBalancerBackendAddressPools = append(loadBalancerBackendAddressPools, network.BackendAddressPool{
			ID: utils.String(id),
		})
	}

	loadBalancerInboundNatRules := make([]network.InboundNatRule, 0)
	for _, id := range info.loadBalancerInboundNatRuleIDs {
		loadBalancerInboundNatRules = append(loadBalancerInboundNatRules, network.InboundNatRule{
			ID: utils.String(id),
		})
	}

	for _, config := range *output {
		if config.InterfaceIPConfigurationPropertiesFormat == nil {
			continue
		}

		if config.InterfaceIPConfigurationPropertiesFormat.PrivateIPAddressVersion != network.IPv4 {
			continue
		}

		config.ApplicationSecurityGroups = &applicationSecurityGroups
		config.ApplicationGatewayBackendAddressPools = &applicationGatewayBackendAddressPools
		config.LoadBalancerBackendAddressPools = &loadBalancerBackendAddressPools
		config.LoadBalancerInboundNatRules = &loadBalancerInboundNatRules
	}

	return output
}
