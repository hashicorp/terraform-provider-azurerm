// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type Registration struct{}

var (
	_ sdk.TypedServiceRegistrationWithAGitHubLabel   = Registration{}
	_ sdk.UntypedServiceRegistrationWithAGitHubLabel = Registration{}
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
		ManagerDataSource{},
		ManagerNetworkGroupDataSource{},
		ManagerConnectivityConfigurationDataSource{},
	}
}

func (r Registration) Resources() []sdk.Resource {
	return []sdk.Resource{
		CustomIpPrefixResource{},
		ManagerAdminRuleResource{},
		ManagerAdminRuleCollectionResource{},
		ManagerDeploymentResource{},
		ManagerConnectivityConfigurationResource{},
		ManagerManagementGroupConnectionResource{},
		ManagerNetworkGroupResource{},
		ManagerResource{},
		ManagerScopeConnectionResource{},
		ManagerSecurityAdminConfigurationResource{},
		ManagerStaticMemberResource{},
		ManagerSubscriptionConnectionResource{},
		PrivateEndpointApplicationSecurityGroupAssociationResource{},
		RouteMapResource{},
		VirtualHubRoutingIntentResource{},
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_application_gateway":                       dataSourceApplicationGateway(),
		"azurerm_application_security_group":                dataSourceApplicationSecurityGroup(),
		"azurerm_bastion_host":                              dataSourceBastionHost(),
		"azurerm_express_route_circuit":                     dataSourceExpressRouteCircuit(),
		"azurerm_express_route_circuit_peering":             dataSourceExpressRouteCircuitPeering(),
		"azurerm_ip_group":                                  dataSourceIpGroup(),
		"azurerm_ip_groups":                                 dataSourceIpGroups(),
		"azurerm_nat_gateway":                               dataSourceNatGateway(),
		"azurerm_network_ddos_protection_plan":              dataSourceNetworkDDoSProtectionPlan(),
		"azurerm_network_interface":                         dataSourceNetworkInterface(),
		"azurerm_network_security_group":                    dataSourceNetworkSecurityGroup(),
		"azurerm_network_watcher":                           dataSourceNetworkWatcher(),
		"azurerm_private_endpoint_connection":               dataSourcePrivateEndpointConnection(),
		"azurerm_private_link_service":                      dataSourcePrivateLinkService(),
		"azurerm_private_link_service_endpoint_connections": dataSourcePrivateLinkServiceEndpointConnections(),
		"azurerm_public_ip":                                 dataSourcePublicIP(),
		"azurerm_public_ips":                                dataSourcePublicIPs(),
		"azurerm_public_ip_prefix":                          dataSourcePublicIpPrefix(),
		"azurerm_route_filter":                              dataSourceRouteFilter(),
		"azurerm_route_table":                               dataSourceRouteTable(),
		"azurerm_network_service_tags":                      dataSourceNetworkServiceTags(),
		"azurerm_subnet":                                    dataSourceSubnet(),
		"azurerm_virtual_hub":                               dataSourceVirtualHub(),
		"azurerm_virtual_hub_connection":                    dataSourceVirtualHubConnection(),
		"azurerm_virtual_hub_route_table":                   dataSourceVirtualHubRouteTable(),
		"azurerm_virtual_network_gateway":                   dataSourceVirtualNetworkGateway(),
		"azurerm_virtual_network_gateway_connection":        dataSourceVirtualNetworkGatewayConnection(),
		"azurerm_virtual_network":                           dataSourceVirtualNetwork(),
		"azurerm_web_application_firewall_policy":           dataWebApplicationFirewallPolicy(),
		"azurerm_virtual_wan":                               dataSourceVirtualWan(),
		"azurerm_local_network_gateway":                     dataSourceLocalNetworkGateway(),
		"azurerm_vpn_gateway":                               dataSourceVPNGateway(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*pluginsdk.Resource {
	return map[string]*pluginsdk.Resource{
		"azurerm_application_gateway":                      resourceApplicationGateway(),
		"azurerm_application_security_group":               resourceApplicationSecurityGroup(),
		"azurerm_bastion_host":                             resourceBastionHost(),
		"azurerm_express_route_circuit_connection":         resourceExpressRouteCircuitConnection(),
		"azurerm_express_route_circuit_authorization":      resourceExpressRouteCircuitAuthorization(),
		"azurerm_express_route_circuit_peering":            resourceExpressRouteCircuitPeering(),
		"azurerm_express_route_circuit":                    resourceExpressRouteCircuit(),
		"azurerm_express_route_connection":                 resourceExpressRouteConnection(),
		"azurerm_express_route_gateway":                    resourceExpressRouteGateway(),
		"azurerm_express_route_port_authorization":         resourceExpressRoutePortAuthorization(),
		"azurerm_express_route_port":                       resourceArmExpressRoutePort(),
		"azurerm_ip_group":                                 resourceIpGroup(),
		"azurerm_ip_group_cidr":                            resourceIpGroupCidr(),
		"azurerm_local_network_gateway":                    resourceLocalNetworkGateway(),
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

		"azurerm_network_packet_capture":                    resourceNetworkPacketCapture(),
		"azurerm_network_profile":                           resourceNetworkProfile(),
		"azurerm_point_to_site_vpn_gateway":                 resourcePointToSiteVPNGateway(),
		"azurerm_private_endpoint":                          resourcePrivateEndpoint(),
		"azurerm_private_link_service":                      resourcePrivateLinkService(),
		"azurerm_public_ip":                                 resourcePublicIp(),
		"azurerm_public_ip_prefix":                          resourcePublicIpPrefix(),
		"azurerm_network_security_group":                    resourceNetworkSecurityGroup(),
		"azurerm_network_security_rule":                     resourceNetworkSecurityRule(),
		"azurerm_network_watcher_flow_log":                  resourceNetworkWatcherFlowLog(),
		"azurerm_network_watcher":                           resourceNetworkWatcher(),
		"azurerm_route_filter":                              resourceRouteFilter(),
		"azurerm_route_table":                               resourceRouteTable(),
		"azurerm_route":                                     resourceRoute(),
		"azurerm_route_server":                              resourceRouteServer(),
		"azurerm_route_server_bgp_connection":               resourceRouteServerBgpConnection(),
		"azurerm_virtual_hub_security_partner_provider":     resourceVirtualHubSecurityPartnerProvider(),
		"azurerm_subnet_service_endpoint_storage_policy":    resourceSubnetServiceEndpointStoragePolicy(),
		"azurerm_subnet_network_security_group_association": resourceSubnetNetworkSecurityGroupAssociation(),
		"azurerm_subnet_route_table_association":            resourceSubnetRouteTableAssociation(),
		"azurerm_subnet_nat_gateway_association":            resourceSubnetNatGatewayAssociation(),
		"azurerm_subnet":                                    resourceSubnet(),
		"azurerm_virtual_hub":                               resourceVirtualHub(),
		"azurerm_virtual_hub_bgp_connection":                resourceVirtualHubBgpConnection(),
		"azurerm_virtual_hub_connection":                    resourceVirtualHubConnection(),
		"azurerm_virtual_hub_ip":                            resourceVirtualHubIP(),
		"azurerm_virtual_hub_route_table":                   resourceVirtualHubRouteTable(),
		"azurerm_virtual_hub_route_table_route":             resourceVirtualHubRouteTableRoute(),
		"azurerm_virtual_machine_packet_capture":            resourceVirtualMachinePacketCapture(),
		"azurerm_virtual_machine_scale_set_packet_capture":  resourceVirtualMachineScaleSetPacketCapture(),
		"azurerm_virtual_network_dns_servers":               resourceVirtualNetworkDnsServers(),
		"azurerm_virtual_network_gateway_connection":        resourceVirtualNetworkGatewayConnection(),
		"azurerm_virtual_network_gateway_nat_rule":          resourceVirtualNetworkGatewayNatRule(),
		"azurerm_virtual_network_gateway":                   resourceVirtualNetworkGateway(),
		"azurerm_virtual_network_peering":                   resourceVirtualNetworkPeering(),
		"azurerm_virtual_network":                           resourceVirtualNetwork(),
		"azurerm_virtual_wan":                               resourceVirtualWan(),
		"azurerm_vpn_gateway":                               resourceVPNGateway(),
		"azurerm_vpn_gateway_connection":                    resourceVPNGatewayConnection(),
		"azurerm_vpn_gateway_nat_rule":                      resourceVPNGatewayNatRule(),
		"azurerm_vpn_server_configuration":                  resourceVPNServerConfiguration(),
		"azurerm_vpn_server_configuration_policy_group":     resourceVPNServerConfigurationPolicyGroup(),
		"azurerm_vpn_site":                                  resourceVpnSite(),
		"azurerm_web_application_firewall_policy":           resourceWebApplicationFirewallPolicy(),
	}
}
