package network

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Network"
}

// WebsiteCategories returns a list of categories which can be used for the sidebar
func (r Registration) WebsiteCategories() []string {
	return []string{
		"Network",
	}
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_application_security_group":                dataSourceApplicationSecurityGroup(),
		"azurerm_express_route_circuit":                     dataSourceExpressRouteCircuit(),
		"azurerm_ip_group":                                  dataSourceIpGroup(),
		"azurerm_nat_gateway":                               dataSourceNatGateway(),
		"azurerm_network_ddos_protection_plan":              dataSourceNetworkDDoSProtectionPlan(),
		"azurerm_network_interface":                         dataSourceNetworkInterface(),
		"azurerm_network_security_group":                    dataSourceArmNetworkSecurityGroup(),
		"azurerm_network_watcher":                           dataSourceArmNetworkWatcher(),
		"azurerm_private_endpoint_connection":               dataSourceArmPrivateEndpointConnection(),
		"azurerm_private_link_service":                      dataSourceArmPrivateLinkService(),
		"azurerm_private_link_service_endpoint_connections": dataSourceArmPrivateLinkServiceEndpointConnections(),
		"azurerm_public_ip":                                 dataSourceArmPublicIP(),
		"azurerm_public_ips":                                dataSourceArmPublicIPs(),
		"azurerm_public_ip_prefix":                          dataSourceArmPublicIpPrefix(),
		"azurerm_route_filter":                              dataSourceRouteFilter(),
		"azurerm_route_table":                               dataSourceRouteTable(),
		"azurerm_network_service_tags":                      dataSourceNetworkServiceTags(),
		"azurerm_subnet":                                    dataSourceArmSubnet(),
		"azurerm_virtual_hub":                               dataSourceArmVirtualHub(),
		"azurerm_virtual_network_gateway":                   dataSourceArmVirtualNetworkGateway(),
		"azurerm_virtual_network_gateway_connection":        dataSourceArmVirtualNetworkGatewayConnection(),
		"azurerm_virtual_network":                           dataSourceArmVirtualNetwork(),
		"azurerm_web_application_firewall_policy":           dataArmWebApplicationFirewallPolicy(),
		"azurerm_virtual_wan":                               dataSourceArmVirtualWan(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_application_gateway":                 resourceApplicationGateway(),
		"azurerm_application_security_group":          resourceApplicationSecurityGroup(),
		"azurerm_bastion_host":                        resourceBastionHost(),
		"azurerm_express_route_circuit_authorization": resourceExpressRouteCircuitAuthorization(),
		"azurerm_express_route_circuit_peering":       resourceExpressRouteCircuitPeering(),
		"azurerm_express_route_circuit":               resourceExpressRouteCircuit(),
		"azurerm_express_route_gateway":               resourceExpressRouteGateway(),
		"azurerm_ip_group":                            resourceIpGroup(),
		"azurerm_local_network_gateway":               resourceLocalNetworkGateway(),
		"azurerm_nat_gateway":                         resourceNatGateway(),
		"azurerm_network_connection_monitor":          resourceNetworkConnectionMonitor(),
		"azurerm_network_ddos_protection_plan":        resourceNetworkDDoSProtectionPlan(),
		"azurerm_network_interface":                   resourceNetworkInterface(),
		"azurerm_network_interface_application_gateway_backend_address_pool_association": resourceNetworkInterfaceApplicationGatewayBackendAddressPoolAssociation(),
		"azurerm_network_interface_application_security_group_association":               resourceNetworkInterfaceApplicationSecurityGroupAssociation(),
		"azurerm_network_interface_backend_address_pool_association":                     resourceNetworkInterfaceBackendAddressPoolAssociation(),
		"azurerm_network_interface_nat_rule_association":                                 resourceNetworkInterfaceNatRuleAssociation(),
		"azurerm_network_interface_security_group_association":                           resourceNetworkInterfaceSecurityGroupAssociation(),
		"azurerm_network_packet_capture":                                                 resourceArmNetworkPacketCapture(),
		"azurerm_network_profile":                                                        resourceArmNetworkProfile(),
		"azurerm_packet_capture":                                                         resourceArmPacketCapture(),
		"azurerm_point_to_site_vpn_gateway":                                              resourceArmPointToSiteVPNGateway(),
		"azurerm_private_endpoint":                                                       resourceArmPrivateEndpoint(),
		"azurerm_private_link_service":                                                   resourceArmPrivateLinkService(),
		"azurerm_public_ip":                                                              resourceArmPublicIp(),
		"azurerm_nat_gateway_public_ip_association":                                      resourceNATGatewayPublicIpAssociation(),
		"azurerm_public_ip_prefix":                                                       resourceArmPublicIpPrefix(),
		"azurerm_network_security_group":                                                 resourceArmNetworkSecurityGroup(),
		"azurerm_network_security_rule":                                                  resourceArmNetworkSecurityRule(),
		"azurerm_network_watcher_flow_log":                                               resourceArmNetworkWatcherFlowLog(),
		"azurerm_network_watcher":                                                        resourceArmNetworkWatcher(),
		"azurerm_route_filter":                                                           resourceRouteFilter(),
		"azurerm_route_table":                                                            resourceRouteTable(),
		"azurerm_route":                                                                  resourceRoute(),
		"azurerm_virtual_hub_security_partner_provider":                                  resourceArmVirtualHubSecurityPartnerProvider(),
		"azurerm_subnet_service_endpoint_storage_policy":                                 resourceArmSubnetServiceEndpointStoragePolicy(),
		"azurerm_subnet_network_security_group_association":                              resourceArmSubnetNetworkSecurityGroupAssociation(),
		"azurerm_subnet_route_table_association":                                         resourceArmSubnetRouteTableAssociation(),
		"azurerm_subnet_nat_gateway_association":                                         resourceArmSubnetNatGatewayAssociation(),
		"azurerm_subnet":                                                                 resourceArmSubnet(),
		"azurerm_virtual_hub":                                                            resourceArmVirtualHub(),
		"azurerm_virtual_hub_bgp_connection":                                             resourceArmVirtualHubBgpConnection(),
		"azurerm_virtual_hub_connection":                                                 resourceArmVirtualHubConnection(),
		"azurerm_virtual_hub_ip":                                                         resourceArmVirtualHubIP(),
		"azurerm_virtual_hub_route_table":                                                resourceArmVirtualHubRouteTable(),
		"azurerm_virtual_network_gateway_connection":                                     resourceArmVirtualNetworkGatewayConnection(),
		"azurerm_virtual_network_gateway":                                                resourceArmVirtualNetworkGateway(),
		"azurerm_virtual_network_peering":                                                resourceArmVirtualNetworkPeering(),
		"azurerm_virtual_network":                                                        resourceArmVirtualNetwork(),
		"azurerm_virtual_wan":                                                            resourceArmVirtualWan(),
		"azurerm_vpn_gateway":                                                            resourceArmVPNGateway(),
		"azurerm_vpn_gateway_connection":                                                 resourceArmVPNGatewayConnection(),
		"azurerm_vpn_server_configuration":                                               resourceArmVPNServerConfiguration(),
		"azurerm_vpn_site":                                                               resourceArmVpnSite(),
		"azurerm_web_application_firewall_policy":                                        resourceArmWebApplicationFirewallPolicy(),
	}
}
