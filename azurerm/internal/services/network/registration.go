package network

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type Registration struct{}

// Name is the name of this Service
func (r Registration) Name() string {
	return "Network"
}

// SupportedDataSources returns the supported Data Sources supported by this Service
func (r Registration) SupportedDataSources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_application_security_group":                dataSourceArmApplicationSecurityGroup(),
		"azurerm_express_route_circuit":                     dataSourceArmExpressRouteCircuit(),
		"azurerm_firewall":                                  dataSourceArmFirewall(),
		"azurerm_lb":                                        dataSourceArmLoadBalancer(),
		"azurerm_lb_backend_address_pool":                   dataSourceArmLoadBalancerBackendAddressPool(),
		"azurerm_nat_gateway":                               dataSourceArmNatGateway(),
		"azurerm_network_ddos_protection_plan":              dataSourceNetworkDDoSProtectionPlan(),
		"azurerm_network_interface":                         dataSourceArmNetworkInterface(),
		"azurerm_network_security_group":                    dataSourceArmNetworkSecurityGroup(),
		"azurerm_network_watcher":                           dataSourceArmNetworkWatcher(),
		"azurerm_private_link_endpoint_connection":          dataSourceArmPrivateLinkEndpointConnection(),
		"azurerm_private_endpoint_connection":               dataSourceArmPrivateEndpointConnection(),
		"azurerm_private_link_service":                      dataSourceArmPrivateLinkService(),
		"azurerm_private_link_service_endpoint_connections": dataSourceArmPrivateLinkServiceEndpointConnections(),
		"azurerm_public_ip":                                 dataSourceArmPublicIP(),
		"azurerm_public_ips":                                dataSourceArmPublicIPs(),
		"azurerm_public_ip_prefix":                          dataSourceArmPublicIpPrefix(),
		"azurerm_route_table":                               dataSourceArmRouteTable(),
		"azurerm_subnet":                                    dataSourceArmSubnet(),
		"azurerm_virtual_hub":                               dataSourceArmVirtualHub(),
		"azurerm_virtual_network_gateway":                   dataSourceArmVirtualNetworkGateway(),
		"azurerm_virtual_network_gateway_connection":        dataSourceArmVirtualNetworkGatewayConnection(),
		"azurerm_virtual_network":                           dataSourceArmVirtualNetwork(),
	}
}

// SupportedResources returns the supported Resources supported by this Service
func (r Registration) SupportedResources() map[string]*schema.Resource {
	return map[string]*schema.Resource{
		"azurerm_application_gateway":                  resourceArmApplicationGateway(),
		"azurerm_application_security_group":           resourceArmApplicationSecurityGroup(),
		"azurerm_bastion_host":                         resourceArmBastionHost(),
		"azurerm_connection_monitor":                   resourceArmConnectionMonitor(),
		"azurerm_ddos_protection_plan":                 resourceArmDDoSProtectionPlan(),
		"azurerm_express_route_circuit_authorization":  resourceArmExpressRouteCircuitAuthorization(),
		"azurerm_express_route_circuit_peering":        resourceArmExpressRouteCircuitPeering(),
		"azurerm_express_route_circuit":                resourceArmExpressRouteCircuit(),
		"azurerm_firewall_application_rule_collection": resourceArmFirewallApplicationRuleCollection(),
		"azurerm_firewall_nat_rule_collection":         resourceArmFirewallNatRuleCollection(),
		"azurerm_firewall_network_rule_collection":     resourceArmFirewallNetworkRuleCollection(),
		"azurerm_firewall":                             resourceArmFirewall(),
		"azurerm_local_network_gateway":                resourceArmLocalNetworkGateway(),
		"azurerm_lb_backend_address_pool":              resourceArmLoadBalancerBackendAddressPool(),
		"azurerm_lb_nat_pool":                          resourceArmLoadBalancerNatPool(),
		"azurerm_lb_nat_rule":                          resourceArmLoadBalancerNatRule(),
		"azurerm_lb_probe":                             resourceArmLoadBalancerProbe(),
		"azurerm_lb_outbound_rule":                     resourceArmLoadBalancerOutboundRule(),
		"azurerm_lb_rule":                              resourceArmLoadBalancerRule(),
		"azurerm_lb":                                   resourceArmLoadBalancer(),
		"azurerm_nat_gateway":                          resourceArmNatGateway(),
		"azurerm_network_connection_monitor":           resourceArmNetworkConnectionMonitor(),
		"azurerm_network_ddos_protection_plan":         resourceArmNetworkDDoSProtectionPlan(),
		"azurerm_network_interface":                    resourceArmNetworkInterface(),
		"azurerm_network_interface_application_gateway_backend_address_pool_association": resourceArmNetworkInterfaceApplicationGatewayBackendAddressPoolAssociation(),
		"azurerm_network_interface_application_security_group_association":               resourceArmNetworkInterfaceApplicationSecurityGroupAssociation(),
		"azurerm_network_interface_backend_address_pool_association":                     resourceArmNetworkInterfaceBackendAddressPoolAssociation(),
		"azurerm_network_interface_nat_rule_association":                                 resourceArmNetworkInterfaceNatRuleAssociation(),
		"azurerm_network_packet_capture":                                                 resourceArmNetworkPacketCapture(),
		"azurerm_network_profile":                                                        resourceArmNetworkProfile(),
		"azurerm_packet_capture":                                                         resourceArmPacketCapture(),
		"azurerm_point_to_site_vpn_gateway":                                              resourceArmPointToSiteVPNGateway(),
		"azurerm_private_link_endpoint":                                                  resourceArmPrivateLinkEndpoint(),
		"azurerm_private_endpoint":                                                       resourceArmPrivateEndpoint(),
		"azurerm_private_link_service":                                                   resourceArmPrivateLinkService(),
		"azurerm_public_ip":                                                              resourceArmPublicIp(),
		"azurerm_public_ip_prefix":                                                       resourceArmPublicIpPrefix(),
		"azurerm_network_security_group":                                                 resourceArmNetworkSecurityGroup(),
		"azurerm_network_security_rule":                                                  resourceArmNetworkSecurityRule(),
		"azurerm_network_watcher_flow_log":                                               resourceArmNetworkWatcherFlowLog(),
		"azurerm_network_watcher":                                                        resourceArmNetworkWatcher(),
		"azurerm_route_table":                                                            resourceArmRouteTable(),
		"azurerm_route":                                                                  resourceArmRoute(),
		"azurerm_subnet_network_security_group_association":                              resourceArmSubnetNetworkSecurityGroupAssociation(),
		"azurerm_subnet_route_table_association":                                         resourceArmSubnetRouteTableAssociation(),
		"azurerm_subnet_nat_gateway_association":                                         resourceArmSubnetNatGatewayAssociation(),
		"azurerm_subnet":                                                                 resourceArmSubnet(),
		"azurerm_virtual_hub":                                                            resourceArmVirtualHub(),
		"azurerm_virtual_network_gateway_connection":                                     resourceArmVirtualNetworkGatewayConnection(),
		"azurerm_virtual_network_gateway":                                                resourceArmVirtualNetworkGateway(),
		"azurerm_virtual_network_peering":                                                resourceArmVirtualNetworkPeering(),
		"azurerm_virtual_network":                                                        resourceArmVirtualNetwork(),
		"azurerm_virtual_wan":                                                            resourceArmVirtualWan(),
		"azurerm_vpn_gateway":                                                            resourceArmVPNGateway(),
		"azurerm_vpn_server_configuration":                                               resourceArmVPNServerConfiguration(),
		"azurerm_web_application_firewall_policy":                                        resourceArmWebApplicationFirewallPolicy(),
	}
}
