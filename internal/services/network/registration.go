// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistrationWithAGitHubLabel   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
	_ sdk.FrameworkServiceRegistration               = Registration{}
)

// Name is the name of this Service
func (r Registration) Name() string {
	return "Network"
}

func (r Registration) AssociatedGitHubLabel() string {
	return "service/network"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Network",
	}
}

func (r Registration) DataSources() []sdk.DataSource {
	return []sdk.DataSource{
		VirtualNetworkPeeringDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		CustomIpPrefixResource{},
		PrivateEndpointApplicationSecurityGroupAssociationResource{},
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_application_security_group":                dataSourceApplicationSecurityGroup(),
		"azurerm_bastion_host":                              dataSourceBastionHost(),
		"azurerm_ip_group":                                  dataSourceIpGroup(),
		"azurerm_ip_groups":                                 dataSourceIpGroups(),
		"azurerm_nat_gateway":                               dataSourceNatGateway(),
		"azurerm_network_ddos_protection_plan":              dataSourceNetworkDDoSProtectionPlan(),
		"azurerm_network_interface":                         dataSourceNetworkInterface(),
		"azurerm_network_security_group":                    dataSourceNetworkSecurityGroup(),
		"azurerm_network_service_tags":                      dataSourceNetworkServiceTags(),
		"azurerm_network_watcher":                           dataSourceNetworkWatcher(),
		"azurerm_private_endpoint_connection":               dataSourcePrivateEndpointConnection(),
		"azurerm_private_link_service":                      dataSourcePrivateLinkService(),
		"azurerm_private_link_service_endpoint_connections": dataSourcePrivateLinkServiceEndpointConnections(),
		"azurerm_public_ip":                                 dataSourcePublicIP(),
		"azurerm_public_ip_prefix":                          dataSourcePublicIpPrefix(),
		"azurerm_public_ips":                                dataSourcePublicIPs(),
		"azurerm_route_table":                               dataSourceRouteTable(),
		"azurerm_subnet":                                    dataSourceSubnet(),
		"azurerm_virtual_network":                           dataSourceVirtualNetwork(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_application_security_group":               resourceApplicationSecurityGroup(),
		"azurerm_bastion_host":                             resourceBastionHost(),
		"azurerm_ip_group":                                 resourceIpGroup(),
		"azurerm_ip_group_cidr":                            resourceIpGroupCidr(),
		"azurerm_nat_gateway":                              resourceNatGateway(),
		"azurerm_nat_gateway_public_ip_association":        resourceNATGatewayPublicIpAssociation(),
		"azurerm_nat_gateway_public_ip_prefix_association": resourceNATGatewayPublicIpPrefixAssociation(),
		"azurerm_network_connection_monitor":               resourceNetworkConnectionMonitor(),
		"azurerm_network_ddos_protection_plan":             resourceNetworkDDoSProtectionPlan(),
		"azurerm_network_interface":                        resourceNetworkInterface(),

		"azurerm_network_interface_application_gateway_backend_address_pool_association": resourceNetworkInterfaceApplicationGatewayBackendAddressPoolAssociation(),
		"azurerm_network_interface_application_security_group_association":               resourceNetworkInterfaceApplicationSecurityGroupAssociation(),
		"azurerm_network_interface_backend_address_pool_association":                     resourceNetworkInterfaceBackendAddressPoolAssociation(),
		"azurerm_network_interface_nat_rule_association":                                 resourceNetworkInterfaceNatRuleAssociation(),
		"azurerm_network_interface_security_group_association":                           resourceNetworkInterfaceSecurityGroupAssociation(),

		"azurerm_network_profile":                           resourceNetworkProfile(),
		"azurerm_network_security_group":                    resourceNetworkSecurityGroup(),
		"azurerm_network_security_rule":                     resourceNetworkSecurityRule(),
		"azurerm_network_watcher":                           resourceNetworkWatcher(),
		"azurerm_network_watcher_flow_log":                  resourceNetworkWatcherFlowLog(),
		"azurerm_private_endpoint":                          resourcePrivateEndpoint(),
		"azurerm_private_link_service":                      resourcePrivateLinkService(),
		"azurerm_public_ip":                                 resourcePublicIp(),
		"azurerm_public_ip_prefix":                          resourcePublicIpPrefix(),
		"azurerm_route":                                     resourceRoute(),
		"azurerm_route_table":                               resourceRouteTable(),
		"azurerm_subnet":                                    resourceSubnet(),
		"azurerm_subnet_nat_gateway_association":            resourceSubnetNatGatewayAssociation(),
		"azurerm_subnet_network_security_group_association": resourceSubnetNetworkSecurityGroupAssociation(),
		"azurerm_subnet_route_table_association":            resourceSubnetRouteTableAssociation(),
		"azurerm_subnet_service_endpoint_storage_policy":    resourceSubnetServiceEndpointStoragePolicy(),
		"azurerm_virtual_machine_packet_capture":            resourceVirtualMachinePacketCapture(),
		"azurerm_virtual_machine_scale_set_packet_capture":  resourceVirtualMachineScaleSetPacketCapture(),
		"azurerm_virtual_network":                           resourceVirtualNetwork(),
		"azurerm_virtual_network_dns_servers":               resourceVirtualNetworkDnsServers(),
		"azurerm_virtual_network_peering":                   resourceVirtualNetworkPeering(),
	}
}

func (r Registration) Actions() []func() action.Action {
	return []func() action.Action{}
}

func (r Registration) ListResources() []sdk.FrameworkListWrappedResource {
	return []sdk.FrameworkListWrappedResource{
		ApplicationSecurityGroupListResource{},
		IpGroupListResource{},
		NatGatewayListResource{},
		NetworkDDoSProtectionPlanListResource{},
		NetworkInterfaceListResource{},
		NetworkProfileListResource{},
		NetworkSecurityGroupListResource{},
		NetworkSecurityRuleListResource{},
		PrivateEndpointListResource{},
		PublicIpListResource{},
		RouteListResource{},
		RouteTableListResource{},
		SubnetListResource{},
		VirtualNetworkListResource{},
		VirtualNetworkPeeringListResource{},
	}
}

func (r Registration) EphemeralResources() []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (r Registration) FrameworkDataSources() []sdk.FrameworkWrappedDataSource {
	return []sdk.FrameworkWrappedDataSource{}
}

func (r Registration) FrameworkResources() []sdk.FrameworkWrappedResource {
	return []sdk.FrameworkWrappedResource{}
}
